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

func TestResume(t *testing.T) {
	type args struct {
		params Params
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "returns error on missing cluster id",
			err:  util.ErrClusterLength,
		},
		{
			name: "returns error on missing API",
			args: args{params: Params{ClusterParams: util.ClusterParams{
				ClusterID: util.ValidClusterID,
			}}},
			err: util.ErrAPIReq,
		},
		{
			name: "succeeds",
			args: args{params: Params{
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
			name: "fails on api error",
			args: args{params: Params{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Resume(tt.args.params); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Resume() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}

func TestPause(t *testing.T) {
	type args struct {
		params Params
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "returns error on missing cluster id",
			err:  util.ErrClusterLength,
		},
		{
			name: "returns error on missing API",
			args: args{params: Params{ClusterParams: util.ClusterParams{
				ClusterID: util.ValidClusterID,
			}}},
			err: util.ErrAPIReq,
		},
		{
			name: "succeeds",
			args: args{params: Params{
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
			name: "fails on api error",
			args: args{params: Params{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Pause(tt.args.params); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Pause() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
