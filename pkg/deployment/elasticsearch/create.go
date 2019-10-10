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

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_elasticsearch"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment/elasticsearch/plan"
	"github.com/elastic/ecctl/pkg/util"
)

// CreateParams is consumed by Create.
type CreateParams struct {
	*api.API
	util.TrackParams

	// Optional settings

	// Plan defines the Elasticsearch cluster plan definition.
	PlanDefinition *models.ElasticsearchClusterPlan

	// ClusterName is the optional name to give to the cluster.
	ClusterName string

	// LegacyPlanParams are used when no Plan definition is passed
	plan.LegacyParams
}

// Validate ensures the parameters are usable by Create.
func (params *CreateParams) Validate() error {
	var err = new(multierror.Error)
	if params.API == nil {
		err = multierror.Append(err, util.ErrAPIReq)
	}

	var legacyCapacity = params.ZoneCount > 0 || params.Capacity > 0
	if params.PlanDefinition != nil && legacyCapacity {
		err = multierror.Append(err, errors.New(
			"cannot specify a plan definition when capacity or zonecount are set",
		))
	}

	err = multierror.Append(err, params.TrackParams.Validate())
	return err.ErrorOrNil()
}

func (params *CreateParams) fillValues() {
	if params.PlanDefinition == nil {
		params.PlanDefinition = plan.NewLegacyPlan(params.LegacyParams)
	}

	if params.LegacyParams.Plugins != nil {
		params.PlanDefinition.Elasticsearch.EnabledBuiltInPlugins = params.Plugins
	}

	if params.LegacyParams.Version != "" {
		params.PlanDefinition.Elasticsearch.Version = params.Version
	}
}

// Create Creates a new Elasticsearch cluster.
func Create(params CreateParams) (*models.ClusterCrudResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}
	params.fillValues()

	resOK, resAccepted, err := params.V1API.ClustersElasticsearch.CreateEsCluster(
		clusters_elasticsearch.NewCreateEsClusterParams().
			WithBody(&models.CreateElasticsearchClusterRequest{
				ClusterName: params.ClusterName,
				Plan:        params.PlanDefinition,
			}),
		params.AuthWriter,
	)

	return util.ParseCUResponse(util.ParseCUResponseParams{
		API:            params.API,
		Err:            err,
		UpdateResponse: resAccepted,
		CreateResponse: resOK,
		TrackParams:    params.TrackParams,
	})
}
