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
	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_elasticsearch"
	"github.com/elastic/cloud-sdk-go/pkg/plan"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

// ShutdownClusterParams used to shut down an elasticsearch cluster.
type ShutdownClusterParams struct {
	util.ClusterParams
	util.TrackParams

	// Hides the cluster from the user console if specified.
	Hide         bool
	SkipSnapshot bool
}

// Validate ensures the parameters are usable by ShutdownCluster.
func (params ShutdownClusterParams) Validate() error {
	var err = new(multierror.Error)

	err = multierror.Append(err, params.ClusterParams.Validate())
	err = multierror.Append(err, params.TrackParams.Validate())

	return err.ErrorOrNil()
}

// ShutdownCluster shuts down an Elasticsearch cluster. If hide is true
// it will also hide the cluster in the UI
func ShutdownCluster(params ShutdownClusterParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	if _, err := params.V1API.ClustersElasticsearch.ShutdownEsCluster(
		clusters_elasticsearch.NewShutdownEsClusterParams().
			WithClusterID(params.ClusterID).
			WithHide(&params.Hide).
			WithSkipSnapshot(&params.SkipSnapshot),
		params.AuthWriter,
	); err != nil {
		return api.UnwrapError(err)
	}

	if !params.Track {
		return nil
	}

	return util.TrackCluster(util.TrackClusterParams{
		Output: params.Output,
		TrackParams: plan.TrackParams{
			API:           params.API,
			ID:            params.ClusterID,
			PollFrequency: params.PollFrequency,
			MaxRetries:    params.MaxRetries,
			Kind:          "elasticsearch",
		},
	})
}
