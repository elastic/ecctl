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

package plan

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_elasticsearch"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/plan/planutil"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"

	depplanutil "github.com/elastic/ecctl/pkg/deployment/planutil"
	"github.com/elastic/ecctl/pkg/util"
)

var (
	errInvalidTransientParameters = errors.New(`current transient settings could lead to data loss`)
)

// ReapplyParams contains the parameters required to call a plan reapply action
type ReapplyParams struct {
	*api.API
	depplanutil.ReapplyParams

	Output io.Writer
	planutil.TrackChangeParams
	Track                bool
	SkipSnapshot         bool `kebabcase:"Skips snapshot on the reapplied plan"`
	SkipDataMigration    bool `kebabcase:"Bypasses the need to wait for data to be migrated from old instances to new instances before continuing the plan (potentially deleting the old instances and losing data)"`
	SkipPostUpgradeSteps bool `kebabcase:"Bypasses 2.x->5.x operations for any plan change ending with a 5.x cluster (eg apply a cluster license, ensure Monitoring is configured)"`
	SkipUpgradeChecker   bool `kebabcase:"Bypasses issue checks that should be resolved before migration (eg contains old Lucene segments)"`
}

// Validate ensures that the parameters are usable
func (params ReapplyParams) Validate() error {
	var err = new(multierror.Error)
	if params.API == nil {
		err = multierror.Append(err, util.ErrAPIReq)
	}

	err = multierror.Append(err, params.ReapplyParams.Validate())

	return err.ErrorOrNil()
}

func newDefaultTransientElasticsearchPlanConfiguration() *models.ElasticsearchPlanControlConfiguration {
	return &models.ElasticsearchPlanControlConfiguration{
		ExtendedMaintenance:  ec.Bool(false),
		OverrideFailsafe:     ec.Bool(false),
		ReallocateInstances:  ec.Bool(false),
		SkipDataMigration:    ec.Bool(false),
		SkipPostUpgradeSteps: ec.Bool(false),
		SkipSnapshot:         ec.Bool(false),
		SkipUpgradeChecker:   ec.Bool(false),
	}
}

type computeTransientParams struct {
	plan      models.ElasticsearchClusterPlan
	transient ReapplyParams
}

// validate ensures that is only dereferenced when the property is not nil
func (tp computeTransientParams) validate() error {
	var transient = tp.transient

	var reallocateInvalid = transient.SkipDataMigration && transient.Reallocate
	if reallocateInvalid {
		return errInvalidTransientParameters
	}

	return nil
}

// Reapply will obtain the latest plan attempt and reapply it resetting all the
// transient settings, Any setting can be overridden if specified in the params.
func Reapply(params ReapplyParams) (*models.ClusterCrudResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	r, err := GetHistory(GetHistoryParams{ClusterParams: util.ClusterParams{
		API:       params.API,
		ClusterID: params.ID,
	}})
	if err != nil {
		return nil, err
	}

	latestAttempt := getLatestAttempt(r)
	if latestAttempt == nil {
		return nil, fmt.Errorf("unable to obtain latest plan attempt")
	}

	body, err := computeTransientSettings(computeTransientParams{
		plan:      *latestAttempt,
		transient: params,
	})
	if err != nil {
		return nil, err
	}

	if !params.HidePlan {
		enc := json.NewEncoder(params.Output)
		enc.SetIndent("", "  ")
		if err := enc.Encode(body); err != nil {
			return nil, err
		}
	}

	res, res2, err := params.V1API.ClustersElasticsearch.UpdateEsClusterPlan(
		clusters_elasticsearch.NewUpdateEsClusterPlanParams().
			WithClusterID(params.ID).
			WithBody(body),
		params.AuthWriter,
	)

	return util.ParseCUResponse(util.ParseCUResponseParams{
		UpdateResponse:    res2,
		CreateResponse:    res,
		Err:               err,
		Track:             params.Track,
		TrackChangeParams: params.TrackChangeParams,
	})
}

// computeTransientSettings resets the transient.PlanConfiguration to a set of
// sane defaults, then, applies the values passed in params.transient
// It will always unset "restore_snapshot" to avoid restoring a snapshot from a
// previous plan that did restore it.
func computeTransientSettings(params computeTransientParams) (*models.ElasticsearchClusterPlan, error) {
	if params.plan.Transient == nil {
		params.plan.Transient = new(models.TransientElasticsearchPlanConfiguration)
	}

	var transient = params.plan.Transient

	// If the strategy is overridden, obtain it from the parameters
	transient.Strategy = params.transient.Strategy()
	transient.RestoreSnapshot = nil
	transient.PlanConfiguration = newDefaultTransientElasticsearchPlanConfiguration()
	transient.PlanConfiguration.SkipSnapshot = ec.Bool(params.transient.SkipSnapshot)
	transient.PlanConfiguration.ReallocateInstances = ec.Bool(params.transient.Reallocate)
	transient.PlanConfiguration.ExtendedMaintenance = ec.Bool(params.transient.ExtendedMaintenance)
	transient.PlanConfiguration.OverrideFailsafe = ec.Bool(params.transient.OverrideFailsafe)
	transient.PlanConfiguration.SkipDataMigration = ec.Bool(params.transient.SkipDataMigration)
	transient.PlanConfiguration.SkipPostUpgradeSteps = ec.Bool(params.transient.SkipPostUpgradeSteps)
	transient.PlanConfiguration.SkipUpgradeChecker = ec.Bool(params.transient.SkipUpgradeChecker)

	if err := params.validate(); err != nil {
		return nil, err
	}
	return &params.plan, nil
}

// getLatestAttempt ensures that we get the latest attempt
// the history that we get back from the API is ordered.
func getLatestAttempt(history []*models.ElasticsearchClusterPlanInfo) *models.ElasticsearchClusterPlan {
	// It is highly unlikely to happen, but just in case it does
	// we don't want a panic with a nil pointer error
	if len(history) == 0 {
		return nil
	}

	return history[len(history)-1].Plan
}
