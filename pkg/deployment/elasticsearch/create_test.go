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
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/output"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment/elasticsearch/plan"
	"github.com/elastic/ecctl/pkg/util"
)

func TestCreateParams_Validate(t *testing.T) {
	type fields struct {
		API            *api.API
		TrackParams    util.TrackParams
		PlanDefinition *models.ElasticsearchClusterPlan
		Name           string
		LegacyParams   plan.LegacyParams
	}
	tests := []struct {
		name   string
		fields fields
		err    error
	}{
		{
			name: "Succeeds",
			fields: fields{
				API: new(api.API),
			},
		},
		{
			name:   "Fails with empty params",
			fields: fields{},
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
			}},
		},
		{
			name: "Fails with empty params",
			fields: fields{
				API:         new(api.API),
				TrackParams: util.TrackParams{Track: true},
			},
			err: &multierror.Error{Errors: []error{
				errors.New("track params: output device cannot be empty"),
			}},
		},
		{
			name: "Fails with incompatible legacy settings",
			fields: fields{
				API: new(api.API),
				LegacyParams: plan.LegacyParams{
					Capacity:  2048,
					ZoneCount: 2,
				},
				PlanDefinition: new(models.ElasticsearchClusterPlan),
			},
			err: &multierror.Error{Errors: []error{
				errors.New("cannot specify a plan definition when capacity or zonecount are set"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := &CreateParams{
				API:            tt.fields.API,
				TrackParams:    tt.fields.TrackParams,
				PlanDefinition: tt.fields.PlanDefinition,
				ClusterName:    tt.fields.Name,
				LegacyParams:   tt.fields.LegacyParams,
			}
			if err := params.Validate(); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("CreateParams.Validate() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}

func TestCreateParams_fillValues(t *testing.T) {
	lp := plan.NewLegacyPlan(plan.LegacyParams{
		Version: "5.6.0",
		Plugins: []string{"plug", "in"},
	})
	type fields struct {
		API            *api.API
		TrackParams    util.TrackParams
		PlanDefinition *models.ElasticsearchClusterPlan
		Name           string
		LegacyParams   plan.LegacyParams
	}
	tests := []struct {
		name   string
		fields fields
		want   *CreateParams
	}{
		{
			name: "fills defaults when a planDefinition is already specified",
			fields: fields{
				LegacyParams: plan.LegacyParams{
					Version: "5.6.0",
					Plugins: []string{"plug", "in"},
				},
				PlanDefinition: &models.ElasticsearchClusterPlan{
					Elasticsearch: &models.ElasticsearchConfiguration{
						DockerImage: "some",
						Version:     "5.0.0",
					},
				},
			},
			want: &CreateParams{
				LegacyParams: plan.LegacyParams{
					Version: "5.6.0",
					Plugins: []string{"plug", "in"},
				},
				PlanDefinition: &models.ElasticsearchClusterPlan{
					Elasticsearch: &models.ElasticsearchConfiguration{
						DockerImage:           "some",
						Version:               "5.6.0",
						EnabledBuiltInPlugins: []string{"plug", "in"},
					},
				},
			},
		},
		{
			name: "fills defaults with an empty PlanDefinition",
			fields: fields{
				LegacyParams: plan.LegacyParams{
					Version: "5.6.0",
					Plugins: []string{"plug", "in"},
				},
			},
			want: &CreateParams{
				LegacyParams: plan.LegacyParams{
					Version: "5.6.0",
					Plugins: []string{"plug", "in"},
				},
				PlanDefinition: lp,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := &CreateParams{
				API:            tt.fields.API,
				TrackParams:    tt.fields.TrackParams,
				PlanDefinition: tt.fields.PlanDefinition,
				ClusterName:    tt.fields.Name,
				LegacyParams:   tt.fields.LegacyParams,
			}
			params.fillValues()
			if !reflect.DeepEqual(params, tt.want) {
				t.Errorf("CreateParams.fillValues() = %v, want %v", params, tt.want)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	var success = `{
		"credentials": {
		  "password": "emulr1rBMrxCP8q5Nh5fcjpJ",
		  "username": "elastic"
		},
		"elasticsearch_cluster_id": "a933b600b2174bb79eecb977f02a9460"
	}`
	type args struct {
		params CreateParams
	}
	tests := []struct {
		name string
		args args
		want *models.ClusterCrudResponse
		err  error
	}{
		{
			name: "fails on parameter validation",
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
			}},
		},
		{
			name: "succeeds",
			args: args{params: CreateParams{
				API: api.NewMock(mock.Response{
					Response: http.Response{
						StatusCode: 200,
						Body:       mock.NewStringBody(success),
					},
				}),
			}},
			want: &models.ClusterCrudResponse{
				ElasticsearchClusterID: "a933b600b2174bb79eecb977f02a9460",
				Credentials: &models.ClusterCredentials{
					Password: "emulr1rBMrxCP8q5Nh5fcjpJ",
					Username: "elastic",
				},
			},
		},
		{
			name: "fails when the API returns an error",
			args: args{params: CreateParams{
				API: api.NewMock(mock.New500Response(mock.NewStringBody(`{}`))),
			}},
			err: errors.New("{}"),
		},
		{
			name: "succeeds with tracking",
			args: args{params: CreateParams{
				TrackParams: util.TrackParams{
					Track:         true,
					Output:        output.NewDevice(new(bytes.Buffer)),
					MaxRetries:    1,
					PollFrequency: time.Nanosecond,
				},
				API: api.NewMock(mock.Response{
					Response: http.Response{
						StatusCode: 200,
						Body:       mock.NewStringBody(success),
					},
				},
					mock.Response{Response: http.Response{
						StatusCode: 404,
						Body:       mock.NewStringBody(success),
					}},
					mock.Response{Response: http.Response{
						StatusCode: 404,
						Body:       mock.NewStringBody(success),
					}},
					mock.Response{Response: http.Response{
						StatusCode: 200,
						Body: mock.NewStructBody(&models.ElasticsearchClusterPlansInfo{
							Current: &models.ElasticsearchClusterPlanInfo{
								PlanAttemptLog: []*models.ClusterPlanStepInfo{},
							},
						}),
					}}),
			}},
			want: &models.ClusterCrudResponse{
				ElasticsearchClusterID: "a933b600b2174bb79eecb977f02a9460",
				Credentials: &models.ClusterCredentials{
					Password: "emulr1rBMrxCP8q5Nh5fcjpJ",
					Username: "elastic",
				},
			},
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
