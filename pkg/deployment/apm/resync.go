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

package apm

import (
	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_apm"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment/deputil"
	"github.com/elastic/ecctl/pkg/util"
)

// ResyncParams represents parameters used to resynchronize an APM cluster.
type ResyncParams struct {
	// API embeds API clients.
	*api.API

	// ClusterID contains an APM cluster ID.
	ClusterID string
}

// Validate ensures that the parameters such as APM cluster ID and
// the API client object are valid before the API call is made.
func (params ResyncParams) Validate() error {
	var err = multierror.Append(new(multierror.Error),
		deputil.ValidateParams(&params),
	)
	return err.ErrorOrNil()
}

// Resync forces indexer to immediately resynchronize the search index
// and cache for a given APM cluster.
func Resync(params ResyncParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	return util.ReturnErrOnly(
		params.API.V1API.ClustersApm.ResyncApmCluster(
			clusters_apm.NewResyncApmClusterParams().
				WithClusterID(params.ClusterID),
			params.AuthWriter,
		),
	)
}
