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

package snaprepo

import (
	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/platform_configuration_snapshots"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"

	"github.com/elastic/ecctl/pkg/util"
)

// Get obtains the specified snapshot repository configuration
func Get(params GetParams) (*models.RepositoryConfig, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	repo, err := params.V1API.PlatformConfigurationSnapshots.GetSnapshotRepository(
		platform_configuration_snapshots.NewGetSnapshotRepositoryParams().
			WithRepositoryName(params.Name),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return repo.Payload, nil
}

// List obtains all the configured platform snapshot repositories
func List(params Params) (*models.RepositoryConfigs, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	repo, err := params.V1API.PlatformConfigurationSnapshots.GetSnapshotRepositories(
		platform_configuration_snapshots.NewGetSnapshotRepositoriesParams(),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return repo.Payload, nil
}

// Delete removes a specified snapshot repository
func Delete(params DeleteParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	_, _, err := params.V1API.PlatformConfigurationSnapshots.DeleteSnapshotRepository(
		platform_configuration_snapshots.NewDeleteSnapshotRepositoryParams().
			WithRepositoryName(params.Name),
		params.AuthWriter,
	)

	return api.UnwrapError(err)
}

// Set adds or updates a snapshot repository from a config
func Set(params SetParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	return util.ReturnErrOnly(
		params.V1API.PlatformConfigurationSnapshots.SetSnapshotRepository(
			platform_configuration_snapshots.NewSetSnapshotRepositoryParams().
				WithRepositoryName(params.Name).
				WithBody(&models.SnapshotRepositoryConfiguration{
					Type:     ec.String(params.Type),
					Settings: params.Config,
				}),
			params.AuthWriter,
		),
	)
}
