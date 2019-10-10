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

package apm

import (
	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_apm"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment/deputil"
	"github.com/elastic/ecctl/pkg/util"
)

// ListParams is used by the List apm function
type ListParams struct {
	*api.API

	// Optional parameters
	deputil.QueryParams
}

// Validate ensures that the parameters are usable by the consuming function.
func (params ListParams) Validate() error {
	return deputil.ValidateParams(&params)
}

// List returns all of the matching APM clusters.
func List(params ListParams) (*models.ApmsInfo, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	if params.Size < 1 {
		params.Size = *clusters_apm.NewGetApmClustersParams().Size
	}

	res, err := params.API.V1API.ClustersApm.GetApmClusters(
		clusters_apm.NewGetApmClustersParams().
			WithQ(ec.String(params.Query)).
			WithSize(ec.Int64(params.Size)).
			WithShowPlans(ec.Bool(params.ShowPlans)).
			WithShowPlanDefaults(ec.Bool(params.ShowPlanDefaults)).
			WithShowMetadata(ec.Bool(params.ShowMetadata)).
			WithShowHidden(ec.Bool(params.ShowHidden)),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}
	return res.Payload, nil
}

// ShowParams is used by the Show apm function
type ShowParams struct {
	*api.API
	ID string

	// Optional parameters
	deputil.QueryParams
}

// Validate ensures that the parameters are usable by the consuming function.
func (params ShowParams) Validate() error {
	var merr = new(multierror.Error)
	if params.API == nil {
		merr = multierror.Append(merr, util.ErrAPIReq)
	}

	if len(params.ID) != 32 {
		merr = multierror.Append(merr, deputil.NewInvalidDeploymentIDError(params.ID))
	}
	return merr.ErrorOrNil()
}

// Show returns a deployment with the specified settings
func Show(params ShowParams) (*models.ApmInfo, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.API.V1API.ClustersApm.GetApmCluster(
		clusters_apm.NewGetApmClusterParams().
			WithClusterID(params.ID).
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
