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

package apm

import (
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func TestResync(t *testing.T) {
	type args struct {
		params ResyncParams
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Fails due to parameter validation (Cluster ID)",
			args: args{params: ResyncParams{}},
			wantErr: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				errors.New(`id "" is invalid`),
			}},
		},
		{
			name: "Fails due to parameter validation (API)",
			args: args{params: ResyncParams{
				ClusterID: "d324608c97154bdba2dff97511d40368",
			}},
			wantErr: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
			}},
		},
		{
			name: "Fails due to unknown API response",
			args: args{params: ResyncParams{
				ClusterID: "d324608c97154bdba2dff97511d40368",
				API: api.NewMock(mock.Response{Response: http.Response{
					StatusCode: http.StatusForbidden,
					Body:       mock.NewStringBody(`{}`),
				}}),
			}},
			wantErr: errors.New("{}"),
		},
		{
			name: "Fails due to API error",
			args: args{params: ResyncParams{
				ClusterID: "d324608c97154bdba2dff97511d40368",
				API: api.NewMock(mock.Response{
					Error: errors.New("error with API"),
				}),
			}},
			wantErr: &url.Error{
				Op:  "Post",
				URL: "https://mock-host/mock-path/clusters/apm/d324608c97154bdba2dff97511d40368/_resync",
				Err: errors.New("error with API"),
			},
		},
		{
			name: "Succeeds to resynchronize APM cluster without errors",
			args: args{params: ResyncParams{
				ClusterID: "d324608c97154bdba2dff97511d40368",
				API: api.NewMock(mock.Response{Response: http.Response{
					StatusCode: http.StatusOK,
					Body:       mock.NewStringBody(`{}`),
				}}),
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Resync(tt.args.params); !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Resync() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResyncAll(t *testing.T) {
	type args struct {
		params ResyncAllParams
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		want    *models.ModelVersionIndexSynchronizationResults
	}{
		{
			name:    "Fails due to parameter validation (API)",
			args:    args{params: ResyncAllParams{}},
			wantErr: errors.New("api reference is required for command"),
		},
		{
			name: "Fails due to unknown API response",
			args: args{params: ResyncAllParams{
				API: api.NewMock(mock.Response{Response: http.Response{
					StatusCode: http.StatusForbidden,
					Body:       mock.NewStringBody(`{"error": "some forbidden error"}`),
				}}),
			}},
			wantErr: errors.New(`{"error": "some forbidden error"}`),
		},
		{
			name: "Fails due to API error",
			args: args{params: ResyncAllParams{
				API: api.NewMock(mock.Response{
					Error: errors.New("error with API"),
				}),
			}},
			wantErr: &url.Error{
				Op:  "Post",
				URL: "https://mock-host/mock-path/clusters/apm/_resync?skip_matching_version=true",
				Err: errors.New("error with API"),
			},
		},
		{
			name: "Succeeds to resynchronize APM instance without errors",
			args: args{params: ResyncAllParams{
				API: api.NewMock(mock.Response{Response: http.Response{
					StatusCode: http.StatusAccepted,
					Body:       mock.NewStringBody(`{}`),
				}}),
			}},
			want: &models.ModelVersionIndexSynchronizationResults{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ResyncAll(tt.args.params)
			if !reflect.DeepEqual(tt.wantErr, err) {
				t.Errorf("ResyncAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ResyncAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
