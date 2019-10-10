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

package enrollmenttoken

import (
	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/platform_configuration_security"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"

	"github.com/elastic/ecctl/pkg/util"
)

// Create creates the token for the specific roles
func Create(params CreateParams) (*models.RequestEnrollmentTokenReply, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	var persistent = params.Duration.Seconds() <= 0
	var tokenConfig = models.EnrollmentTokenRequest{
		Persistent:        ec.Bool(persistent),
		Roles:             params.Roles,
		ValidityInSeconds: int32(params.Duration.Seconds()),
	}

	res, err := params.API.V1API.PlatformConfigurationSecurity.CreateEnrollmentToken(
		platform_configuration_security.NewCreateEnrollmentTokenParams().WithBody(&tokenConfig),
		params.AuthWriter,
	)

	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}

// Delete deletes a persistent token
func Delete(params DeleteParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	return util.ReturnErrOnly(
		params.API.V1API.PlatformConfigurationSecurity.DeleteEnrollmentToken(
			platform_configuration_security.NewDeleteEnrollmentTokenParams().
				WithToken(params.Token),
			params.AuthWriter,
		),
	)
}

// List lists all persistent tokens
func List(params ListParams) (*models.ListEnrollmentTokenReply, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.API.V1API.PlatformConfigurationSecurity.GetEnrollmentTokens(
		platform_configuration_security.NewGetEnrollmentTokensParams(),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}
	return res.Payload, nil
}
