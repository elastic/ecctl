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

package instanceconfig

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"

	"github.com/elastic/ecctl/pkg/util"
)

func TestList(t *testing.T) {
	type args struct {
		params ListParams
	}
	tests := []struct {
		name string
		args args
		want []*models.InstanceConfiguration
		err  error
	}{
		{
			name: "List succeeds",
			args: args{
				params: ListParams{
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(listInstanceConfigsSuccess),
						StatusCode: 200,
					}}),
				},
			},
			want: []*models.InstanceConfiguration{
				{
					ID:                "data.highstorage",
					Description:       "Instance configuration to be used for a higher disk/memory ratio",
					Name:              ec.String("data.highstorage"),
					InstanceType:      ec.String("elasticsearch"),
					StorageMultiplier: float64(32),
					NodeTypes:         []string{"data", "ingest", "master"},
					DiscreteSizes: &models.DiscreteSizes{
						DefaultSize: ec.Int32(1024),
						Resource:    ec.String("memory"),
						Sizes: []int32{
							1024,
							2048,
							4096,
							8192,
							16384,
							32768,
							65536,
							131072,
							262144,
						},
					},
				},
				{
					ID:                "kibana",
					Description:       "Instance configuration to be used for Kibana",
					Name:              ec.String("kibana"),
					InstanceType:      ec.String("kibana"),
					StorageMultiplier: float64(4),
					DiscreteSizes: &models.DiscreteSizes{
						DefaultSize: ec.Int32(1024),
						Resource:    ec.String("memory"),
						Sizes: []int32{
							1024,
							2048,
							4096,
							8192,
						},
					},
				},
			},
		},
		{
			name: "List fails on API error",
			args: args{
				params: ListParams{
					API: api.NewMock(mock.New500Response(mock.NewStringBody(`{"error": "some error"}`))),
				},
			},
			err: errors.New(`{"error": "some error"}`),
		},
		{
			name: "List fails on parameter validation failure",
			args: args{},
			err:  util.ErrAPIReq,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := List(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("List() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGet(t *testing.T) {
	type args struct {
		params GetParams
	}
	tests := []struct {
		name string
		args args
		want *models.InstanceConfiguration
		err  error
	}{
		{
			name: "Get succeeds",
			args: args{
				params: GetParams{
					ID: "data.highstorage",
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(getInstanceConfigsSuccess),
						StatusCode: 200,
					}}),
				},
			},
			want: &models.InstanceConfiguration{
				ID:                "data.highstorage",
				Description:       "Instance configuration to be used for a higher disk/memory ratio",
				Name:              ec.String("data.highstorage"),
				InstanceType:      ec.String("elasticsearch"),
				StorageMultiplier: float64(32),
				NodeTypes:         []string{"data", "ingest", "master"},
				DiscreteSizes: &models.DiscreteSizes{
					DefaultSize: ec.Int32(1024),
					Resource:    ec.String("memory"),
					Sizes: []int32{
						1024,
						2048,
						4096,
						8192,
						16384,
						32768,
						65536,
						131072,
						262144,
					},
				},
			},
		},
		{
			name: "Get succeeds on kibana ID",
			args: args{
				params: GetParams{
					ID: "kibana",
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(getInstanceConfigsSuccessKibana),
						StatusCode: 200,
					}}),
				},
			},
			want: &models.InstanceConfiguration{
				ID:                "kibana",
				Description:       "Instance configuration to be used for Kibana",
				Name:              ec.String("kibana"),
				InstanceType:      ec.String("kibana"),
				StorageMultiplier: float64(4),
				DiscreteSizes: &models.DiscreteSizes{
					DefaultSize: ec.Int32(1024),
					Resource:    ec.String("memory"),
					Sizes: []int32{
						1024,
						2048,
						4096,
						8192,
					},
				},
			},
		},
		{
			name: "Get fails on API error",
			args: args{
				params: GetParams{
					ID:  "kibana",
					API: api.NewMock(mock.New500Response(mock.NewStringBody(`{"error": "some error"}`))),
				},
			},
			err: errors.New(`{"error": "some error"}`),
		},
		{
			name: "Get fails on parameter validation failure (API missing)",
			args: args{},
			err:  util.ErrAPIReq,
		},
		{
			name: "Get fails on parameter validation failure (ID missing)",
			args: args{params: GetParams{API: new(api.API)}},
			err:  errors.New("get: id must not be empty"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	type args struct {
		params CreateParams
	}
	tests := []struct {
		name string
		args args
		want *models.IDResponse
		err  error
	}{
		{
			name: "Create Succeeds",
			args: args{
				params: CreateParams{
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(`{"id": "an autogenerated id"}`),
						StatusCode: 201,
					}}),
					Config: &models.InstanceConfiguration{
						Description:       "Instance configuration to be used for Kibana",
						Name:              ec.String("kibana"),
						InstanceType:      ec.String("kibana"),
						StorageMultiplier: float64(4),
						NodeTypes:         []string{},
						DiscreteSizes: &models.DiscreteSizes{
							DefaultSize: ec.Int32(1024),
							Resource:    ec.String("memory"),
							Sizes: []int32{
								1024,
								2048,
								4096,
								8192,
							},
						},
					},
				},
			},
			want: &models.IDResponse{ID: ec.String("an autogenerated id")},
		},
		{
			name: "Create Succeeds specifying an instance configuration ID",
			args: args{
				params: CreateParams{
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(`{"id": "kibana"}`),
						StatusCode: 200,
					}}),
					Config: &models.InstanceConfiguration{
						ID:                "kibana",
						Description:       "Instance configuration to be used for Kibana",
						Name:              ec.String("kibana"),
						InstanceType:      ec.String("kibana"),
						StorageMultiplier: float64(4),
						NodeTypes:         []string{},
						DiscreteSizes: &models.DiscreteSizes{
							DefaultSize: ec.Int32(1024),
							Resource:    ec.String("memory"),
							Sizes: []int32{
								1024,
								2048,
								4096,
								8192,
							},
						},
					},
				},
			},
			want: &models.IDResponse{ID: ec.String("kibana")},
		},
		{
			name: "Create fails on API error",
			args: args{
				params: CreateParams{
					Config: &models.InstanceConfiguration{
						ID:                "kibana",
						Description:       "Instance configuration to be used for Kibana",
						Name:              ec.String("kibana"),
						InstanceType:      ec.String("kibana"),
						StorageMultiplier: float64(4),
						NodeTypes:         []string{},
						DiscreteSizes: &models.DiscreteSizes{
							DefaultSize: ec.Int32(1024),
							Resource:    ec.String("memory"),
							Sizes: []int32{
								1024,
								2048,
								4096,
								8192,
							},
						},
					},
					API: api.NewMock(mock.New500Response(mock.NewStringBody(`{"error": "some error"}`))),
				},
			},
			err: errors.New(`{"error": "some error"}`),
		},
		{
			name: "Create fails on parameter validation failure (API missing)",
			args: args{},
			err:  util.ErrAPIReq,
		},
		{
			name: "Create fails on parameter validation failure (Config missing)",
			args: args{params: CreateParams{API: new(api.API)}},
			err:  errors.New("create: request needs to have a config set"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Create(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	type args struct {
		params UpdateParams
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Update Succeeds",
			args: args{
				params: UpdateParams{
					ID: "kibana",
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(`{"id": "an autogenerated id"}`),
						StatusCode: 200,
					}}),
					Config: &models.InstanceConfiguration{
						ID:                "kibana",
						Description:       "Instance configuration to be used for Kibana",
						Name:              ec.String("kibana"),
						InstanceType:      ec.String("kibana"),
						StorageMultiplier: float64(4),
						NodeTypes:         []string{},
						DiscreteSizes: &models.DiscreteSizes{
							DefaultSize: ec.Int32(1024),
							Resource:    ec.String("memory"),
							Sizes: []int32{
								1024,
								2048,
								4096,
								8192,
							},
						},
					},
				},
			},
		},
		{
			name: "Update fails on API error",
			args: args{
				params: UpdateParams{
					ID: "kibana",
					Config: &models.InstanceConfiguration{
						ID:                "kibana",
						Description:       "Instance configuration to be used for Kibana",
						Name:              ec.String("kibana"),
						InstanceType:      ec.String("kibana"),
						StorageMultiplier: float64(4),
						NodeTypes:         []string{},
						DiscreteSizes: &models.DiscreteSizes{
							DefaultSize: ec.Int32(1024),
							Resource:    ec.String("memory"),
							Sizes: []int32{
								1024,
								2048,
								4096,
								8192,
							},
						},
					},
					API: api.NewMock(mock.New500Response(mock.NewStringBody(`{"error": "some error"}`))),
				},
			},
			err: errors.New(`{"error": "some error"}`),
		},
		{
			name: "Update fails on parameter validation failure (API missing)",
			args: args{},
			err:  util.ErrAPIReq,
		},
		{
			name: "Update fails on parameter validation failure (ID missing)",
			args: args{params: UpdateParams{API: new(api.API)}},
			err:  errors.New("update: id must not be empty"),
		},
		{
			name: "Update fails on parameter validation failure (Config missing)",
			args: args{params: UpdateParams{API: new(api.API), ID: "id"}},
			err:  errors.New("update: request needs to have a config set"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Update(tt.args.params); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type args struct {
		params DeleteParams
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Delete succeeds",
			args: args{
				params: DeleteParams{
					ID: "data.highstorage",
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(getInstanceConfigsSuccess),
						StatusCode: 200,
					}}),
				},
			},
		},
		{
			name: "Delete succeeds on kibana ID",
			args: args{
				params: DeleteParams{
					ID: "kibana",
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(getInstanceConfigsSuccessKibana),
						StatusCode: 200,
					}}),
				},
			},
		},
		{
			name: "Delete fails on API error",
			args: args{
				params: DeleteParams{
					ID:  "kibana",
					API: api.NewMock(mock.New500Response(mock.NewStringBody(`{"error": "some error"}`))),
				},
			},
			err: errors.New(`{"error": "some error"}`),
		},
		{
			name: "Delete fails on parameter validation failure (API missing)",
			args: args{},
			err:  util.ErrAPIReq,
		},
		{
			name: "Delete fails on parameter validation failure (ID missing)",
			args: args{params: DeleteParams{API: new(api.API)}},
			err:  errors.New("delete: id must not be empty"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Delete(tt.args.params); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}

func TestInstanceConfigurationOperationOutput(t *testing.T) {
	tests := []struct {
		instanceConfigurationOperation InstanceConfigurationOperation
		output                         string
	}{
		{
			instanceConfigurationOperation: InstanceConfigurationOperationNone,
			output:                         "None",
		},
		{
			instanceConfigurationOperation: InstanceConfigurationOperationUpdate,
			output:                         "Update",
		},
		{
			instanceConfigurationOperation: InstanceConfigurationOperationCreate,
			output:                         "Create",
		},
		{
			instanceConfigurationOperation: 10,
			output:                         "Invalid",
		},
	}
	for _, tt := range tests {
		if tt.instanceConfigurationOperation.String() != tt.output {
			t.Errorf("instanceConfigurationOperation got %v, expected %v", tt.instanceConfigurationOperation, tt.output)
		}
	}
}
