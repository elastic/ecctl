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

package kibana

import (
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func TestRestart(t *testing.T) {
	type args struct {
		params DeploymentParams
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Fails restarting due to parameter validation",
			args: args{params: DeploymentParams{}},
			wantErr: &multierror.Error{Errors: []error{
				errors.New("api reference is required for command"),
				errors.New(`id "" is invalid`),
			}},
		},
		{
			name: "Fails restarting",
			args: args{params: DeploymentParams{
				ID: "2c221bd86b7f48959a59ee3128d5c5e8",
				API: api.NewMock(mock.Response{
					Error: errors.New("error with API"),
				}),
			}},
			wantErr: &url.Error{
				Op:  "Post",
				URL: "https://mock-host/mock-path/clusters/kibana/2c221bd86b7f48959a59ee3128d5c5e8/_restart?cancel_pending=false",
				Err: errors.New("error with API"),
			},
		},
		{
			name: "Succeeds restarting without tracking",
			args: args{params: DeploymentParams{
				ID: "d324608c97154bdba2dff97511d40368",
				API: api.NewMock(mock.Response{Response: http.Response{
					StatusCode: 202,
					Body:       mock.NewStringBody(`{}`),
				}}),
			}},
		},
		{
			name: "Succeeds restarting with tracking",
			args: args{params: DeploymentParams{
				TrackChangeParams: util.NewMockTrackChangeParams(""),
				Track:             true,
				ID:                "d324608c97154bdba2dff97511d40368",
				API: api.NewMock(util.AppendTrackResponses(mock.Response{Response: http.Response{
					StatusCode: 202,
					Body:       mock.NewStringBody(`{}`),
				}})...),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Restart(tt.args.params); !reflect.DeepEqual(tt.wantErr, err) {
				t.Errorf("Restart() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
