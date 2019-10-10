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
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/go-openapi/strfmt"
)

const (
	formatErrType = "got error type = %+v, want %+v"
	// ValidClusterID holds a valid cluster id value
	ValidClusterID = "320b7b540dfc967a7a649c18e2fce4ed"
	// InvalidClusterID holds an invalid valid cluster id value
	InvalidClusterID = "!23"
	// ValidDuration holds a valid duration value
	ValidDuration = "3 seconds"
)

var (
	// PlanNotFound represents a mocked response resulting from a plan not found.
	PlanNotFound = mock.Response{Response: http.Response{
		StatusCode: 404,
		Body:       mock.NewStringBody(`{}`),
	}}
)

// NewSuccessfulPlan returns a mocked response from a successful plan.
func NewSuccessfulPlan() mock.Response {
	return mock.Response{Response: http.Response{
		StatusCode: 200,
		Body: mock.NewStructBody(&models.ElasticsearchClusterPlansInfo{
			Current: &models.ElasticsearchClusterPlanInfo{
				PlanAttemptLog: []*models.ClusterPlanStepInfo{
					{
						Status: ec.String("success"),
						StepID: ec.String("step-1"),
					},
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
		PlanNotFound,
		PlanNotFound,
		NewSuccessfulPlan(),
	)

	return response
}

// CheckErrType receives two errors, if either one is not nil, it triggers a
// comparison between the two, returning an error if the errors are not equal
// This function is useful for testing to ensure specific error types.
func CheckErrType(got, want error) error {
	var emptyErrors = want == nil && got == nil

	if !emptyErrors && !reflect.DeepEqual(got, want) {
		return fmt.Errorf(formatErrType, got, want)
	}
	return nil
}

// ParseDate parses a string to date and if parsing generates an error,
// it fails the given testing suite, otherwise it returns the
// strfmt.Datetime parsed value
func ParseDate(t *testing.T, date string) strfmt.DateTime {
	dt, err := strfmt.ParseDateTime(date)
	if err != nil {
		t.Fatal(err)
	}

	return dt
}
