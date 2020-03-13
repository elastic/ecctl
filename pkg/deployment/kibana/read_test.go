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
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api/mock"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"
)

func TestListParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  ListParams
		wantErr bool
		err     error
	}{
		{
			name:   "validate should return all possible errors",
			params: ListParams{},
			err: &multierror.Error{
				Errors: []error{
					errors.New("api reference is required for command"),
				},
			},
			wantErr: true,
		},
		{
			name: "validate should pass if all params are properly set",
			params: ListParams{
				API:     &api.API{},
				Version: "6.2.2",
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

func TestGet(t *testing.T) {
	apiClusterResponse := `{
  "cluster_id": "2c221bd86b7f48959a59ee3128d5c5e8",
  "cluster_name": "soteria",
  "elasticsearch_cluster": {
    "elasticsearch_id": "5c641576747442eba0ebd67944ccbe10"
  },
  "plan_info": {
    "healthy": true
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
        }
      }
    ]
  }
}`
	tests := []struct {
		name     string
		params   ClusterParams
		wantErr  bool
		err      error
		expected models.KibanaClusterInfo
	}{
		{
			name:   "should fail if params validation fails",
			params: ClusterParams{},
			err: &multierror.Error{
				Errors: []error{
					errors.New(`api reference is required for command`),
					errors.New(`id "" is invalid`),
				},
			},
			wantErr: true,
		},
		{
			name: "should succeed if parameters are properly set and API call is successful",
			params: ClusterParams{
				DeploymentParams: DeploymentParams{
					ID: "5c641576747442eba0ebd67944ccbe10",
					API: api.NewMock(mock.Response{Response: http.Response{
						StatusCode: 200,
						Body:       mock.NewStringBody(apiClusterResponse),
					}}),
				},
			},
			wantErr: false,
			expected: models.KibanaClusterInfo{
				ClusterID:   ec.String("2c221bd86b7f48959a59ee3128d5c5e8"),
				ClusterName: ec.String("soteria"),
				ElasticsearchCluster: &models.TargetElasticsearchCluster{
					ElasticsearchID: ec.String("5c641576747442eba0ebd67944ccbe10"),
				},
				PlanInfo: &models.KibanaClusterPlansInfo{
					Healthy: ec.Bool(true),
				},
				Region: "us-east-1",
				Status: ec.String("started"),
				Topology: &models.ClusterTopologyInfo{
					Healthy: ec.Bool(true),
					Instances: []*models.ClusterInstanceInfo{
						{
							AllocatorID:      "i-01c866ac29bf57d4d",
							ContainerStarted: ec.Bool(true),
							Healthy:          ec.Bool(true),
							InstanceConfiguration: &models.ClusterInstanceConfigurationInfo{
								ID:       ec.String("aws.kibana.r4"),
								Name:     ec.String("aws.kibana.r4"),
								Resource: ec.String("memory"),
							},
							InstanceName:    ec.String("instance-0000000002"),
							MaintenanceMode: ec.Bool(false),
							Memory:          &models.ClusterInstanceMemoryInfo{InstanceCapacity: ec.Int32(1024)},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := Get(tt.params)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.err == nil {
				t.Errorf("Get() expected errors = '%v' but no errors returned", tt.err)
			}

			if tt.wantErr && !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Get() expected errors = '%v' but got %v", tt.err, err)
			}

			if !tt.wantErr && !reflect.DeepEqual(tt.expected, *resp) {
				t.Errorf("Get() expected response = '%+v' but got %+v", tt.expected, resp)
			}
		})
	}
}

func TestList(t *testing.T) {
	tests := []struct {
		name    string
		params  ListParams
		want    *models.KibanaClustersInfo
		err     error
		wantErr bool
	}{
		{
			name:   "should fail if params validation fails",
			params: ListParams{},
			err: &multierror.Error{
				Errors: []error{
					errors.New("api reference is required for command"),
				},
			},
			wantErr: true,
		},
		{
			name: "should return error on API error",
			params: ListParams{
				API: api.NewMock(mock.Response{
					Error: errors.New("an error"),
				}),
			},
			err: &url.Error{
				Op:  "Get",
				URL: "https://mock-host/mock-path/clusters/kibana?from=0&show_hidden=false&show_metadata=false&show_plan_defaults=false&show_plans=true&show_settings=false&size=0",
				Err: errors.New("an error"),
			},
			wantErr: true,
		},
		{
			name: "list succeeds",
			params: ListParams{
				API: api.NewMock(mock.Response{Response: http.Response{
					StatusCode: 200,
					Body: mock.NewStructBody(models.KibanaClustersInfo{
						ReturnCount: ec.Int32(2),
						KibanaClusters: []*models.KibanaClusterInfo{
							{
								ClusterID: ec.String("86d2ec6217774eedb93ba38483141997"),
								ElasticsearchCluster: &models.TargetElasticsearchCluster{
									ElasticsearchID: ec.String("d324608c97154bdba2dff97511d40368"),
								},
								Status: ec.String("started"),
							},
							{
								ClusterID: ec.String("d324608c97154bdba2dff97511d40368"),
								ElasticsearchCluster: &models.TargetElasticsearchCluster{
									ElasticsearchID: ec.String("86d2ec6217774eedb93ba38483141997"),
								},
								Status: ec.String("stopped"),
							},
						},
					}),
				}}),
			},
			want: &models.KibanaClustersInfo{
				ReturnCount: ec.Int32(2),
				KibanaClusters: []*models.KibanaClusterInfo{
					{
						ClusterID: ec.String("86d2ec6217774eedb93ba38483141997"),
						ElasticsearchCluster: &models.TargetElasticsearchCluster{
							ElasticsearchID: ec.String("d324608c97154bdba2dff97511d40368"),
						},
						Status: ec.String("started"),
					},
					{
						ClusterID: ec.String("d324608c97154bdba2dff97511d40368"),
						ElasticsearchCluster: &models.TargetElasticsearchCluster{
							ElasticsearchID: ec.String("86d2ec6217774eedb93ba38483141997"),
						},
						Status: ec.String("stopped"),
					},
				},
			},
		},
		{
			name: "list succeeds when a version is supplied",
			params: ListParams{
				Version: "6.2.2",
				API: api.NewMock(mock.Response{Response: http.Response{
					StatusCode: 200,
					Body: mock.NewStructBody(models.KibanaClustersInfo{
						ReturnCount: ec.Int32(2),
						KibanaClusters: []*models.KibanaClusterInfo{
							{
								ClusterID: ec.String("86d2ec6217774eedb93ba38483141997"),
								PlanInfo: &models.KibanaClusterPlansInfo{
									Current: &models.KibanaClusterPlanInfo{
										Plan: &models.KibanaClusterPlan{
											Kibana: &models.KibanaConfiguration{
												Version: "7.0.0",
											},
										},
									},
								},
							},
							{
								ClusterID: ec.String("d324608c97154bdba2dff97511d40368"),
								PlanInfo: &models.KibanaClusterPlansInfo{
									Current: &models.KibanaClusterPlanInfo{
										Plan: &models.KibanaClusterPlan{
											Kibana: &models.KibanaConfiguration{
												Version: "6.2.2",
											},
										},
									},
								},
							},
						},
					}),
				}}),
			},
			want: &models.KibanaClustersInfo{
				ReturnCount: ec.Int32(1),
				KibanaClusters: []*models.KibanaClusterInfo{
					{
						ClusterID: ec.String("d324608c97154bdba2dff97511d40368"),
						PlanInfo: &models.KibanaClusterPlansInfo{
							Current: &models.KibanaClusterPlanInfo{
								Plan: &models.KibanaClusterPlan{
									Kibana: &models.KibanaConfiguration{
										Version: "6.2.2",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := List(tt.params)
			if !reflect.DeepEqual(tt.err, err) {
				t.Errorf("List() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterByVersion(t *testing.T) {
	type args struct {
		payload *models.KibanaClustersInfo
		version string
	}
	tests := []struct {
		name string
		args args
		want *models.KibanaClustersInfo
	}{
		{
			name: "output should match input if no version is supplied",
			args: args{
				payload: &models.KibanaClustersInfo{
					ReturnCount: ec.Int32(2),
					KibanaClusters: []*models.KibanaClusterInfo{
						{
							ClusterID: ec.String("86d2ec6217774eedb93ba38483141997"),
							PlanInfo: &models.KibanaClusterPlansInfo{
								Current: &models.KibanaClusterPlanInfo{
									Plan: &models.KibanaClusterPlan{
										Kibana: &models.KibanaConfiguration{
											Version: "7.0.0",
										},
									},
								},
							},
						},
						{
							ClusterID: ec.String("d324608c97154bdba2dff97511d40368"),
							PlanInfo: &models.KibanaClusterPlansInfo{
								Current: &models.KibanaClusterPlanInfo{
									Plan: &models.KibanaClusterPlan{
										Kibana: &models.KibanaConfiguration{
											Version: "6.2.2",
										},
									},
								},
							},
						},
					},
				},
			},
			want: &models.KibanaClustersInfo{
				ReturnCount: ec.Int32(2),
				KibanaClusters: []*models.KibanaClusterInfo{
					{
						ClusterID: ec.String("86d2ec6217774eedb93ba38483141997"),
						PlanInfo: &models.KibanaClusterPlansInfo{
							Current: &models.KibanaClusterPlanInfo{
								Plan: &models.KibanaClusterPlan{
									Kibana: &models.KibanaConfiguration{
										Version: "7.0.0",
									},
								},
							},
						},
					},
					{
						ClusterID: ec.String("d324608c97154bdba2dff97511d40368"),
						PlanInfo: &models.KibanaClusterPlansInfo{
							Current: &models.KibanaClusterPlanInfo{
								Plan: &models.KibanaClusterPlan{
									Kibana: &models.KibanaConfiguration{
										Version: "6.2.2",
									},
								},
							},
						},
					},
				}},
		},
		{
			name: "should only return version filtered results",
			args: args{
				payload: &models.KibanaClustersInfo{
					ReturnCount: ec.Int32(2),
					KibanaClusters: []*models.KibanaClusterInfo{
						{
							ClusterID: ec.String("86d2ec6217774eedb93ba38483141997"),
							PlanInfo: &models.KibanaClusterPlansInfo{
								Current: &models.KibanaClusterPlanInfo{
									Plan: &models.KibanaClusterPlan{
										Kibana: &models.KibanaConfiguration{
											Version: "7.0.0",
										},
									},
								},
							},
						},
						{
							ClusterID: ec.String("d324608c97154bdba2dff97511d40368"),
							PlanInfo: &models.KibanaClusterPlansInfo{
								Current: &models.KibanaClusterPlanInfo{
									Plan: &models.KibanaClusterPlan{
										Kibana: &models.KibanaConfiguration{
											Version: "6.2.2",
										},
									},
								},
							},
						},
					},
				},
				version: "6.2.2",
			},
			want: &models.KibanaClustersInfo{
				ReturnCount: ec.Int32(1),
				KibanaClusters: []*models.KibanaClusterInfo{
					{
						ClusterID: ec.String("d324608c97154bdba2dff97511d40368"),
						PlanInfo: &models.KibanaClusterPlansInfo{
							Current: &models.KibanaClusterPlanInfo{
								Plan: &models.KibanaClusterPlan{
									Kibana: &models.KibanaConfiguration{
										Version: "6.2.2",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterByVersion(tt.args.payload, tt.args.version)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterByVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
