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

package deployment

import (
	"errors"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/deployments"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

// CreateParams is consumed by Create.
type CreateParams struct {
	*api.API

	Request *models.DeploymentCreateRequest

	RequestID string

	// deployment Overrides
	Overrides *CreateOverrides
}

// CreateOverrides represent the override settings to
type CreateOverrides struct {
	// If set, it will override the deployment name.
	Name string

	// If set, it will override the region when not present in the
	// DeploymentCreateRequest.
	// Note this behaviour is different from the rest of overrides
	// since this field tends to be populated by the global region
	// field which is implicit (by config) rather than explicit by flag.
	Region string

	// If set, it'll override all versions to match this one.
	Version string
}

// Validate ensures the parameters are usable by Shutdown.
func (params CreateParams) Validate() error {
	var merr = new(multierror.Error)

	if params.API == nil {
		merr = multierror.Append(merr, util.ErrAPIReq)
	}

	if params.Request == nil {
		merr = multierror.Append(merr, errors.New("deployment create: request payload cannot be empty"))
	}

	return merr.ErrorOrNil()
}

// Create performs a Create using the specified Request against the API.
func Create(params CreateParams) (*models.DeploymentCreateResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	setOverrides(params.Request, params.Overrides)

	var id *string
	if params.RequestID != "" {
		id = &params.RequestID
	}

	res, res2, err := params.V1API.Deployments.CreateDeployment(
		deployments.NewCreateDeploymentParams().
			WithRequestID(id).
			WithBody(params.Request),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	if res == nil {
		return res2.Payload, nil
	}

	return res.Payload, nil
}

// setOverrides sets a series of overrides
// nolint
func setOverrides(req *models.DeploymentCreateRequest, overrides *CreateOverrides) {
	if req == nil || overrides == nil || req.Resources == nil {
		return
	}

	if overrides.Name != "" {
		req.Name = overrides.Name
	}

	for _, resource := range req.Resources.Apm {
		if resource.Region == nil && overrides.Region != "" {
			resource.Region = &overrides.Region
		}

		if overrides.Version != "" {
			if resource.Plan != nil && resource.Plan.Apm != nil {
				resource.Plan.Apm.Version = overrides.Version
			}
		}
	}

	for _, resource := range req.Resources.Appsearch {
		if resource.Region == nil && overrides.Region != "" {
			resource.Region = &overrides.Region
		}
		if overrides.Version != "" {
			if resource.Plan != nil && resource.Plan.Appsearch != nil {
				resource.Plan.Appsearch.Version = overrides.Version
			}
		}
	}

	for _, resource := range req.Resources.Elasticsearch {
		if resource.Region == nil && overrides.Region != "" {
			resource.Region = &overrides.Region
		}
		if overrides.Version != "" {
			if resource.Plan != nil && resource.Plan.Elasticsearch != nil {
				resource.Plan.Elasticsearch.Version = overrides.Version
			}
		}
	}

	for _, resource := range req.Resources.Kibana {
		if resource.Region == nil && overrides.Region != "" {
			resource.Region = &overrides.Region
		}
		if overrides.Version != "" {
			if resource.Plan != nil && resource.Plan.Kibana != nil {
				resource.Plan.Kibana.Version = overrides.Version
			}
		}
	}
}
