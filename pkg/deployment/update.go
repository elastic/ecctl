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

package deployment

import (
	"errors"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/deployments"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

// UpdateParams is consumed by Update.
type UpdateParams struct {
	*api.API

	DeploymentID string
	Request      *models.DeploymentUpdateRequest

	// Optional values
	SkipSnapshot      bool
	HidePrunedOrphans bool

	// Region is an optional value which if set and the the Request.Resources
	// are missing a region, it populates that field with the value of Region.
	Region string
}

// Validate ensures the parameters are usable by Update.
func (params UpdateParams) Validate() error {
	var merr = new(multierror.Error)

	if params.API == nil {
		merr = multierror.Append(merr, util.ErrAPIReq)
	}

	if params.Request == nil {
		merr = multierror.Append(merr, errors.New("deployment update: request payload cannot be empty"))
	}

	if len(params.DeploymentID) != 32 {
		merr = multierror.Append(merr, util.ErrDeploymentID)
	}

	return merr.ErrorOrNil()
}

// Update receives an update payload with an optional region override in case
// the region isn't specified in the update request payload. Additionally if
// Request.PruneOrphans is false then any omitted resources aren't shutdown.
// The opposite behavior can be expected when the flag is true since the update
// request is treated as the single source of truth and the complete desired
// deployment definition.
func Update(params UpdateParams) (*models.DeploymentUpdateResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	setOverrides(params.Request, &PayloadOverrides{Region: params.Region})

	res, err := params.V1API.Deployments.UpdateDeployment(
		deployments.NewUpdateDeploymentParams().
			WithDeploymentID(params.DeploymentID).
			WithBody(params.Request).
			WithSkipSnapshot(&params.SkipSnapshot).
			WithHidePrunedOrphans(&params.HidePrunedOrphans),
		params.AuthWriter,
	)

	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}
