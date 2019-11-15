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

package monitoring

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func TestEnable(t *testing.T) {
	type args struct {
		params EnableParams
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "succeeds",
			args: args{params: EnableParams{
				ClusterParams: util.ClusterParams{
					API: api.NewMock(mock.Response{
						Response: http.Response{
							StatusCode: http.StatusAccepted,
							Status:     http.StatusText(http.StatusAccepted),
							Body:       mock.NewStringBody(`{}`),
						},
					}),
					ClusterID: util.ValidClusterID,
				},
				TargetID: util.ValidClusterID,
			}},
		},
		{
			name: "fails due to API error",
			args: args{params: EnableParams{
				ClusterParams: util.ClusterParams{
					API:       api.NewMock(mock.New500Response(mock.NewStringBody(`{"error": "some error"}`))),
					ClusterID: util.ValidClusterID,
				},
				TargetID: util.ValidClusterID,
			}},
			err: errors.New(`{"error": "some error"}`),
		},
		{
			name: "fails due to empty parameters",
			args: args{},
			err: &multierror.Error{Errors: []error{
				errors.New(`target id "" is invalid`),
				errors.New(`api reference is required for command`),
				errors.New(`id "" is invalid`),
			}},
		},
		{
			name: "fails due to missing API",
			args: args{params: EnableParams{
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
				},
				TargetID: util.ValidClusterID,
			}},
			err: &multierror.Error{Errors: []error{
				errors.New(`api reference is required for command`),
			}},
		},
		{
			name: "fails due to missing cluster ID",
			args: args{params: EnableParams{
				ClusterParams: util.ClusterParams{
					API: new(api.API),
				},
				TargetID: util.ValidClusterID,
			}},
			err: &multierror.Error{Errors: []error{
				errors.New(`id "" is invalid`),
			}},
		},
		{
			name: "fails due to missing target cluster ID",
			args: args{params: EnableParams{
				ClusterParams: util.ClusterParams{
					API:       new(api.API),
					ClusterID: util.ValidClusterID,
				},
			}},
			err: &multierror.Error{Errors: []error{
				errors.New(`target id "" is invalid`),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Enable(tt.args.params); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Enable() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}

func TestDisable(t *testing.T) {
	type args struct {
		params DisableParams
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "succeeds",
			args: args{params: DisableParams{
				ClusterParams: util.ClusterParams{
					API: api.NewMock(mock.Response{
						Response: http.Response{
							StatusCode: http.StatusAccepted,
							Status:     http.StatusText(http.StatusAccepted),
							Body:       mock.NewStringBody(`{}`),
						},
					}),
					ClusterID: util.ValidClusterID,
				},
			}},
		},
		{
			name: "fails due to API error",
			args: args{params: DisableParams{
				ClusterParams: util.ClusterParams{
					API:       api.NewMock(mock.New500Response(mock.NewStringBody(`{"error": "some error"}`))),
					ClusterID: util.ValidClusterID,
				},
			}},
			err: errors.New(`{"error": "some error"}`),
		},
		{
			name: "fails due to empty parameters",
			args: args{},
			err: &multierror.Error{Errors: []error{
				errors.New(`api reference is required for command`),
				errors.New(`id "" is invalid`),
			}},
		},
		{
			name: "fails due to missing API",
			args: args{params: DisableParams{
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
				},
			}},
			err: &multierror.Error{Errors: []error{
				errors.New(`api reference is required for command`),
			}},
		},
		{
			name: "fails due to missing cluster ID",
			args: args{params: DisableParams{
				ClusterParams: util.ClusterParams{
					API: new(api.API),
				},
			}},
			err: &multierror.Error{Errors: []error{
				errors.New(`id "" is invalid`),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Disable(tt.args.params); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Disable() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
