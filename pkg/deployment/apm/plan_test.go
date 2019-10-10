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
	"bytes"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/output"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment/planutil"
	"github.com/elastic/ecctl/pkg/util"
)

var reapplyHistoryResponse = `{
"healthy": true,
"history": [
	{
	  "attempt_end_time": "2018-12-17T04:53:48.338Z",
	  "attempt_start_time": "2018-12-17T04:53:36.381Z",
	  "healthy": true,
	  "plan": {
		"apm": {
		  "system_settings": {}
		},
		"cluster_topology": [
		  {
			"size": {
			  "resource": "memory",
			  "value": 512
			},
			"zone_count": 1
		  }
		]
	  },
	  "plan_attempt_id": "a0f16c25-5dd4-46a5-87cc-55cca060fae4",
	  "plan_attempt_log": [
		{
		  "completed": "2018-12-17T04:53:48.338Z",
		  "info_log": [
			{
			  "message": "Plan successfully constructed: [PlanSuccessful()]",
			  "stage": "completed",
			  "timestamp": "2018-12-17T04:53:48.338Z"
			}
		  ],
		  "stage": "completed",
		  "started": "2018-12-17T04:53:48.338Z",
		  "status": "success",
		  "step_id": "plan-completed"
		}
	  ],
	  "plan_attempt_name": "attempt-0000000000",
	  "plan_end_time": "2018-12-17T04:53:51.035Z",
	  "source": {
		"action": "apm.create-cluster",
		"admin_id": "admin",
		"date": "2018-12-17T04:53:35.819Z",
		"facilitator": "adminconsole",
		"remote_addresses": [
		  "52.28.156.28"
		  ]
		}
	  }
	]
	}`

var apmCrudResponse = `	{
	"apm_id": "181a0cc28c9143b5a0bda51cd65676b3",
	"secret_token": "wPYRlBJauhJGEZU6Ug"	
}`

func TestGetPlan(t *testing.T) {
	type args struct {
		params PlanParams
	}
	tests := []struct {
		name    string
		args    args
		want    *models.ApmPlansInfo
		wantErr error
	}{
		{
			name: "Fails due to parameter validation",
			args: args{params: PlanParams{}},
			wantErr: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				errors.New(`id "" is invalid`),
			}},
		},
		{
			name: "Succeeds getting plan without tracking",
			args: args{params: PlanParams{
				ID: "86d2ec6217774eedb93ba38483141997",
				API: api.NewMock(mock.Response{Response: http.Response{
					Body: mock.NewStructBody(models.ApmPlansInfo{
						Healthy: ec.Bool(true),
					}),
					StatusCode: 200,
				}}),
			}},
			want: &models.ApmPlansInfo{
				Healthy: ec.Bool(true),
			},
		},
		{
			name: "Succeeds getting plan with tracking",
			args: args{params: PlanParams{
				TrackParams: util.TrackParams{
					Track:         true,
					Output:        output.NewDevice(new(bytes.Buffer)),
					PollFrequency: time.Millisecond,
					MaxRetries:    1,
				},
				ID: "86d2ec6217774eedb93ba38483141997",
				API: api.NewMock(util.AppendTrackResponses(mock.Response{Response: http.Response{
					Body: mock.NewStructBody(models.ApmPlansInfo{
						Healthy: ec.Bool(true),
					}),
					StatusCode: 200,
				}})...),
			}},
			want: &models.ApmPlansInfo{
				Healthy: ec.Bool(true),
			},
		},
		{
			name: "Fails getting plan",
			args: args{params: PlanParams{
				ID: "86d2ec6217774eedb93ba38483141997",
				API: api.NewMock(mock.Response{Response: http.Response{
					Body: mock.NewStructBody(models.BasicFailedReply{
						Errors: []*models.BasicFailedReplyElement{{
							Code:    ec.String("example.cluster_code"),
							Message: ec.String("Oh no, things went wrong"),
						}},
					}),
					StatusCode: 404,
				}}),
			}},
			wantErr: errors.New(marshalIndent(&models.BasicFailedReply{
				Errors: []*models.BasicFailedReplyElement{
					{
						Code:    ec.String("example.cluster_code"),
						Message: ec.String("Oh no, things went wrong"),
					},
				},
			})),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPlan(tt.args.params)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("GetPlan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPlan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCancelPlan(t *testing.T) {
	type args struct {
		params PlanParams
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Fails due to parameter validation",
			args: args{params: PlanParams{}},
			wantErr: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				errors.New(`id "" is invalid`),
			}},
		},
		{
			name: "Succeeds without tracking",
			args: args{params: PlanParams{
				ID: "d324608c97154bdba2dff97511d40368",
				API: api.NewMock(mock.Response{Response: http.Response{
					Body: mock.NewStructBody(models.ApmInfo{
						Status: "success",
					}),
					StatusCode: 200,
				}}),
			}},
		},
		{
			name: "Returns an error when the API returns an error",
			args: args{params: PlanParams{
				ID:  "d324608c97154bdba2dff97511d40368",
				API: api.NewMock(mock.Response{Error: errors.New("an error")}),
			}},
			wantErr: &url.Error{
				Op:  "Delete",
				URL: "https://mock-host/mock-path/clusters/apm/d324608c97154bdba2dff97511d40368/plan/pending?force_delete=false&ignore_missing=false",
				Err: errors.New("an error"),
			},
		},
		{
			name: "Succeeds with tracking",
			args: args{params: PlanParams{
				TrackParams: util.TrackParams{
					Track:         true,
					Output:        output.NewDevice(new(bytes.Buffer)),
					PollFrequency: time.Millisecond,
					MaxRetries:    1,
				},
				ID: "d324608c97154bdba2dff97511d40368",
				API: api.NewMock(util.AppendTrackResponses(mock.Response{Response: http.Response{
					Body: mock.NewStructBody(models.ApmInfo{
						Status: "success",
					}),
					StatusCode: 200,
				}})...),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CancelPlan(tt.args.params); !reflect.DeepEqual(got, tt.wantErr) {
				t.Errorf("CancelPlan() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestListPlanHistory(t *testing.T) {
	type args struct {
		params PlanParams
	}
	tests := []struct {
		name    string
		args    args
		want    []*models.ApmPlanInfo
		wantErr error
	}{
		{
			name: "Fails due to parameter validation",
			args: args{params: PlanParams{}},
			wantErr: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				errors.New(`id "" is invalid`),
			}},
		},
		{
			name: "List plan history succeeds",
			args: args{params: PlanParams{
				ID: "86d2ec6217774eedb93ba38483141997",
				API: api.NewMock(mock.Response{Response: http.Response{
					Body: mock.NewStructBody(models.ApmPlansInfo{
						History: []*models.ApmPlanInfo{
							{PlanAttemptID: "5b498bcf-f3e3-4be6-87fd-826d4c3dd845"},
							{PlanAttemptID: "d34e216b-f975-43e9-9457-3361a22a5459"},
						},
					}),
					StatusCode: 200,
				}}),
			}},
			want: []*models.ApmPlanInfo{
				{PlanAttemptID: "5b498bcf-f3e3-4be6-87fd-826d4c3dd845"},
				{PlanAttemptID: "d34e216b-f975-43e9-9457-3361a22a5459"},
			},
		},
		{
			name: "Returns an error when the API returns an error",
			args: args{params: PlanParams{
				ID:  "d324608c97154bdba2dff97511d40368",
				API: api.NewMock(mock.Response{Error: errors.New("an error")}),
			}},
			wantErr: &url.Error{
				Op:  "Get",
				URL: "https://mock-host/mock-path/clusters/apm/d324608c97154bdba2dff97511d40368/plan/activity?show_plan_defaults=false&show_plan_logs=true",
				Err: errors.New("an error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListPlanHistory(tt.args.params)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("ListPlanHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListPlanHistory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReapplyLatestPlanAttempt(t *testing.T) {
	type args struct {
		params   PlanParams
		reparams planutil.ReapplyParams
	}
	tests := []struct {
		name    string
		args    args
		want    *models.ApmCrudResponse
		wantErr error
	}{
		{
			name:    "Fails due to parameter validation",
			wantErr: errors.New("cluster id cannot be empty"),
		},
		{
			name: "Succeeds reapplying a plan with no tracking",
			args: args{params: PlanParams{
				ID: "181a0cc28c9143b5a0bda51cd65676b3",
				API: api.NewMock(
					mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(reapplyHistoryResponse),
						StatusCode: 200,
					}},
					mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(apmCrudResponse),
						StatusCode: 202,
					}},
				),
			},
				reparams: planutil.ReapplyParams{
					ID:                   "181a0cc28c9143b5a0bda51cd65676b3",
					Default:              true,
					Rolling:              false,
					GrowAndShrink:        false,
					RollingGrowAndShrink: false,
					RollingAll:           false,
					HidePlan:             true,
				},
			},
			want: &models.ApmCrudResponse{
				ApmID:       "181a0cc28c9143b5a0bda51cd65676b3",
				SecretToken: "wPYRlBJauhJGEZU6Ug",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReapplyLatestPlanAttempt(tt.args.params, tt.args.reparams)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("ReapplyLatestPlanAttempt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReapplyLatestPlanAttempt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestObtainLatestAttempt(t *testing.T) {
	type args struct {
		history []*models.ApmPlanInfo
	}
	tests := []struct {
		name string
		args args
		want *models.ApmPlan
	}{
		{
			"History is empty",
			args{
				[]*models.ApmPlanInfo{},
			},
			nil,
		},
		{
			"Retrieves the latest plan attempt",
			args{
				[]*models.ApmPlanInfo{
					{
						Plan: &models.ApmPlan{
							Transient: &models.TransientApmPlanConfiguration{
								PlanConfiguration: &models.ApmPlanControlConfiguration{
									ExtendedMaintenance: ec.Bool(false),
									ReallocateInstances: ec.Bool(false),
								},
							},
						},
					},
					{
						Plan: &models.ApmPlan{
							Transient: &models.TransientApmPlanConfiguration{
								PlanConfiguration: &models.ApmPlanControlConfiguration{
									ExtendedMaintenance: ec.Bool(true),
									ReallocateInstances: ec.Bool(false),
								},
							},
						},
					},
				},
			},
			&models.ApmPlan{
				Transient: &models.TransientApmPlanConfiguration{
					PlanConfiguration: &models.ApmPlanControlConfiguration{
						ExtendedMaintenance: ec.Bool(true),
						ReallocateInstances: ec.Bool(false),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getLatestAttempt(tt.args.history); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getLatestAttempt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComputeTransientSettings(t *testing.T) {
	type args struct {
		params computeTransientParams
	}
	tests := []struct {
		name string
		args args
		want *models.ApmPlan
	}{
		{
			name: "Transient is completely empty",
			args: args{
				params: computeTransientParams{
					plan: models.ApmPlan{},
				},
			},
			want: &models.ApmPlan{
				Transient: &models.TransientApmPlanConfiguration{
					PlanConfiguration: newDefaultTransientPlanConfiguration(),
				},
			},
		},
		{
			name: "Transient is completely empty and ExtendedMaintenance is true",
			args: args{
				params: computeTransientParams{
					plan: models.ApmPlan{},
					transient: planutil.ReapplyParams{
						ExtendedMaintenance: true,
					},
				},
			},
			want: &models.ApmPlan{
				Transient: &models.TransientApmPlanConfiguration{
					PlanConfiguration: &models.ApmPlanControlConfiguration{
						ExtendedMaintenance: ec.Bool(true),
						ReallocateInstances: ec.Bool(false),
					},
				},
			},
		},
		{
			name: "PlanConfiguration is empty",
			args: args{
				params: computeTransientParams{
					plan: models.ApmPlan{
						Transient: &models.TransientApmPlanConfiguration{},
					},
				},
			},
			want: &models.ApmPlan{
				Transient: &models.TransientApmPlanConfiguration{
					PlanConfiguration: &models.ApmPlanControlConfiguration{
						ExtendedMaintenance: ec.Bool(false),
						ReallocateInstances: ec.Bool(false),
					},
				},
			},
		},
		{
			name: "PlanConfiguration is empty and ExtendedMaintenance is true",
			args: args{
				params: computeTransientParams{
					plan: models.ApmPlan{
						Transient: &models.TransientApmPlanConfiguration{},
					},
					transient: planutil.ReapplyParams{
						ExtendedMaintenance: true,
					},
				},
			},
			want: &models.ApmPlan{
				Transient: &models.TransientApmPlanConfiguration{
					PlanConfiguration: &models.ApmPlanControlConfiguration{
						ExtendedMaintenance: ec.Bool(true),
						ReallocateInstances: ec.Bool(false),
					},
				},
			},
		},
		{
			name: "PlanConfiguration contains settings and ExtendedMaintenance is true, the existing options get overridden by defaults",
			args: args{
				params: computeTransientParams{
					plan: models.ApmPlan{
						Transient: &models.TransientApmPlanConfiguration{
							PlanConfiguration: &models.ApmPlanControlConfiguration{
								ReallocateInstances: ec.Bool(true),
							},
						},
					},
					transient: planutil.ReapplyParams{
						ExtendedMaintenance: true,
					},
				},
			},
			want: &models.ApmPlan{
				Transient: &models.TransientApmPlanConfiguration{
					PlanConfiguration: &models.ApmPlanControlConfiguration{
						ExtendedMaintenance: ec.Bool(true),
						ReallocateInstances: ec.Bool(false),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := computeTransientSettings(tt.args.params)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("computeTransientSettings() = %v, want %v", got, tt.want)
			}
		})
	}
}
