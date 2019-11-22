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
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

// GetKeystoreParams is consumed by GetKeystore.
type GetKeystoreParams struct {
	*api.API
	ClusterID string
}

// Validate ensures the parameters are usabl by the consuming function.
func (params GetKeystoreParams) Validate() error {
	var merr = new(multierror.Error)

	if params.API == nil {
		merr = multierror.Append(merr, util.ErrAPIReq)
	}

	if len(params.ClusterID) != 32 {
		merr = multierror.Append(merr, util.ErrClusterLength)
	}

	return merr.ErrorOrNil()
}

// GetKeystore obtains the Elasticsearch keystore values from a ClusterID.
func GetKeystore(params GetKeystoreParams) (*models.KeystoreContents, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.V1API.ClustersElasticsearch.GetEsClusterKeystore(
		clusters_elasticsearch.NewGetEsClusterKeystoreParams().
			WithClusterID(params.ClusterID),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}
