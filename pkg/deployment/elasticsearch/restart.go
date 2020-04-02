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
	"time"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_elasticsearch"
	"github.com/elastic/cloud-sdk-go/pkg/plan/planutil"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"

	depplanutil "github.com/elastic/ecctl/pkg/deployment/planutil"
	"github.com/elastic/ecctl/pkg/util"
)

// RestartClusterParams is used to restart clusters with specific parameters / options
type RestartClusterParams struct {
	util.ClusterParams

	ShardInitWaitTime time.Duration

	RollingByName, RollingByZone, SkipSnapshot, RestoreSnapshot bool

	Track bool
	planutil.TrackChangeParams
}

// Validate ensures the parameters are usable by ShutdownCluster.
func (params RestartClusterParams) Validate() error {
	var err = new(multierror.Error)

	err = multierror.Append(err, params.ClusterParams.Validate())

	return err.ErrorOrNil()
}

// RestartCluster restarts an Elasticsearch cluster, if rolling is set to true
// the restart will happen on one logical zone at a time, meaning the cluster won't
// experience downtime, else, all of the nodes will be restarted at the same time
// causing downtime in the cluster.
func RestartCluster(params RestartClusterParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	groupingStrategy := depplanutil.AllGroupAttribute
	if params.RollingByName {
		groupingStrategy = depplanutil.NameGroupAttribute
	}
	if params.RollingByZone {
		groupingStrategy = depplanutil.ZoneGroupAttribute
	}

	shardInitWaitTime := int64(params.ShardInitWaitTime.Seconds())

	if _, err := params.V1API.ClustersElasticsearch.RestartEsCluster(
		clusters_elasticsearch.NewRestartEsClusterParams().
			WithClusterID(params.ClusterID).
			WithGroupAttribute(&groupingStrategy).
			WithSkipSnapshot(ec.Bool(params.SkipSnapshot)).
			WithRestoreSnapshot(ec.Bool(params.RestoreSnapshot)).
			WithShardInitWaitTime(&shardInitWaitTime),
		params.AuthWriter,
	); err != nil {
		return api.UnwrapError(err)
	}

	if !params.Track {
		return nil
	}

	return planutil.TrackChange(util.SetClusterTracking(
		params.TrackChangeParams, params.ClusterID, util.Elasticsearch,
	))
}
