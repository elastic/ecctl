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
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
)

const (
	// DefaultZoneCount is the number of zones that a Kibana has by default.
	DefaultZoneCount = 1
	// DefaultMemoryPerNode is the default memory per node.
	DefaultMemoryPerNode = 1024
	// DefaultNodeCountPerZone is the default node count per zone set.
	DefaultNodeCountPerZone = 1
)

// NewKibanaBodyParams is used by NewKibanaBody.
type NewKibanaBodyParams struct {
	ID               string
	Name             string
	ZoneCount        int32
	MemoryPerNode    int32
	NodeCountPerZone int32
}

// Fill sets any unset values to its default ones.
func (params *NewKibanaBodyParams) Fill() {
	if params.ZoneCount == 0 {
		params.ZoneCount = DefaultZoneCount
	}
	if params.MemoryPerNode == 0 {
		params.MemoryPerNode = DefaultMemoryPerNode
	}
	if params.NodeCountPerZone == 0 {
		params.NodeCountPerZone = DefaultNodeCountPerZone
	}
}

// NewKibanaBody constructs a kibana body for simple cases, abstracting the
// structure's complexity away from the user.
func NewKibanaBody(params NewKibanaBodyParams) *models.CreateKibanaRequest {
	params.Fill()

	return &models.CreateKibanaRequest{
		ClusterName:            params.Name,
		ElasticsearchClusterID: ec.String(params.ID),
		Plan: &models.KibanaClusterPlan{
			ZoneCount: params.ZoneCount,
			Kibana:    new(models.KibanaConfiguration),
			ClusterTopology: []*models.KibanaClusterTopologyElement{
				{
					MemoryPerNode:    params.MemoryPerNode,
					NodeCountPerZone: params.NodeCountPerZone,
				},
			},
		},
	}
}
