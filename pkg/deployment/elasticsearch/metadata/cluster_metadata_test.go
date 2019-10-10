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

package metadata

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

var metadataSettingsResponse = `{
	"owner_id": "123456",
	"hidden": false,
	"resources": {
		"cpu": {
			"boost": true,
			"hard_limit": false
		}
	}
}`

var errorResponse = `{
	"errors": [{
		"code": "clusters.cluster_not_found",
		"message": "Cluster [66cb4698b42e4ec697470c4a66981534] not found"
	}]
}`

func TestGetSettings(t *testing.T) {
	type args struct {
		params util.ClusterParams
	}
	tests := []struct {
		name    string
		args    args
		want    *models.ClusterMetadataSettings
		wantErr bool
		error   error
	}{
		{
			name: "Get cluster metadata settings succeeds",
			args: args{params: util.ClusterParams{
				ClusterID: util.ValidClusterID,
				API: api.NewMock(mock.Response{Response: http.Response{
					Body:       mock.NewStringBody(metadataSettingsResponse),
					StatusCode: 200,
				}}),
			}},
			want: &models.ClusterMetadataSettings{
				OwnerID: "123456",
				Hidden:  ec.Bool(false),
				Resources: &models.ClusterMetadataResourcesSettings{
					CPU: &models.ClusterMetadataCPUResourcesSettings{
						Boost:     ec.Bool(true),
						HardLimit: ec.Bool(false),
					},
				},
			},
			wantErr: false,
			error:   nil,
		},
		{
			name: "Get cluster metadata settings fails if api returns an error",
			args: args{params: util.ClusterParams{
				ClusterID: util.ValidClusterID,
				API: api.NewMock(mock.Response{Response: http.Response{
					Body:       mock.NewStringBody(errorResponse),
					StatusCode: 400,
				}}),
			}},
			want:    nil,
			wantErr: true,
			error:   errors.New("unknown error (status 400)"),
		},
		{
			name: "Get cluster metadata settings fails if api reference is empty",
			args: args{params: util.ClusterParams{
				ClusterID: util.ValidClusterID,
			}},
			want:    nil,
			wantErr: true,
			error:   errors.New("api reference is required for command"),
		},
		{
			name: "Get cluster metadata settings fails if cluster ID is invalid",
			args: args{params: util.ClusterParams{
				ClusterID: util.InvalidClusterID,
			}},
			want:    nil,
			wantErr: true,
			error:   errors.New("cluster id should have a length of 32 characters"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSettings(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClusterMetadataSettings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.error != nil && !reflect.DeepEqual(err, tt.error) {
				t.Errorf("GetClusterMetadataSettings() actual error = '%v', want error '%v'", err, tt.error)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetClusterMetadataSettings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateClusterMetadataSettings(t *testing.T) {
	type args struct {
		params UpdateSettingsParams
	}
	tests := []struct {
		name    string
		args    args
		want    *models.ClusterMetadataSettings
		wantErr bool
		error   error
	}{
		{
			name: "Update cluster metadata settings succeeds",
			args: args{params: UpdateSettingsParams{
				SettingsParams: &models.ClusterMetadataSettings{},
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(metadataSettingsResponse),
						StatusCode: 200,
					}}),
				},
			}},
			want: &models.ClusterMetadataSettings{
				OwnerID: "123456",
				Hidden:  ec.Bool(false),
				Resources: &models.ClusterMetadataResourcesSettings{
					CPU: &models.ClusterMetadataCPUResourcesSettings{
						Boost:     ec.Bool(true),
						HardLimit: ec.Bool(false),
					},
				},
			},
			wantErr: false,
			error:   nil,
		},
		{
			name: "Get cluster metadata settings fails if api returns an error",
			args: args{params: UpdateSettingsParams{
				SettingsParams: &models.ClusterMetadataSettings{},
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(errorResponse),
						StatusCode: 400,
					}}),
				},
			}},
			want:    nil,
			wantErr: true,
			error:   errors.New("unknown error (status 400)"),
		},
		{
			name:    "Get cluster metadata settings fails if settings reference is empty",
			args:    args{params: UpdateSettingsParams{}},
			want:    nil,
			wantErr: true,
			error: &multierror.Error{Errors: []error{
				errors.New("metadata settings params cannot be empty"),
				util.ErrClusterLength,
			}},
		},
		{
			name: "Get cluster metadata settings fails if api reference is empty",
			args: args{params: UpdateSettingsParams{
				SettingsParams: &models.ClusterMetadataSettings{},
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
				},
			}},
			want:    nil,
			wantErr: true,
			error: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
			}},
		},
		{
			name: "Get cluster metadata settings fails if cluster ID is invalid",
			args: args{params: UpdateSettingsParams{
				SettingsParams: &models.ClusterMetadataSettings{},
				ClusterParams: util.ClusterParams{
					ClusterID: util.InvalidClusterID,
				},
			}},
			want:    nil,
			wantErr: true,
			error: &multierror.Error{Errors: []error{
				util.ErrClusterLength,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpdateSettings(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateClusterMetadataSettings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.error != nil && !reflect.DeepEqual(err, tt.error) {
				t.Errorf("UpdateClusterMetadataSettings() actual error = '%v', want error '%v'", err, tt.error)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateClusterMetadataSettings() = %v, want %v", got, tt.want)
			}
		})
	}
}
