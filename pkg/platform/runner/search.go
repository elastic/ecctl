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

package runner

import (
	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/platform_infrastructure"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/multierror"
	"github.com/go-openapi/strfmt"
)

// SearchParams contains parameters used to search runner's data using Query DSL
type SearchParams struct {
	Params
	Request models.SearchRequest
}

// Validate is the implementation for the ecctl.Validator interface
func (params SearchParams) Validate() error {
	var merr = multierror.NewPrefixed("runner search")

	merr = merr.Append(params.Params.Validate())

	merr = merr.Append(params.Request.Validate(strfmt.Default))

	return merr.ErrorOrNil()
}

// Search searches all the runners using Query DSL
func Search(params SearchParams) (*models.RunnerOverview, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.API.V1API.PlatformInfrastructure.SearchRunners(
		platform_infrastructure.NewSearchRunnersParams().
			WithBody(&params.Request),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}
	return res.Payload, nil
}
