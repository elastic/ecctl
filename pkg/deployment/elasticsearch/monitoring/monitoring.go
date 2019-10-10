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

package monitoring

import (
	"fmt"

	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_elasticsearch"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment/deputil"
	"github.com/elastic/ecctl/pkg/util"
)

// EnableParams is consumed by Enable.
type EnableParams struct {
	util.ClusterParams
	TargetID string
}

// Validate ensures the parameters are consumable by Enable.
func (params EnableParams) Validate() error {
	var err = new(multierror.Error)
	if len(params.TargetID) != 32 {
		err = multierror.Append(err,
			fmt.Errorf(`target id "%s" is invalid`, params.TargetID),
		)
	}

	err = multierror.Append(err, deputil.ValidateParams(params))
	return err.ErrorOrNil()
}

// DisableParams is consumed by Disable.
type DisableParams struct {
	util.ClusterParams
}

// Validate ensures the parameters are consumable by Disable.
func (params DisableParams) Validate() error {
	return deputil.ValidateParams(params)
}

// Enable turns monitoring on the cluster specified deployment.
func Enable(params EnableParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	return util.ReturnErrOnly(
		params.V1API.ClustersElasticsearch.SetEsClusterMonitoring(
			clusters_elasticsearch.NewSetEsClusterMonitoringParams().
				WithClusterID(params.ClusterID).
				WithDestClusterID(params.TargetID),
			params.AuthWriter,
		),
	)
}

// Disable disables monitoring on the specified cluster
func Disable(params DisableParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	return util.ReturnErrOnly(
		params.V1API.ClustersElasticsearch.CancelEsClusterMonitoring(
			clusters_elasticsearch.NewCancelEsClusterMonitoringParams().
				WithClusterID(params.ClusterID),
			params.AuthWriter,
		),
	)
}
