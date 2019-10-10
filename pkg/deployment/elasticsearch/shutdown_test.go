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

package elasticsearch

import (
	"bytes"
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/output"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func TestShutdownCluster(t *testing.T) {
	type args struct {
		params ShutdownClusterParams
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "succeeds without tracking",
			args: args{params: ShutdownClusterParams{
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API: api.NewMock(mock.Response{
						Response: http.Response{
							StatusCode: 202,
							Body:       mock.NewStringBody(`{}`),
						},
					}),
				},
			}},
		},
		{
			name: "succeeds with tracking",
			args: args{params: ShutdownClusterParams{
				TrackParams: util.TrackParams{
					Track:         true,
					PollFrequency: time.Nanosecond,
					MaxRetries:    1,
					Output:        output.NewDevice(new(bytes.Buffer)),
				},
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API: api.NewMock(util.AppendTrackResponses(mock.Response{
						Response: http.Response{
							StatusCode: 202,
							Body:       mock.NewStringBody(`{}`),
						},
					})...),
				},
			}},
		},
		{
			name: "fails due to API error",
			args: args{params: ShutdownClusterParams{
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API: api.NewMock(mock.Response{
						Response: http.Response{
							StatusCode: 500,
							Body:       mock.NewStringBody(`{}`),
						},
					}),
				},
			}},
			err: errors.New("unknown error (status 500)"),
		},
		{
			name: "fails due to parameter validation (Cluster ID)",
			args: args{},
			err: &multierror.Error{Errors: []error{
				util.ErrClusterLength,
			}},
		},
		{
			name: "fails due to parameter validation (API)",
			args: args{params: ShutdownClusterParams{
				ClusterParams: util.ClusterParams{ClusterID: util.ValidClusterID},
			}},
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
			}},
		},
		{
			name: "fails due to parameter validation (track params)",
			args: args{params: ShutdownClusterParams{
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API:       new(api.API),
				},
				TrackParams: util.TrackParams{Track: true},
			}},
			err: &multierror.Error{Errors: []error{
				errors.New("track params: output device cannot be empty"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ShutdownCluster(tt.args.params); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("ShutdownCluster() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
