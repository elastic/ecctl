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

	"github.com/elastic/ecctl/pkg/util"
)

var kibanaTemplateResponse = models.DeploymentTemplateInfo{
	ID: "default",
	ClusterTemplate: &models.DeploymentTemplateDefinitionRequest{
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

func TestNewKibana(t *testing.T) {
	var internalError = models.BasicFailedReply{
		Errors: []*models.BasicFailedReplyElement{
			{},
		},
	}
	internalErrorBytes, _ := json.MarshalIndent(internalError, "", "  ")

	var getResponse = models.DeploymentGetResponse{
		Resources: &models.DeploymentResources{
			Elasticsearch: []*models.ElasticsearchResourceInfo{{
				RefID: ec.String("main-elasticsearch"),
				Info: &models.ElasticsearchClusterInfo{
					PlanInfo: &models.ElasticsearchClusterPlansInfo{
						Current: &models.ElasticsearchClusterPlanInfo{
							Plan: &models.ElasticsearchClusterPlan{
								DeploymentTemplate: &models.DeploymentTemplateReference{
									ID: ec.String("an ID"),
								},
							},
						},
					},
				},
			}},
		},
	}

	type args struct {
		params NewStateless
	}
	tests := []struct {
		name string
		args args
		want *models.KibanaPayload
		err  error
	}{
		{
			name: "fails due to parameter validation",
			args: args{params: NewStateless{}},
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				util.ErrDeploymentID,
				errors.New("deployment topology: region cannot be empty"),
			}},
		},
		{
			name: "fails obtaining the deployment info",
			args: args{params: NewStateless{
				DeploymentID: util.ValidClusterID,
				API:          api.NewMock(mock.New500Response(mock.NewStructBody(internalError))),
				Region:       "ece-region",
			}},
			err: errors.New(string(internalErrorBytes)),
		},
		{
			name: "obtains the deployment info but fails getting the template ID info",
			args: args{params: NewStateless{
				DeploymentID: util.ValidClusterID,
				API: api.NewMock(
					mock.New200Response(mock.NewStructBody(models.DeploymentGetResponse{
						Resources: &models.DeploymentResources{
							Elasticsearch: []*models.ElasticsearchResourceInfo{{
								Info: &models.ElasticsearchClusterInfo{
									PlanInfo: &models.ElasticsearchClusterPlansInfo{},
								},
							}},
						},
					})),
				),
				Region: "ece-region",
			}},
			err: errors.New("unable to obtain deployment template ID from existing deployment ID, please specify a one"),
		},
		{
			name: "obtains the deployment info but fails getting the template ID info",
			args: args{params: NewStateless{
				DeploymentID: util.ValidClusterID,
				API: api.NewMock(
					mock.New200Response(mock.NewStructBody(getResponse)),
					mock.New404Response(mock.NewStructBody(internalError)),
				),
				Region: "ece-region",
			}},
			err: errors.New(string(internalErrorBytes)),
		},
		{
			name: "succeeds with no argument override",
			args: args{params: NewStateless{
				DeploymentID: util.ValidClusterID,
				API: api.NewMock(
					mock.New200Response(mock.NewStructBody(getResponse)),
					mock.New200Response(mock.NewStructBody(kibanaTemplateResponse)),
				),
				Region: "ece-region",
			}},
			want: &models.KibanaPayload{
				ElasticsearchClusterRefID: ec.String("main-elasticsearch"),
				Region:                    ec.String("ece-region"),
				RefID:                     ec.String("main-kibana"),
				Plan: &models.KibanaClusterPlan{
					Kibana: &models.KibanaConfiguration{},
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
		},
		{
			name: "succeeds with argument overrides",
			args: args{params: NewStateless{
				Size:         4096,
				ZoneCount:    3,
				DeploymentID: util.ValidClusterID,
				API: api.NewMock(
					mock.New200Response(mock.NewStructBody(getResponse)),
					mock.New200Response(mock.NewStructBody(kibanaTemplateResponse)),
				),
				Region: "ece-region",
			}},
			want: &models.KibanaPayload{
				ElasticsearchClusterRefID: ec.String("main-elasticsearch"),
				Region:                    ec.String("ece-region"),
				RefID:                     ec.String("main-kibana"),
				Plan: &models.KibanaClusterPlan{
					Kibana: &models.KibanaConfiguration{},
					ClusterTopology: []*models.KibanaClusterTopologyElement{
						{
							Size: &models.TopologySize{
								Resource: ec.String("memory"),
								Value:    ec.Int32(4096),
							},
							ZoneCount: 3,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewKibana(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("NewKibana() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				g, _ := json.Marshal(got)
				w, _ := json.Marshal(tt.want)
				println(string(g))
				println(string(w))
				t.Errorf("NewKibana() = %v, want %v", got, tt.want)
			}
		})
	}
}
