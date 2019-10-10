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

package elasticsearch

import (
	"errors"

	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_elasticsearch"

	"github.com/elastic/ecctl/pkg/util"
)

// stoppedState represents the stopped state string for a cluster
const stoppedState = "stopped"

// DeleteClusterParams ensures the parameters are usable by DeleteCluster.
type DeleteClusterParams struct {
	util.ClusterParams
}

// DeleteCluster deletes an Elasticsearch cluster.
func DeleteCluster(params DeleteClusterParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	res, err := GetCluster(GetClusterParams{
		ClusterParams: params.ClusterParams,
	})
	if err != nil {
		return err
	}

	// Even though the cluster state can change in between querying the state
	// and performing the shutdown, it's not that relevant and thus accepted.
	if res.Status != stoppedState {
		return errors.New("elasticsearch delete: deployment must be stopped")
	}

	return util.ReturnErrOnly(
		params.V1API.ClustersElasticsearch.DeleteEsCluster(
			clusters_elasticsearch.NewDeleteEsClusterParams().
				WithClusterID(params.ClusterID),
			params.AuthWriter,
		),
	)
}
