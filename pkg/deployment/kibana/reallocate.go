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

package kibana

import (
	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_kibana"
	"github.com/elastic/cloud-sdk-go/pkg/plan/planutil"

	"github.com/elastic/ecctl/pkg/deployment/deputil"
	"github.com/elastic/ecctl/pkg/util"
)

// ReallocateParams is used by Reallocate as a config struct
type ReallocateParams struct {
	DeploymentParams
	Instances []string
}

// Reallocate will reallocate the Kibana instance instances, if no
// Instances are specified, all of the instances will be moved.
func Reallocate(params ReallocateParams) error {
	if err := params.Validate(); err != nil {
		return err
	}
	if len(params.Instances) == 0 || params.Instances[0] == "" {
		params.Instances = make([]string, 0)

		kibana, err := Get(ClusterParams{
			DeploymentParams: DeploymentParams{
				API: params.API,
				ID:  params.ID},
			QueryParams: deputil.QueryParams{
				ShowMetadata:     false,
				ShowPlanDefaults: false,
				ShowPlans:        false,
				ShowPlanLogs:     false,
			}})

		if err != nil {
			return err
		}

		for _, i := range kibana.Topology.Instances {
			params.Instances = append(params.Instances, *i.InstanceName)
		}
	}

	if _, err := params.API.V1API.ClustersKibana.MoveKibanaClusterInstances(
		clusters_kibana.NewMoveKibanaClusterInstancesParams().
			WithClusterID(params.ID).
			WithInstanceIds(params.Instances),
		params.API.AuthWriter,
	); err != nil {
		return api.UnwrapError(err)
	}

	return planutil.TrackChange(util.SetClusterTracking(
		params.TrackChangeParams, params.ID, util.Kibana,
	))
}
