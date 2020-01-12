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
	"errors"
	"time"

	"github.com/elastic/cloud-sdk-go/pkg/output"
	"github.com/elastic/cloud-sdk-go/pkg/plan"
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

// TrackParams are intended to be used as a field in other params structs
// that aim to perform plan tracking.
type TrackParams struct {
	Output        *output.Device
	PollFrequency time.Duration
	Track         bool
	MaxRetries    uint8
}

// TrackClusterParams is consumed by TrackCluster
type TrackClusterParams struct {
	plan.TrackParams
	Output *output.Device
}

// Validate ensures that the parameters are usable by the consuming function.
func (params TrackClusterParams) Validate() error {
	if params.Output == nil {
		return errors.New("track: Output cannot be empty")
	}
	return params.TrackParams.Validate()
}

// Validate ensures that the parameters are usable by the consuming function.
func (params TrackParams) Validate() error {
	if params.Track && params.Output == nil {
		return errors.New("track params: output device cannot be empty")
	}
	return nil
}

// TrackCluster tracks a cluster by its ID and kind to a specified output.
func TrackCluster(params TrackClusterParams) error {
	if params.MaxRetries == 0 {
		params.MaxRetries = DefaultRetries
	}

	if params.PollFrequency.Nanoseconds() == 0 {
		params.PollFrequency = DefaultPollFrequency
	}

	if err := params.Validate(); err != nil {
		return err
	}

	c, err := plan.Track(params.TrackParams)
	if err != nil {
		return err
	}

	return plan.Stream(c, params.Output)
}
