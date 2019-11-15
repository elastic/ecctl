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

	"github.com/elastic/ecctl/pkg/util"
)

func TestList(t *testing.T) {
	const listUsersResponse = `{
  "users": [{
    "user_name": "admin",
    "builtin": true
  }, {
    "user_name": "readonly",
    "builtin": true
  }]
}`

	type args struct {
		params ListParams
	}
	tests := []struct {
		name    string
		args    args
		want    *models.UserList
		wantErr bool
		err     error
	}{
		{
			name:    "List fails due to parameter validation failure (missing API)",
			args:    args{},
			wantErr: true,
			err:     util.ErrAPIReq,
		},
		{
			name: "List fails due to API failure",
			args: args{
				params: ListParams{
					API: api.NewMock(mock.New500Response(mock.NewStringBody(`{"error": "some error"}`))),
				},
			},
			wantErr: true,
			err:     errors.New(`{"error": "some error"}`),
		},
		{
			name: "List succeeds",
			args: args{
				params: ListParams{
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(listUsersResponse),
						StatusCode: 200,
					}}),
				},
			},
			want: &models.UserList{
				Users: []*models.User{
					{
						Builtin:  ec.Bool(true),
						UserName: ec.String("admin"),
					},
					{
						Builtin:  ec.Bool(true),
						UserName: ec.String("readonly"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := List(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.err != nil && !reflect.DeepEqual(err, tt.err) {
				t.Errorf("List() actual error = '%v', want error '%v'", err, tt.err)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}
