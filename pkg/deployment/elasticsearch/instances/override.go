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

package instances

import (
	"fmt"
	"strings"

	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_elasticsearch"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/elastic/cloud-sdk-go/pkg/util/slice"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/elastic/ecctl/pkg/deployment/elasticsearch"
	"github.com/elastic/ecctl/pkg/util"
)

const (
	// maxInstanceCapacity is used to determine what the maximum memory size
	// for an instance is.
	maxInstanceCapacity = 65536
)

var (
	// errMaxInstanceCapacityReached is thrown when the instance size goes over
	// the defined limit in maxInstanceCapacity.
	errMaxInstanceCapacityReached = fmt.Errorf("capacity must not exceed %d", maxInstanceCapacity)
)

// OverrideCapacityParams is used by OverrideCapacity
type OverrideCapacityParams struct {
	util.ClusterParams
	Instances   []string
	Value       uint64
	BoostFactor uint8
	Reset       bool
}

// OverrideRequest is used in a SetEsClusterInstancesSettingsOverrides call
// to override any ES settings.
type OverrideRequest struct {
	Capacity  int32
	Instances []string
}

// OverrideResponse is returned by OverrideCapacity
type OverrideResponse OverrideRequest

// TopologySize defines an instance configuration ID with the current size.
type TopologySize struct {
	ID   string
	Size int32
}

// Validate ensures that the parameters are usable by the consuming function
func (params OverrideCapacityParams) Validate() error {
	var err = new(multierror.Error)
	err = multierror.Append(params.ClusterParams.Validate())

	if params.BoostFactor <= 1 && params.Value < 1024 && !params.Reset {
		err = multierror.Append(errors.New("you must specify an absolute value greater than 1024 or a boost factor greater than 1 or reset to the original value"), err)
	}

	if params.BoostFactor > 1 && params.Value > 0 && !params.Reset {
		err = multierror.Append(errors.New("you cannot specify at the same time an absolute value and a boost factor"), err)
	}

	if (params.BoostFactor > 1 || params.Value > 0) && params.Reset {
		err = multierror.Append(errors.New("you cannot specify at the same time an absolute value or a boost factor and set the reset flag to true"), err)
	}

	return err.ErrorOrNil()
}

// OverrideCapacity performs a series of capacity overrides to a
// series of Elasticsearch cluster instances. If no Instances are specified, the
// whole cluster's instances are boosted from the passed multiplier.
// Takes into account different topologies inside the same cluster, individually
// boosting different capacities.
func OverrideCapacity(params OverrideCapacityParams) ([]OverrideResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	cluster, err := elasticsearch.GetCluster(elasticsearch.GetClusterParams{
		ClusterParams: params.ClusterParams,
		Plans:         true,
		PlanDefaults:  true,
	})

	if err != nil {
		return nil, err
	}

	var overrides []OverrideRequest
	for _, topology := range cluster.PlanInfo.Current.Plan.ClusterTopology {
		size, err := GetSizeFromTopology(topology)
		if err != nil {
			continue
		}

		override := NewOverrideRequest(cluster.Topology.Instances,
			size, newSize(params, size.Size), params.Instances,
		)

		if override.Capacity > maxInstanceCapacity {
			return nil, errors.Wrap(errMaxInstanceCapacityReached, strings.Join(override.Instances, " "))
		}

		if !slice.IsEmpty(override.Instances) {
			overrides = append(overrides, override)
		}
	}

	var merr = new(multierror.Error)
	var response []OverrideResponse
	for _, override := range overrides {
		res, err := params.API.V1API.ClustersElasticsearch.SetEsClusterInstancesSettingsOverrides(
			clusters_elasticsearch.NewSetEsClusterInstancesSettingsOverridesParams().
				WithClusterID(params.ClusterID).
				WithInstanceIds(override.Instances).
				WithRestartAfterUpdate(ec.Bool(true)).
				WithBody(&models.ElasticsearchClusterInstanceSettingsOverrides{
					InstanceCapacity: override.Capacity,
				}),
			params.AuthWriter,
		)
		if err != nil {
			merr = multierror.Append(merr, err)
		} else {
			response = append(response, OverrideResponse{
				Instances: override.Instances,
				Capacity:  res.Payload.InstanceCapacity,
			})
		}
	}

	return response, merr.ErrorOrNil()
}

func newSize(params OverrideCapacityParams, size int32) int32 {
	var newSize int32
	if params.Reset {
		return size
	}
	if params.BoostFactor >= 1 {
		return size * int32(params.BoostFactor)
	}
	if params.Value > 0 {
		return int32(params.Value)
	}
	return newSize
}

// GetSizeFromTopology attempts to obtain the defined topology size of an
// Elasticsearch cluster, returning an error if the size cannot be obtained.
func GetSizeFromTopology(topology *models.ElasticsearchClusterTopologyElement) (TopologySize, error) {
	var size = TopologySize{
		// Handles legacy case.
		ID:   strings.Replace(topology.NodeConfiguration, "legacy", "classic", 1),
		Size: topology.MemoryPerNode,
	}
	// ECE case
	if size.ID == "default" {
		size.ID = strings.Replace(topology.NodeConfiguration, "default", "data.default", 1)
	}
	if size.ID == "" {
		size.ID = topology.InstanceConfigurationID
		if size.Size == 0 {
			if *topology.Size.Value == 0 {
				return TopologySize{}, fmt.Errorf("couldn't obtain size from %s", size.ID)
			}
			size.Size = *topology.Size.Value
		}
	}

	return size, nil
}

// NewOverrideRequest receives a list of instances, the Topology size and the new value to be used
// factor and returns the Override request that has to be sent to the API.
// Any instances that don't match the TopologySize will be ignored.
func NewOverrideRequest(instances []*models.ClusterInstanceInfo, size TopologySize, newValue int32, filter []string) OverrideRequest {
	var override = OverrideRequest{Capacity: newValue}
	for _, instance := range instances {
		legacyMatches := instance.InstanceConfiguration.Name == size.ID
		instanceConfigMatches := instance.InstanceConfiguration.ID == size.ID
		if legacyMatches || instanceConfigMatches {
			if slice.HasString(filter, *instance.InstanceName) {
				override.Instances = append(override.Instances, *instance.InstanceName)
			}
		}
	}
	return override
}
