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

package util

import (
	"time"

	"github.com/elastic/cloud-sdk-go/pkg/plan/planutil"
)

const (
	// DefaultRetries retries for the plan.TrackParams which accounts
	// for the Pending plan not being present in the backend, it will retry
	// the request the times specified here. 2 is the default anectdotically
	// because it provides a maximum sleeping time of (PollFrequency * 2)^2
	// or math.Exp2(2). Increasing this value will cause the PlanTracker to
	// sleep for more time than it's required, thus making ecctl less efficient.
	DefaultRetries = 2

	// DefaultPollFrequency is frequency on which the API is polled for updates
	// on pending plans. This value is also used as the cooldown time when used
	// with MaxRetries > 0.
	DefaultPollFrequency = time.Second * 10
)

// SetClusterTracking modifies the TrackChangeParams to track a specific id and
// kind ignoring downstream changes.
func SetClusterTracking(params planutil.TrackChangeParams, id, kind string) planutil.TrackChangeParams {
	params.ResourceID = id
	params.Kind = kind
	params.IgnoreDownstream = true
	return params
}
