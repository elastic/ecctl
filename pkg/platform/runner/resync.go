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

	"github.com/elastic/ecctl/pkg/util"
)

// ResyncParams is consumed by Resync
type ResyncParams struct {
	Params
	ID string
}

// Validate ensures the parameters are usable by the consuming function.
func (params ResyncParams) Validate() error {
	var merr = multierror.NewPrefixed("runner resync")
	if params.ID == "" {
		merr = merr.Append(util.ErrIDCannotBeEmpty)
	}

	merr = merr.Append(params.Params.Validate())

	return merr.ErrorOrNil()
}

// Resync forces indexer to immediately resynchronize the search index
// and cache for a given runner.
func Resync(params ResyncParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	return util.ReturnErrOnly(
		params.API.V1API.PlatformInfrastructure.ResyncRunner(
			platform_infrastructure.NewResyncRunnerParams().
				WithRunnerID(params.ID),
			params.API.AuthWriter,
		),
	)
}

// ResyncAll asynchronously resynchronizes the search index for all runners.
func ResyncAll(params Params) (*models.ModelVersionIndexSynchronizationResults, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.API.V1API.PlatformInfrastructure.ResyncRunners(
		platform_infrastructure.NewResyncRunnersParams(),
		params.API.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}
