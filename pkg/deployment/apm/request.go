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
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
)

const (
	// DefaultZoneCount is the number of zones that a Apm has by default.
	DefaultZoneCount = 1
	// DefaultMemoryPerNode is the default memory per node.
	DefaultMemoryPerNode = 512
)

// NewApmBodyParams is used by NewApmBody.
type NewApmBodyParams struct {
	ID            string
	Name          string
	ZoneCount     int32
	MemoryPerNode int32
}

// Fill sets any unset values to its default ones.
func (params *NewApmBodyParams) Fill() {
	if params.ZoneCount == 0 {
		params.ZoneCount = DefaultZoneCount
	}
	if params.MemoryPerNode == 0 {
		params.MemoryPerNode = DefaultMemoryPerNode
	}
}

// NewApmBody constructs a apm body for simple cases, abstracting the
// structure's complexity away from the user.
func NewApmBody(params NewApmBodyParams) *models.CreateApmRequest {
	params.Fill()

	return &models.CreateApmRequest{
		DisplayName:            params.Name,
		ElasticsearchClusterID: ec.String(params.ID),
		Plan: &models.ApmPlan{
			Apm: new(models.ApmConfiguration),
			ClusterTopology: []*models.ApmTopologyElement{
				{
					ZoneCount: params.ZoneCount,
					Size: &models.TopologySize{
						Resource: ec.String("memory"),
						Value:    ec.Int32(params.MemoryPerNode),
					},
				},
			},
		},
	}
}
