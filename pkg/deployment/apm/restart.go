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
	"github.com/elastic/cloud-sdk-go/pkg/plan/planutil"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment/deputil"
	"github.com/elastic/ecctl/pkg/util"
)

// RestartParams is used by Restart.
type RestartParams struct {
	*api.API
	ID string

	// Force cancels any pending plans that might be in progress.
	Force bool

	Track bool
	planutil.TrackChangeParams
}

// Validate ensures that the parameters are usable by the consuming function.
func (params RestartParams) Validate() error {
	var err = multierror.Append(new(multierror.Error),
		deputil.ValidateParams(&params),
	)
	return err.ErrorOrNil()
}

// Restart an Apm cluster. There's two possible scenarios:
// * Cluster is running: re-applies the existing plan.
// * Cluster is stopped: starts it up with the most recent successful plan.
func Restart(params RestartParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	_, err := params.API.V1API.ClustersApm.RestartApm(
		clusters_apm.NewRestartApmParams().
			WithClusterID(params.ID).
			WithCancelPending(ec.Bool(params.Force)),
		params.AuthWriter,
	)
	if err != nil {
		return api.UnwrapError(err)
	}

	if !params.Track {
		return nil
	}

	return planutil.TrackChange(util.SetClusterTracking(
		params.TrackChangeParams, params.ID, util.Apm,
	))
}
