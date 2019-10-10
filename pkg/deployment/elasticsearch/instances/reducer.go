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

package instances

import (
	"github.com/elastic/ecctl/pkg/deployment/elasticsearch"
	"github.com/elastic/ecctl/pkg/util"
)

// List returns a slice with the Elasticsearch instance names in it.
func List(params util.ClusterParams) ([]string, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := elasticsearch.GetCluster(elasticsearch.GetClusterParams{
		ClusterParams: params,
	})
	if err != nil {
		return nil, err
	}

	var instanceSlice = make([]string, 0, len(res.Topology.Instances))
	for _, ins := range res.Topology.Instances {
		instanceSlice = append(instanceSlice, *ins.InstanceName)
	}

	return instanceSlice, nil
}
