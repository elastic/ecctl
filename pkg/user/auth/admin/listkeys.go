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

package userauthadmin

import (
	"errors"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/authentication"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

// ListKeysParams is consumed by ListKeys
type ListKeysParams struct {
	*api.API

	UserID string
	All    bool
}

// Validate ensures the parameters are usable by the consuming function.
func (params ListKeysParams) Validate() error {
	var merr = new(multierror.Error)
	if params.API == nil {
		merr = multierror.Append(merr, util.ErrAPIReq)
	}

	if allAndUserIDSpecified := params.All && params.UserID != ""; allAndUserIDSpecified {
		merr = multierror.Append(merr, errors.New("userauthadmin: list keys requires a user ID or the all bool set, not both"))
	}

	if noAllOrUserIDSpecified := !params.All && params.UserID == ""; noAllOrUserIDSpecified {
		merr = multierror.Append(merr, errors.New("userauthadmin: list keys requires a user ID or all bool set"))
	}

	return merr.ErrorOrNil()
}

// ListKeys returns the API keys for either the specified user or all the
// platform users
func ListKeys(params ListKeysParams) (*models.APIKeysResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	return listUserOrAllKeys(params)
}

func listUserOrAllKeys(params ListKeysParams) (*models.APIKeysResponse, error) {
	if params.All {
		res, err := params.V1API.Authentication.GetUsersAPIKeys(
			authentication.NewGetUsersAPIKeysParams(),
			params.AuthWriter,
		)

		if err != nil {
			return nil, api.UnwrapError(err)
		}
		return res.Payload, nil
	}

	res, err := params.V1API.Authentication.GetUserAPIKeys(
		authentication.NewGetUserAPIKeysParams().
			WithUserID(params.UserID),
		params.AuthWriter,
	)

	if err != nil {
		return nil, api.UnwrapError(err)
	}
	return res.Payload, nil
}
