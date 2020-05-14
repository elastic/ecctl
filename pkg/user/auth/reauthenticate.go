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

package userauth

import (
	"errors"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/authentication"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/multierror"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"

	"github.com/elastic/ecctl/pkg/util"
)

// ReAuthenticateParams is consumed by ReAuthenticate
type ReAuthenticateParams struct {
	*api.API
	Password []byte
}

// Validate ensures the parameters are usable by the consuming function.
func (params ReAuthenticateParams) Validate() error {
	var merr = multierror.NewPrefixed("user auth")
	if params.API == nil {
		merr = merr.Append(util.ErrAPIReq)
	}

	if len(params.Password) == 0 {
		merr = merr.Append(errors.New("reauthenticate requires a password"))
	}

	return merr.ErrorOrNil()
}

// ReAuthenticate reauthenticates against the API by requiring the user's
// password on the request payload.
func ReAuthenticate(params ReAuthenticateParams) (string, error) {
	if err := params.Validate(); err != nil {
		return "", err
	}

	res, err := params.V1API.Authentication.ReAuthenticate(
		authentication.NewReAuthenticateParams().
			WithBody(&models.ReAuthenticationRequest{
				Password: ec.String(string(params.Password)),
			}),
		params.AuthWriter,
	)
	if err != nil {
		return "", api.UnwrapError(err)
	}

	return *res.Payload.SecurityToken, nil
}
