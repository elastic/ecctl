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

package enrollmenttoken

import (
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"

	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func TestCreate(t *testing.T) {
	type args struct {
		params CreateParams
	}
	tests := []struct {
		name string
		args args
		want *models.RequestEnrollmentTokenReply
		err  error
	}{
		{
			name: "Create fails due to incorrect duration",
			args: args{params: CreateParams{
				API:      new(api.API),
				Duration: time.Hour * 999999,
			}},
			err: &multierror.Error{Errors: []error{
				errors.New("validity value 3599996400 exceeds max allowed 2147483647 value in seconds"),
			}},
		},
		{
			name: "Create fails due to missing API",
			args: args{params: CreateParams{}},
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
			}},
		},
		{
			name: "Create Succeeds with persistent token",
			args: args{params: CreateParams{
				API: api.NewMock(mock.Response{
					Response: http.Response{
						StatusCode: 200,
						Body: mock.NewStructBody(models.RequestEnrollmentTokenReply{
							Token:   ec.String("some token"),
							TokenID: "some-token-id",
						}),
					},
				}),
			}},
			want: &models.RequestEnrollmentTokenReply{
				Token:   ec.String("some token"),
				TokenID: "some-token-id",
			},
		},
		{
			name: "Create fails due to API error",
			args: args{params: CreateParams{
				API: api.NewMock(mock.Response{
					Error: errors.New("error"),
				}),
			}},
			err: &url.Error{
				Op:  "Post",
				URL: "https://mock-host/mock-path/platform/configuration/security/enrollment-tokens",
				Err: errors.New("error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Create(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type args struct {
		params DeleteParams
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Create fails due to missing token",
			args: args{params: DeleteParams{
				API: new(api.API),
			}},
			err: &multierror.Error{Errors: []error{
				errors.New("token cannot be empty"),
			}},
		},
		{
			name: "Create fails due to missing API",
			args: args{params: DeleteParams{
				Token: "token",
			}},
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
			}},
		},
		{
			name: "Delete fails due to API error",
			args: args{params: DeleteParams{
				Token: "atoken",
				API: api.NewMock(mock.Response{
					Error: errors.New("error"),
				}),
			}},
			err: &url.Error{
				Op:  "Delete",
				URL: "https://mock-host/mock-path/platform/configuration/security/enrollment-tokens/atoken",
				Err: errors.New("error"),
			},
		},
		{
			name: "Delete Succeeds with persistent token",
			args: args{params: DeleteParams{
				Token: "atoken",
				API: api.NewMock(mock.Response{
					Response: http.Response{
						StatusCode: 200,
						Body:       mock.NewStructBody(struct{}{}),
					},
				}),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Delete(tt.args.params); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}

func TestList(t *testing.T) {
	type args struct {
		params ListParams
	}
	tests := []struct {
		name string
		args args
		want *models.ListEnrollmentTokenReply
		err  error
	}{
		{
			name: "List Succeeds",
			args: args{params: ListParams{
				API: api.NewMock(mock.Response{
					Response: http.Response{
						StatusCode: 200,
						Body: mock.NewStructBody(models.ListEnrollmentTokenReply{
							Tokens: []*models.ListEnrollmentTokenElement{
								{TokenID: ec.String("token-1")},
								{TokenID: ec.String("token-2"), Roles: []string{"role"}},
							},
						}),
					},
				}),
			}},
			want: &models.ListEnrollmentTokenReply{
				Tokens: []*models.ListEnrollmentTokenElement{
					{TokenID: ec.String("token-1")},
					{TokenID: ec.String("token-2"), Roles: []string{"role"}},
				},
			},
		},
		{
			name: "List fails due to API error",
			args: args{params: ListParams{
				API: api.NewMock(mock.Response{
					Error: errors.New("error"),
				}),
			}},
			err: &url.Error{
				Op:  "Get",
				URL: "https://mock-host/mock-path/platform/configuration/security/enrollment-tokens",
				Err: errors.New("error"),
			},
		},
		{
			name: "Create fails due to missing API",
			args: args{params: ListParams{}},
			err:  util.ErrAPIReq,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := List(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("List() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}
