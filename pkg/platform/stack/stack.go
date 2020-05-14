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
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/stack"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/multierror"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/go-openapi/runtime"

	"github.com/elastic/ecctl/pkg/util"
)

// Get obtains a stackpack to the current installation
func Get(params GetParams) (*models.StackVersionConfig, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.API.V1API.Stack.GetVersionStack(
		stack.NewGetVersionStackParams().
			WithVersion(params.Version),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}

// List lists all stackpacks in the current installation
func List(params ListParams) (*models.StackVersionConfigs, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.API.V1API.Stack.GetVersionStacks(
		stack.NewGetVersionStacksParams().
			WithShowDeleted(ec.Bool(params.Deleted)),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	stacks := res.Payload.Stacks
	sort.Slice(stacks, func(i, j int) bool {
		return compareVersions(stacks[i].Version, stacks[j].Version)
	})

	return res.Payload, nil
}

// Upload uploads a stackpack from a location
func Upload(params UploadParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	res, err := params.V1API.Stack.UpdateStackPacks(
		stack.NewUpdateStackPacksParams().
			WithFile(runtime.NamedReader("StackPack", params.StackPack)),
		params.AuthWriter,
	)
	if err != nil {
		return api.UnwrapError(err)
	}

	var merr = multierror.NewPrefixed("stack upload")
	for _, e := range res.Payload.Errors {
		for _, ee := range e.Errors.Errors {
			// ECE stack packs seem to have a __MACOSX packed file which is
			// causing the command to return an error. Error thrown is:
			// This version cannot be parsed: [__MACOSX] because:
			// Unknown version string: [__MACOSX]
			if !strings.Contains(*ee.Message, "__MACOSX") {
				merr = merr.Append(fmt.Errorf("%s: %s", *ee.Code, *ee.Message))
			}
		}
	}
	return merr.ErrorOrNil()
}

// Delete deletes a stackpack
func Delete(params DeleteParams) error {
	if err := params.Validate(); err != nil {
		return err
	}
	return util.ReturnErrOnly(
		params.API.V1API.Stack.DeleteVersionStack(
			stack.NewDeleteVersionStackParams().
				WithVersion(params.Version),
			params.AuthWriter,
		),
	)
}

func compareVersions(version1, version2 string) bool {
	v1 := strings.Split(version1, ".")
	major1 := v1[0]
	minor1 := v1[1]
	patch1 := v1[2]
	p1, err1 := strconv.Atoi(patch1)

	v2 := strings.Split(version2, ".")
	major2 := v2[0]
	minor2 := v2[1]
	patch2 := v2[2]
	p2, err2 := strconv.Atoi(patch2)

	var patchComp = patch1 > patch2
	if err1 == nil && err2 == nil {
		patchComp = p1 > p2
	}

	return major1 > major2 ||
		(major1 == major2 && minor1 > minor2) ||
		(major1 == major2 && minor1 == minor2 && patchComp)
}
