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

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_elasticsearch"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

// SetKeystoreParams is consumed by SetKeystore.
type SetKeystoreParams struct {
	*api.API
	ClusterID string

	Request *models.KeystoreContents
}

// Validate ensures the parameters are usabl by the consuming function.
func (params SetKeystoreParams) Validate() error {
	var merr = new(multierror.Error)

	if params.API == nil {
		merr = multierror.Append(merr, util.ErrAPIReq)
	}

	if len(params.ClusterID) != 32 {
		merr = multierror.Append(merr, util.ErrClusterLength)
	}

	if params.Request == nil {
		merr = multierror.Append(merr,
			errors.New("elasticsarch keystore: set requires the keystore contents to be specified"),
		)
	}

	return merr.ErrorOrNil()
}

// SetKeystore updates the Elasticsearch keystore values from a ClusterID.
// The specified request is treated as a partial definition and the values in
// it are used to either create / update the keystore values. Existing keystore
// values will not be affected if missing from the request.
func SetKeystore(params SetKeystoreParams) (*models.KeystoreContents, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.V1API.ClustersElasticsearch.SetEsClusterKeystore(
		clusters_elasticsearch.NewSetEsClusterKeystoreParams().
			WithClusterID(params.ClusterID).
			WithBody(params.Request),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}
