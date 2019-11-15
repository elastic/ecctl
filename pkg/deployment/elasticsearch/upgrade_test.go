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
	"bytes"
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/output"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment/planutil"
	"github.com/elastic/ecctl/pkg/util"
)

func TestUpgrade(t *testing.T) {
	type args struct {
		params UpgradeParams
	}
	tests := []struct {
		name string
		args args
		want *models.ClusterCrudResponse
		err  error
	}{
		{
			name: "fails due to empty params",
			args: args{},
			err: &multierror.Error{Errors: []error{
				errors.New("cluster id should have a length of 32 characters"),
				errors.New("version string empty"),
			}},
		},
		{
			name: "fails due to erroneus version params",
			args: args{params: UpgradeParams{
				Version: "xxxx",
			}},
			err: &multierror.Error{Errors: []error{
				errors.New("cluster id should have a length of 32 characters"),
				errors.New("no major.minor.patch elements found"),
			}},
		},
		{
			name: "getting the current cluster plan fails",
			args: args{params: UpgradeParams{
				Version: "5.0.0",
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API: api.NewMock(mock.Response{
						Response: http.Response{
							StatusCode: 200,
							Body:       mock.NewStringBody("{}"),
						},
					}),
				},
			}},
			err: errors.New("cluster has no plan information"),
		},
		{
			name: "getting the current cluste plan fails due to API error",
			args: args{params: UpgradeParams{
				Version: "5.0.0",
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API:       api.NewMock(mock.New500Response(mock.NewStringBody(`{"error": "some error"}`))),
				},
			}},
			err: errors.New(`{"error": "some error"}`),
		},
		{
			name: "succeeds without tracking",
			args: args{params: UpgradeParams{
				Version: "5.0.0",
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API: api.NewMock(mock.Response{
						Response: http.Response{
							StatusCode: 200,
							Body: mock.NewStructBody(
								newElasticSearchClusterInfoFromVersion("2.4.5"),
							),
						},
					}, mock.Response{
						Response: http.Response{
							StatusCode: 202,
							Body: mock.NewStructBody(models.ClusterCrudResponse{
								ElasticsearchClusterID: util.ValidClusterID,
							}),
						},
					}),
				},
			}},
			want: &models.ClusterCrudResponse{
				ElasticsearchClusterID: util.ValidClusterID,
			},
		},
		{
			name: "succeeds with tracking",
			args: args{params: UpgradeParams{
				TrackParams: util.TrackParams{
					Track:         true,
					Output:        output.NewDevice(new(bytes.Buffer)),
					PollFrequency: time.Millisecond,
					MaxRetries:    1,
				},
				Version: "5.0.0",
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API: api.NewMock(util.AppendTrackResponses(
						mock.Response{
							Response: http.Response{
								StatusCode: 200,
								Body: mock.NewStructBody(
									newElasticSearchClusterInfoFromVersion("2.4.5"),
								),
							},
						},
						mock.Response{
							Response: http.Response{
								StatusCode: 202,
								Body: mock.NewStructBody(models.ClusterCrudResponse{
									ElasticsearchClusterID: util.ValidClusterID,
								}),
							},
						},
						util.PlanNotFound,
					)...),
				},
			}},
			want: &models.ClusterCrudResponse{
				ElasticsearchClusterID: util.ValidClusterID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Upgrade(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Upgrade() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Upgrade() = %v, want %v", got, tt.want)
			}
		})
	}
}

func newInstanceTopology(p, c, d int) *models.ClusterInstanceInfo {
	return &models.ClusterInstanceInfo{
		Memory: &models.ClusterInstanceMemoryInfo{
			MemoryPressure:   int32(p),
			InstanceCapacity: ec.Int32(int32(c)),
		},
		Disk: &models.ClusterInstanceDiskInfo{
			DiskSpaceUsed: ec.Int64(int64(d)),
		},
	}
}

func newElasticSearchClusterInfoWithVersionAndTopology(v string, t ...*models.ClusterInstanceInfo) *models.ElasticsearchClusterInfo {
	ci := newElasticSearchClusterInfoFromVersion(v)
	ci.Topology = &models.ClusterTopologyInfo{
		Instances: []*models.ClusterInstanceInfo{},
	}
	ci.Topology.Instances = append(ci.Topology.Instances, t...)

	return ci
}

func TestAppcomputeStrategyForClusterUpgrade(t *testing.T) {
	type args struct {
		cluster *models.ElasticsearchClusterInfo
		version string
	}
	tests := []struct {
		name string
		args args
		want *models.ElasticsearchClusterPlan
	}{
		{
			"Upgrade version 2.4.5 to 5.5.0",
			args{
				newElasticSearchClusterInfoFromVersion("2.4.5"),
				"5.5.0",
			},
			&models.ElasticsearchClusterPlan{
				Elasticsearch: &models.ElasticsearchConfiguration{
					Version: "5.5.0",
					SystemSettings: &models.ElasticsearchSystemSettings{
						DefaultShardsPerIndex: 0,
					},
				},
				Transient: &models.TransientElasticsearchPlanConfiguration{
					Strategy: planutil.MajorUpgradeStrategy,
				},
			},
		},
		{
			"Upgrade version 5.4.1 to 5.5.0",
			args{
				newElasticSearchClusterInfoWithVersionAndTopology(
					"5.4.1",
					newInstanceTopology(10, 1024, 200),
				),
				"5.5.0",
			},
			&models.ElasticsearchClusterPlan{
				Elasticsearch: &models.ElasticsearchConfiguration{
					Version: "5.5.0",
					SystemSettings: &models.ElasticsearchSystemSettings{
						Scripting: newDefault5ScriptingSettings(),
					},
				},
				Transient: &models.TransientElasticsearchPlanConfiguration{
					Strategy: planutil.DefaultPlanStrategy,
				},
			},
		},
		{
			"Upgrade version 5.6.4 to 6.0.0",
			args{
				newElasticSearchClusterInfoFromVersion("5.6.4"),
				"6.0.0",
			},
			&models.ElasticsearchClusterPlan{
				Elasticsearch: &models.ElasticsearchConfiguration{
					Version: "6.0.0",
					SystemSettings: &models.ElasticsearchSystemSettings{
						Scripting: newDefault6ScriptingSettings(),
					},
				},
				Transient: &models.TransientElasticsearchPlanConfiguration{
					Strategy: planutil.MajorUpgradeStrategy,
				},
			},
		},
		{
			"Upgrade version 5.4.3 to 5.5.0 With High Memory Pressure",
			args{
				newElasticSearchClusterInfoWithVersionAndTopology(
					"5.4.3",
					newInstanceTopology(70, 1024, 200),
				),
				"5.5.0",
			},
			&models.ElasticsearchClusterPlan{
				Elasticsearch: &models.ElasticsearchConfiguration{
					Version: "5.5.0",
					SystemSettings: &models.ElasticsearchSystemSettings{
						Scripting: newDefault5ScriptingSettings(),
					},
				},
				Transient: &models.TransientElasticsearchPlanConfiguration{
					Strategy: planutil.RollingByNameStrategy,
				},
			},
		},
		{
			"Upgrade version 5.4.3 to 5.5.0 With High % of Disk used",
			args{
				newElasticSearchClusterInfoWithVersionAndTopology(
					"5.4.3",
					newInstanceTopology(10, 1024, 7168),
				),
				"5.5.0",
			},
			&models.ElasticsearchClusterPlan{
				Elasticsearch: &models.ElasticsearchConfiguration{
					Version: "5.5.0",
					SystemSettings: &models.ElasticsearchSystemSettings{
						Scripting: newDefault5ScriptingSettings(),
					},
				},
				Transient: &models.TransientElasticsearchPlanConfiguration{
					Strategy: planutil.RollingByNameStrategy,
				},
			},
		},
		{
			"Upgrade version 5.4.3 to 5.5.0 with more than 100GB used",
			args{
				newElasticSearchClusterInfoWithVersionAndTopology(
					"5.4.3",
					newInstanceTopology(10, 32768, 102400),
				),
				"5.5.0",
			},
			&models.ElasticsearchClusterPlan{
				Elasticsearch: &models.ElasticsearchConfiguration{
					Version: "5.5.0",
					SystemSettings: &models.ElasticsearchSystemSettings{
						Scripting: newDefault5ScriptingSettings(),
					},
				},
				Transient: &models.TransientElasticsearchPlanConfiguration{
					Strategy: planutil.RollingByNameStrategy,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getUpgradePlan(tt.args.cluster, tt.args.version); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getUpgradePlan() = %+v, want %+v", got.Transient.Strategy, tt.want.Transient.Strategy)
			}
		})
	}
}
