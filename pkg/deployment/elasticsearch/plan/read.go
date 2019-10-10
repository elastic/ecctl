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

package plan

import (
	"errors"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_elasticsearch"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
)

// Get returns either the current cluster plan or the pending plan.
func Get(params GetParams) (*models.ElasticsearchClusterPlanInfo, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.V1API.ClustersElasticsearch.GetEsClusterPlanActivity(
		clusters_elasticsearch.NewGetEsClusterPlanActivityParams().
			WithClusterID(params.ClusterID).
			WithShowPlanDefaults(ec.Bool(true)),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	if !params.Pending {
		return res.Payload.Current, nil
	}
	if res.Payload.Pending == nil {
		return nil, errors.New("no pending plan")
	}
	return res.Payload.Pending, nil
}

// GetHistory returns the historic plan list
func GetHistory(params GetHistoryParams) ([]*models.ElasticsearchClusterPlanInfo, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.V1API.ClustersElasticsearch.GetEsClusterPlanActivity(
		clusters_elasticsearch.NewGetEsClusterPlanActivityParams().
			WithClusterID(params.ClusterID).
			WithShowPlanDefaults(ec.Bool(true)),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload.History, nil
}
