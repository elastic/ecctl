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
	"github.com/elastic/cloud-sdk-go/pkg/models"
)

const (
	// AllGroupAttribute represents the rolling strategy property
	// will cause the action to take place to all instances at the
	// same time. Potentially causing downtime in the cluster.
	AllGroupAttribute = "__all__"

	// NameGroupAttribute represents the rolling strategy porperty
	// will cause the action to take place on one instance at a time.
	NameGroupAttribute = "__name__"

	// ZoneGroupAttribute represents the rolling strategy porperty
	// will cause the action to take place on one instance at a time.
	ZoneGroupAttribute = "logical_zone_name"
)

var (
	// DefaultPlanStrategy will use the "default", strategy
	// which defaults to Grow and Shrink
	DefaultPlanStrategy = &models.PlanStrategy{}

	// GrowAndShrinkStrategy will cause the plan to create new instances first
	// and once the changes are finished it will remove the old instances
	GrowAndShrinkStrategy = &models.PlanStrategy{
		GrowAndShrink: new(models.RollingStrategyConfig),
	}

	// RollingGrowAndShrinkStrategy will cause the plan to perform a grow and shrink
	// but instead of creating all the new instances at once, it will do it rolling.
	// This reduces the amount of available capacity for any change on clusters 1>.
	RollingGrowAndShrinkStrategy = &models.PlanStrategy{
		RollingGrowAndShrink: new(models.RollingStrategyConfig),
	}

	// MajorUpgradeStrategy represents the strategy that will
	// be used in Major version upgrades
	MajorUpgradeStrategy = &models.PlanStrategy{
		Rolling: &models.RollingStrategyConfig{
			GroupBy: AllGroupAttribute,
		},
	}

	// RollingByNameStrategy represents the strategy that will
	// be used when the plan wants to be applied one node at a
	// time without causing downtime
	RollingByNameStrategy = &models.PlanStrategy{
		Rolling: &models.RollingStrategyConfig{
			GroupBy: NameGroupAttribute,
		},
	}
)
