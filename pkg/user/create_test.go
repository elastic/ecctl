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
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func TestCreateParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  CreateParams
		wantErr bool
		err     error
	}{
		{
			name: "validate should return all possible errors",
			params: CreateParams{
				Email: "hi",
			},
			err: &multierror.Error{
				Errors: []error{
					errors.New("api reference is required for command"),
					errors.New("user: create requires a username"),
					errors.New("user: create requires a password with a minimum of 8 characters"),
					errors.New("user: create requires at least 1 role"),
					errors.New("user: hi is not a valid email address format"),
				},
			},
			wantErr: true,
		},
		{
			name: "validate should return an error when entered password is too short",
			params: CreateParams{
				API:      &api.API{},
				UserName: "bob",
				Password: []byte("pass"),
				Roles:    []string{platformAdminRole},
			},
			err: &multierror.Error{
				Errors: []error{
					errors.New("user: create requires a password with a minimum of 8 characters"),
				},
			},
			wantErr: true,
		},
		{
			name: "validate should return an error when ece_platform_admin is used along other roles",
			params: CreateParams{
				API:      &api.API{},
				UserName: "bob",
				Password: []byte("supersecretpass"),
				Roles:    []string{platformAdminRole, platformViewerRole},
			},
			err: &multierror.Error{
				Errors: []error{
					errors.New("user: ece_platform_admin cannot be used in conjunction with other roles"),
				},
			},
			wantErr: true,
		},
		{
			name: "validate should return an error when ece_platform_admin is used along other roles",
			params: CreateParams{
				API:      &api.API{},
				UserName: "bob",
				Password: []byte("supersecretpass"),
				Roles:    []string{deploymentsManagerRole, deploymentsViewerRole},
			},
			err: &multierror.Error{
				Errors: []error{
					errors.New("user: only one of ece_deployment_manager or ece_deployment_viewer can be chosen"),
				},
			},
			wantErr: true,
		},
		{
			name: "validate should pass if all params are properly set",
			params: CreateParams{
				API:      &api.API{},
				UserName: "bob",
				Email:    "hi@example.com",
				Password: []byte("supersecretpass"),
				Roles:    []string{platformViewerRole, deploymentsManagerRole},
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

func TestCreate(t *testing.T) {
	successResponse := `{
  "builtin": false,
  "security": {
    "enabled": true,
    "roles": [
      "ece_deployment_viewer"
    ]
  },
  "user_name": "bob"
}`

	type args struct {
		params CreateParams
	}
	tests := []struct {
		name    string
		args    args
		want    *models.User
		wantErr bool
		err     error
	}{
		{
			name: "Create fails due to parameter validation failure (missing API)",
			args: args{
				params: CreateParams{
					UserName: "bob",
					Password: []byte("supersecretpass"),
					Roles:    []string{"ece_platform_admin"},
				},
			},
			wantErr: true,
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
			}},
		},
		{
			name: "Create fails due to API failure",
			args: args{
				params: CreateParams{
					UserName: "bob",
					Password: []byte("supersecretpass"),
					Roles:    []string{"ece_platform_admin"},
					API:      api.NewMock(mock.New404Response(mock.NewStringBody(`{"error": "some error"}`))),
				},
			},
			wantErr: true,
			err:     errors.New(`{"error": "some error"}`),
		},
		{
			name: "Create succeeds",
			args: args{
				params: CreateParams{
					UserName: "bob",
					Password: []byte("supersecretpass"),
					Roles:    []string{"ece_deployment_viewer"},
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(successResponse),
						StatusCode: 200,
					}}),
				},
			},
			want: &models.User{
				Builtin:  ec.Bool(false),
				UserName: ec.String("bob"),
				Security: &models.UserSecurity{
					Enabled: ec.Bool(true),
					Roles:   []string{"ece_deployment_viewer"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Create(tt.args.params)
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

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
