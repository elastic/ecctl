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
	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_apm"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/plan"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment/deputil"
	"github.com/elastic/ecctl/pkg/util"
)

// UpgradeParams are the parameters needed to use Upgrade.
type UpgradeParams struct {
	*api.API
	ID string

	util.TrackParams
}

// Validate ensures that the parameters are usable by the consuming function.
func (params UpgradeParams) Validate() error {
	var err = multierror.Append(new(multierror.Error),
		deputil.ValidateParams(&params),
		params.TrackParams.Validate(),
	)
	return err.ErrorOrNil()
}

// Upgrade bumps the APM instance to match the Elasticsearch deployment one.
func Upgrade(params UpgradeParams) (*models.ClusterUpgradeInfo, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.API.V1API.ClustersApm.UpgradeApm(
		clusters_apm.NewUpgradeApmParams().
			WithClusterID(params.ID),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	if !params.Track {
		return res.Payload, nil
	}

	return res.Payload, util.TrackCluster(util.TrackClusterParams{
		Output: params.Output,
		TrackParams: plan.TrackParams{
			API:           params.API,
			PollFrequency: params.PollFrequency,
			MaxRetries:    params.MaxRetries,
			ID:            params.ID,
			Kind:          "apm",
		},
	})
}
