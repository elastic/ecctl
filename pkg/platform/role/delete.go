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
	"github.com/elastic/cloud-sdk-go/pkg/multierror"

	"github.com/elastic/ecctl/pkg/util"
)

// DeleteParams is consumed by Delete.
type DeleteParams struct {
	*api.API
	ID string
}

// Validate ensures the parameters are usable.
func (params DeleteParams) Validate() error {
	var merr = multierror.NewPrefixed("role delete")
	if params.API == nil {
		merr = merr.Append(util.ErrAPIReq)
	}

	if params.ID == "" {
		merr = merr.Append(errors.New("id cannot be empty"))
	}

	return merr.ErrorOrNil()
}

// Delete delets a role by ID.
func Delete(params DeleteParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	return util.ReturnErrOnly(
		params.V1API.PlatformInfrastructure.DeleteBlueprinterRole(
			platform_infrastructure.NewDeleteBlueprinterRoleParams().
				WithBlueprinterRoleID(params.ID),
			params.AuthWriter,
		),
	)
}
