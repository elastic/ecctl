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

package elasticsearch

import (
	"errors"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func TestList(t *testing.T) {
	type args struct {
		params ListParams
	}
	tests := []struct {
		name string
		args args
		want *models.ElasticsearchClustersInfo
		err  error
	}{
		{
			name: "succeeds",
			args: args{params: ListParams{
				API: api.NewMock(mock.Response{
					Response: http.Response{
						StatusCode: 200,
						Body: mock.NewStructBody(models.ElasticsearchClustersInfo{
							ReturnCount: ec.Int32(2),
							ElasticsearchClusters: []*models.ElasticsearchClusterInfo{
								{ClusterID: ec.String(util.ValidClusterID)},
								{ClusterID: ec.String("123456")},
							},
						}),
					},
				}),
			}},
			want: &models.ElasticsearchClustersInfo{
				ReturnCount: ec.Int32(2),
				ElasticsearchClusters: []*models.ElasticsearchClusterInfo{
					{ClusterID: ec.String(util.ValidClusterID)},
					{ClusterID: ec.String("123456")},
				},
			},
		},
		{
			name: "succeeds with filter",
			args: args{params: ListParams{
				Version: "7.0.0",
				API: api.NewMock(mock.Response{
					Response: http.Response{
						StatusCode: 200,
						Body: mock.NewStructBody(models.ElasticsearchClustersInfo{
							ReturnCount: ec.Int32(2),
							ElasticsearchClusters: []*models.ElasticsearchClusterInfo{
								{
									ClusterID: ec.String(util.ValidClusterID),
									PlanInfo: &models.ElasticsearchClusterPlansInfo{
										Current: &models.ElasticsearchClusterPlanInfo{
											Plan: &models.ElasticsearchClusterPlan{
												Elasticsearch: &models.ElasticsearchConfiguration{
													Version: "7.0.0",
												},
											},
										},
									},
								},
								{ClusterID: ec.String("123456")},
							},
						}),
					},
				}),
			}},
			want: &models.ElasticsearchClustersInfo{
				ReturnCount: ec.Int32(1),
				ElasticsearchClusters: []*models.ElasticsearchClusterInfo{
					{
						ClusterID: ec.String(util.ValidClusterID),
						PlanInfo: &models.ElasticsearchClusterPlansInfo{
							Current: &models.ElasticsearchClusterPlanInfo{
								Plan: &models.ElasticsearchClusterPlan{
									Elasticsearch: &models.ElasticsearchConfiguration{
										Version: "7.0.0",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "fails due to api error",
			args: args{params: ListParams{
				API: api.NewMock(mock.New500Response(mock.NewStringBody(`{"error": "some error"}`))),
			}},
			err: errors.New(`{"error": "some error"}`),
		},
		{
			name: "fails due to parameter validation",
			args: args{},
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := List(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("List() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func newDefault5ScriptingSettings() *models.ElasticsearchScriptingUserSettings {
	return &models.ElasticsearchScriptingUserSettings{
		File: &models.ElasticsearchScriptTypeSettings{
			Enabled:     ec.Bool(true),
			SandboxMode: ec.Bool(true),
		},
		Inline: &models.ElasticsearchScriptTypeSettings{
			Enabled:     ec.Bool(true),
			SandboxMode: ec.Bool(true),
		},
		Stored: &models.ElasticsearchScriptTypeSettings{
			Enabled:     ec.Bool(true),
			SandboxMode: ec.Bool(true),
		},
		ExpressionsEnabled: ec.Bool(true),
		MustacheEnabled:    ec.Bool(true),
		PainlessEnabled:    ec.Bool(true),
	}
}

func newDefault6ScriptingSettings() *models.ElasticsearchScriptingUserSettings {
	return &models.ElasticsearchScriptingUserSettings{
		// File: &models.ElasticsearchScriptTypeSettings{},
		Inline: &models.ElasticsearchScriptTypeSettings{
			Enabled: ec.Bool(true),
		},
		Stored: &models.ElasticsearchScriptTypeSettings{
			Enabled: ec.Bool(true),
		},
	}
}

func newElasticSearchClusterInfoFromVersion(v string) *models.ElasticsearchClusterInfo {
	var systemSettings *models.ElasticsearchSystemSettings
	// Resembles default 2.x settings
	if strings.Split(v, ".")[0] == "2" {
		systemSettings = &models.ElasticsearchSystemSettings{
			DefaultShardsPerIndex: 1,
		}
	}
	// Resembles 5.x default settings
	if strings.Split(v, ".")[0] == "5" {
		systemSettings = &models.ElasticsearchSystemSettings{
			Scripting: newDefault5ScriptingSettings(),
		}
	}
	// Resembles 6.x default settings
	if strings.Split(v, ".")[0] == "6" {
		systemSettings = &models.ElasticsearchSystemSettings{
			Scripting: newDefault6ScriptingSettings(),
		}
	}
	return &models.ElasticsearchClusterInfo{
		PlanInfo: &models.ElasticsearchClusterPlansInfo{
			Current: &models.ElasticsearchClusterPlanInfo{
				Plan: &models.ElasticsearchClusterPlan{
					Elasticsearch: &models.ElasticsearchConfiguration{
						Version:        v,
						SystemSettings: systemSettings,
					},
				},
			},
		},
	}
}

func TestAppfilterVersion(t *testing.T) {
	type args struct {
		clusters *models.ElasticsearchClustersInfo
		v        string
	}
	tests := []struct {
		name string
		args args
		want *models.ElasticsearchClustersInfo
	}{
		{
			"Empty version filter returns the same information",
			args{
				&models.ElasticsearchClustersInfo{},
				"",
			},
			&models.ElasticsearchClustersInfo{},
		},
		{
			"specific version filter returns a specific set of clusters",
			args{
				&models.ElasticsearchClustersInfo{
					ElasticsearchClusters: []*models.ElasticsearchClusterInfo{
						newElasticSearchClusterInfoFromVersion("5.5.0"),
						newElasticSearchClusterInfoFromVersion("2.4.5"),
						newElasticSearchClusterInfoFromVersion("1.7.6"),
						newElasticSearchClusterInfoFromVersion("5.5.0"),
						newElasticSearchClusterInfoFromVersion("5.4.3"),
						newElasticSearchClusterInfoFromVersion("5.5.0"),
					},
				},
				"5.5.0",
			},
			&models.ElasticsearchClustersInfo{
				ElasticsearchClusters: []*models.ElasticsearchClusterInfo{
					newElasticSearchClusterInfoFromVersion("5.5.0"),
					newElasticSearchClusterInfoFromVersion("5.5.0"),
					newElasticSearchClusterInfoFromVersion("5.5.0"),
				},
				ReturnCount: ec.Int32(3),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterVersion(tt.args.clusters, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}
