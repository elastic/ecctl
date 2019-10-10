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

package user

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/elastic/ecctl/pkg/util"
)

func TestGetParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  GetParams
		wantErr bool
		err     error
	}{
		{
			name:   "validate should return all possible errors",
			params: GetParams{},
			err: &multierror.Error{
				Errors: []error{
					errors.New("api reference is required for command"),
					errors.New("user: get requires a username"),
				},
			},
			wantErr: true,
		},
		{
			name: "validate should pass if all params are properly set",
			params: GetParams{
				API:      &api.API{},
				UserName: "hermenelgilda",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.err == nil {
				t.Errorf("Validate() expected errors = '%v' but no errors returned", tt.err)
			}

			if tt.wantErr && err.Error() != tt.err.Error() {
				t.Errorf("Validate() expected errors = '%v' but got %v", tt.err, err)
			}
		})
	}
}

func TestGet(t *testing.T) {
	const getUserResponse = `{
    "user_name": "admin",
    "builtin": true
}`

	type args struct {
		params GetParams
	}
	tests := []struct {
		name    string
		args    args
		want    *models.User
		wantErr bool
		err     error
	}{
		{
			name:    "Get fails due to parameter validation failure",
			args:    args{},
			wantErr: true,
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				errors.New("user: get requires a username"),
			}},
		},
		{
			name: "Get fails due to API failure",
			args: args{
				params: GetParams{
					UserName: "hermenelgilda",
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(""),
						StatusCode: 500,
					}}),
				},
			},
			wantErr: true,
			err:     errors.New("unknown error (status 500)"),
		},
		{
			name: "Get succeeds",
			args: args{
				params: GetParams{
					UserName: "admin",
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(getUserResponse),
						StatusCode: 200,
					}}),
				},
			},
			want: &models.User{
				Builtin:  ec.Bool(true),
				UserName: ec.String("admin"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.err != nil && err.Error() != tt.err.Error() {
				t.Errorf("Get() actual error = '%v', want error '%v'", err, tt.err)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCurrentParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  GetCurrentParams
		wantErr bool
		err     error
	}{
		{
			name:    "validate should return all possible errors",
			params:  GetCurrentParams{},
			err:     errors.New("api reference is required for command"),
			wantErr: true,
		},
		{
			name: "validate should pass if all params are properly set",
			params: GetCurrentParams{
				API: &api.API{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.err == nil {
				t.Errorf("Validate() expected errors = '%v' but no errors returned", tt.err)
			}

			if tt.wantErr && err.Error() != tt.err.Error() {
				t.Errorf("Validate() expected errors = '%v' but got %v", tt.err, err)
			}
		})
	}
}

func TestGetCurrent(t *testing.T) {
	const getCurrentResponse = `{
    "user_name": "admin",
    "builtin": true
}`

	const getErrorResponse = `{
  "errors": null
}`

	type args struct {
		params GetCurrentParams
	}
	tests := []struct {
		name    string
		args    args
		want    *models.User
		wantErr bool
		err     error
	}{
		{
			name:    "Get fails due to parameter validation failure",
			args:    args{},
			wantErr: true,
			err:     errors.New("api reference is required for command"),
		},
		{
			name: "Get fails due to API failure",
			args: args{
				params: GetCurrentParams{
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(""),
						StatusCode: 500,
					}}),
				},
			},
			wantErr: true,
			err:     errors.New(getErrorResponse),
		},
		{
			name: "Get succeeds",
			args: args{
				params: GetCurrentParams{
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(getCurrentResponse),
						StatusCode: 200,
					}}),
				},
			},
			want: &models.User{
				Builtin:  ec.Bool(true),
				UserName: ec.String("admin"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCurrent(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCurrent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.err != nil && err.Error() != tt.err.Error() {
				t.Errorf("GetCurrent() actual error = '%v', want error '%v'", err, tt.err)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCurrent() = %v, want %v", got, tt.want)
			}
		})
	}
}
