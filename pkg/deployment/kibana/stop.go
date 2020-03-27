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
	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_kibana"
	"github.com/elastic/cloud-sdk-go/pkg/plan/planutil"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

// StopParams is used by the Stop function
type StopParams struct {
	Hide bool
	DeploymentParams
}

// Validate ensures that the parameters are usable by the consuming function.
func (params StopParams) Validate() error {
	var merr = new(multierror.Error)

	merr = multierror.Append(merr, params.DeploymentParams.Validate())

	return merr.ErrorOrNil()
}

// Stop stops the cluster.
func Stop(params StopParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	_, err := params.API.V1API.ClustersKibana.ShutdownKibanaCluster(
		clusters_kibana.NewShutdownKibanaClusterParams().
			WithClusterID(params.ID).
			WithHide(&params.Hide),
		params.API.AuthWriter,
	)
	if err != nil {
		return api.UnwrapError(err)
	}

	if !params.Track {
		return nil
	}

	return planutil.TrackChange(util.SetClusterTracking(
		params.TrackChangeParams, params.ID, util.Kibana,
	))
}
