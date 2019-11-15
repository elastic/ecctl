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

package kibana

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
	"github.com/elastic/cloud-sdk-go/pkg/output"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func TestReallocateParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  ReallocateParams
		wantErr bool
		err     error
	}{
		{
			name:   "validate should return all possible errors",
			params: ReallocateParams{},
			err: &multierror.Error{
				Errors: []error{
					errors.New(`api reference is required for command`),
					errors.New(`id "" is invalid`),
					errors.New(`track: Output cannot be empty`),
				},
			},
			wantErr: true,
		},
		{
			name: "validate should pass if all params are properly set",
			params: ReallocateParams{
				DeploymentParams: DeploymentParams{
					ID:  "5c641576747442eba0ebd67944ccbe10",
					API: &api.API{},
				},
				Output: new(output.Device),
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

			if tt.wantErr && !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Validate() expected errors = '%v' but got %v", tt.err, err)
			}
		})
	}
}

func TestReallocate(t *testing.T) {
	resp := `{
  "cluster_id": "2c221bd86b7f48959a59ee3128d5c5e8",
  "cluster_name": "soteria",
  "elasticsearch_cluster": {
    "elasticsearch_id": "5c641576747442eba0ebd67944ccbe10"
  },
  "external_links": [],
  "healthy": true,
  "metadata": {
    "endpoint": "2c221bd86b7f48959a59ee3128d5c5e8.us-east-1.aws.staging.foundit.no",
    "last_modified": "2018-10-29T12:24:10.210Z",
    "version": 12
  },
  "plan_info": {
    "healthy": true,
    "history": []
  },
  "region": "us-east-1",
  "status": "started",
  "topology": {
    "healthy": true,
    "instances": [
      {
        "allocator_id": "i-01c866ac29bf57d4d",
        "container_started": true,
        "healthy": true,
        "instance_configuration": {
          "id": "aws.kibana.r4",
          "name": "aws.kibana.r4",
          "resource": "memory"
        },
        "instance_name": "instance-0000000002",
        "maintenance_mode": false,
        "memory": {
          "instance_capacity": 1024
        },
        "service_roles": null,
        "service_running": true,
        "service_version": "6.4.1",
        "zone": "us-east-1a"
      }
    ]
  }
}`
	tests := []struct {
		name    string
		params  ReallocateParams
		wantErr bool
		err     error
	}{
		{
			name:   "should fail if params validation fails",
			params: ReallocateParams{},
			err: &multierror.Error{
				Errors: []error{
					errors.New(`api reference is required for command`),
					errors.New(`id "" is invalid`),
					errors.New(`track: Output cannot be empty`),
				},
			},
			wantErr: true,
		},
		{
			name: "should fail if get kibana instance fails",
			params: ReallocateParams{
				DeploymentParams: DeploymentParams{
					ID: "853cfe89c4a74fb6a6477574d3c03771",
					API: api.NewMock(mock.Response{
						Error: errors.New("kibana not found"),
					}),
				},
				Output: new(output.Device),
			},
			wantErr: true,
			err: &url.Error{
				Op:  "Get",
				URL: "https://mock-host/mock-path/clusters/kibana/853cfe89c4a74fb6a6477574d3c03771?convert_legacy_plans=false&show_metadata=false&show_plan_defaults=false&show_plan_logs=false&show_plans=false&show_settings=false",
				Err: errors.New("kibana not found"),
			},
		},
		{
			name: "should fail if get reallocate fails",
			params: ReallocateParams{
				DeploymentParams: DeploymentParams{
					ID: "853cfe89c4a74fb6a6477574d3c03771",
					API: api.NewMock(
						mock.Response{Response: http.Response{
							StatusCode: 200,
							Body:       mock.NewStringBody(resp),
						}},
						mock.Response{
							Error: errors.New("kibana not found"),
						},
					),
				},
				Output: new(output.Device),
			},
			wantErr: true,
			err: &url.Error{
				Op:  "Post",
				URL: "https://mock-host/mock-path/clusters/kibana/853cfe89c4a74fb6a6477574d3c03771/instances/instance-0000000002/_move?force_update=false&ignore_missing=false&validate_only=false",
				Err: errors.New("kibana not found"),
			},
		},
		{
			name: "tracking is enabled",
			params: ReallocateParams{
				DeploymentParams: DeploymentParams{
					ID: "853cfe89c4a74fb6a6477574d3c03771",
					TrackParams: util.TrackParams{
						Track:         true,
						Output:        output.NewDevice(new(bytes.Buffer)),
						PollFrequency: time.Millisecond,
						MaxRetries:    1,
					},
					API: api.NewMock(util.AppendTrackResponses(
						mock.Response{Response: http.Response{
							StatusCode: 200,
							Body:       mock.NewStringBody(resp),
						}},
						mock.Response{
							Error: errors.New("kibana not found"),
						},
					)...),
				},
				Output: new(output.Device),
			},
			wantErr: true,
			err: &url.Error{
				Op:  "Post",
				URL: "https://mock-host/mock-path/clusters/kibana/853cfe89c4a74fb6a6477574d3c03771/instances/instance-0000000002/_move?force_update=false&ignore_missing=false&validate_only=false",
				Err: errors.New("kibana not found"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Reallocate(tt.params)

			if (err != nil) != tt.wantErr {
				t.Errorf("Reallocate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.err == nil {
				t.Errorf("Reallocate() expected errors = '%v' but no errors returned", tt.err)
			}

			if tt.wantErr && !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Reallocate() expected errors = '%v' but got %v", tt.err, err)
			}
		})
	}
}
