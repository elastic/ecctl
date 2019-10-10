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
	"errors"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_apm"
	"github.com/elastic/cloud-sdk-go/pkg/plan"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment/deputil"
	"github.com/elastic/ecctl/pkg/util"
)

// DeleteParams are the parameters needed to use Delete.
type DeleteParams struct {
	*api.API
	ID string
}

// Validate ensures that the parameters are usable by the consuming function.
func (params DeleteParams) Validate() error {
	return deputil.ValidateParams(&params)
}

// Delete deletes an APM.
func Delete(params DeleteParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	dep, err := Show(ShowParams{
		API: params.API,
		ID:  params.ID,
	})
	if err != nil {
		return api.UnwrapError(err)
	}

	if dep.Status != "stopped" {
		return errors.New("apm delete: deployment must be stopped")
	}

	return util.ReturnErrOnly(params.API.V1API.ClustersApm.DeleteApm(
		clusters_apm.NewDeleteApmParams().
			WithClusterID(params.ID),
		params.AuthWriter,
	))
}

// ShutdownParams are the parameters needed to use Shutdown.
type ShutdownParams struct {
	*api.API
	ID string

	// Stops the cluster and hides it from the user.
	Hide bool

	util.TrackParams
}

// Validate ensures that the parameters are usable by the consuming function.
func (params ShutdownParams) Validate() error {
	var merr = new(multierror.Error)
	if params.API == nil {
		merr = multierror.Append(merr, util.ErrAPIReq)
	}

	if len(params.ID) != 32 {
		merr = multierror.Append(merr, deputil.NewInvalidDeploymentIDError(params.ID))
	}

	merr = multierror.Append(merr, params.TrackParams.Validate())
	return merr.ErrorOrNil()
}

// Shutdown stops a running APM.
func Shutdown(params ShutdownParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	err := util.ReturnErrOnly(params.API.V1API.ClustersApm.ShutdownApm(
		clusters_apm.NewShutdownApmParams().
			WithClusterID(params.ID).
			WithHide(ec.Bool(params.Hide)),
		params.AuthWriter,
	))
	if err != nil {
		return api.UnwrapError(err)
	}

	if !params.Track {
		return nil
	}

	return util.TrackCluster(util.TrackClusterParams{
		Output: params.Output,
		TrackParams: plan.TrackParams{
			API:           params.API,
			PollFrequency: params.PollFrequency,
			MaxRetries:    params.MaxRetries,
			ID:            params.ID,
			Kind:          "apm",
		},
	})
}
