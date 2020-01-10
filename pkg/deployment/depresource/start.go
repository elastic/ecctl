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
)

// StartParams is consumed by start.
type StartParams struct {
	deployment.ResourceParams

	All bool
}

// StartInstancesParams is consumed by StartInstances.
type StartInstancesParams struct {
	StartParams
	IgnoreMissing *bool
	InstanceIDs   []string
}

// Validate ensures the parameters are usable by StartInstances.
func (params *StartInstancesParams) Validate() error {
	var merr = new(multierror.Error)

	if len(params.InstanceIDs) == 0 {
		merr = multierror.Append(merr, errors.New("deployment start: at least 1 instance ID must be provided"))
	}

	merr = multierror.Append(merr, params.StartParams.Validate())

	return merr.ErrorOrNil()
}

// Start starts all instances belonging to a deployment resource type.
func Start(params StartParams) (models.DeploymentResourceCommandResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.V1API.Deployments.StartDeploymentResourceInstancesAll(
		deployments.NewStartDeploymentResourceInstancesAllParams().
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

// StartInstances starts defined instances belonging to a deployment resource.
func StartInstances(params StartInstancesParams) (models.DeploymentResourceCommandResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.V1API.Deployments.StartDeploymentResourceInstances(
		deployments.NewStartDeploymentResourceInstancesParams().
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

// StartAllOrSpecified starts all or defined instances belonging to a deployment resource.
func StartAllOrSpecified(params StartInstancesParams) (models.DeploymentResourceCommandResponse, error) {
	if params.All {
		res, err := Start(params.StartParams)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	res, err := StartInstances(params)
	if err != nil {
		return nil, err
	}
	return res, nil
}
