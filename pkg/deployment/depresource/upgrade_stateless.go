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
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment"
	"github.com/elastic/ecctl/pkg/deployment/deputil"
	"github.com/elastic/ecctl/pkg/util"
)

// UpgradeStatelessParams is consumed by UpgradeStateless
type UpgradeStatelessParams struct {
	*api.API

	DeploymentID string
	Type         string
	RefID        string
}

// Validate ensures the parameters are usable by the consuming function.
func (params UpgradeStatelessParams) Validate() error {
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

func (params *UpgradeStatelessParams) fillDefaults() error {
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

// UpgradeStateless upgrades a stateless deployment resource like APM, Kibana
// and AppSearch.
func UpgradeStateless(params UpgradeStatelessParams) (*models.DeploymentResourceUpgradeResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	if err := params.fillDefaults(); err != nil {
		return nil, err
	}

	res, err := params.V1API.Deployments.UpgradeDeploymentStatelessResource(
		deployments.NewUpgradeDeploymentStatelessResourceParams().
			WithStatelessResourceKind(params.Type).
			WithDeploymentID(params.DeploymentID).
			WithRefID(params.RefID),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}
