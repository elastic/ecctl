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

package depresource

import (
	"errors"
	"fmt"
	"io"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/platform/stack"
	"github.com/elastic/ecctl/pkg/util"
)

// LatestStackVersionParams is consumed by LatestStackVersion.
type LatestStackVersionParams struct {
	*api.API

	// When specified, the consuming function will return that version.
	Version string

	// If spceified, it's the io.Writer where info messages are written to.
	Writer io.Writer
}

// Validate ensures the parameters are usable by the consuming function.
func (params LatestStackVersionParams) Validate() error {
	var merr = new(multierror.Error)

	if params.API == nil {
		merr = multierror.Append(merr, util.ErrAPIReq)
	}

	return merr.ErrorOrNil()
}

// LatestStackVersion obtains the latest stack version from the API given an
// empty version from the parameters, if the passed version is not empty, then
// it will be returned.
func LatestStackVersion(params LatestStackVersionParams) (string, error) {
	if err := params.Validate(); err != nil {
		return "", err
	}

	if params.Version != "" {
		return params.Version, nil
	}

	res, err := stack.List(stack.ListParams{
		API: params.API,
	})
	if err != nil {
		return "", errors.New("version discovery: failed to obtain stack list, please specify a version")
	}

	// This check is probably a bit over the top, but you never know.
	if len(res.Stacks) == 0 {
		return "", errors.New("version discovery: stack list is seemingly empty, something is terribly wrong")
	}
	var version = res.Stacks[0].Version

	if params.Writer != nil {
		fmt.Fprintln(params.Writer, "Obtained latest stack version:", version)
	}

	return version, nil
}
