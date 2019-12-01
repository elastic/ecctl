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
	"os"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/hashicorp/go-multierror"
)

func TestParseElasticsearchInput(t *testing.T) {
	var rawClusterTopology = []string{
		`{"name": "data", "size": 2048, "zone_count": 2}`,
		`{"name": "ml", "size": 4096, "zone_count": 1}`,
		`{"name": "master", "size": 1024, "zone_count": 1}`,
	}
	var clusterTopology = []*models.ElasticsearchClusterTopologyElement{
		{
			ZoneCount:               2,
			InstanceConfigurationID: "default.data",
			Size: &models.TopologySize{
				Resource: ec.String("memory"),
				Value:    ec.Int32(2048),
			},
			NodeType: &models.ElasticsearchNodeType{
				Data: ec.Bool(true),
			},
		},
		{
			ZoneCount:               1,
			InstanceConfigurationID: "default.ml",
			Size: &models.TopologySize{
				Resource: ec.String("memory"),
				Value:    ec.Int32(4096),
			},
			NodeType: &models.ElasticsearchNodeType{
				Ml: ec.Bool(true),
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
	}
	type args struct {
		params ParseElasticsearchInputParams
	}
	tests := []struct {
		name string
		args args
		want *models.ElasticsearchPayload
		err  error
	}{
		{
			name: "returns the payload directly when it's specified",
			args: args{params: ParseElasticsearchInputParams{
				Payload: &models.ElasticsearchPayload{DisplayName: "somename"},
			}},
			want: &models.ElasticsearchPayload{DisplayName: "somename"},
		},
		{
			name: "returns the payload from a set of raw topology elements",
			args: args{params: ParseElasticsearchInputParams{
				NewElasticsearchParams: NewElasticsearchParams{
					API:     api.NewMock(mock.New200Response(mock.NewStructBody(elasticsearchTemplateResponse))),
					Region:  "ece-region",
					Version: "7.4.2",
					Name:    "mycluster",
				},
				TopologyElements: rawClusterTopology,
			}},
			want: &models.ElasticsearchPayload{
				DisplayName: "mycluster",
				Region:      ec.String("ece-region"),
				RefID:       ec.String(DefaultElasticsearchRefID),
				Plan: &models.ElasticsearchClusterPlan{
					Elasticsearch: &models.ElasticsearchConfiguration{
						Version: "7.4.2",
					},
					DeploymentTemplate: &models.DeploymentTemplateReference{
						ID: ec.String(DefaultTemplateID),
					},
					ClusterTopology: clusterTopology,
				},
			},
		},
		{
			name: "returns the payload from a set of raw topology elements and auto-discovers version",
			args: args{params: ParseElasticsearchInputParams{
				NewElasticsearchParams: NewElasticsearchParams{
					API: api.NewMock(
						mock.New200Response(mock.NewStructBody(models.StackVersionConfigs{
							Stacks: []*models.StackVersionConfig{
								{Version: "6.4.2"},
								{Version: "7.4.2"},
								{Version: "5.4.2"},
							},
						})),
						mock.New200Response(mock.NewStructBody(elasticsearchTemplateResponse)),
					),
					Region: "ece-region",
					Name:   "mycluster",
				},
				TopologyElements: rawClusterTopology,
			}},
			want: &models.ElasticsearchPayload{
				DisplayName: "mycluster",
				Region:      ec.String("ece-region"),
				RefID:       ec.String(DefaultElasticsearchRefID),
				Plan: &models.ElasticsearchClusterPlan{
					Elasticsearch: &models.ElasticsearchConfiguration{
						Version: "7.4.2",
					},
					DeploymentTemplate: &models.DeploymentTemplateReference{
						ID: ec.String(DefaultTemplateID),
					},
					ClusterTopology: clusterTopology,
				},
			},
		},
		{
			name: "returns the payload from a set of raw topology elements and fails the version auto-discover",
			args: args{params: ParseElasticsearchInputParams{
				NewElasticsearchParams: NewElasticsearchParams{
					API: api.NewMock(
						mock.New200Response(mock.NewStringBody("failed to get the version error")),
					),
					Region: "ece-region",
					Name:   "mycluster",
				},
				TopologyElements: rawClusterTopology,
			}},
			err: errors.New("version discovery: failed to obtain stack list, please specify a version"),
		},
		{
			name: "returns the payload from size and zonecount elements and auto-discovers version",
			args: args{params: ParseElasticsearchInputParams{
				NewElasticsearchParams: NewElasticsearchParams{
					API: api.NewMock(
						mock.New200Response(mock.NewStructBody(models.StackVersionConfigs{
							Stacks: []*models.StackVersionConfig{
								{Version: "6.4.2"},
								{Version: "7.4.2"},
								{Version: "5.4.2"},
							},
						})),
						mock.New200Response(mock.NewStructBody(elasticsearchTemplateResponse)),
					),
					Region: "ece-region",
					Name:   "mycluster",
				},
				Size:      2048,
				ZoneCount: 3,
			}},
			want: &models.ElasticsearchPayload{
				DisplayName: "mycluster",
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
							ZoneCount:               3,
							InstanceConfigurationID: "default.data",
							Size: &models.TopologySize{
								Resource: ec.String("memory"),
								Value:    ec.Int32(2048),
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
			name: "returns the payload from size and zonecount elements and auto-discovers version",
			args: args{params: ParseElasticsearchInputParams{
				NewElasticsearchParams: NewElasticsearchParams{
					API:    api.NewMock(),
					Region: "ece-region",
					Name:   "mycluster",
				},
				TopologyElements: []string{
					`{"name": ""}`,
				},
			}},
			err: &multierror.Error{Errors: []error{
				errors.New("deployment topology: name cannot be empty"),
				errors.New("deployment topology: size cannot be empty"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseElasticsearchInput(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("ParseElasticsearchInput() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				json.NewEncoder(os.Stdout).Encode(got)
				println("WANT")
				json.NewEncoder(os.Stdout).Encode(tt.want)
				t.Errorf("ParseElasticsearchInput() = %v, want %v", got, tt.want)
			}
		})
	}
}
