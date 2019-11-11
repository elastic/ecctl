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

// SearchParams is consumed by Search.
type SearchParams struct {
	*api.API

	Request *models.SearchRequest
}

// Validate ensures the parameters are usable by Shutdown.
func (params SearchParams) Validate() error {
	var merr = new(multierror.Error)

	if params.API == nil {
		merr = multierror.Append(merr, util.ErrAPIReq)
	}

	if params.Request == nil {
		merr = multierror.Append(merr, errors.New("deployment search: request cannot be empty"))
	}

	return merr.ErrorOrNil()
}

// Search performs a search using the specified Request against the API.
func Search(params SearchParams) (*models.DeploymentsSearchResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.V1API.Deployments.SearchDeployments(
		deployments.NewSearchDeploymentsParams().
			WithBody(params.Request),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}
