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
	"strings"

	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
)

const (
	// DefaultNodeConfiguration applied on DNT enabled platforms
	DefaultNodeConfiguration = "default"

	// DefaultNodeCountPerZone is the legacy default nodecount per zone
	DefaultNodeCountPerZone = int32(1)
)

// DefaultScriptingSettings defaults to the scripting settings to be used.
var DefaultScriptingSettings = map[string]*models.ElasticsearchScriptingUserSettings{
	"1": nil,
	"2": nil,
	"5": {
		Inline: &models.ElasticsearchScriptTypeSettings{
			Enabled:     ec.Bool(true),
			SandboxMode: ec.Bool(true),
		},
		Stored: &models.ElasticsearchScriptTypeSettings{
			Enabled:     ec.Bool(true),
			SandboxMode: ec.Bool(true),
		},
		File: &models.ElasticsearchScriptTypeSettings{
			Enabled:     ec.Bool(false),
			SandboxMode: ec.Bool(false),
		},
		ExpressionsEnabled: ec.Bool(true),
		MustacheEnabled:    ec.Bool(true),
		PainlessEnabled:    ec.Bool(true),
	},
	"6": nil,
	"7": nil,
}

// LegacyParams represents the settings to create a simple cluster without
// a complex definition or split between data, master and ingest nodes.
type LegacyParams struct {
	Plugins   []string
	Version   string
	ZoneCount int32
	Capacity  int32
}

// newDefaultESPlan creates a new default elasticsearch plan. This will change
// with the addition of DNT. FIXME: Add DNT configurations.
func newDefaultESPlan() *models.ElasticsearchClusterPlan {
	return &models.ElasticsearchClusterPlan{
		Elasticsearch: &models.ElasticsearchConfiguration{},
		Transient: &models.TransientElasticsearchPlanConfiguration{
			Strategy: &models.PlanStrategy{},
		},
		ClusterTopology: []*models.ElasticsearchClusterTopologyElement{
			{
				NodeType: &models.ElasticsearchNodeType{
					Data:   ec.Bool(true),
					Ingest: ec.Bool(true),
					Master: ec.Bool(true),
				},
				NodeCountPerZone: DefaultNodeCountPerZone,
			},
		},
	}
}

// NewLegacyPlan creates a new legacy plan from the parameters.
func NewLegacyPlan(params LegacyParams) *models.ElasticsearchClusterPlan {
	var plan = newDefaultESPlan()

	plan.Elasticsearch.Version, plan.Elasticsearch.EnabledBuiltInPlugins = params.Version, params.Plugins
	plan.ZoneCount, plan.ClusterTopology[0].MemoryPerNode = params.ZoneCount, params.Capacity
	plan.ClusterTopology[0].NodeConfiguration = DefaultNodeConfiguration

	// Set defaults for new cluster
	plan.Elasticsearch.SystemSettings = &models.ElasticsearchSystemSettings{
		Scripting: DefaultScriptingSettings[strings.Split(params.Version, ".")[0]],
	}

	return plan
}
