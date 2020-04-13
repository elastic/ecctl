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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/output"
	sdkSync "github.com/elastic/cloud-sdk-go/pkg/sync"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/go-openapi/strfmt"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment/depresource"
	"github.com/elastic/ecctl/pkg/util"
)

const planStepLogErrorMessage = "Unexpected error during step: [perform-snapshot]: [no.found.constructor.models.TimeoutException: Timeout]"

type vacateCase struct {
	topology     []vacateCaseClusters
	skipTracking bool
}

type vacateCaseClusters struct {
	// Allocator ID
	Allocator string
	// Elasticsearch clusters that will be simulated
	elasticsearch []vacateCaseClusterConfig
	// Kibana instances that will be simulated
	kibana []vacateCaseClusterConfig
	// APM clusters that will be simulated
	apm []vacateCaseClusterConfig
	// App Search clusters that will be simulated
	appsearch []vacateCaseClusterConfig
}

type vacateCaseClusterConfig struct {
	// ID of the cluster
	ID string
	// Fails the `_move` api call to actually move a cluster
	fail bool
	// Plan steps that are fetched by the
	steps [][]*models.ClusterPlanStepInfo
	// Current plan steps the move
	plan []*models.ClusterPlanStepInfo
}

func newBody(t *testing.T, v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func newDeploymentDiscovery() mock.Response {
	return mock.Response{Response: http.Response{
		StatusCode: 200,
		Body: mock.NewStructBody(models.DeploymentsSearchResponse{
			Deployments: []*models.DeploymentSearchResponse{
				{ID: ec.String("DISCOVERED_DEPLOYMENT_ID")},
			},
		}),
	}}
}

func newElasticsearchMove(t *testing.T, clusterID, allocator string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(newBody(t,
		models.MoveClustersCommandResponse{
			Moves: &models.MoveClustersDetails{
				ElasticsearchClusters: []*models.MoveElasticsearchClusterDetails{
					{
						ClusterID: ec.String(clusterID),
						CalculatedPlan: &models.TransientElasticsearchPlanConfiguration{
							PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
								MoveAllocators: []*models.AllocatorMoveRequest{
									{
										From: ec.String(allocator),
									},
								},
							},
						},
					},
				},
			},
		},
	)))
}

func newKibanaMove(t *testing.T, clusterID, allocator string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(newBody(t,
		models.MoveClustersCommandResponse{
			Moves: &models.MoveClustersDetails{
				KibanaClusters: []*models.MoveKibanaClusterDetails{
					{
						ClusterID: ec.String(clusterID),
						CalculatedPlan: &models.TransientKibanaPlanConfiguration{
							PlanConfiguration: &models.KibanaPlanControlConfiguration{
								MoveAllocators: []*models.AllocatorMoveRequest{
									{
										From: ec.String(allocator),
									},
								},
							},
						},
					},
				},
			},
		},
	)))
}

func newApmMove(t *testing.T, clusterID, allocator string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(newBody(t,
		models.MoveClustersCommandResponse{
			Moves: &models.MoveClustersDetails{
				ApmClusters: []*models.MoveApmClusterDetails{
					{
						ClusterID: ec.String(clusterID),
						CalculatedPlan: &models.TransientApmPlanConfiguration{
							PlanConfiguration: &models.ApmPlanControlConfiguration{
								MoveAllocators: []*models.AllocatorMoveRequest{
									{
										From: ec.String(allocator),
									},
								},
							},
						},
					},
				},
			},
		},
	)))
}

func newAppsearchMove(t *testing.T, clusterID, allocator string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(newBody(t,
		models.MoveClustersCommandResponse{
			Moves: &models.MoveClustersDetails{
				AppsearchClusters: []*models.MoveAppSearchDetails{{
					ClusterID: ec.String(clusterID),
					CalculatedPlan: &models.TransientAppSearchPlanConfiguration{
						PlanConfiguration: &models.AppSearchPlanControlConfiguration{
							MoveAllocators: []*models.AllocatorMoveRequest{
								{From: ec.String(allocator)},
							},
						},
					},
				}},
			},
		},
	)))
}

func newKibanaMoveFailure(t *testing.T, clusterID, allocator string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(newBody(t,
		models.MoveClustersCommandResponse{
			Failures: &models.MoveClustersDetails{
				KibanaClusters: []*models.MoveKibanaClusterDetails{
					{
						ClusterID: ec.String(clusterID),
						CalculatedPlan: &models.TransientKibanaPlanConfiguration{
							PlanConfiguration: &models.KibanaPlanControlConfiguration{
								MoveAllocators: []*models.AllocatorMoveRequest{
									{
										From: ec.String(allocator),
									},
								},
							},
						},
						Errors: []*models.BasicFailedReplyElement{
							{
								Code:    ec.String("some code"),
								Message: ec.String("failed for reason"),
							},
						},
					},
				},
			},
		},
	)))
}

func newMulipleMoves(t *testing.T, allocator string, es, kibana, apm, appsearch []string) io.ReadCloser {
	var res = models.MoveClustersCommandResponse{Moves: &models.MoveClustersDetails{}}
	for _, id := range es {
		res.Moves.ElasticsearchClusters = append(res.Moves.ElasticsearchClusters,
			&models.MoveElasticsearchClusterDetails{
				ClusterID: ec.String(id),
				CalculatedPlan: &models.TransientElasticsearchPlanConfiguration{
					PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
						MoveAllocators: []*models.AllocatorMoveRequest{
							{From: ec.String(allocator)},
						},
					},
				},
			},
		)
	}

	for _, id := range kibana {
		res.Moves.KibanaClusters = append(res.Moves.KibanaClusters,
			&models.MoveKibanaClusterDetails{
				ClusterID: ec.String(id),
				CalculatedPlan: &models.TransientKibanaPlanConfiguration{
					PlanConfiguration: &models.KibanaPlanControlConfiguration{
						MoveAllocators: []*models.AllocatorMoveRequest{
							{From: ec.String(allocator)},
						},
					},
				},
			},
		)
	}
	for _, id := range apm {
		res.Moves.ApmClusters = append(res.Moves.ApmClusters,
			&models.MoveApmClusterDetails{
				ClusterID: ec.String(id),
				CalculatedPlan: &models.TransientApmPlanConfiguration{
					PlanConfiguration: &models.ApmPlanControlConfiguration{
						MoveAllocators: []*models.AllocatorMoveRequest{
							{From: ec.String(allocator)},
						},
					},
				},
			},
		)
	}

	for _, id := range appsearch {
		res.Moves.AppsearchClusters = append(res.Moves.AppsearchClusters,
			&models.MoveAppSearchDetails{
				ClusterID: ec.String(id),
				CalculatedPlan: &models.TransientAppSearchPlanConfiguration{
					PlanConfiguration: &models.AppSearchPlanControlConfiguration{
						MoveAllocators: []*models.AllocatorMoveRequest{
							{From: ec.String(allocator)},
						},
					},
				},
			},
		)
	}
	return ioutil.NopCloser(strings.NewReader(newBody(t, res)))
}

func newPollerBody(t *testing.T, id string, pending, current *models.ElasticsearchClusterPlanInfo) io.ReadCloser {
	payload := models.DeploymentGetResponse{
		ID: ec.String("DISCOVERED_DEPLOYMENT_ID"),
		Resources: &models.DeploymentResources{
			Elasticsearch: []*models.ElasticsearchResourceInfo{
				{
					ID:    &id,
					RefID: ec.String(depresource.DefaultElasticsearchRefID),
					Info: &models.ElasticsearchClusterInfo{PlanInfo: &models.ElasticsearchClusterPlansInfo{
						Current: current, Pending: pending,
					}},
				},
			},
		},
	}
	return newPayloadFromStruct(t, payload)
}

func newApmPollerBody(t *testing.T, id string, pending, current *models.ApmPlanInfo) io.ReadCloser {
	payload := models.DeploymentGetResponse{
		ID: ec.String("DISCOVERED_DEPLOYMENT_ID"),
		Resources: &models.DeploymentResources{
			Apm: []*models.ApmResourceInfo{
				{
					ID:    &id,
					RefID: ec.String(depresource.DefaultApmRefID),
					Info: &models.ApmInfo{PlanInfo: &models.ApmPlansInfo{
						Current: current, Pending: pending,
					}},
				},
			},
		},
	}
	return newPayloadFromStruct(t, payload)
}

func newAppSearchPollerBody(t *testing.T, id string, pending, current *models.AppSearchPlanInfo) io.ReadCloser {
	payload := models.DeploymentGetResponse{
		ID: ec.String("DISCOVERED_DEPLOYMENT_ID"),
		Resources: &models.DeploymentResources{
			Appsearch: []*models.AppSearchResourceInfo{
				{
					ID:    &id,
					RefID: ec.String(depresource.DefaultAppSearchRefID),
					Info: &models.AppSearchInfo{PlanInfo: &models.AppSearchPlansInfo{
						Current: current, Pending: pending,
					}},
				},
			},
		},
	}
	return newPayloadFromStruct(t, payload)
}

func newKibanaPollerBody(t *testing.T, id string, pending, current *models.KibanaClusterPlanInfo) io.ReadCloser {
	payload := models.DeploymentGetResponse{
		ID: ec.String("DISCOVERED_DEPLOYMENT_ID"),
		Resources: &models.DeploymentResources{
			Kibana: []*models.KibanaResourceInfo{
				{
					ID:    &id,
					RefID: ec.String(depresource.DefaultKibanaRefID),
					Info: &models.KibanaClusterInfo{PlanInfo: &models.KibanaClusterPlansInfo{
						Current: current, Pending: pending,
					}},
				},
			},
		},
	}
	return newPayloadFromStruct(t, payload)
}

func newPayloadFromStruct(t *testing.T, payload interface{}) io.ReadCloser {
	var response, err = json.MarshalIndent(payload, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	response = append(response, []byte("\n\n")...)
	return ioutil.NopCloser(bytes.NewReader(response))
}

func newPlanStep(name, status string) *models.ClusterPlanStepInfo {
	started := strfmt.DateTime(time.Now())
	return &models.ClusterPlanStepInfo{
		StepID:  &name,
		Started: &started,
		Status:  &status,
	}
}

func newPlanStepWithDetails(name, status string, details []*models.ClusterPlanStepLogMessageInfo) *models.ClusterPlanStepInfo {
	step := newPlanStep(name, status)
	step.InfoLog = details
	return step
}

func newOutputResponses(res ...string) string {
	var responses = new(bytes.Buffer)
	for _, response := range res {
		responses.WriteString(fmt.Sprintln(response))
	}
	return responses.String()
}

func newMultierror(errs ...error) string {
	var merr = multierror.Error{Errors: errs}
	return merr.Error()
}

func newAllocator(t *testing.T, id, clusterID, kind string) io.ReadCloser {
	res := models.AllocatorInfo{
		AllocatorID: ec.String(id),
		Capacity: &models.AllocatorCapacity{Memory: &models.AllocatorCapacityMemory{
			Total: ec.Int32(236544),
		}},
		Features: []string{
			"apm",
			"elasticsearch",
			"elasticsearch_data",
			"ssd",
			"templates",
			"tinyauth",
		},
		HostIP: ec.String("172.25.61.100"),
		Instances: []*models.AllocatedInstanceStatus{
			{
				ClusterHealthy: ec.Bool(true),
				ClusterID:      ec.String(clusterID),
				ClusterName:    "mycluster",
				ClusterType:    ec.String(kind),
				Healthy:        ec.Bool(true),
				InstanceName:   ec.String("instance-0000000000"),
				NodeMemory:     ec.Int32(1024),
			},
		},
		Metadata: []*models.MetadataItem{
			{
				Key:   ec.String("version"),
				Value: ec.String("2017-09-30"),
			},
			{
				Key:   ec.String("instanceId"),
				Value: ec.String("i-09a0e797fb3af6864"),
			},
			{
				Key:   ec.String("architecture"),
				Value: ec.String("x86_64"),
			},
			{
				Key:   ec.String("instanceType"),
				Value: ec.String("i3.8xlarge"),
			},
			{
				Key:   ec.String("availabilityZone"),
				Value: ec.String("us-east-1a"),
			},
			{
				Key:   ec.String("pendingTime"),
				Value: ec.String("2018-05-18T13:24:21Z"),
			},
			{
				Key:   ec.String("imageId"),
				Value: ec.String("ami-ba0a51c0"),
			},
			{
				Key:   ec.String("privateIp"),
				Value: ec.String("172.25.61.100"),
			},
			{
				Key:   ec.String("region"),
				Value: ec.String("us-east-1"),
			},
		},
		PublicHostname: ec.String("172.25.61.100"),
		Settings:       &models.AllocatorSettings{},
		Status: &models.AllocatorHealthStatus{
			Connected:       ec.Bool(true),
			Healthy:         ec.Bool(true),
			MaintenanceMode: ec.Bool(false),
		},
		ZoneID: ec.String("us-east-1a"),
	}

	return ioutil.NopCloser(strings.NewReader(newBody(t, res)))
}

// newVacateTestCase builds a test case from a vacateCase topology
func newVacateTestCase(t *testing.T, tc vacateCase) *VacateParams {
	var responses = make([]mock.Response, 0, len(tc.topology))
	var allocators = make([]string, 0, len(tc.topology))

	// Do top level Vacate() Calls
	for i := range tc.topology {
		var topology = tc.topology[i]
		var alloc = topology.Allocator
		allocators = append(allocators, alloc)
		var esmoves []string
		var kibanaMoves []string
		var apmMoves []string
		var appsearchMoves []string

		// Get all moves
		for ii := range topology.elasticsearch {
			esmoves = append(esmoves, topology.elasticsearch[ii].ID)
		}
		for ii := range topology.kibana {
			kibanaMoves = append(kibanaMoves, topology.kibana[ii].ID)
		}
		for ii := range topology.apm {
			apmMoves = append(apmMoves, topology.apm[ii].ID)
		}
		for ii := range topology.appsearch {
			appsearchMoves = append(appsearchMoves, topology.appsearch[ii].ID)
		}

		responses = append(responses, mock.Response{Response: http.Response{
			Body:       newMulipleMoves(t, alloc, esmoves, kibanaMoves, apmMoves, appsearchMoves),
			StatusCode: 202,
		}})
	}

	// This calls happen inside the pool.Pool for each cluster
	for i := range tc.topology {
		var topology = tc.topology[i]
		var alloc = topology.Allocator

		for ii := range topology.elasticsearch {
			_, r := newElasticsearchVacateMove(t, alloc, topology.elasticsearch[ii])
			responses = append(responses, r...)
		}

		for ii := range topology.kibana {
			_, r := newKibanaVacateMove(t, alloc, topology.kibana[ii])
			responses = append(responses, r...)
		}

		for ii := range topology.apm {
			_, r := newAPMVacateMove(t, alloc, topology.apm[ii])
			responses = append(responses, r...)
		}

		for ii := range topology.appsearch {
			_, r := newAppsearchVacateMove(t, alloc, topology.appsearch[ii])
			responses = append(responses, r...)
		}
	}
	return &VacateParams{
		SkipTracking:   tc.skipTracking,
		Output:         output.NewDevice(sdkSync.NewBuffer()),
		Allocators:     allocators,
		API:            api.NewMock(responses...),
		OutputFormat:   "text",
		Concurrency:    1,
		MaxPollRetries: 1,
		TrackFrequency: time.Nanosecond,
	}
}

func newAPMVacateMove(t *testing.T, alloc string, move vacateCaseClusterConfig) (*api.API, []mock.Response) {
	var responses = make([]mock.Response, 0, 4)
	responses = append(responses, mock.Response{Response: http.Response{
		Body:       newAllocator(t, alloc, move.ID, util.Apm),
		StatusCode: 200,
	}}, mock.Response{Response: http.Response{
		Body:       newApmMove(t, move.ID, alloc),
		StatusCode: 202,
	}})

	if move.fail {
		body := ioutil.NopCloser(strings.NewReader(newBody(t, &models.MoveClustersCommandResponse{
			Failures: &models.MoveClustersDetails{
				ApmClusters: []*models.MoveApmClusterDetails{{
					ClusterID: ec.String(move.ID),
					CalculatedPlan: &models.TransientApmPlanConfiguration{
						PlanConfiguration: &models.ApmPlanControlConfiguration{
							MoveAllocators: []*models.AllocatorMoveRequest{{From: ec.String(alloc)}},
						},
					},
					Errors: []*models.BasicFailedReplyElement{
						{
							Code:    ec.String("a code"),
							Message: ec.String("a message"),
						},
					},
				}},
			}})))
		// Return a response with a failed move
		responses = append(responses, mock.Response{Response: http.Response{
			Body:       body,
			StatusCode: 202,
		}})
		// No extra responses should be given back for this cluster
		// when a move failures happens.
		return api.NewMock(responses...), responses
	}

	responses = append(responses, mock.Response{Response: http.Response{
		Body:       newApmMove(t, move.ID, alloc),
		StatusCode: 202,
	}}, newDeploymentDiscovery())

	// Define steps
	for iii := range move.steps {
		var step = move.steps[iii]
		responses = append(responses, mock.Response{Response: http.Response{
			StatusCode: 200,
			Body: newApmPollerBody(t, move.ID,
				&models.ApmPlanInfo{PlanAttemptLog: step},
				nil,
			),
		}})
	}

	// Plan finished
	responses = append(responses,
		util.PlanNotFound,
		util.PlanNotFound,
		mock.Response{Response: http.Response{
			StatusCode: 200,
			Body: newApmPollerBody(t, move.ID,
				nil,
				&models.ApmPlanInfo{PlanAttemptLog: move.plan},
			),
		}},
	)

	return api.NewMock(responses...), responses
}

func newKibanaVacateMove(t *testing.T, alloc string, move vacateCaseClusterConfig) (*api.API, []mock.Response) {
	var responses = make([]mock.Response, 0, 4)
	responses = append(responses, mock.Response{Response: http.Response{
		Body:       newAllocator(t, alloc, move.ID, "kibana"),
		StatusCode: 200,
	}}, mock.Response{Response: http.Response{
		Body:       newKibanaMove(t, move.ID, alloc),
		StatusCode: 202,
	}})

	if move.fail {
		body := ioutil.NopCloser(strings.NewReader(newBody(t, &models.MoveClustersCommandResponse{
			Failures: &models.MoveClustersDetails{
				KibanaClusters: []*models.MoveKibanaClusterDetails{{
					ClusterID: ec.String(move.ID),
					CalculatedPlan: &models.TransientKibanaPlanConfiguration{
						PlanConfiguration: &models.KibanaPlanControlConfiguration{
							MoveAllocators: []*models.AllocatorMoveRequest{{From: ec.String(alloc)}},
						},
					},
					Errors: []*models.BasicFailedReplyElement{
						{
							Code:    ec.String("a code"),
							Message: ec.String("a message"),
						},
					},
				}},
			}})))
		// Return a response with a failed move
		responses = append(responses, mock.Response{Response: http.Response{
			Body:       body,
			StatusCode: 202,
		}})

		// No extra responses should be given back for this cluster
		// when a move failures happens.
		return api.NewMock(responses...), responses
	}

	responses = append(responses, mock.Response{Response: http.Response{
		Body:       newKibanaMove(t, move.ID, alloc),
		StatusCode: 202,
	}}, newDeploymentDiscovery())

	// Define steps
	for iii := range move.steps {
		var step = move.steps[iii]
		responses = append(responses, mock.Response{Response: http.Response{
			StatusCode: 200,
			Body: newKibanaPollerBody(t, move.ID,
				&models.KibanaClusterPlanInfo{PlanAttemptLog: step},
				nil,
			),
		}})
	}

	responses = append(responses,
		util.PlanNotFound,
		util.PlanNotFound,
		mock.Response{Response: http.Response{
			StatusCode: 200,
			Body: newKibanaPollerBody(t, move.ID,
				nil,
				&models.KibanaClusterPlanInfo{PlanAttemptLog: move.plan},
			),
		}},
	)

	return api.NewMock(responses...), responses
}

func newElasticsearchVacateMove(t *testing.T, alloc string, move vacateCaseClusterConfig) (*api.API, []mock.Response) {
	var responses = make([]mock.Response, 0, 4)
	responses = append(responses, mock.Response{Response: http.Response{
		Body:       newAllocator(t, alloc, move.ID, "elasticsearch"),
		StatusCode: 200,
	}}, mock.Response{Response: http.Response{
		Body:       newElasticsearchMove(t, move.ID, alloc),
		StatusCode: 202,
	}})

	if move.fail {
		body := ioutil.NopCloser(strings.NewReader(newBody(t, &models.MoveClustersCommandResponse{
			Failures: &models.MoveClustersDetails{
				ElasticsearchClusters: []*models.MoveElasticsearchClusterDetails{{
					ClusterID: ec.String(move.ID),
					CalculatedPlan: &models.TransientElasticsearchPlanConfiguration{
						PlanConfiguration: &models.ElasticsearchPlanControlConfiguration{
							MoveAllocators: []*models.AllocatorMoveRequest{{From: ec.String(alloc)}},
						},
					},
					Errors: []*models.BasicFailedReplyElement{{
						Code:    ec.String("a code"),
						Message: ec.String("a message"),
					}},
				}},
			}})))
		// Return a response with a failed move
		responses = append(responses, mock.Response{Response: http.Response{
			Body:       body,
			StatusCode: 202,
		}})
		// No extra responses should be given back for this cluster
		// when a move failures happens.
		return api.NewMock(responses...), responses
	}

	// do the actual cluster move from the calculated plan
	responses = append(responses, mock.Response{Response: http.Response{
		Body:       newElasticsearchMove(t, move.ID, alloc),
		StatusCode: 202,
	}}, newDeploymentDiscovery())

	// Define steps
	for iii := range move.steps {
		var step = move.steps[iii]
		responses = append(responses, mock.Response{Response: http.Response{
			StatusCode: 200,
			Body: newPollerBody(t, move.ID,
				&models.ElasticsearchClusterPlanInfo{PlanAttemptLog: step},
				nil,
			),
		}})
	}

	// Plan finished
	responses = append(responses,
		util.PlanNotFound,
		util.PlanNotFound,
		mock.Response{Response: http.Response{
			StatusCode: 200,
			Body: newPollerBody(t, move.ID,
				nil,
				&models.ElasticsearchClusterPlanInfo{PlanAttemptLog: move.plan},
			),
		}},
	)

	return api.NewMock(responses...), responses
}

func newAppsearchVacateMove(t *testing.T, alloc string, move vacateCaseClusterConfig) (*api.API, []mock.Response) {
	var responses = make([]mock.Response, 0, 4)
	responses = append(responses, mock.Response{Response: http.Response{
		Body:       newAllocator(t, alloc, move.ID, "appsearch"),
		StatusCode: 200,
	}}, mock.Response{Response: http.Response{
		Body:       newAppsearchMove(t, move.ID, alloc),
		StatusCode: 202,
	}})

	if move.fail {
		body := ioutil.NopCloser(strings.NewReader(newBody(t, &models.MoveClustersCommandResponse{
			Failures: &models.MoveClustersDetails{
				ApmClusters: []*models.MoveApmClusterDetails{{
					ClusterID: ec.String(move.ID),
					CalculatedPlan: &models.TransientApmPlanConfiguration{
						PlanConfiguration: &models.ApmPlanControlConfiguration{
							MoveAllocators: []*models.AllocatorMoveRequest{{From: ec.String(alloc)}},
						},
					},
					Errors: []*models.BasicFailedReplyElement{
						{
							Code:    ec.String("a code"),
							Message: ec.String("a message"),
						},
					},
				}},
			}})))
		// Return a response with a failed move
		responses = append(responses, mock.Response{Response: http.Response{
			Body:       body,
			StatusCode: 202,
		}})
		// No extra responses should be given back for this cluster
		// when a move failures happens.
		return api.NewMock(responses...), responses
	}

	responses = append(responses, mock.Response{Response: http.Response{
		Body:       newAppsearchMove(t, move.ID, alloc),
		StatusCode: 202,
	}}, newDeploymentDiscovery())

	// Define steps
	for iii := range move.steps {
		var step = move.steps[iii]
		responses = append(responses, mock.Response{Response: http.Response{
			StatusCode: 200,
			Body: newAppSearchPollerBody(t, move.ID,
				&models.AppSearchPlanInfo{PlanAttemptLog: step},
				nil,
			),
		}})
	}

	// Plan finished
	responses = append(responses,
		util.PlanNotFound,
		util.PlanNotFound,
		mock.Response{Response: http.Response{
			StatusCode: 200,
			Body: newAppSearchPollerBody(t, move.ID,
				nil,
				&models.AppSearchPlanInfo{PlanAttemptLog: move.plan},
			),
		}},
	)

	return api.NewMock(responses...), responses
}
