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

package instances

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"

	"github.com/elastic/ecctl/pkg/util"
)

func TestStartRouting(t *testing.T) {
	type args struct {
		params Params
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "succeeds with explicit instances",
			args: args{params: Params{
				Instances: []string{
					"instance-000000",
					"instance-000001",
				},
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API: api.NewMock(mock.Response{
						Response: http.Response{
							Status:     http.StatusText(http.StatusAccepted),
							StatusCode: http.StatusAccepted,
							Body:       mock.NewStringBody(`{}`),
						},
					}),
				},
			}},
		},
		{
			name: "fails due to API error",
			args: args{params: Params{
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API: api.NewMock(mock.Response{
						Response: http.Response{
							Status:     http.StatusText(http.StatusInternalServerError),
							StatusCode: http.StatusInternalServerError,
							Body:       mock.NewStringBody(`{}`),
						},
					}),
				},
			}},
			err: errors.New(`unknown error (status 500)`),
		},
		{
			name: "fails due to parameter validation",
			args: args{},
			err:  util.ErrClusterLength,
		},
		{
			name: "fails due to parameter validation (missing API)",
			args: args{params: Params{ClusterParams: util.ClusterParams{
				ClusterID: util.ValidClusterID},
			}},
			err: util.ErrAPIReq,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StartRouting(tt.args.params); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("StartRouting() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}

func TestStopRouting(t *testing.T) {
	type args struct {
		params Params
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "succeeds with explicit instances",
			args: args{params: Params{
				Instances: []string{
					"instance-000000",
					"instance-000001",
				},
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API: api.NewMock(mock.Response{
						Response: http.Response{
							Status:     http.StatusText(http.StatusAccepted),
							StatusCode: http.StatusAccepted,
							Body:       mock.NewStringBody(`{}`),
						},
					}),
				},
			}},
		},
		{
			name: "fails due to API error",
			args: args{params: Params{
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API: api.NewMock(mock.Response{
						Response: http.Response{
							Status:     http.StatusText(http.StatusInternalServerError),
							StatusCode: http.StatusInternalServerError,
							Body:       mock.NewStringBody(`{}`),
						},
					}),
				},
			}},
			err: errors.New(`unknown error (status 500)`),
		},
		{
			name: "fails due to parameter validation",
			args: args{},
			err:  util.ErrClusterLength,
		},
		{
			name: "fails due to parameter validation (missing API)",
			args: args{params: Params{ClusterParams: util.ClusterParams{
				ClusterID: util.ValidClusterID},
			}},
			err: util.ErrAPIReq,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StopRouting(tt.args.params); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("StopRouting() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
