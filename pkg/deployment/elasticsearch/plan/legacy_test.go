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

package plan

import (
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
)

func TestNewLegacyPlan(t *testing.T) {
	type args struct {
		params LegacyParams
	}
	tests := []struct {
		name string
		args args
		want *models.ElasticsearchClusterPlan
	}{
		{
			"New Plan with plugins",
			args{params: LegacyParams{
				Version:   "5.5.0",
				ZoneCount: int32(2),
				Capacity:  int32(2048),
				Plugins:   []string{"analysis-icu", "analysis-kuromoji", "analysis-phonetic", "analysis-smartcn", "analysis-stempel", "analysis-ukrainian", "ingest-attachment", "ingest-geoip", "ingest-user-agent", "mapper-murmur3", "mapper-size", "repository-azure", "repository-gcs"},
			}},
			&models.ElasticsearchClusterPlan{
				Elasticsearch: &models.ElasticsearchConfiguration{
					Version:               "5.5.0",
					EnabledBuiltInPlugins: []string{"analysis-icu", "analysis-kuromoji", "analysis-phonetic", "analysis-smartcn", "analysis-stempel", "analysis-ukrainian", "ingest-attachment", "ingest-geoip", "ingest-user-agent", "mapper-murmur3", "mapper-size", "repository-azure", "repository-gcs"},
					SystemSettings: &models.ElasticsearchSystemSettings{
						Scripting: DefaultScriptingSettings["5"],
					},
				},
				ZoneCount: 2,
				ClusterTopology: []*models.ElasticsearchClusterTopologyElement{
					{
						NodeType: &models.ElasticsearchNodeType{
							Data:   ec.Bool(true),
							Ingest: ec.Bool(true),
							Master: ec.Bool(true),
						},
						MemoryPerNode:     int32(2048),
						NodeCountPerZone:  DefaultNodeCountPerZone,
						NodeConfiguration: DefaultNodeConfiguration,
					},
				},
				Transient: &models.TransientElasticsearchPlanConfiguration{
					Strategy: &models.PlanStrategy{},
				},
			},
		},
		{
			"New Plan without plugins",
			args{params: LegacyParams{
				Version:   "5.5.0",
				ZoneCount: int32(2),
				Capacity:  int32(2048),
			}},
			&models.ElasticsearchClusterPlan{
				Elasticsearch: &models.ElasticsearchConfiguration{
					Version:               "5.5.0",
					EnabledBuiltInPlugins: nil,
					SystemSettings: &models.ElasticsearchSystemSettings{
						Scripting: DefaultScriptingSettings["5"],
					},
				},
				ZoneCount: 2,
				ClusterTopology: []*models.ElasticsearchClusterTopologyElement{
					{
						NodeType: &models.ElasticsearchNodeType{
							Data:   ec.Bool(true),
							Ingest: ec.Bool(true),
							Master: ec.Bool(true),
						},
						MemoryPerNode:     int32(2048),
						NodeCountPerZone:  DefaultNodeCountPerZone,
						NodeConfiguration: DefaultNodeConfiguration,
					},
				},
				Transient: &models.TransientElasticsearchPlanConfiguration{
					Strategy: &models.PlanStrategy{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLegacyPlan(tt.args.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLegacyPlan() = %v, want %v", got, tt.want)
			}
		})
	}
}
