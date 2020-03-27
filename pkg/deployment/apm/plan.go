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
	"encoding/json"
	"fmt"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_apm"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/plan/planutil"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment/deputil"
	deploymentplanutil "github.com/elastic/ecctl/pkg/deployment/planutil"
	"github.com/elastic/ecctl/pkg/util"
)

// PlanParams is used by the Plan command.
type PlanParams struct {
	*api.API
	// ID represents the deployment ID.
	ID           string
	PlanDefaults bool

	Track bool
	planutil.TrackChangeParams
}

// Validate ensures that the parameters are usable by the consuming function.
func (params PlanParams) Validate() error {
	var err = multierror.Append(new(multierror.Error),
		deputil.ValidateParams(&params),
	)
	return err.ErrorOrNil()
}

// GetPlan returns the plan information (if any) of a cluster
func GetPlan(params PlanParams) (*models.ApmPlansInfo, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.API.V1API.ClustersApm.GetApmClusterPlanActivity(
		clusters_apm.NewGetApmClusterPlanActivityParams().
			WithClusterID(params.ID).
			WithShowPlanDefaults(&params.PlanDefaults),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	if !params.Track {
		return res.Payload, nil
	}

	return res.Payload, planutil.TrackChange(util.SetClusterTracking(
		params.TrackChangeParams, params.ID, util.Apm,
	))
}

// CancelPlan cancels the pending plan on the specified cluster.
func CancelPlan(params PlanParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	err := util.ReturnErrOnly(
		params.API.V1API.ClustersApm.CancelApmPendingPlan(
			clusters_apm.NewCancelApmPendingPlanParams().WithClusterID(params.ID),
			params.AuthWriter,
		))
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

// ListPlanHistory returns the historic plan list
func ListPlanHistory(params PlanParams) ([]*models.ApmPlanInfo, error) {
	res, err := GetPlan(params)
	if err != nil {
		return nil, err
	}

	return res.History, nil
}

func newDefaultTransientPlanConfiguration() *models.ApmPlanControlConfiguration {
	return &models.ApmPlanControlConfiguration{
		ExtendedMaintenance: ec.Bool(false),
		ReallocateInstances: ec.Bool(false),
	}
}

type computeTransientParams struct {
	plan      models.ApmPlan
	transient deploymentplanutil.ReapplyParams
}

// ReapplyLatestPlanAttempt will obtain the latest plan attempt and reapply it
// resetting all of the transient settings, Any setting can be overridden if
// specified in the params.
func ReapplyLatestPlanAttempt(params PlanParams, reparams deploymentplanutil.ReapplyParams) (*models.ApmCrudResponse, error) {
	if err := reparams.Validate(); err != nil {
		return nil, err
	}

	r, err := GetPlan(params)
	if err != nil {
		return nil, err
	}

	latestAttempt := getLatestAttempt(r.History)
	if latestAttempt == nil {
		return nil, fmt.Errorf("unable to obtain latest failed plan")
	}

	body := computeTransientSettings(computeTransientParams{
		plan:      *latestAttempt,
		transient: reparams,
	})
	if !reparams.HidePlan {
		enc := json.NewEncoder(params.TrackChangeParams.Writer)
		enc.SetIndent("", "  ")

		if err := enc.Encode(latestAttempt); err != nil {
			return nil, err
		}
	}

	_, res, err := params.API.V1API.ClustersApm.UpdateApmPlan(
		clusters_apm.NewUpdateApmPlanParams().
			WithClusterID(reparams.ID).
			WithBody(body),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	if !params.Track {
		return res.Payload, nil
	}

	return res.Payload, planutil.TrackChange(util.SetClusterTracking(
		params.TrackChangeParams, params.ID, util.Apm,
	))
}

// computeTransientSettings resets the transient.PlanConfiguration to a set of
// sane defaults, then, applies the values passed in params.transient
// It will always unset "restore_snapshot" to avoid restoring a snapshot from a
// previous plan that did restore it.
func computeTransientSettings(params computeTransientParams) *models.ApmPlan {
	if params.plan.Transient == nil {
		params.plan.Transient = new(models.TransientApmPlanConfiguration)
	}

	var transient = params.plan.Transient
	// If the strategy is overridden, obtain it from the parameters
	transient.Strategy = params.transient.Strategy()
	transient.PlanConfiguration = newDefaultTransientPlanConfiguration()
	transient.PlanConfiguration.ReallocateInstances = ec.Bool(params.transient.Reallocate)
	transient.PlanConfiguration.ExtendedMaintenance = ec.Bool(params.transient.ExtendedMaintenance)

	return &params.plan
}

// getLatestAttempt ensures that we get the latest attempt
// the history that we get back from the API is ordered.
func getLatestAttempt(history []*models.ApmPlanInfo) *models.ApmPlan {
	// It is highly unlikely to happen, but just in case it does
	// we don't want a panic to bubble up.
	if len(history) == 0 {
		return nil
	}

	return history[len(history)-1].Plan
}
