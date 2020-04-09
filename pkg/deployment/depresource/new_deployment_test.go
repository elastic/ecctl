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
)

var apmKibanaTemplateResponse = models.DeploymentTemplateInfo{
	ID: "default",
	ClusterTemplate: &models.DeploymentTemplateDefinitionRequest{
		Apm: &models.CreateApmInCreateElasticsearchRequest{
			Plan: &models.ApmPlan{
				ClusterTopology: []*models.ApmTopologyElement{
					{
						Size: &models.TopologySize{
							Resource: ec.String("memory"),
							Value:    ec.Int32(1024),
						},
						ZoneCount: 1,
					},
				},
			},
		},
		Kibana: &models.CreateKibanaInCreateElasticsearchRequest{
			Plan: &models.KibanaClusterPlan{
				ClusterTopology: []*models.KibanaClusterTopologyElement{
					{
						Size: &models.TopologySize{
							Resource: ec.String("memory"),
							Value:    ec.Int32(1024),
						},
						ZoneCount: 1,
					},
				},
			},
		},
		Plan: &models.ElasticsearchClusterPlan{
			ClusterTopology: defaultESTopologies,
		},
	},
}

var appsearchKibanaTemplateResponse = models.DeploymentTemplateInfo{
	ID: "default",
	ClusterTemplate: &models.DeploymentTemplateDefinitionRequest{
		Appsearch: &models.CreateAppSearchRequest{
			Plan: &models.AppSearchPlan{
				ClusterTopology: []*models.AppSearchTopologyElement{
					{
						Size: &models.TopologySize{
							Resource: ec.String("memory"),
							Value:    ec.Int32(1024),
						},
						ZoneCount: 1,
					},
				},
			},
		},
		Kibana: &models.CreateKibanaInCreateElasticsearchRequest{
			Plan: &models.KibanaClusterPlan{
				ClusterTopology: []*models.KibanaClusterTopologyElement{
					{
						Size: &models.TopologySize{
							Resource: ec.String("memory"),
							Value:    ec.Int32(1024),
						},
						ZoneCount: 1,
					},
				},
			},
		},
		Plan: &models.ElasticsearchClusterPlan{
			ClusterTopology: defaultESTopologies,
		},
	},
}

func TestNew(t *testing.T) {
	type args struct {
		params NewParams
	}
	tests := []struct {
		name    string
		args    args
		want    *models.DeploymentCreateRequest
		wantErr error
	}{
		{
			name: "Fails due to API error",
			args: args{params: NewParams{
				Version: "7.6.1",
				Region:  "ece-region",
				ElasticsearchInstance: InstanceParams{
					RefID:     "main-elasticsearch",
					Size:      1024,
					ZoneCount: 1,
				},
				KibanaInstance: InstanceParams{
					RefID:     "main-kibana",
					Size:      1024,
					ZoneCount: 1,
				},
				DeploymentTemplateID: "default",
				API: api.NewMock(mock.New500Response(
					mock.NewStringBody("error"),
				)),
			}},
			wantErr: errors.New("error"),
		},
		{
			name: "Fails to create a deployment payload with ES and Kibana instances",
			args: args{params: NewParams{
				Version: "7.6.1",
				Region:  "ece-region",
				ElasticsearchInstance: InstanceParams{
					RefID:     "main-elasticsearch",
					Size:      1024,
					ZoneCount: 1,
				},
				KibanaInstance: InstanceParams{
					RefID:     "main-kibana",
					Size:      1024,
					ZoneCount: 1,
				},
				DeploymentTemplateID: "default",
				API: api.NewMock(
					mock.New200Response(mock.NewStructBody(defaultTemplateResponse)),
					mock.New200Response(mock.NewStructBody(defaultTemplateResponse)),
				),
			}},
			wantErr: errors.New("deployment: the default template is not configured for Kibana. Please use another template if you wish to start Kibana instances"),
		},
		{
			name: "Fails to create a deployment payload with ES, Kibana and APM instances",
			args: args{params: NewParams{
				Version: "7.6.1",
				Region:  "ece-region",
				ElasticsearchInstance: InstanceParams{
					RefID:     "main-elasticsearch",
					Size:      1024,
					ZoneCount: 1,
				},
				KibanaInstance: InstanceParams{
					RefID:     "main-kibana",
					Size:      1024,
					ZoneCount: 1,
				},
				ApmInstance: InstanceParams{
					RefID:     "main-apm",
					Size:      1024,
					ZoneCount: 1,
				},
				DeploymentTemplateID: "default",
				ApmEnable:            true,
				API: api.NewMock(
					mock.New200Response(mock.NewStructBody(appsearchKibanaTemplateResponse)),
					mock.New200Response(mock.NewStructBody(appsearchKibanaTemplateResponse)),
					mock.New200Response(mock.NewStructBody(appsearchKibanaTemplateResponse)),
				),
			}},
			wantErr: errors.New("deployment: the default template is not configured for APM. Please use another template if you wish to start APM instances"),
		},
		{
			name: "Fails to create a deployment payload with ES, Kibana and App Search instances",
			args: args{params: NewParams{
				Version: "7.6.1",
				Region:  "ece-region",
				ElasticsearchInstance: InstanceParams{
					RefID:     "main-elasticsearch",
					Size:      1024,
					ZoneCount: 1,
				},
				KibanaInstance: InstanceParams{
					RefID:     "main-kibana",
					Size:      1024,
					ZoneCount: 1,
				},
				AppsearchInstance: InstanceParams{
					RefID:     "main-appsearch",
					Size:      1024,
					ZoneCount: 1,
				},
				DeploymentTemplateID: "default",
				AppsearchEnable:      true,
				API: api.NewMock(
					mock.New200Response(mock.NewStructBody(apmKibanaTemplateResponse)),
					mock.New200Response(mock.NewStructBody(apmKibanaTemplateResponse)),
					mock.New200Response(mock.NewStructBody(apmKibanaTemplateResponse)),
				),
			}},
			wantErr: errors.New("deployment: the default template is not configured for App Search. Please use another template if you wish to start App Search instances"),
		},
		{
			name: "Succeeds to create a deployment payload with ES and Kibana instances",
			args: args{params: NewParams{
				Version: "7.6.1",
				Region:  "ece-region",
				ElasticsearchInstance: InstanceParams{
					RefID:     "main-elasticsearch",
					Size:      1024,
					ZoneCount: 1,
				},
				KibanaInstance: InstanceParams{
					RefID:     "main-kibana",
					Size:      1024,
					ZoneCount: 1,
				},
				DeploymentTemplateID: "default",
				API: api.NewMock(
					mock.New200Response(mock.NewStructBody(kibanaTemplateResponse)),
					mock.New200Response(mock.NewStructBody(kibanaTemplateResponse)),
				),
			}},
			want: &models.DeploymentCreateRequest{Resources: &models.DeploymentCreateResources{
				Elasticsearch: []*models.ElasticsearchPayload{{
					RefID:  ec.String("main-elasticsearch"),
					Region: ec.String("ece-region"),
					Plan: &models.ElasticsearchClusterPlan{
						Elasticsearch: &models.ElasticsearchConfiguration{
							Version: "7.6.1",
						},
						DeploymentTemplate: &models.DeploymentTemplateReference{
							ID: ec.String("default"),
						},
						ClusterTopology: []*models.ElasticsearchClusterTopologyElement{{
							ZoneCount:               1,
							InstanceConfigurationID: "default.data",
							Size: &models.TopologySize{
								Resource: ec.String("memory"),
								Value:    ec.Int32(1024),
							},
							NodeType: &models.ElasticsearchNodeType{
								Data: ec.Bool(true),
							},
						}},
					}},
				},
				Kibana: []*models.KibanaPayload{{
					ElasticsearchClusterRefID: ec.String("main-elasticsearch"),
					Region:                    ec.String("ece-region"),
					RefID:                     ec.String("main-kibana"),
					Plan: &models.KibanaClusterPlan{
						Kibana: &models.KibanaConfiguration{
							Version: "7.6.1",
						},
						ClusterTopology: []*models.KibanaClusterTopologyElement{
							{
								Size: &models.TopologySize{
									Resource: ec.String("memory"),
									Value:    ec.Int32(1024),
								},
								ZoneCount: 1,
							},
						},
					},
				}},
			}},
		},
		{
			name: "Succeeds to create a deployment payload with ES, Kibana and APM instances",
			args: args{params: NewParams{
				Version: "7.6.1",
				Region:  "ece-region",
				ElasticsearchInstance: InstanceParams{
					RefID:     "main-elasticsearch",
					Size:      1024,
					ZoneCount: 1,
				},
				KibanaInstance: InstanceParams{
					RefID:     "main-kibana",
					Size:      1024,
					ZoneCount: 1,
				},
				ApmInstance: InstanceParams{
					RefID:     "main-apm",
					Size:      1024,
					ZoneCount: 1,
				},
				DeploymentTemplateID: "default",
				ApmEnable:            true,
				API: api.NewMock(
					mock.New200Response(mock.NewStructBody(apmKibanaTemplateResponse)),
					mock.New200Response(mock.NewStructBody(apmKibanaTemplateResponse)),
					mock.New200Response(mock.NewStructBody(apmKibanaTemplateResponse)),
				),
			}},
			want: &models.DeploymentCreateRequest{Resources: &models.DeploymentCreateResources{
				Elasticsearch: []*models.ElasticsearchPayload{{
					RefID:  ec.String("main-elasticsearch"),
					Region: ec.String("ece-region"),
					Plan: &models.ElasticsearchClusterPlan{
						Elasticsearch: &models.ElasticsearchConfiguration{
							Version: "7.6.1",
						},
						DeploymentTemplate: &models.DeploymentTemplateReference{
							ID: ec.String("default"),
						},
						ClusterTopology: []*models.ElasticsearchClusterTopologyElement{{
							ZoneCount:               1,
							InstanceConfigurationID: "default.data",
							Size: &models.TopologySize{
								Resource: ec.String("memory"),
								Value:    ec.Int32(1024),
							},
							NodeType: &models.ElasticsearchNodeType{
								Data: ec.Bool(true),
							},
						}},
					}},
				},
				Kibana: []*models.KibanaPayload{{
					ElasticsearchClusterRefID: ec.String("main-elasticsearch"),
					Region:                    ec.String("ece-region"),
					RefID:                     ec.String("main-kibana"),
					Plan: &models.KibanaClusterPlan{
						Kibana: &models.KibanaConfiguration{
							Version: "7.6.1",
						},
						ClusterTopology: []*models.KibanaClusterTopologyElement{
							{
								Size: &models.TopologySize{
									Resource: ec.String("memory"),
									Value:    ec.Int32(1024),
								},
								ZoneCount: 1,
							},
						},
					},
				}},
				Apm: []*models.ApmPayload{{
					ElasticsearchClusterRefID: ec.String("main-elasticsearch"),
					Region:                    ec.String("ece-region"),
					RefID:                     ec.String("main-apm"),
					Plan: &models.ApmPlan{
						Apm: &models.ApmConfiguration{
							Version: "7.6.1",
						},
						ClusterTopology: []*models.ApmTopologyElement{
							{
								Size: &models.TopologySize{
									Resource: ec.String("memory"),
									Value:    ec.Int32(1024),
								},
								ZoneCount: 1,
							},
						},
					},
				}},
			}},
		},
		{
			name: "Succeeds to create a deployment payload with ES, Kibana and Appsearch instances",
			args: args{params: NewParams{
				Version: "7.6.1",
				Region:  "ece-region",
				ElasticsearchInstance: InstanceParams{
					RefID:     "main-elasticsearch",
					Size:      1024,
					ZoneCount: 1,
				},
				KibanaInstance: InstanceParams{
					RefID:     "main-kibana",
					Size:      1024,
					ZoneCount: 1,
				},
				AppsearchInstance: InstanceParams{
					RefID:     "main-appsearch",
					Size:      1024,
					ZoneCount: 1,
				},
				DeploymentTemplateID: "default",
				AppsearchEnable:      true,
				API: api.NewMock(
					mock.New200Response(mock.NewStructBody(appsearchKibanaTemplateResponse)),
					mock.New200Response(mock.NewStructBody(appsearchKibanaTemplateResponse)),
					mock.New200Response(mock.NewStructBody(appsearchKibanaTemplateResponse)),
				),
			}},
			want: &models.DeploymentCreateRequest{Resources: &models.DeploymentCreateResources{
				Elasticsearch: []*models.ElasticsearchPayload{{
					RefID:  ec.String("main-elasticsearch"),
					Region: ec.String("ece-region"),
					Plan: &models.ElasticsearchClusterPlan{
						Elasticsearch: &models.ElasticsearchConfiguration{
							Version: "7.6.1",
						},
						DeploymentTemplate: &models.DeploymentTemplateReference{
							ID: ec.String("default"),
						},
						ClusterTopology: []*models.ElasticsearchClusterTopologyElement{{
							ZoneCount:               1,
							InstanceConfigurationID: "default.data",
							Size: &models.TopologySize{
								Resource: ec.String("memory"),
								Value:    ec.Int32(1024),
							},
							NodeType: &models.ElasticsearchNodeType{
								Data: ec.Bool(true),
							},
						}},
					}},
				},
				Kibana: []*models.KibanaPayload{{
					ElasticsearchClusterRefID: ec.String("main-elasticsearch"),
					Region:                    ec.String("ece-region"),
					RefID:                     ec.String("main-kibana"),
					Plan: &models.KibanaClusterPlan{
						Kibana: &models.KibanaConfiguration{
							Version: "7.6.1",
						},
						ClusterTopology: []*models.KibanaClusterTopologyElement{
							{
								Size: &models.TopologySize{
									Resource: ec.String("memory"),
									Value:    ec.Int32(1024),
								},
								ZoneCount: 1,
							},
						},
					},
				}},
				Appsearch: []*models.AppSearchPayload{{
					ElasticsearchClusterRefID: ec.String("main-elasticsearch"),
					Region:                    ec.String("ece-region"),
					RefID:                     ec.String("main-appsearch"),
					Plan: &models.AppSearchPlan{
						Appsearch: &models.AppSearchConfiguration{
							Version: "7.6.1",
						},
						ClusterTopology: []*models.AppSearchTopologyElement{
							{
								Size: &models.TopologySize{
									Resource: ec.String("memory"),
									Value:    ec.Int32(1024),
								},
								ZoneCount: 1,
							},
						},
					},
				}},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.params)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
