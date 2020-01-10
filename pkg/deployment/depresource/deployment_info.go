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
	"errors"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/deployments"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

// GetDeploymentInfoParams is consumed by GetDeploymentInfo.
type GetDeploymentInfoParams struct {
	*api.API
	DeploymentID string
}

// GetDeploymentInfoResponse is returned by GetDeploymentInfo.
type GetDeploymentInfoResponse struct {
	// RefID is the ElasticsearchRefID.
	RefID string

	// Deployment template which the Elasticsearch resource is using.
	DeploymentTemplate string
}

// Validate ensures the parameters are usable by the consuming function.
func (params GetDeploymentInfoParams) Validate() error {
	var merr = new(multierror.Error)

	if params.API == nil {
		merr = multierror.Append(merr, util.ErrAPIReq)
	}

	if len(params.DeploymentID) != 32 {
		merr = multierror.Append(merr, util.ErrDeploymentID)
	}

	return merr.ErrorOrNil()
}

// GetDeploymentInfo obtains the Deployment's elasticsearch relevant info and
// packs it into GetDeploymentInfoResponse. It is convenient for actions which
// need to auto-discover information of their deployment.
func GetDeploymentInfo(params GetDeploymentInfoParams) (GetDeploymentInfoResponse, error) {
	var emptyRes GetDeploymentInfoResponse
	res, err := params.V1API.Deployments.GetDeployment(
		deployments.NewGetDeploymentParams().
			WithDeploymentID(params.DeploymentID).
			WithEnrichWithTemplate(ec.Bool(true)).
			WithConvertLegacyPlans(ec.Bool(true)).
			WithShowPlans(ec.Bool(true)),
		params.AuthWriter,
	)
	if err != nil {
		return emptyRes, api.UnwrapError(err)
	}

	for _, resource := range res.Payload.Resources.Elasticsearch {
		var planInfo = resource.Info.PlanInfo
		var hasPlan = planInfo.Current != nil && planInfo.Current.Plan != nil
		var refID string
		if hasPlan {
			refID = *resource.RefID
		}

		if hasPlan && planInfo.Current.Plan.DeploymentTemplate != nil {
			return GetDeploymentInfoResponse{
				RefID:              refID,
				DeploymentTemplate: *planInfo.Current.Plan.DeploymentTemplate.ID,
			}, nil
		}
	}

	return emptyRes, errors.New("unable to obtain deployment template ID from existing deployment ID, please specify a one")
}
