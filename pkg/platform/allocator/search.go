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

package allocator

import (
	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/platform_infrastructure"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/go-openapi/strfmt"

	"github.com/elastic/ecctl/pkg/util"
)

// SearchParams contains parameters used to search allocator's data using Query DSL
type SearchParams struct {
	Request models.SearchRequest
	*api.API
}

// Validate is the implementation for the ecctl.Validator interface
func (params SearchParams) Validate() error {
	if params.API == nil {
		return util.ErrAPIReq
	}

	return params.Request.Validate(strfmt.Default)
}

// Search searches all the allocators using Query DSL
func Search(params SearchParams) (*models.AllocatorOverview, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.API.V1API.PlatformInfrastructure.SearchAllocators(
		platform_infrastructure.NewSearchAllocatorsParams().
			WithBody(&params.Request),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}
	return res.Payload, nil
}
