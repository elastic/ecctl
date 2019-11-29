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
	"fmt"

	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/hashicorp/go-multierror"
)

const (
	// DataNode identifies the node type which stores data.
	DataNode = "data"

	// MasterNode identifies the node type which is master elegible.
	MasterNode = "master"

	// MLNode identifies the node type performing the Machine Learning.
	MLNode = "ml"
)

var (
	// DefaultTopology will be used when no topology is specified.
	DefaultTopology = []ElasticsearchTopologyElement{
		DefaultTopologyElement,
	}

	// DefaultTopologyElement defines the element used in DefaultTopology
	DefaultTopologyElement = ElasticsearchTopologyElement{
		Name:      DataNode,
		Size:      DefaultDataSize,
		ZoneCount: DefaultDataZoneCount,
	}
)

// BuildElasticsearchTopologyParams is consumed by BuildElasticsearchTopology.
type BuildElasticsearchTopologyParams struct {
	// Deployment Template ID from which the ClusterTopology comes from.
	TemplateID string

	// API-obtained slice of the Elasticsearch Deployment topology.
	ClusterTopology []*models.ElasticsearchClusterTopologyElement

	// User specified desired topology from which a new topology will be built.
	Topology []ElasticsearchTopologyElement
}

// ElasticsearchTopologyElement is a single cluster topology element, meaning
// a number of instances (controlled by ZoneCount) for a single NodeType.
type ElasticsearchTopologyElement struct {
	// Name can be one of "data", "master" or "ml".
	Name string `json:"name"`

	// Number of zones to span the cluster on.
	ZoneCount int32 `json:"zone_count,omitempty"`

	// Memory size of the cluster node.
	Size int32 `json:"size"`
}

// Sets ZoneCount to a DefaultDataZoneCount.
func (element *ElasticsearchTopologyElement) fillDefaults() {
	if element.ZoneCount == 0 {
		element.ZoneCount = DefaultDataZoneCount
	}
}

// Validate ensures the parameters are usable by the consuming function.
func (element *ElasticsearchTopologyElement) Validate() error {
	var merr = new(multierror.Error)

	if element.Name == "" {
		merr = multierror.Append(merr, errors.New("deployment topology: name cannot be empty"))
	}

	if element.Size == 0 {
		merr = multierror.Append(merr, errors.New("deployment topology: size cannot be empty"))
	}

	return merr.ErrorOrNil()
}

// NewElasticsearchTopology creates a []ElasticsearchTopologyElement from a
// slice of raw strings which are then unmarshaled into the desired slice type.
// If any of the topology elements is not usable, an error is returned.
func NewElasticsearchTopology(topology []string) ([]ElasticsearchTopologyElement, error) {
	var t = make([]ElasticsearchTopologyElement, 0, len(topology))
	for _, rawElement := range topology {
		var element ElasticsearchTopologyElement
		element.fillDefaults()
		if err := json.Unmarshal([]byte(rawElement), &element); err != nil {
			return nil, fmt.Errorf("depresource: failed unpacking raw topology: %s", err)
		}

		if err := element.Validate(); err != nil {
			return nil, err
		}

		t = append(t, element)
	}
	return t, nil
}

// NewElasticsearchTopologyElement creates a new topology element given a zone
// count and node size. Using DefaultTopologyElement as the blueprint.
func NewElasticsearchTopologyElement(size, zoneCount int32) ElasticsearchTopologyElement {
	var element = DefaultTopologyElement
	if size > 0 {
		element.Size = size
	}

	if zoneCount > 0 {
		element.ZoneCount = zoneCount
	}

	return element
}

// BuildElasticsearchTopology receives an ElasticsearchCluster topology from
// the deployment template definition and what the user has specified through
// a simplified version of the deployment topology, iterating over the template
// topology and triying to match the types that the user can specify: ["ml",
// "data", "master"], with their instance configuration ID. the matchNodeType
// function takes care of the matching.
// Once the NodeType has beeen matched, any overrides coming from the user-
// specified topology settings are set, overriding the default deployment
// template defined defaults.
func BuildElasticsearchTopology(params BuildElasticsearchTopologyParams) ([]*models.ElasticsearchClusterTopologyElement, error) {
	var topologyList []*models.ElasticsearchClusterTopologyElement
	for _, desired := range params.Topology {
		for _, t := range params.ClusterTopology {
			if matchNodeType(*t.NodeType, desired) {
				// Override the desired topology if values are non zero
				if desired.Size > 0 {
					t.Size.Value = ec.Int32(desired.Size)
				}
				if desired.ZoneCount > 0 {
					t.ZoneCount = desired.ZoneCount
				}

				topologyList = append(topologyList, t)
			}
		}
	}

	if len(topologyList) == 0 {
		return nil, fmt.Errorf(
			"deployment topology: failed to obtain desired topology names (%+v) in deployment template id \"%s\"",
			params.Topology, params.TemplateID,
		)
	}

	return topologyList, nil
}

// matchNodeType compares  ElasticsearchTopologyElement name (NodeType) to the
// actual NodeTypes specified in a deployment template cluster topology. The
// Name field can be ["data", "master", "ml"].
func matchNodeType(got models.ElasticsearchNodeType, want ElasticsearchTopologyElement) bool {
	if want.Name == DataNode {
		return got.Data != nil && *got.Data
	}

	if want.Name == MasterNode {
		var dataFalse = (got.Data != nil && !*got.Data) || got.Data == nil
		var masterTrue = got.Master != nil && *got.Master
		return dataFalse && masterTrue
	}

	if want.Name == MLNode {
		var dataFalse = (got.Data != nil && !*got.Data) || got.Data == nil
		var mlTrue = got.Ml != nil && *got.Ml
		return dataFalse && mlTrue
	}

	return false
}
