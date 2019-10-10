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

package constructor

import (
	"errors"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/platform_infrastructure"
	"github.com/elastic/cloud-sdk-go/pkg/models"

	"github.com/elastic/ecctl/pkg/util"
)

var (
	errAPICannotBeNil  = errors.New("api field cannot be nil")
	errIDCannotBeEmpty = errors.New("id field cannot be empty")
)

// Params is the generic set of parameters used for any constructor call
type Params struct {
	*api.API
}

// Validate checks the parameters
func (params Params) Validate() error {
	if params.API == nil {
		return errAPICannotBeNil
	}

	return nil
}

// GetParams is the set of parameters required for
type GetParams struct {
	Params
	ID string
}

// Validate checks the parameters
func (params GetParams) Validate() error {
	if params.ID == "" {
		return errIDCannotBeEmpty
	}

	return params.Params.Validate()
}

// MaintenanceParams is the set of parameters required for EnableMaintenace and
// DisableMaintenance
type MaintenanceParams struct {
	Params
	ID string
}

// Validate checks the parameters
func (params MaintenanceParams) Validate() error {
	if params.ID == "" {
		return errIDCannotBeEmpty
	}

	return params.Params.Validate()
}

// List gets the list of constuctors for a region
func List(params Params) (*models.ConstructorOverview, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	constructors, err := params.API.V1API.PlatformInfrastructure.GetConstructors(
		platform_infrastructure.NewGetConstructorsParams(),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return constructors.Payload, nil
}

// Get returns information about a specific constructor
func Get(params GetParams) (*models.ConstructorInfo, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	constructor, err := params.API.V1API.PlatformInfrastructure.GetConstructor(
		platform_infrastructure.NewGetConstructorParams().WithConstructorID(params.ID),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return constructor.Payload, nil
}

// EnableMaintenace sets the constructor to operational mode
func EnableMaintenace(params MaintenanceParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	return util.ReturnErrOnly(
		params.API.V1API.PlatformInfrastructure.StartConstructorMaintenanceMode(
			platform_infrastructure.NewStartConstructorMaintenanceModeParams().
				WithConstructorID(params.ID),
			params.AuthWriter,
		),
	)
}

// DisableMaintenance unsets the constructor to operational mode
func DisableMaintenance(params MaintenanceParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	return util.ReturnErrOnly(
		params.API.V1API.PlatformInfrastructure.StopConstructorMaintenanceMode(
			platform_infrastructure.NewStopConstructorMaintenanceModeParams().
				WithConstructorID(params.ID),
			params.AuthWriter,
		),
	)
}
