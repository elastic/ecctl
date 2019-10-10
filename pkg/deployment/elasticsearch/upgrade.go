// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package elasticsearch

import (
	"errors"
	"strings"

	"github.com/blang/semver"
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_elasticsearch"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment/planutil"
	"github.com/elastic/ecctl/pkg/util"
)

const (
	// defaultMultiplier is the default memory / storage ratio
	// used to establish the smallInstanceDiskTreshold
	defaultMultiplier = int32(24)

	// highMemoryThreshold is used to determine if an instance
	// is suffering from high memory pressure
	highMemoryThreshold = 65

	// hardDiskThreshold is an upper bound used to determine
	// wether an instance has more than 100GB of disk used
	hardDiskThreshold = int64(100 * 1024)
)

// UpgradeParams is consumed by Upgrade.
type UpgradeParams struct {
	util.ClusterParams
	util.TrackParams

	Version string
}

// Validate ensures that the parameters are usable by Upgrade.
func (params UpgradeParams) Validate() error {
	var merr = new(multierror.Error)
	merr = multierror.Append(merr, params.ClusterParams.Validate())
	merr = multierror.Append(merr, params.TrackParams.Validate())

	if _, e := semver.Parse(params.Version); e != nil {
		merr = multierror.Append(merr, errors.New(strings.ToLower(e.Error())))
	}

	return merr.ErrorOrNil()
}

// Upgrade the Elasticsearch cluster to the specified version.
func Upgrade(params UpgradeParams) (*models.ClusterCrudResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	cluster, err := GetCluster(GetClusterParams{
		ClusterParams: util.ClusterParams{
			API:       params.API,
			ClusterID: params.ClusterID,
		},
		PlanDefaults: true,
		Plans:        true,
	})
	if err != nil {
		return nil, err
	}

	if cluster.PlanInfo == nil || cluster.PlanInfo.Current == nil {
		return nil, errors.New("cluster has no plan information")
	}

	res, res2, err := params.V1API.ClustersElasticsearch.UpdateEsClusterPlan(
		clusters_elasticsearch.NewUpdateEsClusterPlanParams().
			WithClusterID(params.ClusterID).
			WithBody(
				getUpgradePlan(cluster, params.Version),
			),
		params.AuthWriter,
	)

	return util.ParseCUResponse(util.ParseCUResponseParams{
		API:            params.API,
		UpdateResponse: res2,
		CreateResponse: res,
		Err:            err,
		TrackParams:    params.TrackParams,
	})
}

func getUpgradePlan(cluster *models.ElasticsearchClusterInfo, version string) *models.ElasticsearchClusterPlan {
	newPlan := cluster.PlanInfo.Current.Plan
	currentVersion := newPlan.Elasticsearch.Version
	newPlan.Elasticsearch.Version = version

	// Default Plan for Minor upgrades
	newPlan.Transient = &models.TransientElasticsearchPlanConfiguration{
		Strategy: planutil.DefaultPlanStrategy,
	}

	previousVer := semver.MustParse(currentVersion)
	nextVersion := semver.MustParse(version)

	// The only valid Strategy for Major version upgrades is all performing the
	// change to all instances at the same time, resulting in downtime
	if previousVer.Major != nextVersion.Major {
		newPlan.Transient.Strategy = planutil.MajorUpgradeStrategy
		// This part ensures that upgrades from 2.x to 5.x succeed. The setting
		// is enforced from the API to be unused, else it would error.
		if newPlan.Elasticsearch.SystemSettings != nil {
			if newPlan.Elasticsearch.SystemSettings.DefaultShardsPerIndex > 0 {
				newPlan.Elasticsearch.SystemSettings.DefaultShardsPerIndex = 0
			}
		}

		// Elasticsearch 6.x deprecates a few settings in scripting
		if nextVersion.Major == 6 {
			// This is so ugly but necessary it seems until we have the _upgrade endpoint in the API
			scriptingSettings := newPlan.Elasticsearch.SystemSettings.Scripting
			scriptingSettings.File, scriptingSettings.Stored.SandboxMode, scriptingSettings.Inline.SandboxMode = nil, nil, nil
			scriptingSettings.ExpressionsEnabled, scriptingSettings.MustacheEnabled, scriptingSettings.PainlessEnabled = nil, nil, nil
		}

		return newPlan
	}

	// Strategy for Minor upgrades on clusters that have instances with memory or disk pressure.
	// One instance at a time. This is to help reduce the overhead that data migrate causes.
	for _, instance := range cluster.Topology.Instances {
		// The criteria to decide wether the minor upgrade will be performed in place is:
		// * If the instance has at least 1/4 of the disk used.
		// * OR the instance has a higher memory pressure than `highMemoryThreshold`
		// * OR the instance meets the hard limit of 100GB of disk usage.
		// The upgrade will be performed in place one instance at a time. In practical terms:
		// Clusters with instances of >16GB will always default to the Hard limit of 100GB
		// with the defaultMultiplier disk ratio (16 * 24 / 4 = 96).
		smallInstanceDiskTreshold := int64(
			float32(*instance.Memory.InstanceCapacity*defaultMultiplier) * float32(0.25),
		)
		percentageDiskThresholdExceeded := *instance.Disk.DiskSpaceUsed >= smallInstanceDiskTreshold
		hardDiskThresholdExceeded := *instance.Disk.DiskSpaceUsed >= hardDiskThreshold
		memoryPressureThresholdExceeded := instance.Memory.MemoryPressure >= highMemoryThreshold

		// Asserts wether it complies with the criteria
		needsRollingStrategyGroupedByName := memoryPressureThresholdExceeded || percentageDiskThresholdExceeded || hardDiskThresholdExceeded
		if needsRollingStrategyGroupedByName {
			newPlan.Transient.Strategy = planutil.RollingByNameStrategy
		}
	}
	return newPlan
}
