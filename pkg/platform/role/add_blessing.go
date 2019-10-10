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

// AddBlessingParams is consumed by AddBlessing.
type AddBlessingParams struct {
	*api.API

	Blessing *models.Blessing
	RunnerID string
	ID       string
}

// Validate ensures the parameters are usable.
func (params AddBlessingParams) Validate() error {
	var merr = new(multierror.Error)
	if params.API == nil {
		merr = multierror.Append(merr, util.ErrAPIReq)
	}

	if params.Blessing == nil {
		merr = multierror.Append(merr, errors.New("role add blessing: blessing definition cannot be empty"))
	}

	if params.ID == "" {
		merr = multierror.Append(merr, errors.New("role add blessing: id cannot be empty"))
	}

	if params.RunnerID == "" {
		merr = multierror.Append(merr, errors.New("role add blessing: runner id cannot be empty"))
	}

	return merr.ErrorOrNil()
}

// AddBlessing adds a role blessing to a runner ID.
func AddBlessing(params AddBlessingParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	return util.ReturnErrOnly(
		params.V1API.PlatformInfrastructure.AddBlueprinterBlessing(
			platform_infrastructure.NewAddBlueprinterBlessingParams().
				WithBlueprinterRoleID(params.ID).
				WithRunnerID(params.RunnerID).
				WithBody(params.Blessing),
			params.AuthWriter,
		),
	)
}
