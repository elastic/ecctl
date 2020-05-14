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
	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/authentication"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/multierror"

	"github.com/elastic/ecctl/pkg/util"
)

// ListKeysParams is consumed by ListKeys
type ListKeysParams struct {
	*api.API
}

// Validate ensures the parameters are usable by the consuming function.
func (params ListKeysParams) Validate() error {
	var merr = multierror.NewPrefixed("user auth")
	if params.API == nil {
		merr = merr.Append(util.ErrAPIReq)
	}

	return merr.ErrorOrNil()
}

// ListKeys returns the available API keys for the current user.
func ListKeys(params ListKeysParams) (*models.APIKeysResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.V1API.Authentication.GetAPIKeys(
		authentication.NewGetAPIKeysParams(),
		params.AuthWriter,
	)

	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}
