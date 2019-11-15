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

var basicESCluster = `{
  "created": true,
  "id": "0837d2cd080743e9be080bca163c0b92",
  "name": "my example cluster",
  "resources": [{
    "cloud_id": "my_elasticsearch_cluster:MTkyLjE2OC40NC4xMC5pcC5lcy5pbzo5MjQzJDAxZmEyODU4NzZmNTRlNjk5ZGEzZDNkNmZkOGE4NGYxJGNhOGFjNjU1NWYwZDQzZDhiYTEwNDhjOThlYTYwMjY1",
    "credentials": {
      "password": "6n7Q5fXoFZDnpLOVPi5FnVLa",
      "username": "elastic"
    },
    "id": "01fa285876f54e699da3d3d6fd8a84f1",
    "kind": "elasticsearch",
    "ref_id": "my-es-cluster",
    "region": "ece-region"
  }, {
    "elasticsearch_cluster_ref_id": "my-es-cluster",
    "id": "ca8ac6555f0d43d8ba1048c98ea60265",
    "kind": "kibana",
    "ref_id": "my-kibana-instance",
    "region": "ece-region"
  }]
}`

func TestCreate(t *testing.T) {
	type args struct {
		params CreateParams
	}
	tests := []struct {
		name string
		args args
		want *models.DeploymentCreateResponse
		err  error
	}{
		{
			name: "fails on parameter validation",
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				errors.New("deployment create: request payload cannot be empty"),
			}},
		},
		{
			name: "fails on API error",
			args: args{params: CreateParams{
				API:     api.NewMock(mock.New500Response(mock.NewStringBody("error"))),
				Request: &models.DeploymentCreateRequest{},
			}},
			err: errors.New("error"),
		},
		{
			name: "succeeds",
			args: args{params: CreateParams{
				API: api.NewMock(mock.New201Response(mock.NewStringBody(basicESCluster))),
				Request: &models.DeploymentCreateRequest{
					Name: "my example cluster",
					Resources: &models.DeploymentCreateResources{
						Elasticsearch: []*models.ElasticsearchPayload{
							{Plan: &models.ElasticsearchClusterPlan{
								Elasticsearch: &models.ElasticsearchConfiguration{Version: "6.8.4"},
							}},
						},
						Kibana: []*models.KibanaPayload{
							{Plan: &models.KibanaClusterPlan{
								Kibana: &models.KibanaConfiguration{Version: "6.8.4"},
							}},
						},
					},
				},
			}},
			want: &models.DeploymentCreateResponse{
				Created: ec.Bool(true),
				ID:      ec.String("0837d2cd080743e9be080bca163c0b92"),
				Name:    ec.String("my example cluster"),
				Resources: []*models.DeploymentResource{
					{
						CloudID: "my_elasticsearch_cluster:MTkyLjE2OC40NC4xMC5pcC5lcy5pbzo5MjQzJDAxZmEyODU4NzZmNTRlNjk5ZGEzZDNkNmZkOGE4NGYxJGNhOGFjNjU1NWYwZDQzZDhiYTEwNDhjOThlYTYwMjY1",
						Credentials: &models.ClusterCredentials{
							Password: "6n7Q5fXoFZDnpLOVPi5FnVLa",
							Username: "elastic",
						},
						ID:     ec.String("01fa285876f54e699da3d3d6fd8a84f1"),
						Kind:   ec.String("elasticsearch"),
						RefID:  ec.String("my-es-cluster"),
						Region: ec.String("ece-region"),
					},
					{
						ElasticsearchClusterRefID: "my-es-cluster",
						ID:                        ec.String("ca8ac6555f0d43d8ba1048c98ea60265"),
						Kind:                      ec.String("kibana"),
						RefID:                     ec.String("my-kibana-instance"),
						Region:                    ec.String("ece-region"),
					},
				},
			},
		},
		{
			name: "succeeds with idempotency ID",
			args: args{params: CreateParams{
				RequestID: "1232131231",
				API:       api.NewMock(mock.New201Response(mock.NewStringBody(basicESCluster))),
				Request: &models.DeploymentCreateRequest{
					Name: "my example cluster",
					Resources: &models.DeploymentCreateResources{
						Elasticsearch: []*models.ElasticsearchPayload{
							{Plan: &models.ElasticsearchClusterPlan{
								Elasticsearch: &models.ElasticsearchConfiguration{Version: "6.8.4"},
							}},
						},
						Kibana: []*models.KibanaPayload{
							{Plan: &models.KibanaClusterPlan{
								Kibana: &models.KibanaConfiguration{Version: "6.8.4"},
							}},
						},
					},
				},
			}},
			want: &models.DeploymentCreateResponse{
				Created: ec.Bool(true),
				ID:      ec.String("0837d2cd080743e9be080bca163c0b92"),
				Name:    ec.String("my example cluster"),
				Resources: []*models.DeploymentResource{
					{
						CloudID: "my_elasticsearch_cluster:MTkyLjE2OC40NC4xMC5pcC5lcy5pbzo5MjQzJDAxZmEyODU4NzZmNTRlNjk5ZGEzZDNkNmZkOGE4NGYxJGNhOGFjNjU1NWYwZDQzZDhiYTEwNDhjOThlYTYwMjY1",
						Credentials: &models.ClusterCredentials{
							Password: "6n7Q5fXoFZDnpLOVPi5FnVLa",
							Username: "elastic",
						},
						ID:     ec.String("01fa285876f54e699da3d3d6fd8a84f1"),
						Kind:   ec.String("elasticsearch"),
						RefID:  ec.String("my-es-cluster"),
						Region: ec.String("ece-region"),
					},
					{
						ElasticsearchClusterRefID: "my-es-cluster",
						ID:                        ec.String("ca8ac6555f0d43d8ba1048c98ea60265"),
						Kind:                      ec.String("kibana"),
						RefID:                     ec.String("my-kibana-instance"),
						Region:                    ec.String("ece-region"),
					},
				},
			},
		},
		{
			name: "succeeds with idempotency ID returns a 202 when the resource is still creating with the same ID",
			args: args{params: CreateParams{
				RequestID: "1232131231",
				API:       api.NewMock(mock.New202Response(mock.NewStringBody(basicESCluster))),
				Request: &models.DeploymentCreateRequest{
					Name: "my example cluster",
					Resources: &models.DeploymentCreateResources{
						Elasticsearch: []*models.ElasticsearchPayload{
							{Plan: &models.ElasticsearchClusterPlan{
								Elasticsearch: &models.ElasticsearchConfiguration{Version: "6.8.4"},
							}},
						},
						Kibana: []*models.KibanaPayload{
							{Plan: &models.KibanaClusterPlan{
								Kibana: &models.KibanaConfiguration{Version: "6.8.4"},
							}},
						},
					},
				},
			}},
			want: &models.DeploymentCreateResponse{
				Created: ec.Bool(true),
				ID:      ec.String("0837d2cd080743e9be080bca163c0b92"),
				Name:    ec.String("my example cluster"),
				Resources: []*models.DeploymentResource{
					{
						CloudID: "my_elasticsearch_cluster:MTkyLjE2OC40NC4xMC5pcC5lcy5pbzo5MjQzJDAxZmEyODU4NzZmNTRlNjk5ZGEzZDNkNmZkOGE4NGYxJGNhOGFjNjU1NWYwZDQzZDhiYTEwNDhjOThlYTYwMjY1",
						Credentials: &models.ClusterCredentials{
							Password: "6n7Q5fXoFZDnpLOVPi5FnVLa",
							Username: "elastic",
						},
						ID:     ec.String("01fa285876f54e699da3d3d6fd8a84f1"),
						Kind:   ec.String("elasticsearch"),
						RefID:  ec.String("my-es-cluster"),
						Region: ec.String("ece-region"),
					},
					{
						ElasticsearchClusterRefID: "my-es-cluster",
						ID:                        ec.String("ca8ac6555f0d43d8ba1048c98ea60265"),
						Kind:                      ec.String("kibana"),
						RefID:                     ec.String("my-kibana-instance"),
						Region:                    ec.String("ece-region"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Create(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setOverrides(t *testing.T) {
	var eceRegion = "ece-region"
	type args struct {
		req       *models.DeploymentCreateRequest
		overrides *CreateOverrides
	}
	tests := []struct {
		name string
		args args
		want *models.DeploymentCreateRequest
	}{
		{
			name: "set name override",
			args: args{
				req: &models.DeploymentCreateRequest{
					Name:      "Some",
					Resources: &models.DeploymentCreateResources{},
				},
				overrides: &CreateOverrides{Name: "some other"},
			},
			want: &models.DeploymentCreateRequest{
				Name:      "some other",
				Resources: &models.DeploymentCreateResources{},
			},
		},
		{
			name: "set name, version and region override",
			args: args{
				overrides: &CreateOverrides{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req = tt.args.req
			setOverrides(req, tt.args.overrides)

			if !reflect.DeepEqual(req, tt.want) {
				t.Errorf("Create() = %v, want %v", req, tt.want)
			}
		})
	}
}
