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
	"github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment/deputil"
	"github.com/elastic/ecctl/pkg/util"
)

// DeleteStatelessParams is consumed by Delete
type DeleteStatelessParams struct {
	*api.API

	DeploymentID string
	Type         string
	RefID        string
}

// Validate ensures the parameters are usable by the consuming function.
func (params DeleteStatelessParams) Validate() error {
	var merr = new(multierror.Error)

	if params.API == nil {
		merr = multierror.Append(merr, util.ErrAPIReq)
	}

	if len(params.DeploymentID) != 32 {
		merr = multierror.Append(merr, deputil.NewInvalidDeploymentIDError(params.DeploymentID))
	}

	if params.RefID == "" {
		merr = multierror.Append(merr, errors.New("deployment resource ref id cannot be empty"))
	}

	if params.Type == "" {
		merr = multierror.Append(merr, errors.New("deployment resource type cannot be empty"))
	}

	if params.Type == "elasticsearch" {
		merr = multierror.Append(merr, errors.New("deployment resource type \"elasticsearch\" is not supported"))
	}

	return merr.ErrorOrNil()
}

// DeleteStateless upgrades a stateless deployment resource like APM, Kibana
// and AppSearch.
func DeleteStateless(params DeleteStatelessParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	return util.ReturnErrOnly(
		params.V1API.Deployments.DeleteDeploymentStatelessResource(
			deployments.NewDeleteDeploymentStatelessResourceParams().
				WithStatelessResourceKind(params.Type).
				WithDeploymentID(params.DeploymentID).
				WithRefID(params.RefID),
			params.AuthWriter,
		),
	)
}
