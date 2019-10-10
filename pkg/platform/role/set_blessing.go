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

// SetBlessingsParams is consumed by SetBlessings.
type SetBlessingsParams struct {
	*api.API

	Blessings *models.Blessings
	ID        string
}

// Validate ensures the parameters are usable.
func (params SetBlessingsParams) Validate() error {
	var merr = new(multierror.Error)
	if params.API == nil {
		merr = multierror.Append(merr, util.ErrAPIReq)
	}

	if params.Blessings == nil {
		merr = multierror.Append(merr, errors.New("role set blessing: blessing definitions cannot be empty"))
	}

	if params.ID == "" {
		merr = multierror.Append(merr, errors.New("role set blessing: id cannot be empty"))
	}

	return merr.ErrorOrNil()
}

// SetBlessings sets a role blessing definitions, the update uses PUT meaning
// it does not require a Role with the current + updated data, only requiring
// the changes which want to be updated.
func SetBlessings(params SetBlessingsParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	return util.ReturnErrOnly(
		params.V1API.PlatformInfrastructure.SetBlueprinterBlessings(
			platform_infrastructure.NewSetBlueprinterBlessingsParams().
				WithBlueprinterRoleID(params.ID).
				WithBody(params.Blessings),
			params.AuthWriter,
		),
	)
}
