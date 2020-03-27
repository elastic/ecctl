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

package plan

import (
	"bytes"
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment/planutil"
	"github.com/elastic/ecctl/pkg/util"
)

func TestObtainLatestAttempt(t *testing.T) {
	type args struct {
		history []*models.ElasticsearchClusterPlanInfo
	}
	tests := []struct {
		name string
		args args
		want *models.ElasticsearchClusterPlan
	}{
		{
			"History is empty",
			args{
				[]*models.ElasticsearchClusterPlanInfo{},
			},
			nil,
		},
		{
			"Retrieves the latest plan attempt",
			args{
				[]*models.ElasticsearchClusterPlanInfo{
					{
						Plan: &models.ElasticsearchClusterPlan{
							ZoneCount: 3,
						},
					},
					{
						Plan: &models.ElasticsearchClusterPlan{
							ZoneCount: 1,
						},
					},
				},
			},
			&models.ElasticsearchClusterPlan{
				ZoneCount: 1,
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
		name    string
		args    args
		want    *models.ElasticsearchClusterPlan
		wantErr bool
	}{
		{
			name: "Transient is completely empty",
			args: args{
				params: computeTransientParams{
					plan: models.ElasticsearchClusterPlan{},
				},
			},
			want: &models.ElasticsearchClusterPlan{
				Transient: &models.TransientElasticsearchPlanConfiguration{
					PlanConfiguration: newDefaultTransientElasticsearchPlanConfiguration(),
				},
			},
			wantErr: false,
		},
		{
			name: "Transient is completely empty and SkipSnapshot is true",
			args: args{
				params: computeTransientParams{
					plan:      models.ElasticsearchClusterPlan{},
					transient: ReapplyParams{SkipSnapshot: true},
				},
			},
			want: &models.ElasticsearchClusterPlan{
				Transient: &models.TransientElasticsearchPlanConfiguration{
					PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
						ExtendedMaintenance:  ec.Bool(false),
						OverrideFailsafe:     ec.Bool(false),
						ReallocateInstances:  ec.Bool(false),
						SkipDataMigration:    ec.Bool(false),
						SkipPostUpgradeSteps: ec.Bool(false),
						SkipSnapshot:         ec.Bool(true),
						SkipUpgradeChecker:   ec.Bool(false),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "PlanConfiguration is empty",
			args: args{
				params: computeTransientParams{
					plan: models.ElasticsearchClusterPlan{
						Transient: &models.TransientElasticsearchPlanConfiguration{},
					},
				},
			},
			want: &models.ElasticsearchClusterPlan{
				Transient: &models.TransientElasticsearchPlanConfiguration{
					PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
						ExtendedMaintenance:  ec.Bool(false),
						OverrideFailsafe:     ec.Bool(false),
						ReallocateInstances:  ec.Bool(false),
						SkipDataMigration:    ec.Bool(false),
						SkipPostUpgradeSteps: ec.Bool(false),
						SkipSnapshot:         ec.Bool(false),
						SkipUpgradeChecker:   ec.Bool(false),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "PlanConfiguration is empty and SkipSnapshot is true",
			args: args{
				params: computeTransientParams{
					plan: models.ElasticsearchClusterPlan{
						Transient: &models.TransientElasticsearchPlanConfiguration{},
					},
					transient: ReapplyParams{
						SkipSnapshot: true,
					},
				},
			},
			want: &models.ElasticsearchClusterPlan{
				Transient: &models.TransientElasticsearchPlanConfiguration{
					PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
						ExtendedMaintenance:  ec.Bool(false),
						OverrideFailsafe:     ec.Bool(false),
						ReallocateInstances:  ec.Bool(false),
						SkipDataMigration:    ec.Bool(false),
						SkipPostUpgradeSteps: ec.Bool(false),
						SkipSnapshot:         ec.Bool(true),
						SkipUpgradeChecker:   ec.Bool(false),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "PlanConfiguration contains settings and SkipSnapshot is true, the existing options get overridden by defaults",
			args: args{
				params: computeTransientParams{
					plan: models.ElasticsearchClusterPlan{
						Transient: &models.TransientElasticsearchPlanConfiguration{
							PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
								SkipDataMigration: ec.Bool(true),
							},
						},
					},
					transient: ReapplyParams{SkipSnapshot: true},
				},
			},
			want: &models.ElasticsearchClusterPlan{
				Transient: &models.TransientElasticsearchPlanConfiguration{
					PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
						ExtendedMaintenance:  ec.Bool(false),
						OverrideFailsafe:     ec.Bool(false),
						ReallocateInstances:  ec.Bool(false),
						SkipDataMigration:    ec.Bool(false),
						SkipPostUpgradeSteps: ec.Bool(false),
						SkipSnapshot:         ec.Bool(true),
						SkipUpgradeChecker:   ec.Bool(false),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "PlanConfiguration contains settings and SkipSnapshot is true, the existing options get overridden by defaults",
			args: args{
				params: computeTransientParams{
					plan: models.ElasticsearchClusterPlan{
						Transient: &models.TransientElasticsearchPlanConfiguration{
							PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
								SkipDataMigration: ec.Bool(true),
							},
						},
					},
					transient: ReapplyParams{
						SkipSnapshot: true,
						ReapplyParams: planutil.ReapplyParams{
							Reallocate: true,
						},
					},
				},
			},
			want: &models.ElasticsearchClusterPlan{
				Transient: &models.TransientElasticsearchPlanConfiguration{
					PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
						ExtendedMaintenance:  ec.Bool(false),
						OverrideFailsafe:     ec.Bool(false),
						ReallocateInstances:  ec.Bool(true),
						SkipDataMigration:    ec.Bool(false),
						SkipPostUpgradeSteps: ec.Bool(false),
						SkipSnapshot:         ec.Bool(true),
						SkipUpgradeChecker:   ec.Bool(false),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "PlanConfiguration contains settings and SkipSnapshot is true, the existing options get overridden by defaults",
			args: args{
				params: computeTransientParams{
					plan: models.ElasticsearchClusterPlan{
						Transient: &models.TransientElasticsearchPlanConfiguration{
							PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
								SkipDataMigration: ec.Bool(true),
							},
						},
					},
					transient: ReapplyParams{
						SkipSnapshot: true,
						ReapplyParams: planutil.ReapplyParams{
							Reallocate: true,
						},
					},
				},
			},
			want: &models.ElasticsearchClusterPlan{
				Transient: &models.TransientElasticsearchPlanConfiguration{
					PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
						ExtendedMaintenance:  ec.Bool(false),
						OverrideFailsafe:     ec.Bool(false),
						ReallocateInstances:  ec.Bool(true),
						SkipDataMigration:    ec.Bool(false),
						SkipPostUpgradeSteps: ec.Bool(false),
						SkipSnapshot:         ec.Bool(true),
						SkipUpgradeChecker:   ec.Bool(false),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "PlanConfiguration contains potential data loss settings, an error is returned",
			args: args{
				params: computeTransientParams{
					plan: models.ElasticsearchClusterPlan{
						Transient: &models.TransientElasticsearchPlanConfiguration{
							PlanConfiguration: new(models.ElasticsearchPlanControlConfiguration),
						},
					},
					transient: ReapplyParams{
						SkipDataMigration: true,
						ReapplyParams: planutil.ReapplyParams{
							Reallocate: true,
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := computeTransientSettings(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("computeTransientSettings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("computeTransientSettings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestComputeTransientParamsValidate(t *testing.T) {
	type fields struct {
		plan      models.ElasticsearchClusterPlan
		transient ReapplyParams
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "validate fails due to invalid settings",
			fields: fields{
				plan: models.ElasticsearchClusterPlan{
					Transient: &models.TransientElasticsearchPlanConfiguration{
						PlanConfiguration: new(models.ElasticsearchPlanControlConfiguration),
					},
				},
				transient: ReapplyParams{
					ReapplyParams: planutil.ReapplyParams{
						Reallocate: true,
					},
					SkipDataMigration: true,
				},
			},
			wantErr: true,
		},
		{
			name: "validate succeeds due to valid settings",
			fields: fields{
				plan: models.ElasticsearchClusterPlan{
					Transient: &models.TransientElasticsearchPlanConfiguration{
						PlanConfiguration: new(models.ElasticsearchPlanControlConfiguration),
					},
				},
				transient: ReapplyParams{
					SkipDataMigration: true,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := computeTransientParams{
				plan:      tt.fields.plan,
				transient: tt.fields.transient,
			}
			if err := tp.validate(); (err != nil) != tt.wantErr {
				t.Errorf("computeTransientParams.validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReapply(t *testing.T) {
	type args struct {
		params ReapplyParams
	}
	tests := []struct {
		name string
		args args
		want *models.ClusterCrudResponse
		err  error
	}{
		{
			name: "Fails due to parameter validation",
			args: args{params: ReapplyParams{}},
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				errors.New("cluster id cannot be empty"),
			}},
		},
		{
			name: "Fails to obtain latest plan attempt without Tracking",
			args: args{params: ReapplyParams{
				SkipSnapshot: true,
				API: api.NewMock(
					mock.Response{Response: http.Response{
						StatusCode: http.StatusOK,
						Status:     http.StatusText(http.StatusOK),
						Body: mock.NewStructBody(models.ElasticsearchClusterPlansInfo{
							History: []*models.ElasticsearchClusterPlanInfo{
								{PlanAttemptID: "someID"},
							},
						}),
					}},
				),
				ReapplyParams: planutil.ReapplyParams{
					ID:      "b786acd298292c2d521c0e8741761b4d",
					Default: true,
				},
			}},
			err: errors.New("unable to obtain latest plan attempt"),
		},
		{
			name: "Succeeds with Tracking",
			args: args{params: ReapplyParams{
				Output:            new(bytes.Buffer),
				Track:             true,
				TrackChangeParams: util.NewMockTrackChangeParams(""),
				SkipSnapshot:      true,
				API: api.NewMock(
					mock.Response{Response: http.Response{
						StatusCode: http.StatusOK,
						Status:     http.StatusText(http.StatusOK),
						Body: mock.NewStructBody(models.ElasticsearchClusterPlansInfo{
							History: []*models.ElasticsearchClusterPlanInfo{
								{PlanAttemptID: "someID"},
								{PlanAttemptID: "someID", Plan: &models.ElasticsearchClusterPlan{}},
							},
						}),
					}},
					mock.Response{Response: http.Response{
						StatusCode: http.StatusOK,
						Status:     http.StatusText(http.StatusOK),
						Body: mock.NewStructBody(&models.ClusterCrudResponse{
							ElasticsearchClusterID: "b786acd298292c2d521c0e8741761b4d",
						}),
					}},
				),
				ReapplyParams: planutil.ReapplyParams{
					ID:      "b786acd298292c2d521c0e8741761b4d",
					Default: true,
				},
			}},
			want: &models.ClusterCrudResponse{ElasticsearchClusterID: "b786acd298292c2d521c0e8741761b4d"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Reapply(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Reapply() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reapply() = %v, want %v", got, tt.want)
			}
		})
	}
}
