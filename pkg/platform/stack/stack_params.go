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

package stack

import (
	"errors"
	"io"
	"strings"

	"github.com/blang/semver"
	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/multierror"

	"github.com/elastic/ecctl/pkg/util"
)

// GetParams is consumed by Get
type GetParams struct {
	*api.API
	Version string
}

// Validate ensures that the parameters are usable by the consuming
// function
func (params GetParams) Validate() error {
	var merr = multierror.NewPrefixed("stack get")
	if params.API == nil {
		merr = merr.Append(util.ErrAPIReq)
	}

	if _, e := semver.Parse(params.Version); e != nil {
		merr = merr.Append(errors.New(strings.ToLower(e.Error())))
	}

	return merr.ErrorOrNil()
}

// ListParams is consumed by List
type ListParams struct {
	*api.API
	Deleted bool
}

// Validate ensures that the parameters are usable by the consuming
// function
func (params ListParams) Validate() error {
	if params.API == nil {
		return util.ErrAPIReq
	}

	return nil
}

// UploadParams is consumed by Upload
type UploadParams struct {
	*api.API
	StackPack io.Reader
}

// Validate ensures that the parameters are usable by the consuming
// function
func (params UploadParams) Validate() error {
	var merr = multierror.NewPrefixed("stack upload")
	if params.API == nil {
		merr = merr.Append(util.ErrAPIReq)
	}

	if params.StackPack == nil {
		merr = merr.Append(errors.New("stackpack cannot be empty"))
	}

	return merr.ErrorOrNil()
}

// DeleteParams is consumed by Delete
type DeleteParams struct {
	*api.API
	Version string
}

// Validate ensures that the parameters are usable by the consuming
// function
func (params DeleteParams) Validate() error {
	var merr = multierror.NewPrefixed("stack delete")
	if params.API == nil {
		merr = merr.Append(util.ErrAPIReq)
	}

	if _, e := semver.Parse(params.Version); e != nil {
		merr = merr.Append(errors.New(strings.ToLower(e.Error())))
	}

	return merr.ErrorOrNil()
}
