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
	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/deployments"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment/deputil"
	"github.com/elastic/ecctl/pkg/util"
)

var systemAlerts = ec.Int64(5)

// GetParams is consumed by get resource functions
type GetParams struct {
	*api.API
	DeploymentID string

	// Optional parameters
	deputil.QueryParams

	// RefID, when specified, skips auto-discovering the deployment resource
	// RefID and instead uses the one that's passed.
	RefID string
}

// Validate ensures that the parameters are usable by the consuming function.
func (params GetParams) Validate() error {
	var merr = new(multierror.Error)
	if params.API == nil {
		merr = multierror.Append(merr, util.ErrAPIReq)
	}

	if len(params.DeploymentID) != 32 {
		merr = multierror.Append(merr, deputil.NewInvalidDeploymentIDError(params.DeploymentID))
	}

	return merr.ErrorOrNil()
}

// Get returns info about a deployment.
func Get(params GetParams) (*models.DeploymentGetResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.API.V1API.Deployments.GetDeployment(
		deployments.NewGetDeploymentParams().
			WithDeploymentID(params.DeploymentID).
			WithShowPlans(ec.Bool(params.ShowPlans)).
			WithShowPlanDefaults(ec.Bool(params.ShowPlanDefaults)).
			WithShowPlanLogs(ec.Bool(params.ShowPlanLogs)).
			WithShowMetadata(ec.Bool(params.ShowMetadata)).
			WithShowSettings(ec.Bool(params.ShowSettings)).
			WithShowSystemAlerts(systemAlerts),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}

// GetApm returns info about an apm resource belonging to a given deployment.
func GetApm(params GetParams) (*models.ApmResourceInfo, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.API.V1API.Deployments.GetDeploymentApmResourceInfo(
		deployments.NewGetDeploymentApmResourceInfoParams().
			WithDeploymentID(params.DeploymentID).
			WithRefID(params.RefID).
			WithShowPlans(ec.Bool(params.ShowPlans)).
			WithShowPlanDefaults(ec.Bool(params.ShowPlanDefaults)).
			WithShowPlanLogs(ec.Bool(params.ShowPlanLogs)).
			WithShowMetadata(ec.Bool(params.ShowMetadata)).
			WithShowSettings(ec.Bool(params.ShowSettings)),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}
	return res.Payload, nil
}

// GetAppSearch returns info about an appsearch resource belonging to a given deployment.
func GetAppSearch(params GetParams) (*models.AppSearchResourceInfo, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.API.V1API.Deployments.GetDeploymentAppsearchResourceInfo(
		deployments.NewGetDeploymentAppsearchResourceInfoParams().
			WithDeploymentID(params.DeploymentID).
			WithRefID(params.RefID).
			WithShowPlans(ec.Bool(params.ShowPlans)).
			WithShowPlanDefaults(ec.Bool(params.ShowPlanDefaults)).
			WithShowPlanLogs(ec.Bool(params.ShowPlanLogs)).
			WithShowMetadata(ec.Bool(params.ShowMetadata)).
			WithShowSettings(ec.Bool(params.ShowSettings)),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}
	return res.Payload, nil
}

// GetElasticsearch returns info about an elasticsearch resource belonging to a given deployment.
func GetElasticsearch(params GetParams) (*models.ElasticsearchResourceInfo, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.API.V1API.Deployments.GetDeploymentEsResourceInfo(
		deployments.NewGetDeploymentEsResourceInfoParams().
			WithDeploymentID(params.DeploymentID).
			WithRefID(params.RefID).
			WithShowPlans(ec.Bool(params.ShowPlans)).
			WithShowPlanDefaults(ec.Bool(params.ShowPlanDefaults)).
			WithShowPlanLogs(ec.Bool(params.ShowPlanLogs)).
			WithShowMetadata(ec.Bool(params.ShowMetadata)).
			WithShowSettings(ec.Bool(params.ShowSettings)).
			WithShowSystemAlerts(systemAlerts),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}
	return res.Payload, nil
}

// GetKibana returns info about an kibana resource belonging to a given deployment.
func GetKibana(params GetParams) (*models.KibanaResourceInfo, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.API.V1API.Deployments.GetDeploymentKibResourceInfo(
		deployments.NewGetDeploymentKibResourceInfoParams().
			WithDeploymentID(params.DeploymentID).
			WithRefID(params.RefID).
			WithShowPlans(ec.Bool(params.ShowPlans)).
			WithShowPlanDefaults(ec.Bool(params.ShowPlanDefaults)).
			WithShowPlanLogs(ec.Bool(params.ShowPlanLogs)).
			WithShowMetadata(ec.Bool(params.ShowMetadata)).
			WithShowSettings(ec.Bool(params.ShowSettings)),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}
	return res.Payload, nil
}
