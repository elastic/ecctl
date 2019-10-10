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
	"encoding/json"
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

func TestList(t *testing.T) {
	type args struct {
		params ListParams
	}
	tests := []struct {
		name string
		args args
		want *models.ApmsInfo
		err  error
	}{
		{
			name: "list fails due to parameter validation",
			err:  &multierror.Error{Errors: []error{util.ErrAPIReq}},
		},
		{
			name: "list succeeds",
			args: args{params: ListParams{
				API: api.NewMock(mock.Response{Response: http.Response{
					Body: mock.NewStructBody(models.ApmsInfo{
						ReturnCount: ec.Int32(2),
						Apms: []*models.ApmInfo{
							{
								ID: ec.String("86d2ec6217774eedb93ba38483141997"),
								ElasticsearchCluster: &models.TargetElasticsearchCluster{
									ElasticsearchID: ec.String("d324608c97154bdba2dff97511d40368"),
								},
								Status: "started",
							},
							{
								ID: ec.String("d324608c97154bdba2dff97511d40368"),
								ElasticsearchCluster: &models.TargetElasticsearchCluster{
									ElasticsearchID: ec.String("86d2ec6217774eedb93ba38483141997"),
								},
								Status: "stopped",
							},
						},
					}),
					StatusCode: 200,
				}}),
			}},
			want: &models.ApmsInfo{
				ReturnCount: ec.Int32(2),
				Apms: []*models.ApmInfo{
					{
						ID: ec.String("86d2ec6217774eedb93ba38483141997"),
						ElasticsearchCluster: &models.TargetElasticsearchCluster{
							ElasticsearchID: ec.String("d324608c97154bdba2dff97511d40368"),
						},
						Status: "started",
					},
					{
						ID: ec.String("d324608c97154bdba2dff97511d40368"),
						ElasticsearchCluster: &models.TargetElasticsearchCluster{
							ElasticsearchID: ec.String("86d2ec6217774eedb93ba38483141997"),
						},
						Status: "stopped",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := List(tt.args.params)
			if !reflect.DeepEqual(tt.err, err) {
				t.Errorf("List() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShow(t *testing.T) {
	type args struct {
		params ShowParams
	}
	tests := []struct {
		name string
		args args
		want *models.ApmInfo
		err  error
	}{
		{
			name: "Show fails due to parameter validation",
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				errors.New(`id "" is invalid`),
			}},
		},
		{
			name: "Show succeeds",
			args: args{params: ShowParams{
				ID: "86d2ec6217774eedb93ba38483141997",
				API: api.NewMock(mock.Response{Response: http.Response{
					Body: mock.NewStructBody(models.ApmInfo{
						ID: ec.String("86d2ec6217774eedb93ba38483141997"),
						ElasticsearchCluster: &models.TargetElasticsearchCluster{
							ElasticsearchID: ec.String("d324608c97154bdba2dff97511d40368"),
						},
						Status: "started",
					}),
					StatusCode: 200,
				}}),
			}},
			want: &models.ApmInfo{
				ID: ec.String("86d2ec6217774eedb93ba38483141997"),
				ElasticsearchCluster: &models.TargetElasticsearchCluster{
					ElasticsearchID: ec.String("d324608c97154bdba2dff97511d40368"),
				},
				Status: "started",
			},
		},
		{
			name: "Show returns an error the API",
			args: args{params: ShowParams{
				ID: "86d2ec6217774eedb93ba38483141997",
				API: api.NewMock(mock.Response{Response: http.Response{
					Body: mock.NewStructBody(models.BasicFailedReply{
						Errors: []*models.BasicFailedReplyElement{{
							Code:    ec.String("some.code"),
							Message: ec.String("something went wrong"),
						}},
					}),
					StatusCode: 404,
				}}),
			}},
			err: errors.New(marshalIndent(&models.BasicFailedReply{
				Errors: []*models.BasicFailedReplyElement{
					{
						Code:    ec.String("some.code"),
						Message: ec.String("something went wrong"),
					},
				},
			})),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Show(tt.args.params)
			if !reflect.DeepEqual(tt.err, err) {
				t.Errorf("Show() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Show() = %v, want %v", got, tt.want)
			}
		})
	}
}

func marshalIndent(i interface{}) string {
	b, _ := json.MarshalIndent(i, "", "  ")
	return string(b)
}
