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
	"io"

	"github.com/elastic/cloud-sdk-go/pkg/models"
)

// ParseElasticsearchInputParams is consumed by ParseElasticsearchInput.
type ParseElasticsearchInputParams struct {
	NewElasticsearchParams

	Payload          *models.ElasticsearchPayload
	TopologyElements []string
	Size, ZoneCount  int32
	Writer           io.Writer
}

// ParseElasticsearchInput handles all the parameters as optional, providing
// a nicer API when it's used. The bulk of what it does is:
// * If a Payload is already specifide, it returns it.
// * Tries to create an []ElasticsearchTopologyElement from a raw []string.
// * If the previous step returns an empty slice, it uses a default slice which
//   might override the values when Size or ZoneCount are set in the params.
// * Auto-discovers the latest Stack version if Version is not specified.
// When all of those steps are done, it finally calls NewElasticsearch building
// the resulting ElasticsearchPayload.
func ParseElasticsearchInput(params ParseElasticsearchInputParams) (*models.ElasticsearchPayload, error) {
	if params.Payload != nil {
		return params.Payload, nil
	}

	topology, err := NewElasticsearchTopology(params.TopologyElements)
	if err != nil {
		return nil, err
	}

	// On empty topology, use the default one with the size & count specified
	if len(topology) == 0 {
		topology = append(topology, NewElasticsearchTopologyElement(
			params.Size, params.ZoneCount,
		))
	}

	// Version Discovery
	version, err := LatestStackVersion(LatestStackVersionParams{
		Writer:  params.Writer,
		API:     params.API,
		Version: params.Version,
	})
	if err != nil {
		return nil, err
	}

	var NewEsparams = params.NewElasticsearchParams
	NewEsparams.Version = version
	NewEsparams.Topology = topology

	return NewElasticsearch(NewEsparams)
}
