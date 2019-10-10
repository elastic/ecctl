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

package deployment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/platform_configuration_templates"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/pkg/errors"

	"github.com/elastic/ecctl/pkg/util"
)

var (
	errInvalidTemplateID         = errors.New("invalid template ID")
	errDeploymentTemplateMissing = errors.New("deployment template is missing")
)

// TemplateOperation define an operation needed on DT.
type TemplateOperation uint

const (
	// TemplateOperationNone No operation needed on DT
	TemplateOperationNone TemplateOperation = 0
	// TemplateOperationUpdate update DT
	TemplateOperationUpdate TemplateOperation = 1
	// TemplateOperationCreate create DT
	TemplateOperationCreate TemplateOperation = 2
)

// Valid check validity of deployment template action
func (a TemplateOperation) Valid() bool {
	return TemplateOperationNone <= a && a <= TemplateOperationCreate
}

func (a TemplateOperation) String() string {
	switch a {
	case TemplateOperationNone:
		return fmt.Sprint("None")
	case TemplateOperationUpdate:
		return fmt.Sprint("Update")
	case TemplateOperationCreate:
		return fmt.Sprint("Create")
	default:
		return fmt.Sprint("Invalid")
	}
}

// CompareResult contains the diffs between two deployment template after they have
// been compared.
type CompareResult struct {
	Equals    bool
	Diff      string
	DT        *models.DeploymentTemplateInfo
	Operation TemplateOperation
}

// ListTemplates obtains all the configured platform deployment templates
func ListTemplates(params ListTemplateParams) (*platform_configuration_templates.GetDeploymentTemplatesOK, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.V1API.PlatformConfigurationTemplates.GetDeploymentTemplates(
		platform_configuration_templates.NewGetDeploymentTemplatesParams().
			WithStackVersion(ec.String(params.StackVersion)).
			WithMetadata(ec.String(params.Metadata)).
			WithShowInstanceConfigurations(ec.Bool(params.ShowInstanceConfig)),
		params.AuthWriter,
	)

	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res, nil
}

// GetTemplate obtains information about a specific platform deployment template
func GetTemplate(params GetTemplateParams) (*models.DeploymentTemplateInfo, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.V1API.PlatformConfigurationTemplates.GetDeploymentTemplate(
		platform_configuration_templates.NewGetDeploymentTemplateParams().
			WithShowInstanceConfigurations(ec.Bool(params.ShowInstanceConfig)).
			WithTemplateID(params.ID),
		params.AuthWriter,
	)

	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}

// DeleteTemplate deletes a specific platform deployment template
func DeleteTemplate(params GetTemplateParams) error {
	if err := params.Validate(); err != nil {
		return err
	}
	return util.ReturnErrOnly(params.V1API.PlatformConfigurationTemplates.DeleteDeploymentTemplate(
		platform_configuration_templates.NewDeleteDeploymentTemplateParams().
			WithTemplateID(params.ID),
		params.AuthWriter,
	),
	)
}

// CreateTemplate creates a platform deployment template
func CreateTemplate(params CreateTemplateParams) (string, error) {
	if err := params.Validate(); err != nil {
		return "", err
	}

	if params.ID != "" {
		if err := UpdateTemplate(UpdateTemplateParams{
			TemplateParams: TemplateParams{
				API: params.API,
				ID:  params.ID,
			},
			DeploymentTemplateInfo: params.DeploymentTemplateInfo,
		}); err != nil {
			return "", api.UnwrapError(err)
		}
		return params.ID, nil
	}
	resp, err := params.V1API.PlatformConfigurationTemplates.CreateDeploymentTemplate(
		platform_configuration_templates.NewCreateDeploymentTemplateParams().
			WithBody(params.DeploymentTemplateInfo),
		params.AuthWriter,
	)

	if err != nil {
		return "", api.UnwrapError(err)
	}

	return *resp.Payload.ID, nil
}

// UpdateTemplate updates a platform deployment template
func UpdateTemplate(params UpdateTemplateParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	_, _, err := params.V1API.PlatformConfigurationTemplates.SetDeploymentTemplate(
		platform_configuration_templates.NewSetDeploymentTemplateParams().
			WithBody(params.DeploymentTemplateInfo).
			WithTemplateID(params.TemplateParams.ID),
		params.AuthWriter,
	)

	if err != nil {
		return api.UnwrapError(err)
	}

	return nil
}

// PullToFolder downloads deployment templates and save them in a local folder
func PullToFolder(params PullTemplateToFolderParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	res, err := ListTemplates(ListTemplateParams{
		API: params.API,
	})
	if err != nil {
		return err
	}

	return writeDeploymentTemplateToFolder(params.Folder, res.Payload)
}

// writeDeploymentTemplateToFolder this will write all the deployment template to a folder
// following this structure:
//   folder/
//   folder/id.json
func writeDeploymentTemplateToFolder(folder string, deploymentTemplates []*models.DeploymentTemplateInfo) error {
	err := os.MkdirAll(filepath.Dir(folder), os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "error creating dir %s", folder)
	}

	for _, deploymentTemplate := range deploymentTemplates {
		// we do not care the source that created or modified the template.
		deploymentTemplate.Source = nil
		jsonMarshal, err := json.MarshalIndent(deploymentTemplate, "", "  ")
		if err != nil {
			return errors.Wrapf(err, "error converting deployment template %s into json", deploymentTemplate.ID)
		}
		jsonPath := filepath.Join(folder, deploymentTemplate.ID+".json")
		err = ioutil.WriteFile(jsonPath, jsonMarshal, os.ModePerm)
		if err != nil {
			return errors.Wrapf(err, "error writing deployment template %s", jsonPath)
		}
	}
	return nil
}

// GetDeploymentTemplateFromFolder retrieves deployment template from a local folder.
// It can retrieve all DT definition or a single one. Depends on the id value.
func GetDeploymentTemplateFromFolder(params TemplateFromFolderParams) (map[string]*models.DeploymentTemplateInfo, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	if _, err := os.Stat(params.Folder); os.IsNotExist(err) {
		return nil, fmt.Errorf("folder %s does not exist", params.Folder)
	}

	var deploymentTemplate *models.DeploymentTemplateInfo

	if params.ID != "" {
		filePath := fmt.Sprintf("%s%s.json", params.Folder, params.ID)
		f, err := os.Open(filePath)
		if err != nil {
			return nil, fmt.Errorf("error %s while opening %s", err, filePath)
		}
		deploymentTemplate = &models.DeploymentTemplateInfo{}
		jsonParser := json.NewDecoder(f)
		if err = jsonParser.Decode(deploymentTemplate); err != nil {
			return nil, fmt.Errorf("error %s while decoding %s", err, filePath)
		}
		return map[string]*models.DeploymentTemplateInfo{
			params.ID: deploymentTemplate,
		}, nil
	}

	matches, err := filepath.Glob(filepath.Join(params.Folder, "*.json"))
	if err != nil {
		return nil, err
	}

	res := map[string]*models.DeploymentTemplateInfo{}
	for _, deploymentTemplateFile := range matches {
		fileContents, err := ioutil.ReadFile(deploymentTemplateFile)
		if err != nil {
			return nil, err
		}
		deploymentTemplate = &models.DeploymentTemplateInfo{}
		if err := json.Unmarshal(fileContents, deploymentTemplate); err != nil {
			return nil, err
		}
		res[deploymentTemplate.ID] = deploymentTemplate
	}
	return res, nil
}

// GetRemoteDeploymentTemplates retrieves remote deployment templates
func GetRemoteDeploymentTemplates(params GetTemplateParams) (map[string]*models.DeploymentTemplateInfo, error) {
	remoteDeploymentTemplate := map[string]*models.DeploymentTemplateInfo{}
	if params.ID != "" {
		res, err := GetTemplate(params)
		if err != nil {
			res = nil
		} else {
			res.Source = nil
		}
		remoteDeploymentTemplate[params.ID] = res
		return remoteDeploymentTemplate, nil
	}

	res, err := ListTemplates(ListTemplateParams{
		API: params.API,
	})
	if err != nil {
		return nil, err
	}
	for _, deploymentTemplate := range res.Payload {
		deploymentTemplate.Source = nil
		remoteDeploymentTemplate[deploymentTemplate.ID] = deploymentTemplate
	}

	return remoteDeploymentTemplate, nil
}
