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
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	multierror "github.com/hashicorp/go-multierror"
)

type incorrectParams struct{ Some string }

func (params incorrectParams) Validate() error { return nil }

type idParams struct{ ID string }

func (params idParams) Validate() error { return nil }

type idInvalidTypeParams struct{ ID bool }

func (params idInvalidTypeParams) Validate() error { return nil }

type apiParams struct{ *api.API }

func (params apiParams) Validate() error { return nil }

type apiInvalidTypeParams struct{ API bool }

func (params apiInvalidTypeParams) Validate() error { return nil }

type correctParams struct {
	*api.API
	ID string
}

func (params correctParams) Validate() error { return nil }

type alternativeCorrectParams struct {
	*api.API
	ClusterID string
}

func (params alternativeCorrectParams) Validate() error { return nil }

func TestValidateBasicParams(t *testing.T) {
	type args struct {
		params Validator
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Validate fails when params is nil",
			args: args{},
			err:  errors.New("params cannot be nil"),
		},
		{
			name: "Validate fails when params are missing the API or ID field(s)",
			args: args{
				params: &incorrectParams{Some: "a"},
			},
			err: errors.New("params must have one of API or ID"),
		},
		{
			name: "Validate fails when params are missing a correct ID",
			args: args{
				params: &idParams{ID: "someValue"},
			},
			err: &multierror.Error{Errors: []error{
				errors.New(`id "someValue" is invalid`),
			}},
		},
		{
			name: "Validate fails when params are missing a correct ID",
			args: args{
				params: &idParams{},
			},
			err: &multierror.Error{Errors: []error{
				errors.New(`id "" is invalid`),
			}},
		},
		{
			name: "Validate fails when params are missing an API",
			args: args{
				params: &idParams{},
			},
			err: &multierror.Error{Errors: []error{
				errors.New(`id "" is invalid`),
			}},
		},
		{
			name: "Validate fails when params have an ID type but it's of invalid type",
			args: args{
				params: &idInvalidTypeParams{ID: true},
			},
			err: &multierror.Error{Errors: []error{
				errors.New("field ID must be of type string"),
			}},
		},
		{
			name: "Validate fails when params are missing an API",
			args: args{
				params: &apiParams{},
			},
			err: &multierror.Error{Errors: []error{
				errors.New(`api reference is required for command`),
			}},
		},
		{
			name: "Validate fails when params have an API type but it's of invalid type",
			args: args{
				params: &apiInvalidTypeParams{API: true},
			},
			err: &multierror.Error{Errors: []error{
				errors.New("field API must be of type *github.com/elastic/cloud-sdk-go/pkg/api.API"),
			}},
		},
		{
			name: "Validate succeeds",
			args: args{params: &correctParams{
				ID:  "269d4dc158f4457491572a90ac124e6f",
				API: new(api.API),
			}},
		},
		{
			name: "Validate alternative struct succeeds",
			args: args{params: &alternativeCorrectParams{
				ClusterID: "269d4dc158f4457491572a90ac124e6f",
				API:       new(api.API),
			}},
		},
		{
			name: "Validate non pointer alternative struct succeeds",
			args: args{params: alternativeCorrectParams{
				ClusterID: "269d4dc158f4457491572a90ac124e6f",
				API:       new(api.API),
			}},
		},
		{
			name: "Validate fails when params are missing an API",
			args: args{
				params: &alternativeCorrectParams{},
			},
			err: &multierror.Error{Errors: []error{
				errors.New(`api reference is required for command`),
				errors.New(`id "" is invalid`),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateParams(tt.args.params); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("ValidateBasicParams() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
