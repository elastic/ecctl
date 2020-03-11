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
	"github.com/elastic/cloud-sdk-go/pkg/client/deployments"

	"github.com/elastic/ecctl/pkg/deployment"
	"github.com/elastic/ecctl/pkg/util"
)

// ShutdownParams is consumed by Shutdown.
type ShutdownParams struct {
	deployment.ResourceParams

	// Optional Values
	// Skips taking a snapshot before shutting down the deployment resource.
	SkipSnapshot bool

	// Hides the resource. Hidden resources are not listed by default.
	Hide bool
}

// Shutdown stops all the running instances for the specified resource kind ref
// ID on a Deployment. If no refID is specified, it tries to autodiscover it.
func Shutdown(params ShutdownParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	if params.Kind == "elasticsearch" {
		return util.ReturnErrOnly(
			params.V1API.Deployments.ShutdownDeploymentEsResource(
				deployments.NewShutdownDeploymentEsResourceParams().
					WithDeploymentID(params.DeploymentID).
					WithSkipSnapshot(&params.SkipSnapshot).
					WithRefID(params.RefID).
					WithHide(&params.Hide),
				params.AuthWriter,
			),
		)
	}

	return util.ReturnErrOnly(
		params.V1API.Deployments.ShutdownDeploymentStatelessResource(
			deployments.NewShutdownDeploymentStatelessResourceParams().
				WithDeploymentID(params.DeploymentID).
				WithSkipSnapshot(&params.SkipSnapshot).
				WithStatelessResourceKind(params.Kind).
				WithRefID(params.RefID).
				WithHide(&params.Hide),
			params.AuthWriter,
		),
	)
}
