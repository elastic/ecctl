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

// GetKeyParams is consumed by GetKey
type GetKeyParams struct {
	*api.API

	ID     string
	UserID string
}

// Validate ensures the parameters are usable by the consuming function.
func (params GetKeyParams) Validate() error {
	var merr = new(multierror.Error)
	if params.API == nil {
		merr = multierror.Append(merr, util.ErrAPIReq)
	}

	if params.ID == "" {
		merr = multierror.Append(merr, errors.New("userauthadmin: get key requires a key id"))
	}

	if params.UserID == "" {
		merr = multierror.Append(merr, errors.New("userauthadmin: get key requires a user id"))
	}

	return merr.ErrorOrNil()
}

// GetKey returns the API key details for the specified key and user id.
func GetKey(params GetKeyParams) (*models.APIKeyResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.V1API.Authentication.GetUserAPIKey(
		authentication.NewGetUserAPIKeyParams().
			WithAPIKeyID(params.ID).
			WithUserID(params.UserID),
		params.AuthWriter,
	)

	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}
