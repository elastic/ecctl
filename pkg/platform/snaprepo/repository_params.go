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

package snaprepo

import (
	"errors"

	"github.com/elastic/cloud-sdk-go/pkg/api"

	"github.com/elastic/ecctl/pkg/util"
)

var (
	errAPICannotBeNil    = errors.New("api field cannot be nil")
	errNameCannotBeEmpty = errors.New("name field cannot be empty")
	errConfigMustBeSet   = errors.New("config field must be set")
)

// Params is embedded in all of the specific action functions
type Params struct {
	*api.API
}

// Validate ensures that parameters are correct
func (p Params) Validate() error {
	if p.API == nil {
		return errAPICannotBeNil
	}

	return nil
}

// GetParams is used for the Get call
type GetParams struct {
	Params
	Name string
}

// DeleteParams is used for the Delete call
type DeleteParams = GetParams

// Validate ensures that parameters are correct
func (p GetParams) Validate() error {
	if p.Name == "" {
		return errNameCannotBeEmpty
	}

	return p.Params.Validate()
}

// SetParams is used for the Set Call, which will create or update a snapshot
// repository
type SetParams struct {
	Params
	Name   string
	Type   string
	Config util.Validator
}

// Validate ensures that parameters are correct
func (p SetParams) Validate() error {
	if p.Name == "" {
		return errNameCannotBeEmpty
	}

	if p.Config == nil {
		return errConfigMustBeSet
	}

	if err := p.Config.Validate(); err != nil {
		return err
	}

	return p.Params.Validate()
}
