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
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

// ShowParams is consumed by Show.
type ShowParams struct {
	*api.API
	ID string
}

// Validate ensures the parameters are valid
func (params ShowParams) Validate() error {
	var merr = new(multierror.Error)
	if params.API == nil {
		merr = multierror.Append(merr, util.ErrAPIReq)
	}

	if params.ID == "" {
		merr = multierror.Append(merr, errors.New("role show: id cannot be empty"))
	}

	return merr.ErrorOrNil()
}

// Show returns a platform's role by ID
func Show(params ShowParams) (*models.RoleAggregate, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.V1API.PlatformInfrastructure.GetBlueprinterRole(
		platform_infrastructure.NewGetBlueprinterRoleParams().
			WithBlueprinterRoleID(params.ID),
		params.AuthWriter,
	)

	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}
