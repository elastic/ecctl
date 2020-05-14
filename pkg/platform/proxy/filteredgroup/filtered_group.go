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

package filteredgroup

import (
	"errors"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/platform_infrastructure"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/multierror"

	"github.com/elastic/ecctl/pkg/util"
)

var (
	errIDCannotBeEmpty                            = errors.New("proxies filtered group id cannot be empty")
	errFiltersCannotBeEmpty                       = errors.New("proxies filtered group filters cannot be empty")
	errExpectedProxiesCountCannotBeLesserThanZero = errors.New("proxies filtered group expected proxies count must be greater than 0")
	errVersionCannotBeLesserTahZero               = errors.New("proxies filtered group version cannot be less than 0")
)

// CommonParams is the set of parameters required for retrieving a proxies filtered group information
type CommonParams struct {
	*api.API
	ID string
}

// CreateParams is the set of parameters required for creating or updating proxies filtered group
type CreateParams struct {
	CommonParams
	Filters              map[string]string
	ExpectedProxiesCount int32
}

// UpdateParams is the set of parameters required for updating proxies filtered group
type UpdateParams struct {
	CreateParams
	Version int64
}

// Validate parameters for get and delete functions
func (params CommonParams) Validate() error {
	var merr = multierror.NewPrefixed("filtered group")
	if params.ID == "" {
		merr = merr.Append(errIDCannotBeEmpty)
	}

	if params.API == nil {
		merr = merr.Append(util.ErrAPIReq)
	}

	return merr.ErrorOrNil()
}

// Validate parameters for Create function
func (params CreateParams) Validate() error {
	var merr = multierror.NewPrefixed("filtered group")
	merr = merr.Append(params.CommonParams.Validate())

	if len(params.Filters) < 1 {
		merr = merr.Append(errFiltersCannotBeEmpty)
	}

	if params.ExpectedProxiesCount < 1 {
		merr = merr.Append(errExpectedProxiesCountCannotBeLesserThanZero)
	}

	return merr.ErrorOrNil()
}

// Validate parameters for Update function
func (params UpdateParams) Validate() error {
	var merr = multierror.NewPrefixed("filtered group")
	merr = merr.Append(params.CreateParams.Validate())

	if params.Version < 0 {
		merr = merr.Append(errVersionCannotBeLesserTahZero)
	}

	return merr.ErrorOrNil()
}

func createProxiesFilters(filters map[string]string) []*models.ProxiesFilter {
	proxiesFilters := make([]*models.ProxiesFilter, 0)
	for key, value := range filters {
		var k, v = key, value
		proxiesFilters = append(proxiesFilters, &models.ProxiesFilter{
			Key:   &k,
			Value: &v,
		})
	}

	return proxiesFilters
}

// Get returns information about a specific proxies filtered group
func Get(params CommonParams) (*models.ProxiesFilteredGroup, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	proxy, err := params.API.V1API.PlatformInfrastructure.GetProxiesFilteredGroup(
		platform_infrastructure.NewGetProxiesFilteredGroupParams().
			WithProxiesFilteredGroupID(params.ID),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return proxy.Payload, nil
}

// Create creates proxies filtered group with passed parameters
func Create(params CreateParams) (*models.ProxiesFilteredGroup, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	filters := createProxiesFilters(params.Filters)

	proxy, err := params.API.V1API.PlatformInfrastructure.CreateProxiesFilteredGroup(
		platform_infrastructure.NewCreateProxiesFilteredGroupParams().
			WithBody(&models.ProxiesFilteredGroup{
				Filters:              filters,
				ID:                   params.ID,
				ExpectedProxiesCount: &params.ExpectedProxiesCount,
			}),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return proxy.Payload, nil
}

// Delete deletes proxies filtered group by group id
func Delete(params CommonParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	_, err := params.API.V1API.PlatformInfrastructure.DeleteProxiesFilteredGroup(
		platform_infrastructure.NewDeleteProxiesFilteredGroupParams().
			WithProxiesFilteredGroupID(params.ID),
		params.AuthWriter,
	)
	return api.UnwrapError(err)
}

// Update updates information for already existing proxies filtered group
func Update(params UpdateParams) (*models.ProxiesFilteredGroup, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	filters := createProxiesFilters(params.Filters)

	body := &models.ProxiesFilteredGroup{
		Filters:              filters,
		ID:                   params.ID,
		ExpectedProxiesCount: &params.ExpectedProxiesCount,
	}

	proxy, err := params.API.V1API.PlatformInfrastructure.UpdateProxiesFilteredGroup(
		platform_infrastructure.NewUpdateProxiesFilteredGroupParams().
			WithBody(body).WithVersion(&params.Version).WithProxiesFilteredGroupID(params.ID),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return proxy.Payload, nil
}

// List gets the list of proxies filter groups for a region
func List(params CommonParams) ([]*models.ProxiesFilteredGroupHealth, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	proxies, err := params.API.V1API.PlatformInfrastructure.GetProxiesHealth(
		platform_infrastructure.NewGetProxiesHealthParams(),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return proxies.Payload.FilteredGroups, nil
}
