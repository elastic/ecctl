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

package stack

import (
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/multierror"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"

	"github.com/elastic/ecctl/pkg/util"
)

func TestCompareVersions(t *testing.T) {
	type fields struct {
		Version1 string
		Version2 string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
		want    bool
	}{
		{
			name: "Compare versions different major versions",
			fields: fields{
				Version1: "6.0.0",
				Version2: "5.9.9",
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "Compare versions different minor versions",
			fields: fields{
				Version1: "5.5.0",
				Version2: "5.4.9",
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "Compare versions different patch versions",
			fields: fields{
				Version1: "5.5.5",
				Version2: "5.5.4",
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "Compare versions different patch versions greater than 9",
			fields: fields{
				Version1: "5.5.10",
				Version2: "5.5.9",
			},
			wantErr: false,
			want:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := compareVersions(tt.fields.Version1, tt.fields.Version2)
			if res != tt.want {
				t.Errorf("CompareVersions() want = %v, actual %v", tt.want, res)
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
		want *models.StackVersionConfig
		err  error
	}{
		{
			name: "Get Succeeds",
			args: args{params: GetParams{
				Version: "6.0.0",
				API: api.NewMock(mock.Response{
					Response: http.Response{
						StatusCode: http.StatusOK,
						Status:     http.StatusText(http.StatusOK),
						Body: mock.NewStructBody(models.StackVersionConfig{
							Deleted: ec.Bool(false),
							Version: "6.0.0",
							Kibana: &models.StackVersionKibanaConfig{
								CapacityConstraints: &models.StackVersionInstanceCapacityConstraint{
									Max: ec.Int32(8192),
									Min: ec.Int32(1024),
								},
							},
						}),
					},
				}),
			}},
			want: &models.StackVersionConfig{
				Deleted: ec.Bool(false),
				Version: "6.0.0",
				Kibana: &models.StackVersionKibanaConfig{
					CapacityConstraints: &models.StackVersionInstanceCapacityConstraint{
						Max: ec.Int32(8192),
						Min: ec.Int32(1024),
					},
				},
			},
		},
		{
			name: "Get fails due to API error",
			args: args{params: GetParams{
				Version: "6.0.0",
				API: api.NewMock(mock.Response{
					Error: errors.New(`{"error": "some error"}`),
				}),
			}},
			err: &url.Error{
				Op:  "Get",
				URL: "https://mock-host/mock-path/stack/versions/6.0.0",
				Err: errors.New(`{"error": "some error"}`),
			},
		},
		{
			name: "Get fails due to missing api",
			args: args{params: GetParams{
				Version: "6.0.0",
			}},
			err: multierror.NewPrefixed("stack get",
				util.ErrAPIReq,
			),
		},
		{
			name: "Get fails due to missing version",
			args: args{params: GetParams{
				API: new(api.API),
			}},
			err: multierror.NewPrefixed("stack get",
				errors.New("version string empty"),
			),
		},
		{
			name: "Get fails due to empty parameters",
			args: args{params: GetParams{}},
			err: multierror.NewPrefixed("stack get",
				util.ErrAPIReq,
				errors.New("version string empty"),
			),
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

func TestList(t *testing.T) {
	type args struct {
		params ListParams
	}
	tests := []struct {
		name string
		args args
		want *models.StackVersionConfigs
		err  error
	}{
		{
			name: "List succeeds",
			args: args{params: ListParams{
				API: api.NewMock(mock.Response{
					Response: http.Response{
						StatusCode: http.StatusOK,
						Status:     http.StatusText(http.StatusOK),
						Body: mock.NewStructBody(models.StackVersionConfigs{
							Stacks: []*models.StackVersionConfig{
								{
									Deleted: ec.Bool(false),
									Version: "6.0.0",
									Kibana: &models.StackVersionKibanaConfig{
										CapacityConstraints: &models.StackVersionInstanceCapacityConstraint{
											Max: ec.Int32(8192),
											Min: ec.Int32(1024),
										},
									},
								},
								{
									Deleted: ec.Bool(false),
									Version: "6.1.0",
									Kibana: &models.StackVersionKibanaConfig{
										CapacityConstraints: &models.StackVersionInstanceCapacityConstraint{
											Max: ec.Int32(8192),
											Min: ec.Int32(1024),
										},
									},
								},
								{
									Deleted: ec.Bool(false),
									Version: "6.2.0",
									Kibana: &models.StackVersionKibanaConfig{
										CapacityConstraints: &models.StackVersionInstanceCapacityConstraint{
											Max: ec.Int32(8192),
											Min: ec.Int32(1024),
										},
									},
								},
								{
									Deleted: ec.Bool(false),
									Version: "5.6.0",
									Kibana: &models.StackVersionKibanaConfig{
										CapacityConstraints: &models.StackVersionInstanceCapacityConstraint{
											Max: ec.Int32(8192),
											Min: ec.Int32(1024),
										},
									},
								},
							},
						}),
					},
				}),
			}},
			want: &models.StackVersionConfigs{
				Stacks: []*models.StackVersionConfig{
					{
						Deleted: ec.Bool(false),
						Version: "6.2.0",
						Kibana: &models.StackVersionKibanaConfig{
							CapacityConstraints: &models.StackVersionInstanceCapacityConstraint{
								Max: ec.Int32(8192),
								Min: ec.Int32(1024),
							},
						},
					},
					{
						Deleted: ec.Bool(false),
						Version: "6.1.0",
						Kibana: &models.StackVersionKibanaConfig{
							CapacityConstraints: &models.StackVersionInstanceCapacityConstraint{
								Max: ec.Int32(8192),
								Min: ec.Int32(1024),
							},
						},
					},
					{
						Deleted: ec.Bool(false),
						Version: "6.0.0",
						Kibana: &models.StackVersionKibanaConfig{
							CapacityConstraints: &models.StackVersionInstanceCapacityConstraint{
								Max: ec.Int32(8192),
								Min: ec.Int32(1024),
							},
						},
					},
					{
						Deleted: ec.Bool(false),
						Version: "5.6.0",
						Kibana: &models.StackVersionKibanaConfig{
							CapacityConstraints: &models.StackVersionInstanceCapacityConstraint{
								Max: ec.Int32(8192),
								Min: ec.Int32(1024),
							},
						},
					},
				},
			},
		},
		{
			name: "List fails due to API error",
			args: args{params: ListParams{
				API: api.NewMock(mock.Response{
					Error: errors.New(`{"error": "some error"}`),
				}),
			}},
			err: &url.Error{
				Op:  "Get",
				URL: "https://mock-host/mock-path/stack/versions?show_deleted=false&show_unusable=false",
				Err: errors.New(`{"error": "some error"}`),
			},
		},
		{
			name: "List deleted fails due to API error",
			args: args{params: ListParams{
				Deleted: true,
				API: api.NewMock(mock.Response{
					Error: errors.New(`{"error": "some error"}`),
				}),
			}},
			err: &url.Error{
				Op:  "Get",
				URL: "https://mock-host/mock-path/stack/versions?show_deleted=true&show_unusable=false",
				Err: errors.New(`{"error": "some error"}`),
			},
		},
		{
			name: "Get fails due to missing api",
			args: args{params: ListParams{}},
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

func TestUpload(t *testing.T) {
	type args struct {
		params UploadParams
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Upload Succeeds",
			args: args{params: UploadParams{
				StackPack: strings.NewReader("aa"),
				API: api.NewMock(mock.Response{
					Response: http.Response{
						StatusCode: http.StatusOK,
						Status:     http.StatusText(http.StatusOK),
						Body:       mock.NewStringBody("{}"),
					},
				}),
			}},
		},
		{
			name: "Upload fails due to API error",
			args: args{params: UploadParams{
				StackPack: strings.NewReader("aa"),
				API: api.NewMock(mock.Response{
					Error: errors.New(`{"error": "some error"}`),
				}),
			}},
			err: &url.Error{
				Op:  "Post",
				URL: "https://mock-host/mock-path/stack/versions",
				Err: errors.New(`{"error": "some error"}`),
			},
		},
		{
			name: "Upload fails due to empty parameters",
			args: args{params: UploadParams{}},
			err: multierror.NewPrefixed("stack upload",
				util.ErrAPIReq,
				errors.New("stackpack cannot be empty"),
			),
		},
		{
			name: "Upload fails due to stackpack upload error",
			args: args{params: UploadParams{
				StackPack: strings.NewReader("aa"),
				API: api.NewMock(mock.Response{
					Response: http.Response{
						StatusCode: 200,
						Body: mock.NewStructBody(models.StackVersionArchiveProcessingResult{
							Errors: []*models.StackVersionArchiveProcessingError{
								{Errors: &models.BasicFailedReply{
									Errors: []*models.BasicFailedReplyElement{
										{
											Code:    ec.String("some.code.error"),
											Message: ec.String("some message"),
										},
									},
								}},
							},
						}),
					},
				}),
			}},
			err: multierror.NewPrefixed("stack upload",
				errors.New("some.code.error: some message"),
			),
		},
		{
			name: "Upload fails due to empty parameters",
			args: args{params: UploadParams{}},
			err: multierror.NewPrefixed("stack upload",
				util.ErrAPIReq,
				errors.New("stackpack cannot be empty"),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Upload(tt.args.params); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Upload() error = %v, wantErr %v", err, tt.err)
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
			name: "Delete Succeeds",
			args: args{params: DeleteParams{
				Version: "5.6.0",
				API: api.NewMock(mock.Response{
					Response: http.Response{
						StatusCode: http.StatusOK,
						Status:     http.StatusText(http.StatusOK),
						Body:       mock.NewStringBody("{}"),
					},
				}),
			}},
		},
		{
			name: "Delete fails due to API error",
			args: args{params: DeleteParams{
				Version: "5.6.0",
				API: api.NewMock(mock.Response{
					Error: errors.New(`{"error": "some error"}`),
				}),
			}},
			err: &url.Error{
				Op:  "Delete",
				URL: "https://mock-host/mock-path/stack/versions/5.6.0",
				Err: errors.New(`{"error": "some error"}`),
			},
		},
		{
			name: "Delete fails due to empty API",
			args: args{params: DeleteParams{
				Version: "5.0.0",
			}},
			err: multierror.NewPrefixed("stack delete",
				util.ErrAPIReq,
			),
		},
		{
			name: "Delete fails due to empty version",
			args: args{params: DeleteParams{
				API: new(api.API),
			}},
			err: multierror.NewPrefixed("stack delete",
				errors.New("version string empty"),
			),
		},
		{
			name: "Delete fails due to empty parameters",
			args: args{params: DeleteParams{}},
			err: multierror.NewPrefixed("stack delete",
				util.ErrAPIReq,
				errors.New("version string empty"),
			),
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
