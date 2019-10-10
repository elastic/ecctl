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
	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_elasticsearch"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/plan"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment/deputil"
	"github.com/elastic/ecctl/pkg/util"
)

// UpdateParams is consumed by Update
type UpdateParams struct {
	*api.API
	ID           string
	ValidateOnly bool
	Plan         models.ElasticsearchClusterPlan

	util.TrackParams
}

// Validate ensures that the parameters are usable by the consuming function.
func (params UpdateParams) Validate() error {
	var err = multierror.Append(new(multierror.Error),
		deputil.ValidateParams(&params),
		params.TrackParams.Validate(),
	)
	return err.ErrorOrNil()
}

// Update applies (or validates) the given plan
func Update(params UpdateParams) (*models.ClusterCrudResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	planOk, planAccepted, err := params.API.V1API.ClustersElasticsearch.UpdateEsClusterPlan(
		clusters_elasticsearch.NewUpdateEsClusterPlanParams().
			WithClusterID(params.ID).
			WithValidateOnly(&params.ValidateOnly).
			WithBody(&params.Plan),
		params.AuthWriter,
	)

	if err != nil {
		return nil, err
	}

	if params.ValidateOnly {
		return planOk.Payload, nil
	}

	if !params.Track {
		return planAccepted.Payload, nil
	}

	return planAccepted.Payload, util.TrackCluster(util.TrackClusterParams{
		Output: params.Output,
		TrackParams: plan.TrackParams{
			API:           params.API,
			PollFrequency: params.PollFrequency,
			MaxRetries:    params.MaxRetries,
			ID:            params.ID,
			Kind:          "elasticsearch",
		},
	})
}
