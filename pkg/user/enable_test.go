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
	"github.com/elastic/cloud-sdk-go/pkg/multierror"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"

	"github.com/elastic/ecctl/pkg/util"
)

func TestEnableParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  EnableParams
		wantErr bool
		err     error
	}{
		{
			name:   "validate should return all possible errors",
			params: EnableParams{},
			err: multierror.NewPrefixed("user",
				errors.New("enable requires a username"),
				errors.New("api reference is required for command"),
			),
			wantErr: true,
		},
		{
			name: "validate should return an error if username is empty",
			params: EnableParams{
				API:     &api.API{},
				Enabled: true,
			},
			err: multierror.NewPrefixed("user",
				errors.New("enable requires a username"),
			),
			wantErr: true,
		},
		{
			name: "validate should return an error when api is missing",
			params: EnableParams{
				UserName: "tiburcio",
				Enabled:  false,
			},
			err: multierror.NewPrefixed("user",
				errors.New("api reference is required for command"),
			),
			wantErr: true,
		},
		{
			name: "validate should pass if all params are properly set",
			params: EnableParams{
				UserName: "tiburcio",
				API:      &api.API{},
				Enabled:  true,
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

func TestEnable(t *testing.T) {
	const successResponse = `{
  "builtin": false,
  "security": {
    "enabled": true,
    "roles": [
      "ece_deployment_viewer"
    ]
  },
  "user_name": "tiburcio"
}`

	type args struct {
		params EnableParams
	}
	tests := []struct {
		name    string
		args    args
		want    *models.User
		wantErr bool
		err     error
	}{
		{
			name: "Enable fails due to parameter validation failure (missing API)",
			args: args{
				params: EnableParams{},
			},
			wantErr: true,
			err: multierror.NewPrefixed("user",
				errors.New("enable requires a username"),
				util.ErrAPIReq,
			),
		},
		{
			name: "Enable fails due to API failure",
			args: args{
				params: EnableParams{
					UserName: "tiburcio",
					Enabled:  true,
					API:      api.NewMock(mock.SampleNotFoundError()),
				},
			},
			wantErr: true,
			err:     mock.MultierrorNotFound,
		},
		{
			name: "Enable succeeds",
			args: args{
				params: EnableParams{
					UserName: "tiburcio",
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(successResponse),
						StatusCode: 200,
					}}),
				},
			},
			want: &models.User{
				Builtin:  ec.Bool(false),
				UserName: ec.String("tiburcio"),
				Security: &models.UserSecurity{
					Enabled: ec.Bool(true),
					Roles:   []string{"ece_deployment_viewer"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Enable(tt.args.params)
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
				t.Errorf("Enable() = %v, want %v", got, tt.want)
			}
		})
	}
}
