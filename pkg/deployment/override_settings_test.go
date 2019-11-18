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

package deployment

import (
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/models"
)

func Test_setOverrides(t *testing.T) {
	var eceRegion = "ece-region"
	var overriddenRegion = "overridden-region"
	type args struct {
		req       interface{}
		overrides *PayloadOverrides
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "set name override",
			args: args{
				req: &models.DeploymentCreateRequest{
					Name:      "Some",
					Resources: &models.DeploymentCreateResources{},
				},
				overrides: &PayloadOverrides{Name: "some other"},
			},
			want: &models.DeploymentCreateRequest{
				Name:      "some other",
				Resources: &models.DeploymentCreateResources{},
			},
		},
		{
			name: "set name, version and region override",
			args: args{
				overrides: &PayloadOverrides{
					Name:    "some other",
					Version: "7.4.1",
					Region:  eceRegion,
				},
				req: &models.DeploymentCreateRequest{
					Name: "Some",
					Resources: &models.DeploymentCreateResources{
						Apm: []*models.ApmPayload{
							{
								Plan: &models.ApmPlan{
									Apm: &models.ApmConfiguration{Version: "1.2.3"},
								},
							},
						},
						Appsearch: []*models.AppSearchPayload{
							{
								Plan: &models.AppSearchPlan{
									Appsearch: &models.AppSearchConfiguration{Version: "1.2.3"},
								},
							},
						},
						Elasticsearch: []*models.ElasticsearchPayload{
							{
								Plan: &models.ElasticsearchClusterPlan{
									Elasticsearch: &models.ElasticsearchConfiguration{Version: "1.2.3"},
								},
							},
						},
						Kibana: []*models.KibanaPayload{
							{
								Plan: &models.KibanaClusterPlan{
									Kibana: &models.KibanaConfiguration{Version: "1.2.3"},
								},
							},
						},
					},
				},
			},
			want: &models.DeploymentCreateRequest{
				Name: "some other",
				Resources: &models.DeploymentCreateResources{
					Apm: []*models.ApmPayload{
						{
							Region: &eceRegion,
							Plan: &models.ApmPlan{
								Apm: &models.ApmConfiguration{Version: "7.4.1"},
							},
						},
					},
					Appsearch: []*models.AppSearchPayload{
						{
							Region: &eceRegion,
							Plan: &models.AppSearchPlan{
								Appsearch: &models.AppSearchConfiguration{Version: "7.4.1"},
							},
						},
					},
					Elasticsearch: []*models.ElasticsearchPayload{
						{
							Region: &eceRegion,
							Plan: &models.ElasticsearchClusterPlan{
								Elasticsearch: &models.ElasticsearchConfiguration{Version: "7.4.1"},
							},
						},
					},
					Kibana: []*models.KibanaPayload{
						{
							Region: &eceRegion,
							Plan: &models.KibanaClusterPlan{
								Kibana: &models.KibanaConfiguration{Version: "7.4.1"},
							},
						},
					},
				},
			},
		},
		{
			name: "set region override on a DeploymentUpdateRequest",
			args: args{
				overrides: &PayloadOverrides{
					Region: "overridden-region",
				},
				req: &models.DeploymentUpdateRequest{
					Resources: &models.DeploymentUpdateResources{
						Apm: []*models.ApmPayload{
							{
								Plan: &models.ApmPlan{
									Apm: &models.ApmConfiguration{Version: "7.4.1"},
								},
							},
						},
						Appsearch: []*models.AppSearchPayload{
							{
								Plan: &models.AppSearchPlan{
									Appsearch: &models.AppSearchConfiguration{Version: "7.4.1"},
								},
							},
						},
						Elasticsearch: []*models.ElasticsearchPayload{
							{
								Plan: &models.ElasticsearchClusterPlan{
									Elasticsearch: &models.ElasticsearchConfiguration{Version: "7.4.1"},
								},
							},
						},
						Kibana: []*models.KibanaPayload{
							{
								Plan: &models.KibanaClusterPlan{
									Kibana: &models.KibanaConfiguration{Version: "7.4.1"},
								},
							},
						},
					},
				},
			},
			want: &models.DeploymentUpdateRequest{
				Resources: &models.DeploymentUpdateResources{
					Apm: []*models.ApmPayload{
						{
							Region: &overriddenRegion,
							Plan: &models.ApmPlan{
								Apm: &models.ApmConfiguration{Version: "7.4.1"},
							},
						},
					},
					Appsearch: []*models.AppSearchPayload{
						{
							Region: &overriddenRegion,
							Plan: &models.AppSearchPlan{
								Appsearch: &models.AppSearchConfiguration{Version: "7.4.1"},
							},
						},
					},
					Elasticsearch: []*models.ElasticsearchPayload{
						{
							Region: &overriddenRegion,
							Plan: &models.ElasticsearchClusterPlan{
								Elasticsearch: &models.ElasticsearchConfiguration{Version: "7.4.1"},
							},
						},
					},
					Kibana: []*models.KibanaPayload{
						{
							Region: &overriddenRegion,
							Plan: &models.KibanaClusterPlan{
								Kibana: &models.KibanaConfiguration{Version: "7.4.1"},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req = tt.args.req
			setOverrides(req, tt.args.overrides)

			if !reflect.DeepEqual(req, tt.want) {
				t.Errorf("setOverrides() = %v, want %v", req, tt.want)
			}
		})
	}
}
