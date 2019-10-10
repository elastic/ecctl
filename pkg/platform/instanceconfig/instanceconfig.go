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

package instanceconfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/platform_configuration_instances"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/pkg/errors"

	"github.com/elastic/ecctl/pkg/util"
)

// InstanceConfigurationOperation define an operation needed for the iC.
type InstanceConfigurationOperation uint

const (
	// InstanceConfigurationOperationNone No operation needed on IC
	InstanceConfigurationOperationNone InstanceConfigurationOperation = 0
	// InstanceConfigurationOperationUpdate update IC
	InstanceConfigurationOperationUpdate InstanceConfigurationOperation = 1
	// InstanceConfigurationOperationCreate create IC
	InstanceConfigurationOperationCreate InstanceConfigurationOperation = 2
)

// Valid check validity of instance configuration action
func (a InstanceConfigurationOperation) Valid() bool {
	return (InstanceConfigurationOperationNone <= a && a <= InstanceConfigurationOperationCreate)
}

func (a InstanceConfigurationOperation) String() string {
	switch a {
	case InstanceConfigurationOperationNone:
		return fmt.Sprint("None")
	case InstanceConfigurationOperationUpdate:
		return fmt.Sprint("Update")
	case InstanceConfigurationOperationCreate:
		return fmt.Sprint("Create")
	default:
		return fmt.Sprint("Invalid")
	}
}

// List returns an array of all instance configurations
func List(params ListParams) ([]*models.InstanceConfiguration, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.API.V1API.PlatformConfigurationInstances.GetInstanceConfigurations(
		platform_configuration_instances.NewGetInstanceConfigurationsParams(),
		params.AuthWriter,
	)

	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}

// Get obtains an instance configuration from an ID
func Get(params GetParams) (*models.InstanceConfiguration, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.API.V1API.PlatformConfigurationInstances.GetInstanceConfiguration(
		platform_configuration_instances.NewGetInstanceConfigurationParams().
			WithID(params.ID),
		params.AuthWriter,
	)

	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}

// Create creates a new instance configuration.
func Create(params CreateParams) (*models.IDResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	if params.Config.ID != "" {
		if err := Update(UpdateParams{
			API:    params.API,
			ID:     params.Config.ID,
			Config: params.Config,
		}); err != nil {
			return nil, api.UnwrapError(err)
		}
		return &models.IDResponse{ID: ec.String(params.Config.ID)}, nil
	}

	res, err := params.API.V1API.PlatformConfigurationInstances.CreateInstanceConfiguration(
		platform_configuration_instances.NewCreateInstanceConfigurationParams().
			WithInstance(params.Config),
		params.AuthWriter,
	)

	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}

// Update overwrites an already existing instance configuration.
func Update(params UpdateParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	_, _, err := params.API.V1API.PlatformConfigurationInstances.SetInstanceConfiguration(
		platform_configuration_instances.NewSetInstanceConfigurationParams().
			WithID(params.ID).
			WithInstance(params.Config),
		params.AuthWriter,
	)

	if err != nil {
		return api.UnwrapError(err)
	}

	return nil
}

// Delete deletes an already existing instance configuration.
func Delete(params DeleteParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	return util.ReturnErrOnly(
		params.API.V1API.PlatformConfigurationInstances.DeleteInstanceConfiguration(
			platform_configuration_instances.NewDeleteInstanceConfigurationParams().
				WithID(params.ID),
			params.AuthWriter,
		),
	)
}

// PullToFolder downloads instance configs and save them in a local folder
func PullToFolder(params PullToFolderParams) error {
	if err := params.Validate(); err != nil {
		return err
	}
	res, err := List(ListParams{API: params.API})
	if err != nil {
		return err
	}
	return writeInstanceConfigToFolder(params.Folder, res)
}

// writeInstanceConfigToFolder this will write all the instance configs to a folder
// following this structure:
//   folder/
//   folder/id.json
func writeInstanceConfigToFolder(folder string, instanceConfigs []*models.InstanceConfiguration) error {
	err := os.MkdirAll(filepath.Dir(folder), os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "error creating dir %s", folder)
	}

	for _, instanceConfig := range instanceConfigs {
		b, err := json.MarshalIndent(instanceConfig, "", "  ")
		if err != nil {
			return errors.Wrapf(err, "error converting instance config %s into json", instanceConfig.ID)
		}
		jsonPath := filepath.Join(folder, instanceConfig.ID+".json")
		err = ioutil.WriteFile(jsonPath, b, os.ModePerm)
		if err != nil {
			return errors.Wrapf(err, "error writing instance config %s", jsonPath)
		}
	}
	return nil
}
