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
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/output"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func TestEnable(t *testing.T) {
	type args struct {
		params DeploymentParams
	}
	tests := []struct {
		name string
		args args
		want *models.ClusterCrudResponse
		err  error
	}{
		{
			name: "Validate parameters fails on empty params",
			args: args{params: DeploymentParams{}},
			want: nil,
			err: &multierror.Error{
				Errors: []error{
					errors.New("api reference is required for command"),
					errors.New("id \"\" is invalid"),
				},
			},
		},
		{
			name: "Validate parameters fails on invalid params",
			args: args{params: DeploymentParams{
				API: new(api.API),
				ID:  "Some invalid ID",
			}},
			want: nil,
			err: &multierror.Error{
				Errors: []error{
					errors.New("id \"Some invalid ID\" is invalid"),
				},
			},
		},
		{
			name: "Fails to enable",
			args: args{params: DeploymentParams{
				API: api.NewMock(mock.Response{
					Error: errors.New("error with API"),
				}),
				ID: "b786acd298292c2d521c0e8741761b4d",
			}},
			err: &url.Error{
				Op:  "Post",
				URL: "https://mock-host/mock-path/clusters/kibana?validate_only=false",
				Err: errors.New("error with API"),
			},
		},
		{
			name: "Succeeds enabling with valid params",
			args: args{params: DeploymentParams{
				API: api.NewMock(mock.Response{Response: http.Response{
					StatusCode: 201,
					Body: mock.NewStructBody(models.ClusterCrudResponse{
						KibanaClusterID: "bae445a42a44e72f2f27e4b149aa496d",
					}),
				}}),
				ID: "b786acd298292c2d521c0e8741761b4d",
			}},
			want: &models.ClusterCrudResponse{
				KibanaClusterID: "bae445a42a44e72f2f27e4b149aa496d",
			},
		},
		{
			name: "Succeeds enabling with tracking",
			args: args{params: DeploymentParams{
				TrackParams: util.TrackParams{
					Track:         true,
					Output:        output.NewDevice(new(bytes.Buffer)),
					PollFrequency: time.Millisecond,
					MaxRetries:    1,
				},
				API: api.NewMock(
					util.AppendTrackResponses(mock.Response{Response: http.Response{
						StatusCode: 201,
						Body: mock.NewStructBody(models.ClusterCrudResponse{
							KibanaClusterID: "bae445a42a44e72f2f27e4b149aa496d",
						}),
					}})...,
				),
				ID: "b786acd298292c2d521c0e8741761b4d",
			}},
			want: &models.ClusterCrudResponse{
				KibanaClusterID: "bae445a42a44e72f2f27e4b149aa496d",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Enable(tt.args.params)
			if !reflect.DeepEqual(tt.err, err) {
				t.Errorf("Enable() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Enable() = %v, want %v", got, tt.want)
			}
		})
	}
}
