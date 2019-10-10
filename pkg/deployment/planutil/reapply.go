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

package planutil

import (
	"errors"

	"github.com/elastic/cloud-sdk-go/pkg/models"
)

var (
	errIDCannotBeEmpty       = errors.New("cluster id cannot be empty")
	errIDCannotMixStrategies = errors.New("cannot specify multi strategies")
)

// ReapplyParams contains the parameters required to call a plan reapply action
type ReapplyParams struct {
	ID                   string
	HidePlan             bool `kebabcase:"Doesn't print the plan before reapplying"`
	Default              bool `kebabcase:"Overwrites the strategy to the default one"`
	Rolling              bool `kebabcase:"Overwrites the strategy to rolling"`
	GrowAndShrink        bool `kebabcase:"Overwrites the strategy to grow and shrink"`
	RollingGrowAndShrink bool `kebabcase:"Overwrites the strategy to rolling grow and shrink (one at a time)"`
	RollingAll           bool `kebabcase:"Overwrites the strategy to apply the change in all the instances at a time (causes downtime)"`
	Reallocate           bool `kebabcase:"Forces creation of new instances"`
	ExtendedMaintenance  bool `kebabcase:"Stops routing to the cluster instances after the plan has been applied"`
	OverrideFailsafe     bool `kebabcase:"Overrides failsafe at the constructor level that prevent bad things from happening"`
}

// Validate returns an error if the parameters are invalid
func (p ReapplyParams) Validate() error {
	if len(p.ID) < 32 {
		return errIDCannotBeEmpty
	}

	// Valid parameters
	var (
		validDefault              = !p.Rolling && !p.GrowAndShrink && !p.RollingGrowAndShrink
		validRolling              = !p.Default && p.Rolling && !p.GrowAndShrink && !p.RollingGrowAndShrink
		validGrowAndShrink        = !p.Default && !p.Rolling && p.GrowAndShrink && !p.RollingGrowAndShrink
		validRollingGrowAndShrink = !p.Default && !p.Rolling && !p.GrowAndShrink && p.RollingGrowAndShrink
	)

	// We don't care if the default strategy is true or false, but care about the others
	if validDefault || validRolling || validGrowAndShrink || validRollingGrowAndShrink {
		return nil
	}

	return errIDCannotMixStrategies
}

// Strategy returns a plan strategy from the specified ReapplyParams
// If all strategies are false, nil will be returned, which won't alter
// the previously specified strategy.
func (p ReapplyParams) Strategy() *models.PlanStrategy {
	if p.Rolling {
		return RollingByNameStrategy
	}

	if p.RollingAll {
		return MajorUpgradeStrategy
	}

	if p.GrowAndShrink {
		return GrowAndShrinkStrategy
	}

	if p.Default {
		return DefaultPlanStrategy
	}

	if p.RollingGrowAndShrink {
		return RollingGrowAndShrinkStrategy
	}

	return nil
}
