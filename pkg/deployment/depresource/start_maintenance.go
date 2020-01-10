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

// StartMaintenanceMode starts maintenance mode of all instances belonging to a deployment resource type.
func StartMaintenanceMode(params StartParams) (models.DeploymentResourceCommandResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.V1API.Deployments.StartDeploymentResourceInstancesAllMaintenanceMode(
		deployments.NewStartDeploymentResourceInstancesAllMaintenanceModeParams().
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

// StartInstancesMaintenanceMode starts maintenance mode of defined instances belonging to a deployment resource.
func StartInstancesMaintenanceMode(params StartInstancesParams) (models.DeploymentResourceCommandResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.V1API.Deployments.StartDeploymentResourceMaintenanceMode(
		deployments.NewStartDeploymentResourceMaintenanceModeParams().
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

// StartMaintenanceModeAllOrSpecified starts all or defined instances belonging to a deployment resource.
func StartMaintenanceModeAllOrSpecified(params StartInstancesParams) (models.DeploymentResourceCommandResponse, error) {
	if params.All {
		res, err := StartMaintenanceMode(params.StartParams)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	res, err := StartInstancesMaintenanceMode(params)
	if err != nil {
		return nil, err
	}
	return res, nil
}
