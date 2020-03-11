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

func TestStopMaintenanceMode(t *testing.T) {
	type args struct {
		params StopParams
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
			args: args{params: StopParams{
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
			args: args{params: StopParams{
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
			args: args{params: StopParams{
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
			args: args{params: StopParams{
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
			got, err := StopMaintenanceMode(tt.args.params)
			if tt.err != nil && err.Error() != tt.err.Error() {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStopInstancesMaintenanceMode(t *testing.T) {
	type args struct {
		params StopInstancesParams
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
				errors.New("deployment stop: at least 1 instance ID must be provided"),
				errors.New("api reference is required for command"),
				errors.New(`id "" is invalid`),
				errors.New("deployment resource kind cannot be empty"),
				errors.New("failed auto-discovering the resource ref id: api reference is required for command"),
				errors.New(`failed auto-discovering the resource ref id: id "" is invalid`),
			}},
		},
		{
			name: "fails due to API error",
			args: args{params: StopInstancesParams{
				StopParams: StopParams{
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
			args: args{params: StopInstancesParams{
				StopParams: StopParams{
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
			args: args{params: StopInstancesParams{
				StopParams: StopParams{
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
			args: args{params: StopInstancesParams{
				StopParams: StopParams{
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
			got, err := StopInstancesMaintenanceMode(tt.args.params)
			if tt.err != nil && err.Error() != tt.err.Error() {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStopMaintenanceModeAllOrSpecified(t *testing.T) {
	type args struct {
		params StopInstancesParams
	}

	tests := []struct {
		name string
		args args
		want models.DeploymentResourceCommandResponse
		err  error
	}{
		{
			name: "fails due to API error when all is not set",
			args: args{params: StopInstancesParams{
				StopParams: StopParams{
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
			args: args{params: StopInstancesParams{
				StopParams: StopParams{
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
			args: args{params: StopInstancesParams{
				StopParams: StopParams{
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
			args: args{params: StopInstancesParams{
				StopParams: StopParams{
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
			got, err := StopMaintenanceModeAllOrSpecified(tt.args.params)
			if tt.err != nil && err.Error() != tt.err.Error() {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stop() = %v, want %v", got, tt.want)
			}
		})
	}
}
