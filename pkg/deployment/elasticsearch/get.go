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
	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_elasticsearch"
	"github.com/elastic/cloud-sdk-go/pkg/models"

	"github.com/elastic/ecctl/pkg/util"
)

var (
	// systemAlerts controls the number of system alerts that
	// are shown by default
	systemAlerts int64 = 5
)

// GetClusterParams contains parameters used to fetch cluster's data
type GetClusterParams struct {
	util.ClusterParams
	Metadata, PlanDefaults, Plans, Logs, Settings bool
}

// GetCluster fetches the cluster information with the settings
func GetCluster(params GetClusterParams) (*models.ElasticsearchClusterInfo, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.V1API.ClustersElasticsearch.GetEsCluster(
		clusters_elasticsearch.NewGetEsClusterParams().
			WithClusterID(params.ClusterID).
			WithShowMetadata(&params.Metadata).
			WithShowPlanDefaults(&params.PlanDefaults).
			WithShowPlans(&params.Plans).
			WithShowPlanLogs(&params.Logs).
			WithShowSystemAlerts(&systemAlerts).
			WithShowSettings(&params.Settings),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}
