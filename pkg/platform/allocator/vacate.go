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

package allocator

import (
	"fmt"
	"strings"
	"time"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/platform_infrastructure"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/plan"
	"github.com/elastic/cloud-sdk-go/pkg/sync/pool"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/elastic/cloud-sdk-go/pkg/util/slice"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

const (
	// DefaultTrackRetries is used when TrackRetries is empty
	DefaultTrackRetries uint8 = 4

	// DefaultTrackFrequency is used when TrackFrequency is 0
	DefaultTrackFrequency = time.Second

	// PlanPendingMessage is used to discard
	PlanPendingMessage = "There is a plan still pending, cancel that or wait for it to complete before restarting"
)

// Vacate drains allocated cluster nodes away from the allocator list either to
// a specific allocator list or we let the constructor decide if that is empty.
// If clusters is set, it will only move the nodes that are part of those IDs.
// If kind is specified, it will only move the clusters that match that kind.
// If none is specified it will add all of the clusters in the allocator.
// The maximum concurrent moves is controlled by the Concurrency parameter.
func Vacate(params *VacateParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	var emptyTimeout pool.Timeout
	if params.PoolTimeout == emptyTimeout {
		params.PoolTimeout = pool.DefaultTimeout
	}

	p, err := pool.NewPool(pool.Params{
		Size:    params.Concurrency,
		Run:     VacateClusterInPool,
		Timeout: params.PoolTimeout,
		Writer:  params.Output,
	})
	if err != nil {
		return err
	}

	// Errors reported here are returned by a dry-run execution of the move api (validation-only flag is used)
	// and we don't want to stop the real vacate.
	// Instead we are returning the validateOnlyErr with the actual vacate validateOnlyErr at the end of the function
	leftovers, hasWork, validateOnlyErr := moveAllocators(params, p)

	if err := p.Start(); err != nil {
		return err
	}

	// If the queue was full prior to starting it, there might be some
	// leftovers, ranging until the leftovers are inexistent as the pool
	// clears out the work items.
	for len(leftovers) > 0 {
		leftovers, _ = p.Add(leftovers...)
	}

	// Wait until all of the items have been processed
	return multierror.Append(waitVacateCompletion(p, hasWork), validateOnlyErr).ErrorOrNil()
}

// moveAllocators ranges over the list of provided allocators and moves the
// nodes off each allocator, finally, returns any leftovers from full pool
// queues, whether or not any work was added to the pool, and potential errors
// returned from API calls.
func moveAllocators(params *VacateParams, p *pool.Pool) ([]pool.Validator, bool, error) {
	var leftovers []pool.Validator
	var merr = new(multierror.Error)
	var hasWork bool
	for _, id := range params.Allocators {
		left, moved, err := moveNodes(id, params, p)
		merr = multierror.Append(merr, err)
		if len(left) > 0 {
			leftovers = append(leftovers, left...)
		}

		if moved {
			hasWork = true
		}
	}
	return leftovers, hasWork, merr.ErrorOrNil()
}

// moveNodes moves all of the nodes off the specified allocator
func moveNodes(id string, params *VacateParams, p *pool.Pool) ([]pool.Validator, bool, error) {
	var merr = new(multierror.Error)
	res, err := params.API.V1API.PlatformInfrastructure.MoveClusters(
		platform_infrastructure.NewMoveClustersParams().
			WithAllocatorID(id).
			WithMoveOnly(params.MoveOnly).
			WithValidateOnly(ec.Bool(true)),
		params.AuthWriter,
	)

	if err != nil {
		var e = errors.Wrapf(err, "allocator %s", id)
		fmt.Fprintln(params.Output, e)
		merr = multierror.Append(merr, e)
		return nil, false, merr.ErrorOrNil()
	}

	if err := CheckVacateFailures(res.Payload.Failures, params.ClusterFilter); err != nil {
		merr = multierror.Append(merr, err)
	}

	work, hasWork := addAllocatorMovesToPool(addAllocatorMovesToPoolParams{
		ID:           id,
		Pool:         p,
		Moves:        res.Payload.Moves,
		VacateParams: params,
	})

	return work, hasWork, merr.ErrorOrNil()
}

// waitVacateCompletion waits for the pool to be finished if there's work
// items that were added and when all of the items have been processed
// stops the pool, and returns a multierror with any leftovers from the
// stopped pool.
func waitVacateCompletion(p *pool.Pool, hasWork bool) error {
	var errs = new(multierror.Error)
	if hasWork {
		if err := p.Wait(); err != nil {
			if e, ok := err.(*multierror.Error); ok {
				errs = e
			}
		}
	}

	if p.Status() < pool.StoppingStatus {
		// Stop the pool once we've finished all the work
		if err := p.Stop(); err != nil && err != pool.ErrStopOperationTimedOut {
			errs = multierror.Append(errs, err)
		}
	}

	leftovers, _ := p.Leftovers()
	for _, lover := range leftovers {
		if params, ok := lover.(*VacateClusterParams); ok {
			errs = multierror.Append(errs, fmt.Errorf(
				"allocator %s: cluster [%s][%s]: %s",
				params.ID, params.ClusterID, params.Kind,
				"was either cancelled or not processed, follow up accordingly",
			))
		}
	}

	return errs.ErrorOrNil()
}

func addAllocatorMovesToPool(params addAllocatorMovesToPoolParams) ([]pool.Validator, bool) {
	var leftovers []pool.Validator
	var vacates = make([]pool.Validator, 0)
	if params.Moves == nil {
		return leftovers, len(vacates) > 0
	}

	var filter = params.VacateParams.ClusterFilter
	var kindFilter = params.VacateParams.KindFilter
	for _, move := range params.Moves.ElasticsearchClusters {
		if len(filter) > 0 && !slice.HasString(filter, *move.ClusterID) {
			continue
		}

		var kind = "elasticsearch"
		if kindFilter != "" && kind != kindFilter {
			break
		}
		vacates = append(vacates, newVacateClusterParams(params, *move.ClusterID, kind))
	}

	for _, move := range params.Moves.KibanaClusters {
		if len(filter) > 0 && !slice.HasString(filter, *move.ClusterID) {
			continue
		}

		var kind = "kibana"
		if kindFilter != "" && kind != kindFilter {
			break
		}
		vacates = append(vacates, newVacateClusterParams(params, *move.ClusterID, kind))
	}

	for _, move := range params.Moves.ApmClusters {
		if len(filter) > 0 && !slice.HasString(filter, *move.ClusterID) {
			continue
		}

		var kind = "apm"
		if kindFilter != "" && kind != kindFilter {
			break
		}
		vacates = append(vacates, newVacateClusterParams(params, *move.ClusterID, kind))
	}

	if leftover, _ := params.Pool.Add(vacates...); len(leftover) > 0 {
		leftovers = append(leftovers, leftover...)
	}

	return leftovers, len(vacates) > 0
}

func newVacateClusterParams(params addAllocatorMovesToPoolParams, id, kind string) *VacateClusterParams {
	return &VacateClusterParams{
		API:                 params.VacateParams.API,
		ID:                  params.ID,
		Kind:                kind,
		ClusterID:           id,
		ClusterFilter:       params.VacateParams.ClusterFilter,
		PreferredAllocators: params.VacateParams.PreferredAllocators,
		MaxPollRetries:      params.VacateParams.MaxPollRetries,
		TrackFrequency:      params.VacateParams.TrackFrequency,
		Output:              params.VacateParams.Output,
		MoveOnly:            params.VacateParams.MoveOnly,
		PlanOverrides:       params.VacateParams.PlanOverrides,
	}
}

// VacateClusterInPool vacates a cluster from an allocator, complying
// with the pool.RunFunc signature.
func VacateClusterInPool(p pool.Validator) error {
	if p == nil {
		return errors.New("allocator vacate: params cannot be nil")
	}

	if params, ok := p.(*VacateClusterParams); ok {
		return VacateCluster(params)
	}

	return errors.New("allocator vacate: failed casting parameters to *VacateClusterParams")
}

// VacateCluster moves a cluster node off an allocator.
func VacateCluster(params *VacateClusterParams) error {
	params, err := fillVacateClusterParams(params)
	if err != nil {
		return err
	}

	if err := moveClusterByType(params); err != nil {
		return err
	}

	return trackMovedCluster(params)
}

// fillVacateClusterParams validates the parameters and fills any missing
// properties that are set to a default if empty. Performs a Get on the
// allocator to discover the allocator health if AllocatorDown is nil.
func fillVacateClusterParams(params *VacateClusterParams) (*VacateClusterParams, error) {
	if params == nil {
		return nil, errors.New("allocator vacate: params cannot be nil")
	}

	if err := params.Validate(); err != nil {
		return nil, errors.Wrap(err, vacateWrapMessage(params, "parameter validation"))
	}

	if params.AllocatorDown == nil {
		alloc, err := Get(GetParams{API: params.API, ID: params.ID})
		if err != nil {
			return nil, errors.Wrap(err, vacateWrapMessage(params, "allocator health autodiscovery"))
		}
		params.AllocatorDown = ec.Bool(!*alloc.Status.Connected || !*alloc.Status.Healthy)
	}

	if params.MaxPollRetries == 0 {
		params.MaxPollRetries = DefaultTrackRetries
	}

	if params.TrackFrequency.Nanoseconds() == 0 {
		params.TrackFrequency = DefaultTrackFrequency
	}

	return params, nil
}

// returns a wrapping message from the VacateClusterParams
func vacateWrapMessage(params *VacateClusterParams, ctx string) string {
	if params == nil {
		return ""
	}
	return fmt.Sprintf("allocator %s: cluster [%s][%s]: %s",
		params.ID, params.ClusterID, params.Kind, ctx,
	)
}

// newMoveClusterParams
func newMoveClusterParams(params *VacateClusterParams) (*platform_infrastructure.MoveClustersByTypeParams, error) {
	res, err := params.API.V1API.PlatformInfrastructure.MoveClusters(
		platform_infrastructure.NewMoveClustersParams().
			WithAllocatorDown(params.AllocatorDown).
			WithMoveOnly(params.MoveOnly).
			WithAllocatorID(params.ID).
			WithValidateOnly(ec.Bool(true)),
		params.AuthWriter,
	)
	if err != nil {
		return nil, errors.Wrap(api.UnwrapError(err), vacateWrapMessage(params, "validate_only"))
	}

	req := ComputeVacateRequest(res.Payload.Moves,
		[]string{params.ClusterID},
		params.PreferredAllocators,
		params.PlanOverrides,
	)

	var moveParams = platform_infrastructure.NewMoveClustersByTypeParams().
		WithAllocatorID(params.ID).
		WithAllocatorDown(params.AllocatorDown).
		WithBody(req)

	if len(req.ElasticsearchClusters) > 0 {
		moveParams.SetClusterType(ec.String("elasticsearch"))
	}

	if len(req.KibanaClusters) > 0 {
		moveParams.SetClusterType(ec.String("kibana"))
	}

	if len(req.ApmClusters) > 0 {
		moveParams.SetClusterType(ec.String("apm"))
	}

	return moveParams, nil
}

// moveClusterByType moves a cluster's node from its allocator
func moveClusterByType(params *VacateClusterParams) error {
	moveParams, err := newMoveClusterParams(params)
	if err != nil {
		return err
	}

	res, err := params.API.V1API.PlatformInfrastructure.MoveClustersByType(
		moveParams,
		params.AuthWriter,
	)

	if err != nil {
		return errors.Wrap(api.UnwrapError(err), vacateWrapMessage(params, "cluster move"))
	}

	return CheckVacateFailures(res.Payload.Failures, params.ClusterFilter)
}

func trackMovedCluster(params *VacateClusterParams) error {
	channel, err := plan.Track(plan.TrackParams{
		API:           params.API,
		ID:            params.ClusterID,
		Kind:          params.Kind,
		PollFrequency: params.TrackFrequency,
		MaxRetries:    params.MaxPollRetries,
	})
	if err != nil {
		return errors.Wrap(err, vacateWrapMessage(params, "track"))
	}

	return plan.Stream(channel, params.Output)
}

// CheckVacateFailures iterates over the list of failures returning
// a multierror with any of the failures found.
func CheckVacateFailures(failures *models.MoveClustersDetails, filter []string) error {
	if failures == nil {
		return nil
	}

	var merr = new(multierror.Error)
	for _, failure := range failures.ElasticsearchClusters {
		if len(filter) > 0 && !slice.HasString(filter, *failure.ClusterID) {
			continue
		}

		var ferr error
		if len(failure.Errors) > 0 {
			var err = failure.Errors[0]
			ferr = fmt.Errorf("code: %s, message: %s", *err.Code, *err.Message)
		}
		if !strings.Contains(ferr.Error(), PlanPendingMessage) {
			merr = multierror.Append(merr,
				fmt.Errorf("cluster [%s][elasticsearch] failed vacating, reason: %s", *failure.ClusterID, ferr),
			)
		}
	}

	for _, failure := range failures.KibanaClusters {
		if len(filter) > 0 && !slice.HasString(filter, *failure.ClusterID) {
			continue
		}

		var ferr error
		if len(failure.Errors) > 0 {
			var err = failure.Errors[0]
			ferr = fmt.Errorf("code: %s, message: %s", *err.Code, *err.Message)
		}

		if !strings.Contains(ferr.Error(), PlanPendingMessage) {
			merr = multierror.Append(merr,
				fmt.Errorf("cluster [%s][kibana] failed vacating, reason: %s", *failure.ClusterID, ferr),
			)
		}
	}

	for _, failure := range failures.ApmClusters {
		if len(filter) > 0 && !slice.HasString(filter, *failure.ClusterID) {
			continue
		}

		var ferr error
		if len(failure.Errors) > 0 {
			var err = failure.Errors[0]
			ferr = fmt.Errorf("code: %s, message: %s", *err.Code, *err.Message)
		}

		if !strings.Contains(ferr.Error(), PlanPendingMessage) {
			merr = multierror.Append(merr,
				fmt.Errorf("cluster [%s][apm] failed vacating, reason: %s", *failure.ClusterID, ferr),
			)
		}
	}

	return merr.ErrorOrNil()
}

// ComputeVacateRequest filters the tentative cluster that would be moved and
// filters those by ID if it's specified, also setting any preferred allocators
// if that is sent. Any cluster plan overrides will be set in this function.
func ComputeVacateRequest(pr *models.MoveClustersDetails, clusters, to []string, overrides PlanOverrides) *models.MoveClustersRequest {
	var req models.MoveClustersRequest
	for _, c := range pr.ElasticsearchClusters {
		if len(clusters) > 0 && !slice.HasString(clusters, *c.ClusterID) {
			continue
		}

		if overrides.SkipSnapshot != nil {
			c.CalculatedPlan.PlanConfiguration.SkipSnapshot = overrides.SkipSnapshot
		}

		if overrides.SkipDataMigration != nil {
			c.CalculatedPlan.PlanConfiguration.SkipDataMigration = overrides.SkipDataMigration
		}

		if overrides.OverrideFailsafe != nil {
			c.CalculatedPlan.PlanConfiguration.OverrideFailsafe = overrides.OverrideFailsafe
		}

		c.CalculatedPlan.PlanConfiguration.PreferredAllocators = to
		req.ElasticsearchClusters = append(req.ElasticsearchClusters,
			&models.MoveElasticsearchClusterConfiguration{
				ClusterIds:   []string{*c.ClusterID},
				PlanOverride: c.CalculatedPlan,
			},
		)
	}

	for _, c := range pr.KibanaClusters {
		if len(clusters) > 0 && !slice.HasString(clusters, *c.ClusterID) {
			continue
		}

		c.CalculatedPlan.PlanConfiguration.PreferredAllocators = to
		req.KibanaClusters = append(req.KibanaClusters,
			&models.MoveKibanaClusterConfiguration{
				ClusterIds:   []string{*c.ClusterID},
				PlanOverride: c.CalculatedPlan,
			},
		)
	}

	for _, c := range pr.ApmClusters {
		if len(clusters) > 0 && !slice.HasString(clusters, *c.ClusterID) {
			continue
		}

		c.CalculatedPlan.PlanConfiguration.PreferredAllocators = to
		req.ApmClusters = append(req.ApmClusters,
			&models.MoveApmClusterConfiguration{
				ClusterIds:   []string{*c.ClusterID},
				PlanOverride: c.CalculatedPlan,
			},
		)
	}

	return &req
}
