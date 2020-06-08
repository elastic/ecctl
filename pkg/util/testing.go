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

package util

import (
	"bytes"
	"net/http"
	"time"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/output"
	"github.com/elastic/cloud-sdk-go/pkg/plan"
	"github.com/elastic/cloud-sdk-go/pkg/plan/planutil"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/go-openapi/strfmt"
)

const (
	// ValidClusterID holds a valid cluster id value
	ValidClusterID = "320b7b540dfc967a7a649c18e2fce4ed"
	// InvalidClusterID holds an invalid valid cluster id value
	InvalidClusterID = "!23"
	// ValidDuration holds a valid duration value
	ValidDuration = "3 seconds"
)

// NewSuccessfulCurrentPlan returns a mocked response from a successful plan.
func NewSuccessfulCurrentPlan(currentTime ...time.Time) mock.Response {
	var now strfmt.DateTime
	if len(currentTime) > 0 {
		now = strfmt.DateTime(currentTime[0])
	}
	return mock.Response{Response: http.Response{
		StatusCode: 200,
		Body: mock.NewStructBody(&models.DeploymentGetResponse{
			ID: ec.String("DEPLOYMENT_TEST_ID"),
			Resources: &models.DeploymentResources{
				Apm: []*models.ApmResourceInfo{
					{ID: ec.String(ValidClusterID), RefID: ec.String("main-apm"), Info: &models.ApmInfo{PlanInfo: &models.ApmPlansInfo{
						Current: &models.ApmPlanInfo{PlanAttemptLog: []*models.ClusterPlanStepInfo{
							{Status: ec.String("success"), StepID: ec.String("plan-completed"), Started: &now},
						}},
					}}},
				},
				Appsearch: []*models.AppSearchResourceInfo{
					{ID: ec.String(ValidClusterID), RefID: ec.String("main-appsearch"), Info: &models.AppSearchInfo{PlanInfo: &models.AppSearchPlansInfo{
						Current: &models.AppSearchPlanInfo{PlanAttemptLog: []*models.ClusterPlanStepInfo{
							{Status: ec.String("success"), StepID: ec.String("plan-completed"), Started: &now},
						}},
					}}},
				},
				Elasticsearch: []*models.ElasticsearchResourceInfo{
					{ID: ec.String(ValidClusterID), RefID: ec.String("main-elasticsearch"), Info: &models.ElasticsearchClusterInfo{PlanInfo: &models.ElasticsearchClusterPlansInfo{
						Current: &models.ElasticsearchClusterPlanInfo{PlanAttemptLog: []*models.ClusterPlanStepInfo{
							{Status: ec.String("success"), StepID: ec.String("plan-completed"), Started: &now},
						}},
					}}},
				},
				Kibana: []*models.KibanaResourceInfo{
					{ID: ec.String(ValidClusterID), RefID: ec.String("main-kibana"), Info: &models.KibanaClusterInfo{PlanInfo: &models.KibanaClusterPlansInfo{
						Current: &models.KibanaClusterPlanInfo{PlanAttemptLog: []*models.ClusterPlanStepInfo{
							{Status: ec.String("success"), StepID: ec.String("plan-completed"), Started: &now},
						}},
					}}},
				},
			},
		}),
	}}
}

// NewFailedPlanUnknown returns a mocked response from a failed plan.
func NewFailedPlanUnknown() mock.Response {
	return mock.Response{Response: http.Response{
		StatusCode: 200,
		Body: mock.NewStructBody(&models.ElasticsearchClusterPlansInfo{
			Current: &models.ElasticsearchClusterPlanInfo{
				PlanAttemptLog: []*models.ClusterPlanStepInfo{
					{
						Status: ec.String("error"),
						StepID: ec.String("step-1"),
					},
				},
			},
		}),
	}}
}

// AppendTrackResponses is aimed to be used while testing to add the tracking
// responses as a successful plan, it assumes that there's 1 retry performed
// by the client.
func AppendTrackResponses(res ...mock.Response) []mock.Response {
	var response = append(res,
		NewSuccessfulCurrentPlan(),
		NewSuccessfulCurrentPlan(),
		NewSuccessfulCurrentPlan(),
	)

	return response
}

// NewMockTrackChangeParams creates new tracking params for test purposes.
func NewMockTrackChangeParams(id string) planutil.TrackChangeParams {
	var res = AppendTrackResponses()
	if id == "" {
		res = AppendTrackResponses(mock.New200Response(mock.NewStructBody(models.DeploymentsSearchResponse{
			Deployments: []*models.DeploymentSearchResponse{
				{ID: ec.String("cbb4bc6c09684c86aa5de54c05ea1d38")},
			},
		})))
	}
	return planutil.TrackChangeParams{
		TrackChangeParams: plan.TrackChangeParams{
			API:          api.NewMock(res...),
			DeploymentID: id,
		},
		Writer: output.NewDevice(new(bytes.Buffer)),
		Format: "text",
	}
}
