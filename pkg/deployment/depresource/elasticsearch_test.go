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
	"fmt"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/hashicorp/go-multierror"
)

var defaultESTopologies = []*models.ElasticsearchClusterTopologyElement{
	{
		InstanceConfigurationID: "default.data",
		Size: &models.TopologySize{
			Resource: ec.String("memory"),
			Value:    ec.Int32(1024),
		},
		NodeType: &models.ElasticsearchNodeType{
			Data: ec.Bool(true),
		},
	},
	{
		InstanceConfigurationID: "default.master",
		Size: &models.TopologySize{
			Resource: ec.String("memory"),
			Value:    ec.Int32(1024),
		},
		NodeType: &models.ElasticsearchNodeType{
			Master: ec.Bool(true),
		},
	},
	{
		InstanceConfigurationID: "default.ml",
		Size: &models.TopologySize{
			Resource: ec.String("memory"),
			Value:    ec.Int32(1024),
		},
		NodeType: &models.ElasticsearchNodeType{
			Ml: ec.Bool(true),
		},
	},
}

var elasticsearchTemplateResponse = models.DeploymentTemplateInfo{
	ID: "default",
	ClusterTemplate: &models.DeploymentTemplateDefinitionRequest{
		Plan: &models.ElasticsearchClusterPlan{
			ClusterTopology: defaultESTopologies,
		},
	},
}

func TestNewElasticsearch(t *testing.T) {
	type args struct {
		params NewElasticsearchParams
	}
	tests := []struct {
		name string
		args args
		want *models.ElasticsearchPayload
		err  error
	}{
		{
			name: "fails due to parameter validation",
			args: args{params: NewElasticsearchParams{
				Topology: []ElasticsearchTopologyElement{
					{},
				},
			}},
			err: &multierror.Error{Errors: []error{
				errors.New("api reference is required for command"),
				errors.New("deployment topology: region cannot be empty"),
				errors.New("deployment topology: version cannot be empty"),
				fmt.Errorf("topology element [0]: %s", &multierror.Error{Errors: []error{
					errors.New("deployment topology: name cannot be empty"),
					errors.New("deployment topology: size cannot be empty"),
				}}),
			}},
		},
		{
			name: "fails due to API error",
			args: args{params: NewElasticsearchParams{
				API:     api.NewMock(mock.New500Response(mock.NewStringBody("error"))),
				Region:  "ece-region",
				Version: "7.4.2",
			}},
			err: errors.New("error"),
		},
		{
			name: "fails due to unknown desired topology",
			args: args{params: NewElasticsearchParams{
				API:        api.NewMock(mock.New200Response(mock.NewStructBody(elasticsearchTemplateResponse))),
				Region:     "ece-region",
				Version:    "7.4.2",
				TemplateID: "default",
				Topology: []ElasticsearchTopologyElement{
					{Name: "some", Size: 1024},
				},
			}},
			err: errors.New(`deployment topology: failed to obtain desired topology names ([{Name:some ZoneCount:0 Size:1024}]) in deployment template id "default"`),
		},
		{
			name: "Returns the default topology",
			args: args{params: NewElasticsearchParams{
				API:        api.NewMock(mock.New200Response(mock.NewStructBody(elasticsearchTemplateResponse))),
				Region:     "ece-region",
				Version:    "7.4.2",
				TemplateID: "default",
			}},
			want: &models.ElasticsearchPayload{
				DisplayName: "",
				Region:      ec.String("ece-region"),
				RefID:       ec.String(DefaultElasticsearchRefID),
				Plan: &models.ElasticsearchClusterPlan{
					Elasticsearch: &models.ElasticsearchConfiguration{
						Version: "7.4.2",
					},
					DeploymentTemplate: &models.DeploymentTemplateReference{
						ID: ec.String(DefaultTemplateID),
					},
					ClusterTopology: []*models.ElasticsearchClusterTopologyElement{
						{
							ZoneCount:               1,
							InstanceConfigurationID: "default.data",
							Size: &models.TopologySize{
								Resource: ec.String("memory"),
								Value:    ec.Int32(4096),
							},
							NodeType: &models.ElasticsearchNodeType{
								Data: ec.Bool(true),
							},
						},
					},
				},
			},
		},
		{
			name: "Returns a custom topology",
			args: args{params: NewElasticsearchParams{
				API:        api.NewMock(mock.New200Response(mock.NewStructBody(elasticsearchTemplateResponse))),
				Region:     "ece-region",
				Version:    "7.4.2",
				TemplateID: "default",
				Topology: []ElasticsearchTopologyElement{
					{Name: DataNode, Size: 8192, ZoneCount: 2},
					{Name: MasterNode, Size: 1024, ZoneCount: 1},
					{Name: MLNode, Size: 2048, ZoneCount: 1},
				},
			}},
			want: &models.ElasticsearchPayload{
				DisplayName: "",
				Region:      ec.String("ece-region"),
				RefID:       ec.String(DefaultElasticsearchRefID),
				Plan: &models.ElasticsearchClusterPlan{
					Elasticsearch: &models.ElasticsearchConfiguration{
						Version: "7.4.2",
					},
					DeploymentTemplate: &models.DeploymentTemplateReference{
						ID: ec.String(DefaultTemplateID),
					},
					ClusterTopology: []*models.ElasticsearchClusterTopologyElement{
						{
							ZoneCount:               2,
							InstanceConfigurationID: "default.data",
							Size: &models.TopologySize{
								Resource: ec.String("memory"),
								Value:    ec.Int32(8192),
							},
							NodeType: &models.ElasticsearchNodeType{
								Data: ec.Bool(true),
							},
						},
						{
							ZoneCount:               1,
							InstanceConfigurationID: "default.master",
							Size: &models.TopologySize{
								Resource: ec.String("memory"),
								Value:    ec.Int32(1024),
							},
							NodeType: &models.ElasticsearchNodeType{
								Master: ec.Bool(true),
							},
						},
						{
							ZoneCount:               1,
							InstanceConfigurationID: "default.ml",
							Size: &models.TopologySize{
								Resource: ec.String("memory"),
								Value:    ec.Int32(2048),
							},
							NodeType: &models.ElasticsearchNodeType{
								Ml: ec.Bool(true),
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewElasticsearch(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("NewElasticsearch() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewElasticsearch() = %v, want %v", got, tt.want)
			}
		})
	}
}
