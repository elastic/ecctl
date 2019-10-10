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

package user

import (
	"errors"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/users"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

// UpdateParams is consumed by Update
type UpdateParams struct {
	*api.API

	Password                  []byte
	Roles                     []string
	FullName, Email, UserName string
}

// Validate ensures the parameters are usable by the consuming function.
func (params UpdateParams) Validate() error {
	var merr = new(multierror.Error)

	if params.UserName == "" {
		merr = multierror.Append(merr, errors.New("user: update requires a username"))
	}

	if params.API == nil {
		merr = multierror.Append(merr, util.ErrAPIReq)
	}

	if params.Password != nil && len(params.Password) < minPasswordLength {
		merr = multierror.Append(merr, errors.New("user: update requires a password with a minimum of 8 characters"))
	}

	if params.Email != "" {
		if err := util.ValidateEmail("user", params.Email); err != nil {
			merr = multierror.Append(merr, err)
		}
	}

	merr = multierror.Append(merr, ValidateRoles(params.Roles))

	return merr.ErrorOrNil()
}

// Update updates an existing user.
func Update(params UpdateParams) (*models.User, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.V1API.Users.UpdateUser(
		users.NewUpdateUserParams().
			WithUserName(params.UserName).
			WithBody(&models.User{
				UserName: &params.UserName,
				FullName: params.FullName,
				Email:    params.Email,
				Security: &models.UserSecurity{
					Password: string(params.Password),
					Roles:    params.Roles,
				},
			}),
		params.AuthWriter,
	)

	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}

// UpdateCurrent updates the current user.
func UpdateCurrent(params UpdateParams) (*models.User, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.V1API.Users.UpdateCurrentUser(
		users.NewUpdateCurrentUserParams().
			WithBody(&models.User{
				UserName: &params.UserName,
				FullName: params.FullName,
				Email:    params.Email,
				Security: &models.UserSecurity{
					Password: string(params.Password),
					Roles:    params.Roles,
				},
			}),
		params.AuthWriter,
	)

	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}
