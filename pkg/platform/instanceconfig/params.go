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
	"errors"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/go-openapi/strfmt"

	"github.com/elastic/ecctl/pkg/util"
)

// ListParams is used to list all of the available instance configurations.
type ListParams struct {
	*api.API
}

// Validate ensures that the parameters are correct.
func (params ListParams) Validate() error {
	if params.API == nil {
		return util.ErrAPIReq
	}
	return nil
}

// GetParams is used to obtain an instance configuration from an ID.
type GetParams struct {
	*api.API
	ID string
}

// Validate ensures that the parameters are correct.
func (params GetParams) Validate() error {
	if params.API == nil {
		return util.ErrAPIReq
	}

	if params.ID == "" {
		return errors.New("get: id must not be empty")
	}
	return nil
}

// CreateParams is used to create a new instance configuration.
type CreateParams struct {
	*api.API
	Config *models.InstanceConfiguration
}

// Validate ensures that the parameters are correct.
func (params CreateParams) Validate() error {
	if params.API == nil {
		return util.ErrAPIReq
	}

	if params.Config == nil {
		return errors.New("create: request needs to have a config set")
	}

	return params.Config.Validate(strfmt.Default)
}

// UpdateParams is used to overwrite an existing instance configuration.
type UpdateParams struct {
	*api.API
	ID     string
	Config *models.InstanceConfiguration
}

// Validate ensures that the parameters are correct.
func (params UpdateParams) Validate() error {
	if params.API == nil {
		return util.ErrAPIReq
	}
	if params.ID == "" {
		return errors.New("update: id must not be empty")
	}
	if params.Config == nil {
		return errors.New("update: request needs to have a config set")
	}

	return params.Config.Validate(strfmt.Default)
}

// DeleteParams is used to delete an instance configuration from its ID.
type DeleteParams struct {
	*api.API
	ID string
}

// Validate ensures that the parameters are correct.
func (params DeleteParams) Validate() error {
	if params.API == nil {
		return util.ErrAPIReq
	}
	if params.ID == "" {
		return errors.New("delete: id must not be empty")
	}
	return nil
}

// PullToFolderParams is used to store all available instance configurations in a local folder.
type PullToFolderParams struct {
	*api.API
	Folder string
}

// Validate ensures that the parameters are correct.
func (params PullToFolderParams) Validate() error {
	if params.API == nil {
		return util.ErrAPIReq
	}
	if params.Folder == "" {
		return errors.New("pull: folder must not be empty")
	}
	return nil
}
