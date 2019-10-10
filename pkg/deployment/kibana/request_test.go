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
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
)

func TestNewKibanaBody(t *testing.T) {
	type args struct {
		params NewKibanaBodyParams
	}
	tests := []struct {
		name string
		args args
		want *models.CreateKibanaRequest
	}{
		{
			name: "Fills in defaults when no params are specified",
			args: args{params: NewKibanaBodyParams{
				ID: "some ES ID",
			}},
			want: &models.CreateKibanaRequest{
				ClusterName:            "",
				ElasticsearchClusterID: ec.String("some ES ID"),
				Plan: &models.KibanaClusterPlan{
					ZoneCount: DefaultZoneCount,
					Kibana:    new(models.KibanaConfiguration),
					ClusterTopology: []*models.KibanaClusterTopologyElement{
						{
							MemoryPerNode:    DefaultMemoryPerNode,
							NodeCountPerZone: DefaultNodeCountPerZone,
						},
					},
				},
			},
		},
		{
			name: "Fills in with the specified parameters",
			args: args{params: NewKibanaBodyParams{
				ID:               "some other ES ID",
				Name:             "a name",
				ZoneCount:        2,
				MemoryPerNode:    2048,
				NodeCountPerZone: 3,
			}},
			want: &models.CreateKibanaRequest{
				ClusterName:            "a name",
				ElasticsearchClusterID: ec.String("some other ES ID"),
				Plan: &models.KibanaClusterPlan{
					ZoneCount: 2,
					Kibana:    new(models.KibanaConfiguration),
					ClusterTopology: []*models.KibanaClusterTopologyElement{
						{
							MemoryPerNode:    2048,
							NodeCountPerZone: 3,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewKibanaBody(tt.args.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKibanaBody() = %v, want %v", got, tt.want)
			}
		})
	}
}
