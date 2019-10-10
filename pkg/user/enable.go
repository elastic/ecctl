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

// EnableParams is consumed by Enable
type EnableParams struct {
	*api.API

	Enabled  bool
	UserName string
}

// Validate ensures the parameters are usable by the consuming function.
func (params EnableParams) Validate() error {
	var merr = new(multierror.Error)

	if params.UserName == "" {
		merr = multierror.Append(merr, errors.New("user: enable requires a username"))
	}

	if params.API == nil {
		merr = multierror.Append(merr, util.ErrAPIReq)
	}

	return merr.ErrorOrNil()
}

// Enable enables or disables an existing user.
func Enable(params EnableParams) (*models.User, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.V1API.Users.UpdateUser(
		users.NewUpdateUserParams().
			WithUserName(params.UserName).
			WithBody(&models.User{
				UserName: &params.UserName,
				Security: &models.UserSecurity{
					Enabled: &params.Enabled,
				},
			}),
		params.AuthWriter,
	)

	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}
