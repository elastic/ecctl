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

package kibana

import (
	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/plan/planutil"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment/deputil"
)

// stoppedState represents the stopped state string for a cluster
const stoppedState = "stopped"

// DeploymentParams is the base struct meant to be embedded in other kibana pkg params
type DeploymentParams struct {
	API *api.API
	// ID represents the deployment ID.
	ID string

	Track bool
	planutil.TrackChangeParams
}

// Validate ensures that the parameters are usable by the consuming function.
func (params *DeploymentParams) Validate() error {
	var err = multierror.Append(new(multierror.Error),
		deputil.ValidateParams(params),
	)
	return err.ErrorOrNil()
}
