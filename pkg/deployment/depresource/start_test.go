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

package depresource

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment"
	"github.com/elastic/ecctl/pkg/util"
)

func TestStartInstancesParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  StartInstancesParams
		wantErr bool
		err     error
	}{
		{
			name:   "validate should return all possible errors",
			params: StartInstancesParams{},
			err: &multierror.Error{
				Errors: []error{
					errors.New("deployment start: at least 1 instance ID must be provided"),
					errors.New("api reference is required for command"),
					errors.New(`id "" is invalid`),
					errors.New("deployment resource kind cannot be empty"),
					errors.New("failed auto-discovering the resource ref id: api reference is required for command"),
					errors.New(`failed auto-discovering the resource ref id: id "" is invalid`),
				},
			},
			wantErr: true,
		},
		{
			name: "validate should return error on missing instance IDs",
			params: StartInstancesParams{
				StartParams: StartParams{
					ResourceParams: deployment.ResourceParams{
						API:          &api.API{},
						DeploymentID: "f1d329b0fb34470ba8b18361cabdd2bc",
						Kind:         "elasticsearch",
						RefID:        "main-elasticsearch",
					},
				},
			},
			err: &multierror.Error{
				Errors: []error{
					errors.New("deployment start: at least 1 instance ID must be provided"),
				},
			},
			wantErr: true,
		},
		{
			name: "validate should pass if all params are properly set",
			params: StartInstancesParams{
				StartParams: StartParams{
					ResourceParams: deployment.ResourceParams{
						API:          &api.API{},
						DeploymentID: "f1d329b0fb34470ba8b18361cabdd2bc",
						Kind:         "elasticsearch",
						RefID:        "main-elasticsearch",
					},
				},
				InstanceIDs: []string{"instance-0000000001", "instance-0000000002"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.err == nil {
				t.Errorf("Validate() expected errors = '%v' but no errors returned", tt.err)
			}

			if tt.wantErr && err.Error() != tt.err.Error() {
				t.Errorf("Validate() expected errors = '%v' but got %v", tt.err, err)
			}
		})
	}
}

func TestStart(t *testing.T) {
	var internalError = models.BasicFailedReply{
		Errors: []*models.BasicFailedReplyElement{
			{},
		},
	}
	internalErrorBytes, _ := json.MarshalIndent(internalError, "", "  ")
	var errGet500 = `failed auto-discovering the resource ref id: {
  "errors": [
    {
      "code": "deployment.missing",
      "fields": null,
      "message": null
    }
  ]
}`
	type args struct {
		params StartParams
	}

	tests := []struct {
		name string
		args args
		want models.DeploymentResourceCommandResponse
		err  error
	}{
		{
			name: "fails due to parameter validation",
			args: args{},
			err: &multierror.Error{Errors: []error{
				errors.New("api reference is required for command"),
				errors.New(`id "" is invalid`),
				errors.New("deployment resource kind cannot be empty"),
				errors.New("failed auto-discovering the resource ref id: api reference is required for command"),
				errors.New(`failed auto-discovering the resource ref id: id "" is invalid`),
			}},
		},
		{
			name: "fails due to API error",
			args: args{params: StartParams{
				ResourceParams: deployment.ResourceParams{
					API:          api.NewMock(mock.New404Response(mock.NewStructBody(internalError))),
					DeploymentID: util.ValidClusterID,
					Kind:         "elasticsearch",
					RefID:        "main-elasticsearch",
				},
			}},
			err: errors.New(string(internalErrorBytes)),
		},
		{
			name: "fails due to RefID discovery",
			args: args{params: StartParams{
				ResourceParams: deployment.ResourceParams{
					API: api.NewMock(mock.New500Response(mock.NewStructBody(&models.BasicFailedReply{
						Errors: []*models.BasicFailedReplyElement{
							{Code: ec.String("deployment.missing")},
						},
					}))),
					DeploymentID: util.ValidClusterID,
					Kind:         "elasticsearch",
				},
			}},
			err: &multierror.Error{Errors: []error{
				errors.New(errGet500),
			}},
		},
		{
			name: "succeeds",
			args: args{params: StartParams{
				ResourceParams: deployment.ResourceParams{
					API:          api.NewMock(mock.New200Response(mock.NewStringBody(""))),
					DeploymentID: util.ValidClusterID,
					Kind:         "elasticsearch",
					RefID:        "main-elasticsearch",
				},
			}},
		},
		{
			name: "succeeds when RefID is not set",
			args: args{params: StartParams{
				ResourceParams: deployment.ResourceParams{
					API: api.NewMock(
						mock.New200Response(mock.NewStructBody(models.DeploymentGetResponse{
							Healthy: ec.Bool(true),
							ID:      ec.String("3531aaf988594efa87c1aabb7caed337"),
							Resources: &models.DeploymentResources{
								Elasticsearch: []*models.ElasticsearchResourceInfo{{
									ID:    ec.String("3531aaf988594efa87c1aabb7caed337"),
									RefID: ec.String("elasticsearch"),
								}},
							},
						})),
						mock.New200Response(mock.NewStringBody("")),
					),
					DeploymentID: util.ValidClusterID,
					Kind:         "elasticsearch",
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Start(tt.args.params)
			if tt.err != nil && err.Error() != tt.err.Error() {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Start() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStartInstances(t *testing.T) {
	var internalError = models.BasicFailedReply{
		Errors: []*models.BasicFailedReplyElement{
			{},
		},
	}
	internalErrorBytes, _ := json.MarshalIndent(internalError, "", "  ")
	var errGet500 = `failed auto-discovering the resource ref id: {
  "errors": [
    {
      "code": "deployment.missing",
      "fields": null,
      "message": null
    }
  ]
}`
	type args struct {
		params StartInstancesParams
	}

	tests := []struct {
		name string
		args args
		want models.DeploymentResourceCommandResponse
		err  error
	}{
		{
			name: "fails due to parameter validation",
			args: args{},
			err: &multierror.Error{Errors: []error{
				errors.New("deployment start: at least 1 instance ID must be provided"),
				errors.New("api reference is required for command"),
				errors.New(`id "" is invalid`),
				errors.New("deployment resource kind cannot be empty"),
				errors.New("failed auto-discovering the resource ref id: api reference is required for command"),
				errors.New(`failed auto-discovering the resource ref id: id "" is invalid`),
			}},
		},
		{
			name: "fails due to API error",
			args: args{params: StartInstancesParams{
				StartParams: StartParams{
					ResourceParams: deployment.ResourceParams{
						API:          api.NewMock(mock.New404Response(mock.NewStructBody(internalError))),
						DeploymentID: util.ValidClusterID,
						Kind:         "elasticsearch",
						RefID:        "main-elasticsearch",
					},
				},
				InstanceIDs: []string{"instance-0000000001", "instance-0000000002"},
			}},
			err: errors.New(string(internalErrorBytes)),
		},
		{
			name: "fails due to RefID discovery",
			args: args{params: StartInstancesParams{
				StartParams: StartParams{
					ResourceParams: deployment.ResourceParams{
						API: api.NewMock(mock.New500Response(mock.NewStructBody(&models.BasicFailedReply{
							Errors: []*models.BasicFailedReplyElement{
								{Code: ec.String("deployment.missing")},
							},
						}))),
						DeploymentID: util.ValidClusterID,
						Kind:         "elasticsearch",
					},
				},
				InstanceIDs: []string{"instance-0000000001", "instance-0000000002"},
			}},
			err: &multierror.Error{Errors: []error{
				errors.New(errGet500),
			}},
		},
		{
			name: "succeeds",
			args: args{params: StartInstancesParams{
				StartParams: StartParams{
					ResourceParams: deployment.ResourceParams{
						API:          api.NewMock(mock.New200Response(mock.NewStringBody(""))),
						DeploymentID: util.ValidClusterID,
						Kind:         "elasticsearch",
						RefID:        "main-elasticsearch",
					},
				},
				InstanceIDs: []string{"instance-0000000001", "instance-0000000002"},
			}},
		},
		{
			name: "succeeds when RefID is not set",
			args: args{params: StartInstancesParams{
				StartParams: StartParams{
					ResourceParams: deployment.ResourceParams{
						API: api.NewMock(
							mock.New200Response(mock.NewStructBody(models.DeploymentGetResponse{
								Healthy: ec.Bool(true),
								ID:      ec.String("3531aaf988594efa87c1aabb7caed337"),
								Resources: &models.DeploymentResources{
									Elasticsearch: []*models.ElasticsearchResourceInfo{{
										ID:    ec.String("3531aaf988594efa87c1aabb7caed337"),
										RefID: ec.String("elasticsearch"),
									}},
								},
							})),
							mock.New200Response(mock.NewStringBody("")),
						),
						DeploymentID: util.ValidClusterID,
						Kind:         "elasticsearch",
					},
				},
				InstanceIDs: []string{"instance-0000000001", "instance-0000000002"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StartInstances(tt.args.params)
			if tt.err != nil && err.Error() != tt.err.Error() {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Start() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStartAllOrSpecified(t *testing.T) {
	var internalError = models.BasicFailedReply{
		Errors: []*models.BasicFailedReplyElement{
			{},
		},
	}
	internalErrorBytes, _ := json.MarshalIndent(internalError, "", "  ")
	type args struct {
		params StartInstancesParams
	}

	tests := []struct {
		name string
		args args
		want models.DeploymentResourceCommandResponse
		err  error
	}{
		{
			name: "fails due to API error when all is not set",
			args: args{params: StartInstancesParams{
				StartParams: StartParams{
					ResourceParams: deployment.ResourceParams{
						API:          api.NewMock(mock.New404Response(mock.NewStructBody(internalError))),
						DeploymentID: util.ValidClusterID,
						Kind:         "elasticsearch",
						RefID:        "main-elasticsearch",
					},
				},
				InstanceIDs: []string{"instance-0000000001", "instance-0000000002"},
			}},
			err: errors.New(string(internalErrorBytes)),
		},
		{
			name: "fails due to API error when all is set to true",
			args: args{params: StartInstancesParams{
				StartParams: StartParams{
					ResourceParams: deployment.ResourceParams{
						API:          api.NewMock(mock.New404Response(mock.NewStructBody(internalError))),
						DeploymentID: util.ValidClusterID,
						Kind:         "elasticsearch",
						RefID:        "main-elasticsearch",
					},
					All: true,
				},
			}},
			err: errors.New(string(internalErrorBytes)),
		},
		{
			name: "succeeds when all is not set",
			args: args{params: StartInstancesParams{
				StartParams: StartParams{
					ResourceParams: deployment.ResourceParams{
						API:          api.NewMock(mock.New200Response(mock.NewStringBody(""))),
						DeploymentID: util.ValidClusterID,
						Kind:         "elasticsearch",
						RefID:        "main-elasticsearch",
					},
				},
				InstanceIDs: []string{"instance-0000000001", "instance-0000000002"},
			}},
		},
		{
			name: "succeeds when all is set to true",
			args: args{params: StartInstancesParams{
				StartParams: StartParams{
					ResourceParams: deployment.ResourceParams{
						API:          api.NewMock(mock.New200Response(mock.NewStringBody(""))),
						DeploymentID: util.ValidClusterID,
						Kind:         "elasticsearch",
						RefID:        "main-elasticsearch",
					},
					All: true,
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StartAllOrSpecified(tt.args.params)
			if tt.err != nil && err.Error() != tt.err.Error() {
				t.Errorf("Start() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Start() = %v, want %v", got, tt.want)
			}
		})
	}
}
