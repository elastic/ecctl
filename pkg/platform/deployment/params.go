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
	"io"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/output"
	"github.com/go-openapi/strfmt"
	"github.com/pkg/errors"

	"github.com/elastic/ecctl/pkg/util"
)

var folderErrorMessage = "folder must not be empty"
var outErrorMessage = "Out io.Writer must be defined"

// TemplateParams is the parameter of template delete sub-command
type TemplateParams struct {
	*api.API
	ID string
}

// Validate is the implementation for the ecctl.Validator interface
func (params TemplateParams) Validate() error {
	if params.API == nil {
		return util.ErrAPIReq
	}

	if params.ID == "" {
		return errInvalidTemplateID
	}
	return nil
}

// ListTemplateParams is the parameter of template list sub-command
type ListTemplateParams struct {
	*api.API
	// If true, will return details for each instance configuration referenced by the template.
	ShowInstanceConfig bool
	// If present, it will cause the returned deployment templates to be adapted to return only the elements allowed
	// in that version.
	StackVersion string
	// An optional key/value pair in the form of (key:value) that will act as a filter and exclude any templates
	// that do not have a matching metadata item associated.
	Metadata string
}

// Validate is the implementation for the ecctl.Validator interface
func (params ListTemplateParams) Validate() error {
	if params.API == nil {
		return util.ErrAPIReq
	}
	return nil
}

// GetTemplateParams is the parameter of template show sub-command
type GetTemplateParams struct {
	TemplateParams
	ShowInstanceConfig bool
}

// CreateTemplateParams is the parameter of template create sub-command
type CreateTemplateParams struct {
	*api.API
	*models.DeploymentTemplateInfo
}

// Validate is the implementation for the ecctl.Validator interface
func (params CreateTemplateParams) Validate() error {
	if params.API == nil {
		return util.ErrAPIReq
	}
	if params.DeploymentTemplateInfo == nil {
		return errDeploymentTemplateMissing
	}
	return params.DeploymentTemplateInfo.Validate(strfmt.Default)
}

// UpdateTemplateParams is the parameter of template update sub-command
type UpdateTemplateParams struct {
	TemplateParams
	*models.DeploymentTemplateInfo
}

// Validate is the implementation for the ecctl.Validator interface
func (params UpdateTemplateParams) Validate() error {
	if params.DeploymentTemplateInfo == nil {
		return errDeploymentTemplateMissing
	}
	if err := params.DeploymentTemplateInfo.Validate(strfmt.Default); err != nil {
		return err
	}

	return params.TemplateParams.Validate()
}

// PushTemplateParams is used to push deployment templates as defined from a local folder.
// Optionally you can target only one deployment template.
type PushTemplateParams struct {
	*api.API
	ID, Folder string
	Out        *output.Device
	Force      bool
}

// Validate ensures that the parameters are correct.
func (params PushTemplateParams) Validate() error {
	if params.API == nil {
		return util.ErrAPIReq
	}
	if params.Folder == "" {
		return errors.New(folderErrorMessage)
	}
	if params.Out == nil {
		return errors.New(outErrorMessage)
	}
	return nil
}

// TemplateToFolderParams is the parameter for deployment template pull to folder sub-command
type TemplateToFolderParams struct {
	*api.API
	Folder string
}

// Validate is the implementation for the ecctl.Validator interface
func (params TemplateToFolderParams) Validate() error {
	if params.API == nil {
		return util.ErrAPIReq
	}

	if params.Folder == "" {
		return errors.New(folderErrorMessage)
	}
	return nil
}

// PullTemplateToFolderParams is the parameter for deployment template pull to folder sub-command
type PullTemplateToFolderParams struct {
	TemplateToFolderParams
}

// Validate is the implementation for the ecctl.Validator interface
func (params PullTemplateToFolderParams) Validate() error {
	return params.TemplateToFolderParams.Validate()
}

// DiffTemplateParams is used to compare deployment templates as defined in a local folder.
// Optionally you can target only one template.
type DiffTemplateParams struct {
	TemplateToFolderParams
	ID  string
	Out io.Writer
}

// Validate is the implementation for the ecctl.Validator interface
func (params DiffTemplateParams) Validate() error {
	if params.Out == nil {
		return errors.New(outErrorMessage)
	}

	return params.TemplateToFolderParams.Validate()
}

// TemplateFromFolderParams retrieve deployment template from folder
// template ID is optional
type TemplateFromFolderParams struct {
	Folder, ID string
}

// Validate ensures that the parameters are correct.
func (params TemplateFromFolderParams) Validate() error {
	if params.Folder == "" {
		return errors.New(folderErrorMessage)
	}
	return nil
}
