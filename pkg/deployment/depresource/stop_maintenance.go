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
)

// StopMaintenanceMode stops maintenance mode of all instances belonging to a deployment resource kind.
func StopMaintenanceMode(params StopParams) (models.DeploymentResourceCommandResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.V1API.Deployments.StopDeploymentResourceInstancesAllMaintenanceMode(
		deployments.NewStopDeploymentResourceInstancesAllMaintenanceModeParams().
			WithDeploymentID(params.DeploymentID).
			WithResourceKind(params.Kind).
			WithRefID(params.RefID),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}

// StopInstancesMaintenanceMode stops maintenance mode of defined instances belonging to a deployment resource.
func StopInstancesMaintenanceMode(params StopInstancesParams) (models.DeploymentResourceCommandResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.V1API.Deployments.StopDeploymentResourceMaintenanceMode(
		deployments.NewStopDeploymentResourceMaintenanceModeParams().
			WithDeploymentID(params.DeploymentID).
			WithResourceKind(params.Kind).
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

// StopMaintenanceModeAllOrSpecified stops all or defined instances belonging to a deployment resource.
func StopMaintenanceModeAllOrSpecified(params StopInstancesParams) (models.DeploymentResourceCommandResponse, error) {
	if params.All {
		res, err := StopMaintenanceMode(params.StopParams)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	res, err := StopInstancesMaintenanceMode(params)
	if err != nil {
		return nil, err
	}
	return res, nil
}
