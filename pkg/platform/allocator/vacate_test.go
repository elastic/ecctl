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

package allocator

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/client/platform_infrastructure"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/output"
	"github.com/elastic/cloud-sdk-go/pkg/sync/pool"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func discardResponses(a *api.API, _ interface{}) *api.API { return a }

func TestComputeVacateRequest(t *testing.T) {
	type args struct {
		pr        *models.MoveClustersDetails
		clusters  []string
		to        []string
		overrides PlanOverrides
	}
	tests := []struct {
		name string
		args args
		want *models.MoveClustersRequest
	}{
		{
			name: "No filters",
			args: args{
				pr: &models.MoveClustersDetails{
					ElasticsearchClusters: []*models.MoveElasticsearchClusterDetails{
						{
							ClusterID: ec.String("63d765d37613423e97b1040257cf20c8"),
							CalculatedPlan: &models.TransientElasticsearchPlanConfiguration{
								PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
									Timeout:              4096,
									ReallocateInstances:  ec.Bool(false),
									ExtendedMaintenance:  ec.Bool(false),
									OverrideFailsafe:     ec.Bool(false),
									SkipDataMigration:    ec.Bool(false),
									SkipPostUpgradeSteps: ec.Bool(false),
									SkipSnapshot:         ec.Bool(false),
								},
							},
						},
					},
					KibanaClusters: []*models.MoveKibanaClusterDetails{
						{
							ClusterID: ec.String("d7ad23ad6f064709bbae7ab87a7e1bc9"),
							CalculatedPlan: &models.TransientKibanaPlanConfiguration{
								PlanConfiguration: &models.KibanaPlanControlConfiguration{
									Timeout:             4096,
									ReallocateInstances: ec.Bool(false),
									ExtendedMaintenance: ec.Bool(false),
								},
							},
						},
					},
					ApmClusters: []*models.MoveApmClusterDetails{
						{
							ClusterID: ec.String("d8ad23ad6f064709bbae7ab87a7e1bc9"),
							CalculatedPlan: &models.TransientApmPlanConfiguration{
								PlanConfiguration: &models.ApmPlanControlConfiguration{
									Timeout:             4096,
									ReallocateInstances: ec.Bool(false),
									ExtendedMaintenance: ec.Bool(false),
								},
							},
						},
					},
					AppsearchClusters: []*models.MoveAppSearchDetails{
						{
							ClusterID: ec.String("d8ad23ad6f064709bbae7ab87a7e1bc9"),
							CalculatedPlan: &models.TransientAppSearchPlanConfiguration{
								PlanConfiguration: &models.AppSearchPlanControlConfiguration{
									Timeout:             4096,
									ReallocateInstances: ec.Bool(false),
									ExtendedMaintenance: ec.Bool(false),
								},
							},
						},
					},
				},
				clusters: nil,
				to:       nil,
			},
			want: &models.MoveClustersRequest{
				ElasticsearchClusters: []*models.MoveElasticsearchClusterConfiguration{
					{
						ClusterIds: []string{
							"63d765d37613423e97b1040257cf20c8",
						},
						PlanOverride: &models.TransientElasticsearchPlanConfiguration{
							PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
								Timeout:              4096,
								ReallocateInstances:  ec.Bool(false),
								ExtendedMaintenance:  ec.Bool(false),
								OverrideFailsafe:     ec.Bool(false),
								SkipDataMigration:    ec.Bool(false),
								SkipPostUpgradeSteps: ec.Bool(false),
								SkipSnapshot:         ec.Bool(false),
							},
						},
					},
				},
				KibanaClusters: []*models.MoveKibanaClusterConfiguration{
					{
						ClusterIds: []string{
							"d7ad23ad6f064709bbae7ab87a7e1bc9",
						},
						PlanOverride: &models.TransientKibanaPlanConfiguration{
							PlanConfiguration: &models.KibanaPlanControlConfiguration{
								Timeout:             4096,
								ReallocateInstances: ec.Bool(false),
								ExtendedMaintenance: ec.Bool(false),
							},
						},
					},
				},
				ApmClusters: []*models.MoveApmClusterConfiguration{
					{
						ClusterIds: []string{
							"d8ad23ad6f064709bbae7ab87a7e1bc9",
						},
						PlanOverride: &models.TransientApmPlanConfiguration{
							PlanConfiguration: &models.ApmPlanControlConfiguration{
								Timeout:             4096,
								ReallocateInstances: ec.Bool(false),
								ExtendedMaintenance: ec.Bool(false),
							},
						},
					},
				},
				AppsearchClusters: []*models.MoveAppSearchConfiguration{
					{
						ClusterIds: []string{
							"d8ad23ad6f064709bbae7ab87a7e1bc9",
						},
						PlanOverride: &models.TransientAppSearchPlanConfiguration{
							PlanConfiguration: &models.AppSearchPlanControlConfiguration{
								Timeout:             4096,
								ReallocateInstances: ec.Bool(false),
								ExtendedMaintenance: ec.Bool(false),
							},
						},
					},
				},
			},
		},
		{
			name: "No filters with SkipSnapshot override",
			args: args{
				pr: &models.MoveClustersDetails{
					ElasticsearchClusters: []*models.MoveElasticsearchClusterDetails{
						{
							ClusterID: ec.String("63d765d37613423e97b1040257cf20c8"),
							CalculatedPlan: &models.TransientElasticsearchPlanConfiguration{
								PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
									Timeout:              4096,
									ReallocateInstances:  ec.Bool(false),
									ExtendedMaintenance:  ec.Bool(false),
									OverrideFailsafe:     ec.Bool(false),
									SkipDataMigration:    ec.Bool(false),
									SkipPostUpgradeSteps: ec.Bool(false),
									SkipSnapshot:         ec.Bool(false),
								},
							},
						},
					},
				},
				clusters:  nil,
				to:        nil,
				overrides: PlanOverrides{SkipSnapshot: ec.Bool(true)},
			},
			want: &models.MoveClustersRequest{
				ElasticsearchClusters: []*models.MoveElasticsearchClusterConfiguration{
					{
						ClusterIds: []string{
							"63d765d37613423e97b1040257cf20c8",
						},
						PlanOverride: &models.TransientElasticsearchPlanConfiguration{
							PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
								Timeout:              4096,
								ReallocateInstances:  ec.Bool(false),
								ExtendedMaintenance:  ec.Bool(false),
								OverrideFailsafe:     ec.Bool(false),
								SkipDataMigration:    ec.Bool(false),
								SkipPostUpgradeSteps: ec.Bool(false),
								SkipSnapshot:         ec.Bool(true),
							},
						},
					},
				},
			},
		},
		{
			name: "No filters with SkipDataMigration override",
			args: args{
				pr: &models.MoveClustersDetails{
					ElasticsearchClusters: []*models.MoveElasticsearchClusterDetails{
						{
							ClusterID: ec.String("63d765d37613423e97b1040257cf20c8"),
							CalculatedPlan: &models.TransientElasticsearchPlanConfiguration{
								PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
									Timeout:              4096,
									ReallocateInstances:  ec.Bool(false),
									ExtendedMaintenance:  ec.Bool(false),
									OverrideFailsafe:     ec.Bool(false),
									SkipDataMigration:    ec.Bool(false),
									SkipPostUpgradeSteps: ec.Bool(false),
									SkipSnapshot:         ec.Bool(false),
								},
							},
						},
					},
				},
				clusters:  nil,
				to:        nil,
				overrides: PlanOverrides{SkipDataMigration: ec.Bool(true)},
			},
			want: &models.MoveClustersRequest{
				ElasticsearchClusters: []*models.MoveElasticsearchClusterConfiguration{
					{
						ClusterIds: []string{
							"63d765d37613423e97b1040257cf20c8",
						},
						PlanOverride: &models.TransientElasticsearchPlanConfiguration{
							PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
								Timeout:              4096,
								ReallocateInstances:  ec.Bool(false),
								ExtendedMaintenance:  ec.Bool(false),
								OverrideFailsafe:     ec.Bool(false),
								SkipDataMigration:    ec.Bool(true),
								SkipPostUpgradeSteps: ec.Bool(false),
								SkipSnapshot:         ec.Bool(false),
							},
						},
					},
				},
			},
		},
		{
			name: "No filters with OverrideFailsafe override",
			args: args{
				pr: &models.MoveClustersDetails{
					ElasticsearchClusters: []*models.MoveElasticsearchClusterDetails{
						{
							ClusterID: ec.String("63d765d37613423e97b1040257cf20c8"),
							CalculatedPlan: &models.TransientElasticsearchPlanConfiguration{
								PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
									Timeout:              4096,
									ReallocateInstances:  ec.Bool(false),
									ExtendedMaintenance:  ec.Bool(false),
									OverrideFailsafe:     ec.Bool(false),
									SkipDataMigration:    ec.Bool(false),
									SkipPostUpgradeSteps: ec.Bool(false),
									SkipSnapshot:         ec.Bool(false),
								},
							},
						},
					},
				},
				clusters:  nil,
				to:        nil,
				overrides: PlanOverrides{OverrideFailsafe: ec.Bool(true)},
			},
			want: &models.MoveClustersRequest{
				ElasticsearchClusters: []*models.MoveElasticsearchClusterConfiguration{
					{
						ClusterIds: []string{
							"63d765d37613423e97b1040257cf20c8",
						},
						PlanOverride: &models.TransientElasticsearchPlanConfiguration{
							PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
								Timeout:              4096,
								ReallocateInstances:  ec.Bool(false),
								ExtendedMaintenance:  ec.Bool(false),
								OverrideFailsafe:     ec.Bool(true),
								SkipDataMigration:    ec.Bool(false),
								SkipPostUpgradeSteps: ec.Bool(false),
								SkipSnapshot:         ec.Bool(false),
							},
						},
					},
				},
			},
		},
		{
			name: "Set target allocator",
			args: args{
				pr: &models.MoveClustersDetails{
					ElasticsearchClusters: []*models.MoveElasticsearchClusterDetails{
						{
							ClusterID: ec.String("63d765d37613423e97b1040257cf20c8"),
							CalculatedPlan: &models.TransientElasticsearchPlanConfiguration{
								PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
									Timeout:              4096,
									ReallocateInstances:  ec.Bool(false),
									ExtendedMaintenance:  ec.Bool(false),
									OverrideFailsafe:     ec.Bool(false),
									SkipDataMigration:    ec.Bool(false),
									SkipPostUpgradeSteps: ec.Bool(false),
									SkipSnapshot:         ec.Bool(false),
								},
							},
						},
					},
					KibanaClusters: []*models.MoveKibanaClusterDetails{
						{
							ClusterID: ec.String("d7ad23ad6f064709bbae7ab87a7e1bc9"),
							CalculatedPlan: &models.TransientKibanaPlanConfiguration{
								PlanConfiguration: &models.KibanaPlanControlConfiguration{
									Timeout:             4096,
									ReallocateInstances: ec.Bool(false),
									ExtendedMaintenance: ec.Bool(false),
								},
							},
						},
					},
				},
				clusters: nil,
				to:       []string{"192.168.44.11"},
			},
			want: &models.MoveClustersRequest{
				ElasticsearchClusters: []*models.MoveElasticsearchClusterConfiguration{
					{
						ClusterIds: []string{
							"63d765d37613423e97b1040257cf20c8",
						},
						PlanOverride: &models.TransientElasticsearchPlanConfiguration{
							PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
								Timeout:              4096,
								ReallocateInstances:  ec.Bool(false),
								ExtendedMaintenance:  ec.Bool(false),
								OverrideFailsafe:     ec.Bool(false),
								SkipDataMigration:    ec.Bool(false),
								SkipPostUpgradeSteps: ec.Bool(false),
								SkipSnapshot:         ec.Bool(false),
								PreferredAllocators:  []string{"192.168.44.11"},
							},
						},
					},
				},
				KibanaClusters: []*models.MoveKibanaClusterConfiguration{
					{
						ClusterIds: []string{
							"d7ad23ad6f064709bbae7ab87a7e1bc9",
						},
						PlanOverride: &models.TransientKibanaPlanConfiguration{
							PlanConfiguration: &models.KibanaPlanControlConfiguration{
								Timeout:             4096,
								ReallocateInstances: ec.Bool(false),
								ExtendedMaintenance: ec.Bool(false),
								PreferredAllocators: []string{"192.168.44.11"},
							},
						},
					},
				},
			},
		},
		{
			name: "Set specific cluster filter",
			args: args{
				pr: &models.MoveClustersDetails{
					ElasticsearchClusters: []*models.MoveElasticsearchClusterDetails{
						{
							ClusterID: ec.String("63d765d37613423e97b1040257cf20c8"),
							CalculatedPlan: &models.TransientElasticsearchPlanConfiguration{
								PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
									Timeout:              4096,
									ReallocateInstances:  ec.Bool(false),
									ExtendedMaintenance:  ec.Bool(false),
									OverrideFailsafe:     ec.Bool(false),
									SkipDataMigration:    ec.Bool(false),
									SkipPostUpgradeSteps: ec.Bool(false),
									SkipSnapshot:         ec.Bool(false),
								},
							},
						},
						{
							ClusterID: ec.String("63d765d37613423e97b1040257cf20cf"),
							CalculatedPlan: &models.TransientElasticsearchPlanConfiguration{
								PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
									Timeout:              4096,
									ReallocateInstances:  ec.Bool(false),
									ExtendedMaintenance:  ec.Bool(false),
									OverrideFailsafe:     ec.Bool(false),
									SkipDataMigration:    ec.Bool(false),
									SkipPostUpgradeSteps: ec.Bool(false),
									SkipSnapshot:         ec.Bool(false),
								},
							},
						},
					},
					KibanaClusters: []*models.MoveKibanaClusterDetails{
						{
							ClusterID: ec.String("d7ad23ad6f064709bbae7ab87a7e1bc9"),
							CalculatedPlan: &models.TransientKibanaPlanConfiguration{
								PlanConfiguration: &models.KibanaPlanControlConfiguration{
									Timeout:             4096,
									ReallocateInstances: ec.Bool(false),
									ExtendedMaintenance: ec.Bool(false),
								},
							},
						},
					},
				},
				clusters: []string{"63d765d37613423e97b1040257cf20c8"},
				to:       nil,
			},
			want: &models.MoveClustersRequest{
				ElasticsearchClusters: []*models.MoveElasticsearchClusterConfiguration{
					{
						ClusterIds: []string{
							"63d765d37613423e97b1040257cf20c8",
						},
						PlanOverride: &models.TransientElasticsearchPlanConfiguration{
							PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
								Timeout:              4096,
								ReallocateInstances:  ec.Bool(false),
								ExtendedMaintenance:  ec.Bool(false),
								OverrideFailsafe:     ec.Bool(false),
								SkipDataMigration:    ec.Bool(false),
								SkipPostUpgradeSteps: ec.Bool(false),
								SkipSnapshot:         ec.Bool(false),
							},
						},
					},
				},
			},
		},
		{
			name: "Set specific cluster filter and allocator target",
			args: args{
				pr: &models.MoveClustersDetails{
					ElasticsearchClusters: []*models.MoveElasticsearchClusterDetails{
						{
							ClusterID: ec.String("63d765d37613423e97b1040257cf20c8"),
							CalculatedPlan: &models.TransientElasticsearchPlanConfiguration{
								PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
									Timeout:              4096,
									ReallocateInstances:  ec.Bool(false),
									ExtendedMaintenance:  ec.Bool(false),
									OverrideFailsafe:     ec.Bool(false),
									SkipDataMigration:    ec.Bool(false),
									SkipPostUpgradeSteps: ec.Bool(false),
									SkipSnapshot:         ec.Bool(false),
								},
							},
						},
						{
							ClusterID: ec.String("63d765d37613423e97b1040257cf20cf"),
							CalculatedPlan: &models.TransientElasticsearchPlanConfiguration{
								PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
									Timeout:              4096,
									ReallocateInstances:  ec.Bool(false),
									ExtendedMaintenance:  ec.Bool(false),
									OverrideFailsafe:     ec.Bool(false),
									SkipDataMigration:    ec.Bool(false),
									SkipPostUpgradeSteps: ec.Bool(false),
									SkipSnapshot:         ec.Bool(false),
								},
							},
						},
					},
					KibanaClusters: []*models.MoveKibanaClusterDetails{
						{
							ClusterID: ec.String("d7ad23ad6f064709bbae7ab87a7e1bc9"),
							CalculatedPlan: &models.TransientKibanaPlanConfiguration{
								PlanConfiguration: &models.KibanaPlanControlConfiguration{
									Timeout:             4096,
									ReallocateInstances: ec.Bool(false),
									ExtendedMaintenance: ec.Bool(false),
								},
							},
						},
					},
				},
				clusters: []string{"63d765d37613423e97b1040257cf20c8"},
				to:       []string{"192.168.44.11"},
			},
			want: &models.MoveClustersRequest{
				ElasticsearchClusters: []*models.MoveElasticsearchClusterConfiguration{
					{
						ClusterIds: []string{
							"63d765d37613423e97b1040257cf20c8",
						},
						PlanOverride: &models.TransientElasticsearchPlanConfiguration{
							PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
								Timeout:              4096,
								ReallocateInstances:  ec.Bool(false),
								ExtendedMaintenance:  ec.Bool(false),
								OverrideFailsafe:     ec.Bool(false),
								SkipDataMigration:    ec.Bool(false),
								SkipPostUpgradeSteps: ec.Bool(false),
								SkipSnapshot:         ec.Bool(false),
								PreferredAllocators:  []string{"192.168.44.11"},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ComputeVacateRequest(
				tt.args.pr, tt.args.clusters, tt.args.to, tt.args.overrides,
			); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ComputeVacateRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckVacateFailures(t *testing.T) {
	type args struct {
		failures      *models.MoveClustersDetails
		clusterFilter []string
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Returns nil on nil failures",
			args: args{failures: nil},
			err:  nil,
		},
		{
			name: "Returns nil on no failures",
			args: args{failures: new(models.MoveClustersDetails)},
			err:  nil,
		},
		{
			name: "Returns an elasticsearch error on ES vacate failure",
			args: args{
				failures: &models.MoveClustersDetails{
					ElasticsearchClusters: []*models.MoveElasticsearchClusterDetails{
						{
							ClusterID: ec.String("123456789"),
							Errors: []*models.BasicFailedReplyElement{
								{
									Code:    ec.String("unknown"),
									Message: ec.String("a message"),
								},
							},
						},
					},
				},
			},
			err: &multierror.Error{Errors: []error{
				errors.New("resource id [123456789][elasticsearch] failed vacating, reason: code: unknown, message: a message"),
			}},
		},
		{
			name: "Returns a kibana error on Kibana vacate failure",
			args: args{
				failures: &models.MoveClustersDetails{
					KibanaClusters: []*models.MoveKibanaClusterDetails{
						{
							ClusterID: ec.String("123456789"),
							Errors: []*models.BasicFailedReplyElement{
								{
									Code:    ec.String("unknown"),
									Message: ec.String("a kibana error message"),
								},
							},
						},
					},
				},
			},
			err: &multierror.Error{Errors: []error{
				errors.New("resource id [123456789][kibana] failed vacating, reason: code: unknown, message: a kibana error message"),
			}},
		},
		{
			name: "Returns an error on APM vacate failure",
			args: args{
				failures: &models.MoveClustersDetails{
					ApmClusters: []*models.MoveApmClusterDetails{
						{
							ClusterID: ec.String("123456789"),
							Errors: []*models.BasicFailedReplyElement{
								{
									Code:    ec.String("unknown"),
									Message: ec.String("an apm error message"),
								},
							},
						},
					},
				},
			},
			err: &multierror.Error{Errors: []error{
				errors.New("resource id [123456789][apm] failed vacating, reason: code: unknown, message: an apm error message"),
			}},
		},
		{
			name: "Returns an error on AppSearch vacate failure",
			args: args{
				failures: &models.MoveClustersDetails{
					AppsearchClusters: []*models.MoveAppSearchDetails{
						{
							ClusterID: ec.String("123456789"),
							Errors: []*models.BasicFailedReplyElement{
								{
									Code:    ec.String("unknown"),
									Message: ec.String("an apm error message"),
								},
							},
						},
					},
				},
			},
			err: &multierror.Error{Errors: []error{
				errors.New("resource id [123456789][appsearch] failed vacating, reason: code: unknown, message: an apm error message"),
			}},
		},
		{
			name: "Returns an elasticsearch & kibana error on multiple ES & Kibana vacate failures",
			args: args{
				failures: &models.MoveClustersDetails{
					ElasticsearchClusters: []*models.MoveElasticsearchClusterDetails{
						{
							ClusterID: ec.String("123456789"),
							Errors: []*models.BasicFailedReplyElement{
								{
									Code:    ec.String("unknown"),
									Message: ec.String("a message"),
								},
							},
						},
					},
					KibanaClusters: []*models.MoveKibanaClusterDetails{
						{
							ClusterID: ec.String("123456789"),
							Errors: []*models.BasicFailedReplyElement{
								{
									Code:    ec.String("unknown"),
									Message: ec.String("a kibana error message"),
								},
							},
						},
					},
				},
			},
			err: &multierror.Error{Errors: []error{
				errors.New("resource id [123456789][elasticsearch] failed vacating, reason: code: unknown, message: a message"),
				errors.New("resource id [123456789][kibana] failed vacating, reason: code: unknown, message: a kibana error message"),
			}},
		},
		{
			name: "Returns only the clusters specified in the ClusterFilter",
			args: args{
				clusterFilter: []string{"1234567890"},
				failures: &models.MoveClustersDetails{
					ElasticsearchClusters: []*models.MoveElasticsearchClusterDetails{
						{
							ClusterID: ec.String("1234567890"),
							Errors: []*models.BasicFailedReplyElement{
								{
									Code:    ec.String("unknown"),
									Message: ec.String("a message"),
								},
							},
						},
					},
					KibanaClusters: []*models.MoveKibanaClusterDetails{
						{
							ClusterID: ec.String("123456789"),
							Errors:    []*models.BasicFailedReplyElement{},
						},
					},
				},
			},
			err: &multierror.Error{Errors: []error{
				errors.New("resource id [1234567890][elasticsearch] failed vacating, reason: code: unknown, message: a message"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckVacateFailures(tt.args.failures, tt.args.clusterFilter); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("CheckVacateFailures() error = %v, wantErr = %v", err, tt.err)
			}
		})
	}
}

func TestVacateCluster(t *testing.T) {
	var errEmptyParams = `allocator : resource id [][]: parameter validation: 5 errors occurred:
	* api reference is required for command
	* vacate cluster: invalid allocator ID 
	* vacate cluster: invalid cluster ID 
	* vacate cluster: invalid kind 
	* output device cannot be nil

`
	var errInvalidParams = `allocator someID: resource id [3ee11eb40eda22cac0cce259625c6734][somethingweird]: parameter validation: 3 errors occurred:
	* api reference is required for command
	* vacate cluster: invalid kind somethingweird
	* output device cannot be nil

`
	type args struct {
		params *VacateClusterParams
		buf    *bytes.Buffer
	}
	tests := []struct {
		name string
		args args
		err  string
		want string
	}{
		{
			name: "Fails due to params being nil",
			args: args{},
			err:  "allocator vacate: params cannot be nil",
		},
		{
			name: "Fails due to param validation",
			args: args{
				params: new(VacateClusterParams),
			},
			err: errEmptyParams,
		},
		{
			name: "Fails due to param validation (invalid kind)",
			args: args{
				params: &VacateClusterParams{
					ID:        "someID",
					ClusterID: "3ee11eb40eda22cac0cce259625c6734",
					Kind:      "somethingweird",
				},
			},
			err: errInvalidParams,
		},
		{
			name: "Succeeds with an elasticsearch cluster",
			args: args{
				buf: new(bytes.Buffer),
				params: &VacateClusterParams{
					ID:             "someID",
					ClusterID:      "3ee11eb40eda22cac0cce259625c6734",
					Kind:           "elasticsearch",
					Output:         new(output.Device),
					TrackFrequency: time.Nanosecond,
					OutputFormat:   "text",
					MaxPollRetries: 1,
					API: discardResponses(newElasticsearchVacateMove(t, "someID", vacateCaseClusterConfig{
						ID: "3ee11eb40eda22cac0cce259625c6734",
						steps: [][]*models.ClusterPlanStepInfo{
							{
								newPlanStep("step1", "success"),
								newPlanStep("step2", "pending"),
							},
						},
						plan: []*models.ClusterPlanStepInfo{
							newPlanStep("step1", "success"),
							newPlanStep("step2", "success"),
							newPlanStep("plan-completed", "success"),
						},
					})),
				},
			},
			want: newOutputResponses(
				`Deployment [DISCOVERED_DEPLOYMENT_ID] - [Elasticsearch][3ee11eb40eda22cac0cce259625c6734]: running step "step2" (Plan duration )...`,
				"\x1b[92;mDeployment [DISCOVERED_DEPLOYMENT_ID] - [Elasticsearch][3ee11eb40eda22cac0cce259625c6734]: finished running all the plan steps\x1b[0m (Total plan duration )",
			),
		},
		{
			name: "Succeeds with an elasticsearch cluster with no tracking",
			args: args{
				buf: new(bytes.Buffer),
				params: &VacateClusterParams{
					ID:             "someID",
					ClusterID:      "3ee11eb40eda22cac0cce259625c6734",
					Kind:           "elasticsearch",
					Output:         new(output.Device),
					TrackFrequency: time.Nanosecond,
					SkipTracking:   true,
					MaxPollRetries: 1,
					API: discardResponses(
						newElasticsearchVacateMove(t, "someID", vacateCaseClusterConfig{}),
					),
				},
			},
		},
		{
			name: "Succeeds with a kibana instance",
			args: args{
				buf: new(bytes.Buffer),
				params: &VacateClusterParams{
					ID:             "someID",
					ClusterID:      "2ee11eb40eda22cac0cce259625c6734",
					Kind:           "kibana",
					Output:         new(output.Device),
					TrackFrequency: time.Nanosecond,
					OutputFormat:   "text",
					MaxPollRetries: 1,
					API: discardResponses(newKibanaVacateMove(t, "someID", vacateCaseClusterConfig{
						ID: "2ee11eb40eda22cac0cce259625c6734",
						steps: [][]*models.ClusterPlanStepInfo{
							{
								newPlanStep("step1", "success"),
								newPlanStep("step2", "pending"),
							},
							{
								newPlanStep("step1", "success"),
								newPlanStep("step2", "success"),
								newPlanStep("step3", "pending"),
							},
						},
						plan: []*models.ClusterPlanStepInfo{
							newPlanStep("step1", "success"),
							newPlanStep("step2", "success"),
							newPlanStep("step3", "success"),
							newPlanStep("plan-completed", "success"),
						},
					})),
				},
			},
			want: newOutputResponses(
				"Deployment [DISCOVERED_DEPLOYMENT_ID] - [Kibana][2ee11eb40eda22cac0cce259625c6734]: running step \"step2\" (Plan duration )...",
				"Deployment [DISCOVERED_DEPLOYMENT_ID] - [Kibana][2ee11eb40eda22cac0cce259625c6734]: running step \"step3\" (Plan duration )...",
				"\x1b[92;mDeployment [DISCOVERED_DEPLOYMENT_ID] - [Kibana][2ee11eb40eda22cac0cce259625c6734]: finished running all the plan steps\x1b[0m (Total plan duration )",
			),
		},
		{
			name: "Succeeds with a kibana instance with no tracking",
			args: args{
				buf: new(bytes.Buffer),
				params: &VacateClusterParams{
					ID:             "someID",
					ClusterID:      "2ee11eb40eda22cac0cce259625c6734",
					Kind:           "kibana",
					Output:         new(output.Device),
					TrackFrequency: time.Nanosecond,
					SkipTracking:   true,
					MaxPollRetries: 1,
					API: discardResponses(
						newKibanaVacateMove(t, "someID", vacateCaseClusterConfig{}),
					),
				},
			},
		},
		{
			name: "Succeeds with an appsearch instance with no tracking",
			args: args{
				buf: new(bytes.Buffer),
				params: &VacateClusterParams{
					ID:             "someID",
					ClusterID:      "2ee11eb40eda22cac0cce259625c6734",
					Kind:           "appsearch",
					Output:         new(output.Device),
					TrackFrequency: time.Nanosecond,
					SkipTracking:   true,
					MaxPollRetries: 1,
					API: discardResponses(
						newKibanaVacateMove(t, "someID", vacateCaseClusterConfig{}),
					),
				},
			},
		},
		{
			name: "Moving kibana instance fails",
			args: args{
				buf: new(bytes.Buffer),
				params: &VacateClusterParams{
					ID:             "someID",
					ClusterID:      "2ee11eb40eda22cac0cce259625c6734",
					Kind:           "kibana",
					Output:         new(output.Device),
					TrackFrequency: time.Nanosecond,
					MaxPollRetries: 1,
					API: discardResponses(newKibanaVacateMove(t, "someID", vacateCaseClusterConfig{
						ID:   "2ee11eb40eda22cac0cce259625c6734",
						fail: true,
					})),
				},
			},
			err: newMultierror(
				errors.New("resource id [2ee11eb40eda22cac0cce259625c6734][kibana] failed vacating, reason: code: a code, message: a message"),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.buf != nil {
				tt.args.params.Output = output.NewDevice(tt.args.buf)
			}
			if err := VacateCluster(tt.args.params); err != nil && err.Error() != tt.err {
				t.Errorf("VacateCluster() error = %v, wantErr %v", err, tt.err)
			}

			var got string
			if tt.args.buf != nil {
				got = regexp.MustCompile(`duration.*\)`).
					ReplaceAllString(tt.args.buf.String(), "duration )")
			}

			if tt.args.buf != nil && tt.want != got {
				t.Errorf("VacateCluster() output = \n%v, want \n%v", got, tt.want)
			}
		})
	}
}

func Test_fillVacateClusterParams(t *testing.T) {
	type args struct {
		params *VacateClusterParams
	}
	tests := []struct {
		name string
		args args
		want *VacateClusterParams
		err  error
	}{
		{
			name: "returns an error when the allocator discovery health can't be obtained",
			args: args{
				params: &VacateClusterParams{
					API:       api.NewMock(mock.Response{Error: errors.New("unauthorized")}),
					ID:        "allocator-1",
					ClusterID: "3ee11eb40eda22cac0cce259625c6734",
					Kind:      "elasticsearch",
					Output:    output.NewDevice(new(bytes.Buffer)),
				},
			},
			err: errors.New(`allocator allocator-1: resource id [3ee11eb40eda22cac0cce259625c6734][elasticsearch]: allocator health autodiscovery: Get https://mock-host/mock-path/platform/infrastructure/allocators/allocator-1: unauthorized`),
		},
		{
			name: "sets defaults on parameters that aren't specified",
			args: args{
				params: &VacateClusterParams{
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       newAllocator(t, "allocator-1", "3ee11eb40eda22cac0cce259625c6734", "elasticsearch"),
						StatusCode: 200,
					}}),
					ID:        "allocator-1",
					ClusterID: "3ee11eb40eda22cac0cce259625c6734",
					Kind:      "elasticsearch",
					Output:    output.NewDevice(new(bytes.Buffer)),
				},
			},
			want: &VacateClusterParams{
				ID:             "allocator-1",
				ClusterID:      "3ee11eb40eda22cac0cce259625c6734",
				Kind:           "elasticsearch",
				Output:         output.NewDevice(new(bytes.Buffer)),
				MaxPollRetries: util.DefaultRetries,
				TrackFrequency: util.DefaultPollFrequency,
				AllocatorDown:  ec.Bool(false),
			},
		},
		{
			name: "respects explicit parameters",
			args: args{
				params: &VacateClusterParams{
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       newAllocator(t, "allocator-1", "3ee11eb40eda22cac0cce259625c6734", "elasticsearch"),
						StatusCode: 200,
					}}),
					ID:             "allocator-1",
					ClusterID:      "3ee11eb40eda22cac0cce259625c6734",
					Kind:           "elasticsearch",
					Output:         output.NewDevice(new(bytes.Buffer)),
					MaxPollRetries: 4,
					AllocatorDown:  ec.Bool(true),
					TrackFrequency: time.Millisecond,
				},
			},
			want: &VacateClusterParams{
				ID:             "allocator-1",
				ClusterID:      "3ee11eb40eda22cac0cce259625c6734",
				Kind:           "elasticsearch",
				Output:         output.NewDevice(new(bytes.Buffer)),
				MaxPollRetries: 4,
				TrackFrequency: time.Millisecond,
				AllocatorDown:  ec.Bool(true),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fillVacateClusterParams(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				var errMesg string
				if err != nil {
					errMesg = err.Error()
				}
				var wantErr string
				if tt.err != nil {
					wantErr = tt.err.Error()
				}
				if errMesg != wantErr {
					t.Errorf("fillVacateClusterParams() error = %v, wantErr %v", errMesg, wantErr)
					return
				}
			}
			if got != nil {
				got.API = nil
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fillVacateClusterParams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newMoveClusterParams(t *testing.T) {
	type args struct {
		params *VacateClusterParams
	}
	tests := []struct {
		name string
		args args
		want *platform_infrastructure.MoveClustersByTypeParams
		err  error
	}{
		{
			name: "when an API error is returned, the error is properly wrapped",
			args: args{params: &VacateClusterParams{
				API:       api.NewMock(mock.Response{Error: errors.New("unauthorized")}),
				ID:        "allocator-1",
				ClusterID: "3ee11eb40eda22cac0cce259625c6734",
				Kind:      "elasticsearch",
				Output:    output.NewDevice(new(bytes.Buffer)),
			}},
			err: errors.New(`allocator allocator-1: resource id [3ee11eb40eda22cac0cce259625c6734][elasticsearch]: validate_only: Post https://mock-host/mock-path/platform/infrastructure/allocators/allocator-1/clusters/_move?force_update=false&validate_only=true: unauthorized`),
		},
		{
			name: "elasticsearch move succeeds to get parameters",
			args: args{params: &VacateClusterParams{
				API: api.NewMock(mock.Response{Response: http.Response{
					Body:       newElasticsearchMove(t, "3ee11eb40eda22cac0cce259625c6734", "allocator-1"),
					StatusCode: 202,
				}}),
				ID:            "allocator-1",
				ClusterID:     "3ee11eb40eda22cac0cce259625c6734",
				Kind:          "elasticsearch",
				Output:        output.NewDevice(new(bytes.Buffer)),
				AllocatorDown: ec.Bool(false),
			}},
			want: platform_infrastructure.NewMoveClustersByTypeParams().
				WithAllocatorID("allocator-1").
				WithClusterType(ec.String("elasticsearch")).
				WithAllocatorDown(ec.Bool(false)).
				WithBody(&models.MoveClustersRequest{
					ElasticsearchClusters: []*models.MoveElasticsearchClusterConfiguration{
						{
							ClusterIds: []string{"3ee11eb40eda22cac0cce259625c6734"},
							PlanOverride: &models.TransientElasticsearchPlanConfiguration{
								PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
									MoveAllocators: []*models.AllocatorMoveRequest{
										{
											From: ec.String("allocator-1"),
										},
									},
								},
							},
						},
					},
				}),
		},
		{
			name: "elasticsearch move succeeds to get parameters",
			args: args{params: &VacateClusterParams{
				API: api.NewMock(mock.Response{Response: http.Response{
					Body:       newApmMove(t, "3ee11eb40eda22cac0cce259625c6734", "allocator-1"),
					StatusCode: 202,
				}}),
				ID:            "allocator-1",
				ClusterID:     "3ee11eb40eda22cac0cce259625c6734",
				Kind:          util.Apm,
				Output:        output.NewDevice(new(bytes.Buffer)),
				AllocatorDown: ec.Bool(false),
			}},
			want: platform_infrastructure.NewMoveClustersByTypeParams().
				WithAllocatorID("allocator-1").
				WithClusterType(ec.String(util.Apm)).
				WithAllocatorDown(ec.Bool(false)).
				WithBody(&models.MoveClustersRequest{
					ApmClusters: []*models.MoveApmClusterConfiguration{
						{
							ClusterIds: []string{"3ee11eb40eda22cac0cce259625c6734"},
							PlanOverride: &models.TransientApmPlanConfiguration{
								PlanConfiguration: &models.ApmPlanControlConfiguration{
									MoveAllocators: []*models.AllocatorMoveRequest{
										{
											From: ec.String("allocator-1"),
										},
									},
								},
							},
						},
					},
				}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newMoveClusterParams(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				var errMesg string
				if err != nil {
					errMesg = err.Error()
				}
				var wantErr string
				if tt.err != nil {
					wantErr = tt.err.Error()
				}
				if errMesg != wantErr {
					t.Errorf("newMoveClusterParams() error = %v, wantErr %v", errMesg, wantErr)
					return
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				_, diff, _ := util.CompareStructs(got, tt.want)
				t.Errorf("newMoveClusterParams() = \n%+v\n, want \n%+v", got, tt.want)
				fmt.Println(diff)
			}
		})
	}
}

func Test_addAllocatorMovesToPool(t *testing.T) {
	type args struct {
		params addAllocatorMovesToPoolParams
	}
	tests := []struct {
		name          string
		args          args
		wantLeftovers []pool.Validator
		wantMoved     bool
	}{
		{
			name: "Move clusters when no filter is specified",
			args: args{params: addAllocatorMovesToPoolParams{
				ID: "allocator-1",
				Pool: func() *pool.Pool {
					p, _ := pool.NewPool(pool.Params{
						Size:    1,
						Run:     VacateClusterInPool,
						Timeout: pool.DefaultTimeout,
					})
					return p
				}(),
				VacateParams: &VacateParams{},
				Moves: &models.MoveClustersDetails{
					ElasticsearchClusters: []*models.MoveElasticsearchClusterDetails{
						{
							ClusterID: ec.String("63d765d37613423e97b1040257cf20c8"),
							CalculatedPlan: &models.TransientElasticsearchPlanConfiguration{
								PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{},
							},
						},
					},
					KibanaClusters: []*models.MoveKibanaClusterDetails{
						{
							ClusterID: ec.String("d7ad23ad6f064709bbae7ab87a7e1bc9"),
							CalculatedPlan: &models.TransientKibanaPlanConfiguration{
								PlanConfiguration: &models.KibanaPlanControlConfiguration{},
							},
						},
					},
				},
			}},
			wantLeftovers: nil,
			wantMoved:     true,
		},
		{
			name: "Move clusters when no filter is specified (Apm)",
			args: args{params: addAllocatorMovesToPoolParams{
				ID: "allocator-1",
				Pool: func() *pool.Pool {
					p, _ := pool.NewPool(pool.Params{
						Size:    1,
						Run:     VacateClusterInPool,
						Timeout: pool.DefaultTimeout,
					})
					return p
				}(),
				VacateParams: &VacateParams{},
				Moves: &models.MoveClustersDetails{
					ApmClusters: []*models.MoveApmClusterDetails{
						{
							ClusterID: ec.String("d7ad23ad6f064709bbae7ab87a7e1bc9"),
							CalculatedPlan: &models.TransientApmPlanConfiguration{
								PlanConfiguration: &models.ApmPlanControlConfiguration{},
							},
						},
					},
				},
			}},
			wantLeftovers: nil,
			wantMoved:     true,
		},
		{
			name: "Move clusters when no filter is specified (AppSearch)",
			args: args{params: addAllocatorMovesToPoolParams{
				ID: "allocator-1",
				Pool: func() *pool.Pool {
					p, _ := pool.NewPool(pool.Params{
						Size:    1,
						Run:     VacateClusterInPool,
						Timeout: pool.DefaultTimeout,
					})
					return p
				}(),
				VacateParams: &VacateParams{},
				Moves: &models.MoveClustersDetails{
					AppsearchClusters: []*models.MoveAppSearchDetails{
						{
							ClusterID: ec.String("d7ad23ad6f064709bbae7ab87a7e1bc9"),
							CalculatedPlan: &models.TransientAppSearchPlanConfiguration{
								PlanConfiguration: &models.AppSearchPlanControlConfiguration{},
							},
						},
					},
				},
			}},
			wantLeftovers: nil,
			wantMoved:     true,
		},
		{
			name: "Move clusters when a matching filter is specified (Kibana)",
			args: args{params: addAllocatorMovesToPoolParams{
				ID: "allocator-1",
				Pool: func() *pool.Pool {
					p, _ := pool.NewPool(pool.Params{
						Size:    1,
						Run:     VacateClusterInPool,
						Timeout: pool.DefaultTimeout,
					})
					return p
				}(),
				VacateParams: &VacateParams{KindFilter: "kibana"},
				Moves: &models.MoveClustersDetails{
					ElasticsearchClusters: []*models.MoveElasticsearchClusterDetails{
						{
							ClusterID: ec.String("63d765d37613423e97b1040257cf20c8"),
							CalculatedPlan: &models.TransientElasticsearchPlanConfiguration{
								PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{},
							},
						},
					},
					KibanaClusters: []*models.MoveKibanaClusterDetails{
						{
							ClusterID: ec.String("d7ad23ad6f064709bbae7ab87a7e1bc9"),
							CalculatedPlan: &models.TransientKibanaPlanConfiguration{
								PlanConfiguration: &models.KibanaPlanControlConfiguration{},
							},
						},
					},
				},
			}},
			wantLeftovers: nil,
			wantMoved:     true,
		},
		{
			name: "Move clusters when a matching filter is specified (Elasticsearch)",
			args: args{params: addAllocatorMovesToPoolParams{
				ID: "allocator-1",
				Pool: func() *pool.Pool {
					p, _ := pool.NewPool(pool.Params{
						Size:    1,
						Run:     VacateClusterInPool,
						Timeout: pool.DefaultTimeout,
					})
					return p
				}(),
				VacateParams: &VacateParams{KindFilter: "elasticsearch"},
				Moves: &models.MoveClustersDetails{
					ElasticsearchClusters: []*models.MoveElasticsearchClusterDetails{
						{
							ClusterID: ec.String("63d765d37613423e97b1040257cf20c8"),
							CalculatedPlan: &models.TransientElasticsearchPlanConfiguration{
								PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{},
							},
						},
					},
					KibanaClusters: []*models.MoveKibanaClusterDetails{
						{
							ClusterID: ec.String("d7ad23ad6f064709bbae7ab87a7e1bc9"),
							CalculatedPlan: &models.TransientKibanaPlanConfiguration{
								PlanConfiguration: &models.KibanaPlanControlConfiguration{},
							},
						},
					},
				},
			}},
			wantLeftovers: nil,
			wantMoved:     true,
		},
		{
			name: "Move clusters when no filter is specified (AppSearch)",
			args: args{params: addAllocatorMovesToPoolParams{
				ID: "allocator-1",
				Pool: func() *pool.Pool {
					p, _ := pool.NewPool(pool.Params{
						Size:    1,
						Run:     VacateClusterInPool,
						Timeout: pool.DefaultTimeout,
					})
					return p
				}(),
				VacateParams: &VacateParams{KindFilter: "appsearch"},
				Moves: &models.MoveClustersDetails{
					AppsearchClusters: []*models.MoveAppSearchDetails{
						{
							ClusterID: ec.String("d7ad23ad6f064709bbae7ab87a7e1bc9"),
							CalculatedPlan: &models.TransientAppSearchPlanConfiguration{
								PlanConfiguration: &models.AppSearchPlanControlConfiguration{},
							},
						},
					},
				},
			}},
			wantLeftovers: nil,
			wantMoved:     true,
		},
		{
			name: "Skip move clusters when a non-matching filter is specified",
			args: args{params: addAllocatorMovesToPoolParams{
				ID: "allocator-1",
				Pool: func() *pool.Pool {
					p, _ := pool.NewPool(pool.Params{
						Size:    1,
						Run:     VacateClusterInPool,
						Timeout: pool.DefaultTimeout,
					})
					return p
				}(),
				VacateParams: &VacateParams{KindFilter: util.Apm},
				Moves: &models.MoveClustersDetails{
					ElasticsearchClusters: []*models.MoveElasticsearchClusterDetails{
						{
							ClusterID: ec.String("63d765d37613423e97b1040257cf20c8"),
							CalculatedPlan: &models.TransientElasticsearchPlanConfiguration{
								PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{},
							},
						},
					},
					KibanaClusters: []*models.MoveKibanaClusterDetails{
						{
							ClusterID: ec.String("d7ad23ad6f064709bbae7ab87a7e1bc9"),
							CalculatedPlan: &models.TransientKibanaPlanConfiguration{
								PlanConfiguration: &models.KibanaPlanControlConfiguration{},
							},
						},
					},
				},
			}},
			wantLeftovers: nil,
			wantMoved:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLeftovers, gotMoved := addAllocatorMovesToPool(tt.args.params)
			if !reflect.DeepEqual(gotLeftovers, tt.wantLeftovers) {
				t.Errorf("addAllocatorMovesToPool() gotLeftovers = %v, want %v", gotLeftovers, tt.wantLeftovers)
			}
			if gotMoved != tt.wantMoved {
				t.Errorf("addAllocatorMovesToPool() gotMoved = %v, want %v", gotMoved, tt.wantMoved)
			}
		})
	}
}
