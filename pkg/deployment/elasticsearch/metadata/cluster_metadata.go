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

package metadata

import (
	"errors"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_elasticsearch"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

// UpdateSettingsParams are the parameters required to update a cluster settings
type UpdateSettingsParams struct {
	util.ClusterParams
	SettingsParams *models.ClusterMetadataSettings
}

// Validate is the implementation for the ecctl.Validator interface
func (cp UpdateSettingsParams) Validate() error {
	var err = new(multierror.Error)

	if cp.SettingsParams == nil {
		err = multierror.Append(err, errors.New("metadata settings params cannot be empty"))
	}

	err = multierror.Append(err, cp.ClusterParams.Validate())

	return err.ErrorOrNil()
}

// UpdateSettings updates the cluster metadata settings
func UpdateSettings(params UpdateSettingsParams) (*models.ClusterMetadataSettings, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.V1API.ClustersElasticsearch.UpdateEsClusterMetadataSettings(
		clusters_elasticsearch.NewUpdateEsClusterMetadataSettingsParams().
			WithClusterID(params.ClusterID).WithBody(params.SettingsParams),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}

// GetSettings fetches the cluster metadata settings
func GetSettings(params util.ClusterParams) (*models.ClusterMetadataSettings, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.V1API.ClustersElasticsearch.GetEsClusterMetadataSettings(
		clusters_elasticsearch.NewGetEsClusterMetadataSettingsParams().
			WithClusterID(params.ClusterID),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}
