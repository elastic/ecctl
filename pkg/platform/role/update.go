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

package role

import (
	"errors"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/platform_infrastructure"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/multierror"

	"github.com/elastic/ecctl/pkg/util"
)

// UpdateParams is consumed by Update.
type UpdateParams struct {
	*api.API

	Role *models.Role
	ID   string
}

// Validate ensures the parameters are usable.
func (params UpdateParams) Validate() error {
	var merr = multierror.NewPrefixed("role update")
	if params.API == nil {
		merr = merr.Append(util.ErrAPIReq)
	}

	if params.Role == nil {
		merr = merr.Append(errors.New("role definition cannot be empty"))
	}

	if params.ID == "" {
		merr = merr.Append(errors.New("id cannot be empty"))
	}

	return merr.ErrorOrNil()
}

// Update updates a role definition, the update uses PUT meaning it does not
// require a Role with the current + updated data, only requiring the changes
// which want to be updated
func Update(params UpdateParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	return util.ReturnErrOnly(
		params.V1API.PlatformInfrastructure.UpdateBlueprinterRole(
			platform_infrastructure.NewUpdateBlueprinterRoleParams().
				WithBlueprinterRoleID(params.ID).
				WithBody(params.Role),
			params.AuthWriter,
		),
	)
}
