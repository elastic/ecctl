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

package apm

import (
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
)

func TestNewApmBody(t *testing.T) {
	type args struct {
		params NewApmBodyParams
	}
	tests := []struct {
		name string
		args args
		want *models.CreateApmRequest
	}{
		{
			name: "Fills in defaults when no params are specified",
			args: args{params: NewApmBodyParams{
				ID: "some APM ID",
			}},
			want: &models.CreateApmRequest{
				ElasticsearchClusterID: ec.String("some APM ID"),
				Plan: &models.ApmPlan{
					Apm: new(models.ApmConfiguration),
					ClusterTopology: []*models.ApmTopologyElement{
						{
							ZoneCount: DefaultZoneCount,
							Size: &models.TopologySize{
								Resource: ec.String("memory"),
								Value:    ec.Int32(DefaultMemoryPerNode),
							},
						},
					},
				},
			},
		},
		{
			name: "Fills in with the specified parameters",
			args: args{params: NewApmBodyParams{
				ID:            "some other APM ID",
				Name:          "a name",
				ZoneCount:     2,
				MemoryPerNode: 2048,
			}},
			want: &models.CreateApmRequest{
				DisplayName:            "a name",
				ElasticsearchClusterID: ec.String("some other APM ID"),
				Plan: &models.ApmPlan{
					Apm: new(models.ApmConfiguration),
					ClusterTopology: []*models.ApmTopologyElement{
						{
							ZoneCount: 2,
							Size: &models.TopologySize{
								Resource: ec.String("memory"),
								Value:    ec.Int32(2048),
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewApmBody(tt.args.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewApmBody() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
