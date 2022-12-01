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
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/output"
	"github.com/elastic/cloud-sdk-go/pkg/plan"
	"github.com/elastic/cloud-sdk-go/pkg/plan/planutil"

	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/formatter"
	"github.com/elastic/ecctl/pkg/util"
)

func TestNewTrackParams(t *testing.T) {
	textFormatter := &formatter.Text{}
	some := models.APIKeyResponse{}
	someOutputDevice := output.NewDevice(io.Discard)
	emptyMockAAPI := api.NewMock()
	ecctlInstance := ecctl.App{
		API:       emptyMockAAPI,
		Formatter: textFormatter,
		Config: ecctl.Config{
			Output:       "text",
			OutputDevice: someOutputDevice,
		},
	}

	type args struct {
		params TrackParamsConfig
	}
	tests := []struct {
		name string
		args args
		want TrackParams
	}{
		{
			name: "empty FrequencyConfig obtains the default one",
			args: args{params: TrackParamsConfig{
				DeploymentID: util.ValidClusterID,
				Track:        true,
				Template:     "some/template",
				Response:     some,
				App:          &ecctlInstance,
			}},
			want: TrackParams{
				TrackChangeParams: planutil.TrackChangeParams{
					TrackChangeParams: plan.TrackChangeParams{
						DeploymentID: util.ValidClusterID,
						API:          emptyMockAAPI,
						Config:       DefaultTrackFrequencyConfig,
					},
					Writer: someOutputDevice,
					Format: "text",
				},
				Template:  "some/template",
				Formatter: textFormatter,
				Track:     true,
				Response:  some,
			},
		},
		{
			name: "Specifies a TrackFrequencyConfig",
			args: args{params: TrackParamsConfig{
				DeploymentID: util.ValidClusterID,
				Track:        true,
				Template:     "some/template",
				Response:     some,
				App:          &ecctlInstance,
				FrequencyConfig: &plan.TrackFrequencyConfig{
					PollFrequency: time.Second * 10,
					MaxRetries:    5,
				},
			}},
			want: TrackParams{
				TrackChangeParams: planutil.TrackChangeParams{
					TrackChangeParams: plan.TrackChangeParams{
						DeploymentID: util.ValidClusterID,
						API:          emptyMockAAPI,
						Config: plan.TrackFrequencyConfig{
							PollFrequency: time.Second * 10,
							MaxRetries:    5,
						},
					},
					Writer: someOutputDevice,
					Format: "text",
				},
				Template:  "some/template",
				Formatter: textFormatter,
				Track:     true,
				Response:  some,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTrackParams(tt.args.params)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTrackParams() = %v, want %v", got, tt.want)
			}
		})
	}
}
