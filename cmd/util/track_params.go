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

package cmdutil

import (
	"time"

	"github.com/elastic/cloud-sdk-go/pkg/plan"
	"github.com/elastic/cloud-sdk-go/pkg/plan/planutil"

	"github.com/elastic/ecctl/pkg/ecctl"
)

// DefaultTrackFrequencyConfig provides sane defaults for Plan change tracking.
var DefaultTrackFrequencyConfig = plan.TrackFrequencyConfig{
	PollFrequency: time.Second * 5,
	MaxRetries:    3,
}

// TrackParamsConfig is used to create TrackParams which print the output / track.
type TrackParamsConfig struct {
	App             *ecctl.App
	DeploymentID    string
	ResourceID      string
	Kind            string
	Template        string
	Response        interface{}
	Track           bool
	FrequencyConfig *plan.TrackFrequencyConfig
}

// NewTrackParams creates a TrackParams structure from the config.
func NewTrackParams(params TrackParamsConfig) TrackParams {
	if params.FrequencyConfig == nil {
		params.FrequencyConfig = &DefaultTrackFrequencyConfig
	}

	return TrackParams{
		TrackChangeParams: planutil.TrackChangeParams{
			Writer: params.App.Config.OutputDevice,
			Format: params.App.Config.Output,
			TrackChangeParams: plan.TrackChangeParams{
				API:          params.App.API,
				DeploymentID: params.DeploymentID,
				ResourceID:   params.ResourceID,
				Kind:         params.Kind,
				Config:       *params.FrequencyConfig,
			},
		},
		Formatter: params.App.Formatter,
		Track:     params.Track,
		Response:  params.Response,
		Template:  params.Template,
	}
}
