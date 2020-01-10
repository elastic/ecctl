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

	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/hashicorp/go-multierror"
)

func TestNewElasticsearchTopology(t *testing.T) {
	type args struct {
		topology []string
	}
	tests := []struct {
		name string
		args args
		want []ElasticsearchTopologyElement
		err  error
	}{
		{
			name: "correctly parses a single element topology",
			args: args{topology: []string{
				`{"name": "data", "size": 1024, "zone_count": 1}`,
			}},
			want: []ElasticsearchTopologyElement{
				{Name: "data", Size: 1024, ZoneCount: 1},
			},
		},
		{
			name: "correctly parses a multi element topology",
			args: args{topology: []string{
				`{"name": "data", "size": 2048, "zone_count": 2}`,
				`{"name": "ml", "size": 4096, "zone_count": 1}`,
				`{"name": "master", "size": 1024, "zone_count": 1}`,
			}},
			want: []ElasticsearchTopologyElement{
				{Name: "data", Size: 2048, ZoneCount: 2},
				{Name: "ml", Size: 4096, ZoneCount: 1},
				{Name: "master", Size: 1024, ZoneCount: 1},
			},
		},
		{
			name: "fails due to invalid json parses a multi element topology",
			args: args{topology: []string{
				`{"name": "data", "size": 2048, "zone_count": 2}`,
				`{"name": "ml", "size": 4096, "zone_count": 1}`,
				`{"aaaaaaaaaa`,
			}},
			err: errors.New("depresource: failed unpacking raw topology: unexpected end of JSON input"),
		},
		{
			name: "fails due to missing name",
			args: args{topology: []string{
				`{"zone_count": 2}`,
			}},
			err: &multierror.Error{Errors: []error{
				errors.New("deployment topology: name cannot be empty"),
				errors.New("deployment topology: size cannot be empty"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewElasticsearchTopology(tt.args.topology)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("NewElasticsearchTopology() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewElasticsearchTopology() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewElasticsearchTopologyElement(t *testing.T) {
	type args struct {
		size      int32
		zoneCount int32
	}
	tests := []struct {
		name string
		args args
		want ElasticsearchTopologyElement
	}{
		{
			name: "empty size and zonecount returns the default values",
			args: args{},
			want: DefaultTopologyElement,
		},
		{
			name: "returns a parametrized ElasticsearchTopologyElement",
			args: args{size: 2048, zoneCount: 3},
			want: ElasticsearchTopologyElement{
				Name:      DataNode,
				Size:      2048,
				ZoneCount: 3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewElasticsearchTopologyElement(tt.args.size, tt.args.zoneCount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewElasticsearchTopologyElement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildElasticsearchTopology(t *testing.T) {
	var topologies = []*models.ElasticsearchClusterTopologyElement{
		{
			Size: &models.TopologySize{
				Resource: ec.String("memory"),
				Value:    ec.Int32(1024),
			},
			NodeType: &models.ElasticsearchNodeType{
				Data: ec.Bool(true),
			},
		},
		{
			Size: &models.TopologySize{
				Resource: ec.String("memory"),
				Value:    ec.Int32(1024),
			},
			NodeType: &models.ElasticsearchNodeType{
				Ml: ec.Bool(true),
			},
		},
		{
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
		params BuildElasticsearchTopologyParams
	}
	tests := []struct {
		name string
		args args
		want []*models.ElasticsearchClusterTopologyElement
		err  error
	}{
		{
			name: "returns error on empty topology",
			args: args{params: BuildElasticsearchTopologyParams{
				TemplateID: "sometemplateid",
			}},
			err: errors.New(`deployment topology: failed to obtain desired topology names ([]) in deployment template id "sometemplateid"`),
		},
		{
			name: "returns error on non matching desired topology",
			args: args{params: BuildElasticsearchTopologyParams{
				TemplateID: "sometemplateid",
				Topology: []ElasticsearchTopologyElement{
					{Name: "something weird", Size: 2048, ZoneCount: 1},
				},
				ClusterTopology: topologies,
			}},
			err: errors.New(`deployment topology: failed to obtain desired topology names ([{Name:something weird ZoneCount:1 Size:2048}]) in deployment template id "sometemplateid"`),
		},
		{
			name: "matches the topologies",
			args: args{params: BuildElasticsearchTopologyParams{
				TemplateID: "sometemplateid",
				Topology: []ElasticsearchTopologyElement{
					{Name: DataNode, Size: 8192, ZoneCount: 2},
					{Name: MasterNode, Size: 1024, ZoneCount: 1},
					{Name: MLNode, Size: 2048, ZoneCount: 1},
				},
				ClusterTopology: topologies,
			}},
			want: []*models.ElasticsearchClusterTopologyElement{
				{
					ZoneCount: 2,
					Size: &models.TopologySize{
						Resource: ec.String("memory"),
						Value:    ec.Int32(8192),
					},
					NodeType: &models.ElasticsearchNodeType{
						Data: ec.Bool(true),
					},
				},
				{
					ZoneCount: 1,
					Size: &models.TopologySize{
						Resource: ec.String("memory"),
						Value:    ec.Int32(1024),
					},
					NodeType: &models.ElasticsearchNodeType{
						Master: ec.Bool(true),
					},
				},
				{
					ZoneCount: 1,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildElasticsearchTopology(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("BuildElasticsearchTopology() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				if len(got) > 0 && len(tt.want) > 0 {
					json.NewEncoder(os.Stdout).Encode(got)
					json.NewEncoder(os.Stdout).Encode(tt.want)
				}
				t.Errorf("BuildElasticsearchTopology() = %v, want %v", got, tt.want)
			}
		})
	}
}
