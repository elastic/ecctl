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

package deputil

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

// NewInvalidDeploymentIDError is returned when a deployment ID is invalid.
func NewInvalidDeploymentIDError(id string) error {
	return fmt.Errorf("id \"%s\" is invalid", id)
}

// NewInvalidResourceTypeError is returned when a deployment ID is invalid.
func NewInvalidResourceTypeError(resourceType string) error {
	return fmt.Errorf(`"%v" is not a valid resource type. Accepted resource types are: %v`,
		resourceType, util.ValidTypes)
}

// Validator wraps the Validate signature
type Validator interface {
	Validate() error
}

// ValidateParams validates that the Validator that is passed complies
// the basic Validate function
func ValidateParams(params Validator) error {
	if params == nil {
		return errors.New("params cannot be nil")
	}

	var valueParams = reflect.ValueOf(params)
	if reflect.Ptr == reflect.ValueOf(params).Kind() {
		valueParams = valueParams.Elem()
	}

	var valid bool
	var merr = new(multierror.Error)
	if field := valueParams.FieldByName("API"); field.IsValid() {
		valid = true
		v, ok := field.Interface().(*api.API)
		if ok && v == nil {
			merr = multierror.Append(merr, util.ErrAPIReq)
		}
		if !ok {
			merr = multierror.Append(merr,
				errors.New("field API must be of type *github.com/elastic/cloud-sdk-go/pkg/api.API"),
			)
		}
	}

	if field := valueParams.FieldByName("ID"); field.IsValid() {
		valid = true
		id, ok := field.Interface().(string)
		if ok && len(id) != 32 {
			merr = multierror.Append(merr, NewInvalidDeploymentIDError(id))
		}
		if !ok {
			merr = multierror.Append(merr, errors.New("field ID must be of type string"))
		}
	}

	if field := valueParams.FieldByName("ClusterID"); field.IsValid() {
		valid = true
		id, ok := field.Interface().(string)
		if ok && len(id) != 32 {
			merr = multierror.Append(merr, NewInvalidDeploymentIDError(id))
		}
		if !ok {
			merr = multierror.Append(merr, errors.New("field ClusterID must be of type string"))
		}
	}

	if !valid {
		return errors.New("params must have one of API or ID")
	}

	return merr.ErrorOrNil()
}

// QueryParams is meant to be embedded in other param structs to
// provide common deployment parameter query settings.
type QueryParams struct {
	// Single part
	ShowPlans        bool
	ShowPlanDefaults bool
	ShowPlanLogs     bool
	ShowMetadata     bool
	ShowSettings     bool

	// List part
	ShowHidden bool
	Query      string
	Size       int64
}
