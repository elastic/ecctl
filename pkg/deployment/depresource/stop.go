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

package depresource

import (
	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/deployments"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/elastic/ecctl/pkg/deployment"
	"github.com/elastic/ecctl/pkg/deployment/deputil"
	"github.com/elastic/ecctl/pkg/util"
)

// StopParams is consumed by Stop.
type StopParams struct {
	*api.API
	DeploymentID string
	Type         string
	RefID        string
	All          bool
}

// Validate ensures the parameters are usable by Stop.
func (params StopParams) Validate() error {
	var merr = new(multierror.Error)

	if params.API == nil {
		merr = multierror.Append(merr, util.ErrAPIReq)
	}

	if len(params.DeploymentID) != 32 {
		merr = multierror.Append(merr, deputil.NewInvalidDeploymentIDError(params.DeploymentID))
	}

	if params.Type == "" {
		merr = multierror.Append(merr, errors.New("deployment resource type cannot be empty"))
	}

	return merr.ErrorOrNil()
}

func (params *StopParams) fillDefaults() error {
	if params.RefID == "" {
		refID, err := deployment.GetTypeRefID(deployment.GetResourceParams{
			GetParams: deployment.GetParams{
				API:          params.API,
				DeploymentID: params.DeploymentID,
			},
			Type: params.Type,
		})
		if err != nil {
			return err
		}

		params.RefID = refID
	}

	return nil
}

// StopInstancesParams is consumed by StopInstances.
type StopInstancesParams struct {
	StopParams
	IgnoreMissing *bool
	InstanceIDs   []string
}

// Validate ensures the parameters are usable by StopInstances.
func (params StopInstancesParams) Validate() error {
	var merr = new(multierror.Error)

	if len(params.InstanceIDs) == 0 {
		merr = multierror.Append(merr, errors.New("deployment stop: at least 1 instance ID must be provided"))
	}

	merr = multierror.Append(merr, params.StopParams.Validate())

	return merr.ErrorOrNil()
}

// Stop stops all instances belonging to a deployment resource type.
func Stop(params StopParams) (models.DeploymentResourceCommandResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	if err := params.fillDefaults(); err != nil {
		return nil, err
	}

	res, err := params.V1API.Deployments.StopDeploymentResourceInstancesAll(
		deployments.NewStopDeploymentResourceInstancesAllParams().
			WithDeploymentID(params.DeploymentID).
			WithResourceKind(params.Type).
			WithRefID(params.RefID),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}

// StopInstances stops defined instances belonging to a deployment resource.
func StopInstances(params StopInstancesParams) (models.DeploymentResourceCommandResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	if err := params.StopParams.fillDefaults(); err != nil {
		return nil, err
	}

	res, err := params.V1API.Deployments.StopDeploymentResourceInstances(
		deployments.NewStopDeploymentResourceInstancesParams().
			WithDeploymentID(params.DeploymentID).
			WithResourceKind(params.Type).
			WithIgnoreMissing(params.IgnoreMissing).
			WithInstanceIds(params.InstanceIDs).
			WithRefID(params.RefID),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}

// StopAllOrSpecified stops all or defined instances belonging to a deployment resource.
func StopAllOrSpecified(params StopInstancesParams) (models.DeploymentResourceCommandResponse, error) {
	if params.All {
		res, err := Stop(params.StopParams)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	res, err := StopInstances(params)
	if err != nil {
		return nil, err
	}
	return res, nil
}
