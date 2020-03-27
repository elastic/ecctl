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
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/output"
	sdkSync "github.com/elastic/cloud-sdk-go/pkg/sync"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
)

func TestVacate(t *testing.T) {
	type args struct {
		params *VacateParams
		buf    *sdkSync.Buffer
	}
	tests := []struct {
		name string
		args args
		err  string
		want string
	}{
		{
			name: "Succeeds moving a single ES cluster from a single allocator",
			args: args{
				buf: sdkSync.NewBuffer(),
				params: newVacateTestCase(t, vacateCase{
					topology: []vacateCaseClusters{
						{
							Allocator: "allocatorID",
							elasticsearch: []vacateCaseClusterConfig{
								{
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
								},
							},
						},
					},
				}),
			},
			want: newOutputResponses(
				"Deployment [DISCOVERED_DEPLOYMENT_ID] - [Elasticsearch][3ee11eb40eda22cac0cce259625c6734]: running step \"step2\" (Plan duration )...",
				"\x1b[92;mDeployment [DISCOVERED_DEPLOYMENT_ID] - [Elasticsearch][3ee11eb40eda22cac0cce259625c6734]: finished running all the plan steps\x1b[0m (Total plan duration )",
			),
		},
		{
			name: "Succeeds moving a single kibana instance from a single allocator",
			args: args{
				buf: sdkSync.NewBuffer(),
				params: newVacateTestCase(t, vacateCase{
					topology: []vacateCaseClusters{
						{
							Allocator: "allocatorID",
							kibana: []vacateCaseClusterConfig{
								{
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
								},
							},
						},
					},
				}),
			},
			want: newOutputResponses(
				"Deployment [DISCOVERED_DEPLOYMENT_ID] - [Kibana][3ee11eb40eda22cac0cce259625c6734]: running step \"step2\" (Plan duration )...",
				"\x1b[92;mDeployment [DISCOVERED_DEPLOYMENT_ID] - [Kibana][3ee11eb40eda22cac0cce259625c6734]: finished running all the plan steps\x1b[0m (Total plan duration )",
			),
		},
		{
			name: "Succeeds moving a single APM cluster from a single allocator",
			args: args{
				buf: sdkSync.NewBuffer(),
				params: newVacateTestCase(t, vacateCase{
					topology: []vacateCaseClusters{
						{
							Allocator: "allocatorID",
							apm: []vacateCaseClusterConfig{
								{
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
								},
							},
						},
					},
				}),
			},
			want: newOutputResponses(
				"Deployment [DISCOVERED_DEPLOYMENT_ID] - [Apm][3ee11eb40eda22cac0cce259625c6734]: running step \"step2\" (Plan duration )...",
				"\x1b[92;mDeployment [DISCOVERED_DEPLOYMENT_ID] - [Apm][3ee11eb40eda22cac0cce259625c6734]: finished running all the plan steps\x1b[0m (Total plan duration )",
			),
		},
		{
			name: "Succeeds moving a single Appsearch cluster from a single allocator (without tracking)",
			args: args{
				buf: sdkSync.NewBuffer(),
				params: newVacateTestCase(t, vacateCase{
					skipTracking: true,
					topology: []vacateCaseClusters{
						{
							Allocator: "allocatorID",
							appsearch: []vacateCaseClusterConfig{
								{
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
								},
							},
						},
					},
				}),
			},
		},
		{
			name: "Succeeds moving a single Appsearch cluster from a single allocator",
			args: args{
				buf: sdkSync.NewBuffer(),
				params: newVacateTestCase(t, vacateCase{
					topology: []vacateCaseClusters{
						{
							Allocator: "allocatorID",
							appsearch: []vacateCaseClusterConfig{
								{
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
								},
							},
						},
					},
				}),
			},
			want: newOutputResponses(
				"Deployment [DISCOVERED_DEPLOYMENT_ID] - [Appsearch][3ee11eb40eda22cac0cce259625c6734]: running step \"step2\" (Plan duration )...",
				"\x1b[92;mDeployment [DISCOVERED_DEPLOYMENT_ID] - [Appsearch][3ee11eb40eda22cac0cce259625c6734]: finished running all the plan steps\x1b[0m (Total plan duration )",
			),
		},
		{
			name: "Succeeds moving a multiple clusters from a single allocator",
			args: args{
				buf: sdkSync.NewBuffer(),
				params: newVacateTestCase(t, vacateCase{topology: []vacateCaseClusters{
					{
						Allocator: "allocatorID",
						elasticsearch: []vacateCaseClusterConfig{
							{
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
							},
						},
						kibana: []vacateCaseClusterConfig{
							{
								ID: "2ee11eb40eda22cac0cce259625c6734",
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
							},
						},
					},
				}}),
			},
			want: newOutputResponses(
				"Deployment [DISCOVERED_DEPLOYMENT_ID] - [Elasticsearch][3ee11eb40eda22cac0cce259625c6734]: running step \"step2\" (Plan duration )...",
				"\x1b[92;mDeployment [DISCOVERED_DEPLOYMENT_ID] - [Elasticsearch][3ee11eb40eda22cac0cce259625c6734]: finished running all the plan steps\x1b[0m (Total plan duration )",
				"Deployment [DISCOVERED_DEPLOYMENT_ID] - [Kibana][2ee11eb40eda22cac0cce259625c6734]: running step \"step2\" (Plan duration )...",
				"\x1b[92;mDeployment [DISCOVERED_DEPLOYMENT_ID] - [Kibana][2ee11eb40eda22cac0cce259625c6734]: finished running all the plan steps\x1b[0m (Total plan duration )",
			),
		},
		{
			name: "Succeeds moving multiple clusters from a multiple allocators",
			args: args{
				buf: sdkSync.NewBuffer(),
				params: newVacateTestCase(t, vacateCase{topology: []vacateCaseClusters{
					{
						Allocator: "allocatorID-1",
						elasticsearch: []vacateCaseClusterConfig{
							{
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
							},
						},
						kibana: []vacateCaseClusterConfig{
							{
								ID: "2ee11eb40eda22cac0cce259625c6734",
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
							},
						},
					},
					{
						Allocator: "allocatorID-2",
						elasticsearch: []vacateCaseClusterConfig{
							{
								ID: "5ee11eb40eda22cac0cce259625c6734",
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
							},
						},
						kibana: []vacateCaseClusterConfig{
							{
								ID: "4ee11eb40eda22cac0cce259625c6734",
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
							},
						},
					},
				}}),
			},
			want: newOutputResponses(
				"Deployment [DISCOVERED_DEPLOYMENT_ID] - [Elasticsearch][3ee11eb40eda22cac0cce259625c6734]: running step \"step2\" (Plan duration )...",
				"\x1b[92;mDeployment [DISCOVERED_DEPLOYMENT_ID] - [Elasticsearch][3ee11eb40eda22cac0cce259625c6734]: finished running all the plan steps\x1b[0m (Total plan duration )",
				"Deployment [DISCOVERED_DEPLOYMENT_ID] - [Kibana][2ee11eb40eda22cac0cce259625c6734]: running step \"step2\" (Plan duration )...",
				"\x1b[92;mDeployment [DISCOVERED_DEPLOYMENT_ID] - [Kibana][2ee11eb40eda22cac0cce259625c6734]: finished running all the plan steps\x1b[0m (Total plan duration )",
				"Deployment [DISCOVERED_DEPLOYMENT_ID] - [Elasticsearch][5ee11eb40eda22cac0cce259625c6734]: running step \"step2\" (Plan duration )...",
				"\x1b[92;mDeployment [DISCOVERED_DEPLOYMENT_ID] - [Elasticsearch][5ee11eb40eda22cac0cce259625c6734]: finished running all the plan steps\x1b[0m (Total plan duration )",
				"Deployment [DISCOVERED_DEPLOYMENT_ID] - [Kibana][4ee11eb40eda22cac0cce259625c6734]: running step \"step2\" (Plan duration )...",
				"\x1b[92;mDeployment [DISCOVERED_DEPLOYMENT_ID] - [Kibana][4ee11eb40eda22cac0cce259625c6734]: finished running all the plan steps\x1b[0m (Total plan duration )",
			),
		},
		{
			name: "Moving multiple clusters from a multiple allocators that contain track failures",
			args: args{
				buf: sdkSync.NewBuffer(),
				params: newVacateTestCase(t, vacateCase{topology: []vacateCaseClusters{
					{
						Allocator: "allocatorID-1",
						elasticsearch: []vacateCaseClusterConfig{
							{
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
							},
						},
						kibana: []vacateCaseClusterConfig{
							{
								ID: "2ee11eb40eda22cac0cce259625c6734",
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
							},
						},
					},
					{
						Allocator: "allocatorID-2",
						elasticsearch: []vacateCaseClusterConfig{
							{
								ID: "5ee11eb40eda22cac0cce259625c6734",
								steps: [][]*models.ClusterPlanStepInfo{
									{
										newPlanStep("step1", "success"),
										newPlanStep("step2", "pending"),
									},
									{
										newPlanStep("step1", "success"),
										newPlanStepWithDetails("step2", "success", nil),
										newPlanStepWithDetails("step3", "error", []*models.ClusterPlanStepLogMessageInfo{{
											Message: ec.String(planStepLogErrorMessage),
										}}),
									},
								},
								plan: []*models.ClusterPlanStepInfo{
									newPlanStep("step1", "success"),
									newPlanStep("step2", "success"),
									newPlanStepWithDetails("last step", "error", []*models.ClusterPlanStepLogMessageInfo{{
										Message: ec.String(planStepLogErrorMessage),
									}}),
								},
							},
						},
						kibana: []vacateCaseClusterConfig{
							{
								ID: "4ee11eb40eda22cac0cce259625c6734",
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
							},
						},
					},
				}}),
			},
			want: newOutputResponses(
				"Deployment [DISCOVERED_DEPLOYMENT_ID] - [Elasticsearch][3ee11eb40eda22cac0cce259625c6734]: running step \"step2\" (Plan duration )...",
				"\x1b[92;mDeployment [DISCOVERED_DEPLOYMENT_ID] - [Elasticsearch][3ee11eb40eda22cac0cce259625c6734]: finished running all the plan steps\x1b[0m (Total plan duration )",
				"Deployment [DISCOVERED_DEPLOYMENT_ID] - [Kibana][2ee11eb40eda22cac0cce259625c6734]: running step \"step2\" (Plan duration )...",
				"\x1b[92;mDeployment [DISCOVERED_DEPLOYMENT_ID] - [Kibana][2ee11eb40eda22cac0cce259625c6734]: finished running all the plan steps\x1b[0m (Total plan duration )",
				"Deployment [DISCOVERED_DEPLOYMENT_ID] - [Elasticsearch][5ee11eb40eda22cac0cce259625c6734]: running step \"step2\" (Plan duration )...",
				`Deployment [DISCOVERED_DEPLOYMENT_ID] - [Elasticsearch][5ee11eb40eda22cac0cce259625c6734]: running step "step3" caught error: "Unexpected error during step: [perform-snapshot]: [no.found.constructor.models.TimeoutException: Timeout]" (Plan duration )...`,
				`Deployment [DISCOVERED_DEPLOYMENT_ID] - [Elasticsearch][5ee11eb40eda22cac0cce259625c6734]: running step "last step" caught error: "Unexpected error during step: [perform-snapshot]: [no.found.constructor.models.TimeoutException: Timeout]" (Plan duration )...`,
				"Deployment [DISCOVERED_DEPLOYMENT_ID] - [Kibana][4ee11eb40eda22cac0cce259625c6734]: running step \"step2\" (Plan duration )...",
				"\x1b[92;mDeployment [DISCOVERED_DEPLOYMENT_ID] - [Kibana][4ee11eb40eda22cac0cce259625c6734]: finished running all the plan steps\x1b[0m (Total plan duration )",
			),
			err: `1 error occurred:
	* resource id [5ee11eb40eda22cac0cce259625c6734][elasticsearch] Unexpected error during step: [perform-snapshot]: [no.found.constructor.models.TimeoutException: Timeout]

`,
		},
		{
			name: "Moving multiple clusters from a multiple allocators that fail to move",
			args: args{
				buf: sdkSync.NewBuffer(),
				params: newVacateTestCase(t, vacateCase{topology: []vacateCaseClusters{
					{
						Allocator: "allocatorID-1",
						elasticsearch: []vacateCaseClusterConfig{
							{
								ID:   "3ee11eb40eda22cac0cce259625c6734",
								fail: true,
							},
						},
						kibana: []vacateCaseClusterConfig{
							{
								ID: "2ee11eb40eda22cac0cce259625c6734",
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
							},
						},
					},
					{
						Allocator: "allocatorID-2",
						elasticsearch: []vacateCaseClusterConfig{
							{
								ID: "5ee11eb40eda22cac0cce259625c6734",
								steps: [][]*models.ClusterPlanStepInfo{
									{
										newPlanStep("step1", "success"),
										newPlanStep("step2", "pending"),
									},
									{
										newPlanStep("step1", "success"),
										newPlanStep("step2", "success"),
										newPlanStepWithDetails("step3", "error", []*models.ClusterPlanStepLogMessageInfo{{
											Message: ec.String(planStepLogErrorMessage),
										}}),
									},
								},
								plan: []*models.ClusterPlanStepInfo{
									newPlanStep("step1", "success"),
									newPlanStep("step2", "success"),
									newPlanStepWithDetails("step3", "error", []*models.ClusterPlanStepLogMessageInfo{{
										Message: ec.String(planStepLogErrorMessage),
									}}),
									newPlanStepWithDetails("plan-completed", "error", []*models.ClusterPlanStepLogMessageInfo{{
										Message: ec.String(planStepLogErrorMessage),
									}}),
								},
							},
						},
						kibana: []vacateCaseClusterConfig{
							{
								ID: "4ee11eb40eda22cac0cce259625c6734",
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
							},
						},
					},
				}}),
			},
			want: newOutputResponses(
				"Deployment [DISCOVERED_DEPLOYMENT_ID] - [Kibana][2ee11eb40eda22cac0cce259625c6734]: running step \"step2\" (Plan duration )...",
				"\x1b[92;mDeployment [DISCOVERED_DEPLOYMENT_ID] - [Kibana][2ee11eb40eda22cac0cce259625c6734]: finished running all the plan steps\x1b[0m (Total plan duration )",
				"Deployment [DISCOVERED_DEPLOYMENT_ID] - [Elasticsearch][5ee11eb40eda22cac0cce259625c6734]: running step \"step2\" (Plan duration )...",
				`Deployment [DISCOVERED_DEPLOYMENT_ID] - [Elasticsearch][5ee11eb40eda22cac0cce259625c6734]: running step "step3" caught error: "Unexpected error during step: [perform-snapshot]: [no.found.constructor.models.TimeoutException: Timeout]" (Plan duration )...`,
				"\x1b[91;1mDeployment [DISCOVERED_DEPLOYMENT_ID] - [Elasticsearch][5ee11eb40eda22cac0cce259625c6734]: caught error: \"Unexpected error during step: [perform-snapshot]: [no.found.constructor.models.TimeoutException: Timeout]\"\x1b[0m (Total plan duration )",
				"Deployment [DISCOVERED_DEPLOYMENT_ID] - [Kibana][4ee11eb40eda22cac0cce259625c6734]: running step \"step2\" (Plan duration )...",
				"\x1b[92;mDeployment [DISCOVERED_DEPLOYMENT_ID] - [Kibana][4ee11eb40eda22cac0cce259625c6734]: finished running all the plan steps\x1b[0m (Total plan duration )",
			),
			err: `2 errors occurred:
	* resource id [3ee11eb40eda22cac0cce259625c6734][elasticsearch] failed vacating, reason: code: a code, message: a message
	* deployment [DISCOVERED_DEPLOYMENT_ID] - [elasticsearch][5ee11eb40eda22cac0cce259625c6734]: caught error: "Unexpected error during step: [perform-snapshot]: [no.found.constructor.models.TimeoutException: Timeout]"

`,
		},
		{
			name: "Vacate has some failures",
			args: args{
				buf: sdkSync.NewBuffer(),
				params: &VacateParams{
					Allocators:     []string{"allocatorID"},
					Concurrency:    1,
					MaxPollRetries: 1,
					TrackFrequency: time.Nanosecond,
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       newKibanaMoveFailure(t, "3ee11eb40eda22cac0cce259625c6734", "allocatorID"),
						StatusCode: 202,
					}}),
				},
			},
			err: `1 error occurred:
	* resource id [3ee11eb40eda22cac0cce259625c6734][kibana] failed vacating, reason: code: some code, message: failed for reason

`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.buf != nil {
				tt.args.params.Output = output.NewDevice(tt.args.buf)
			}

			if err := Vacate(tt.args.params); err != nil && err.Error() != tt.err {
				t.Errorf("Vacate() error = %v, wantErr %v", err, tt.err)
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

func TestVacateInterrupt(t *testing.T) {
	type args struct {
		params *VacateParams
	}
	tests := []struct {
		name string
		args args
		err  string
		want string
	}{
		{
			name: "Interrupts the vacate",
			args: args{
				params: newVacateTestCase(t, vacateCase{topology: []vacateCaseClusters{
					{
						Allocator: "allocatorID",
						elasticsearch: []vacateCaseClusterConfig{
							{
								ID: "3ee11eb40eda22cac0cce259625c6734",
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
									{
										newPlanStep("step1", "success"),
										newPlanStep("step2", "success"),
										newPlanStep("step3", "success"),
										newPlanStep("step4", "pending"),
									},
								},
								plan: []*models.ClusterPlanStepInfo{
									newPlanStep("step1", "success"),
									newPlanStep("step2", "success"),
									newPlanStep("step3", "success"),
									newPlanStep("step4", "success"),
								},
							},
						},
						kibana: []vacateCaseClusterConfig{
							{
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
									{
										newPlanStep("step1", "success"),
										newPlanStep("step2", "success"),
										newPlanStep("step3", "success"),
										newPlanStep("step4", "pending"),
									},
								},
								plan: []*models.ClusterPlanStepInfo{
									newPlanStep("step1", "success"),
									newPlanStep("step2", "success"),
									newPlanStep("step3", "success"),
									newPlanStep("step4", "success"),
								},
							},
						},
					},
				}}),
			},
			want: newOutputResponses(
				"Deployment [DISCOVERED_DEPLOYMENT_ID] - [Elasticsearch][3ee11eb40eda22cac0cce259625c6734]: running step \"step2\" (Plan duration )...",
				"pool: received interrupt, stopping pool...",
				"Deployment [DISCOVERED_DEPLOYMENT_ID] - [Elasticsearch][3ee11eb40eda22cac0cce259625c6734]: running step \"step3\" (Plan duration )...",
				"Deployment [DISCOVERED_DEPLOYMENT_ID] - [Elasticsearch][3ee11eb40eda22cac0cce259625c6734]: running step \"step4\" (Plan duration )...",
				"\x1b[92;mDeployment [DISCOVERED_DEPLOYMENT_ID] - [Elasticsearch][3ee11eb40eda22cac0cce259625c6734]: finished running all the plan steps\x1b[0m (Total plan duration )",
			),
			err: `1 error occurred:
	* allocator allocatorID: resource id [2ee11eb40eda22cac0cce259625c6734][kibana]: was either cancelled or not processed, follow up accordingly

`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var once sync.Once
			// Sends a signal once the first write is received.
			buf := sdkSync.NewBuffer(func() {
				once.Do(func() {
					p, err := os.FindProcess(os.Getpid())
					if err != nil {
						t.Fatal(err)
					}

					if err := p.Signal(os.Interrupt); err != nil {
						t.Fatal(err)
					}
				})
			})
			tt.args.params.Output = output.NewDevice(buf)

			if err := Vacate(tt.args.params); err != nil && err.Error() != tt.err {
				t.Errorf("Vacate() error = %v, wantErr %v", err, tt.err)
			}

			var got string
			if buf != nil {
				got = regexp.MustCompile(`duration.*\)`).
					ReplaceAllString(buf.String(), "duration )")
			}

			if buf != nil && tt.want != got {
				wantLines := strings.Split(tt.want, "\n")
				var matched bool
				for _, want := range wantLines {
					if strings.Contains(got, want) {
						matched = true
					} else {
						matched = false
					}
				}

				if !matched {
					t.Errorf("VacateCluster() output = \n%v, want \n%v", got, tt.want)
				}
			}
		})
	}
}
