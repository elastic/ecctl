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

package userauth

import (
	"errors"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/multierror"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"

	"github.com/elastic/ecctl/pkg/util"
)

func TestCreateKey(t *testing.T) {
	invalidPassErrType := mock.APIError{
		Code:    "auth.invalid_password",
		Fields:  []string{"body.password"},
		Message: "request password doesn't match the user's password",
	}

	securityTokenResponse := models.ReAuthenticationResponse{
		SecurityToken: ec.String("uzcyenzalonopalMyxBx"),
	}
	createdAPIKey := models.APIKeyResponse{
		Key: "somekeyvalue",
		ID:  ec.String("somekey"),
	}
	type args struct {
		params CreateKeyParams
	}
	tests := []struct {
		name string
		args args
		want *models.APIKeyResponse
		err  error
	}{
		{
			name: "fails due to parameter validation",
			args: args{},
			err: multierror.NewPrefixed("user auth",
				util.ErrAPIReq,
				errors.New("reauthenticate requires a password"),
				errors.New("create key requires a key description"),
			),
		},
		{
			name: "fails due to reauthenticate API error",
			args: args{params: CreateKeyParams{
				Description: "some description",
				ReAuthenticateParams: ReAuthenticateParams{
					Password: []byte("somepass"),
					API:      api.NewMock(mock.NewErrorResponse(400, invalidPassErrType)),
				},
			}},
			err: multierror.NewPrefixed("api error",
				errors.New("auth.invalid_password: request password doesn't match the user's password (body.password)"),
			),
		},
		{
			name: "fails due to create API error",
			args: args{params: CreateKeyParams{
				Description: "some description",
				ReAuthenticateParams: ReAuthenticateParams{
					Password: []byte("somepass"),
					API: api.NewMock(
						mock.New200Response(mock.NewStructBody(securityTokenResponse)),
						mock.NewErrorResponse(400, invalidPassErrType),
					),
				},
			}},
			err: multierror.NewPrefixed("api error",
				errors.New("auth.invalid_password: request password doesn't match the user's password (body.password)"),
			),
		},
		{
			name: "succeeds",
			args: args{params: CreateKeyParams{
				Description: "some description",
				ReAuthenticateParams: ReAuthenticateParams{
					Password: []byte("somepass"),
					API: api.NewMock(
						mock.New200Response(mock.NewStructBody(securityTokenResponse)),
						mock.New201Response(mock.NewStructBody(createdAPIKey)),
					),
				},
			}},
			want: &createdAPIKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateKey(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("CreateKey() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
