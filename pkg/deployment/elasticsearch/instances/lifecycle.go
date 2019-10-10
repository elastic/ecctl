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
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_elasticsearch"

	"github.com/elastic/ecctl/pkg/util"
)

// Params is consumed both vby by Resume and Pause.
type Params struct {
	util.ClusterParams
	Instances []string
}

// Resume starts specific cluster instances
func Resume(params Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	return util.ReturnErrOnly(
		params.V1API.ClustersElasticsearch.StartEsClusterInstances(
			clusters_elasticsearch.NewStartEsClusterInstancesParams().
				WithClusterID(params.ClusterID).
				WithInstanceIds(params.Instances),
			params.AuthWriter,
		),
	)
}

// Pause stops specific cluster instances
func Pause(params Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	return util.ReturnErrOnly(
		params.V1API.ClustersElasticsearch.StopEsClusterInstances(
			clusters_elasticsearch.NewStopEsClusterInstancesParams().
				WithClusterID(params.ClusterID).
				WithInstanceIds(params.Instances),
			params.AuthWriter,
		),
	)
}
