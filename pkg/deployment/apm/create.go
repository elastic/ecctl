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
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/go-openapi/strfmt"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment/deputil"
	"github.com/elastic/ecctl/pkg/util"
)

// CreateParams are consumed by the Create function.
type CreateParams struct {
	*api.API

	// ID of the Deployment on which to link it
	ID string
	models.CreateApmRequest
	util.TrackParams
}

// Validate ensures that the parameters are usable by the consuming function.
func (params *CreateParams) Validate() error {
	var err = multierror.Append(new(multierror.Error),
		deputil.ValidateParams(params),
		params.TrackParams.Validate(),
		params.CreateApmRequest.Validate(strfmt.Default),
	)
	return err.ErrorOrNil()
}

// Fill populates any optional fields that have a default value.
func (params *CreateParams) Fill() {
	if params.ID != "" {
		params.CreateApmRequest.ElasticsearchClusterID = ec.String(params.ID)
	}
}

// Create creates a new APM Cluster from a definition
func Create(params CreateParams) (*models.ApmCrudResponse, error) {
	params.Fill()

	if err := params.Validate(); err != nil {
		return nil, err
	}

	_, res, err := params.V1API.ClustersApm.CreateApm(
		clusters_apm.NewCreateApmParams().
			WithBody(&params.CreateApmRequest),
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
			ID:            res.Payload.ApmID,
			Kind:          "apm",
		},
	})
}
