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
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/output"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func TestStopParams_Validate(t *testing.T) {
	type args struct {
		params StopParams
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "Validate should return all possible errors",
			args: args{params: StopParams{
				DeploymentParams: DeploymentParams{
					ID:  "",
					API: nil,
				},
			}},
			err: &multierror.Error{
				Errors: []error{
					errors.New("api reference is required for command"),
					errors.New(`id "" is invalid`),
				},
			},
			wantErr: true,
		},
		{
			name: "Validate should pass if all params are properly set",
			args: args{params: StopParams{
				Hide: true,
				DeploymentParams: DeploymentParams{
					ID: "2c221bd86b7f48959a59ee3128d5c5e8",
					API: api.NewMock(mock.Response{Response: http.Response{
						StatusCode: 202,
						Body:       mock.NewStringBody(`{}`),
					}}),
				}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Stop(tt.args.params)

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.err == nil {
				t.Errorf("Validate() expected errors = '%v' but no errors returned", tt.err)
			}

			if tt.wantErr && !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Validate() expected errors = '%v' but got %v", tt.err, err)
			}
		})
	}
}

func TestStop(t *testing.T) {
	type args struct {
		params StopParams
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Fails stopping",
			args: args{params: StopParams{
				Hide: true,
				DeploymentParams: DeploymentParams{
					ID: "2c221bd86b7f48959a59ee3128d5c5e8",
					API: api.NewMock(mock.Response{
						Error: errors.New("error with API"),
					}),
				},
			}},
			wantErr: &url.Error{
				Op:  "Post",
				URL: "https://mock-host/mock-path/clusters/kibana/2c221bd86b7f48959a59ee3128d5c5e8/_shutdown?hide=true",
				Err: errors.New("error with API"),
			},
		},
		{
			name: "Succeeds stopping",
			args: args{params: StopParams{
				Hide: true,
				DeploymentParams: DeploymentParams{
					ID: "2c221bd86b7f48959a59ee3128d5c5e8",
					API: api.NewMock(mock.Response{Response: http.Response{
						StatusCode: 202,
						Body:       mock.NewStringBody(`{}`),
					}}),
				}},
			},
		},
		{
			name: "Succeeds stopping with tracking",
			args: args{params: StopParams{
				Hide: true,
				DeploymentParams: DeploymentParams{
					ID: "2c221bd86b7f48959a59ee3128d5c5e8",
					TrackParams: util.TrackParams{
						Track:         true,
						Output:        output.NewDevice(new(bytes.Buffer)),
						PollFrequency: time.Millisecond,
						MaxRetries:    1,
					},
					API: api.NewMock(util.AppendTrackResponses(
						mock.Response{Response: http.Response{
							StatusCode: 202,
							Body:       mock.NewStringBody(`{}`),
						}},
					)...),
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Stop(tt.args.params); !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
