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
	"github.com/go-openapi/strfmt"

	"github.com/elastic/ecctl/pkg/util"
)

// SearchClusterParams contains parameters used to search cluster's data using Query DSL
type SearchClusterParams struct {
	Request models.SearchRequest
	*api.API
}

// Validate is the implementation for the ecctl.Validator interface
func (params SearchClusterParams) Validate() error {
	if params.API == nil {
		return util.ErrAPIReq
	}

	return params.Request.Validate(strfmt.Default)
}

// SearchClusters searches all the elasticsearch clusters using Query DSL
func SearchClusters(params SearchClusterParams) (*models.ElasticsearchClustersInfo, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.V1API.ClustersElasticsearch.SearchEsClusters(
		clusters_elasticsearch.NewSearchEsClustersParams().WithBody(&params.Request),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}
	return res.Payload, nil
}
