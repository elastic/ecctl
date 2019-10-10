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
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func TestDelete(t *testing.T) {
	type args struct {
		params DeleteParams
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name:    "Delete fails due to parameter validation failure",
			args:    args{},
			wantErr: true,
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				errors.New("user: delete requires a username"),
			}},
		},
		{
			name: "Delete fails due to API failure",
			args: args{
				params: DeleteParams{
					UserName: "user bob",
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(""),
						StatusCode: 400,
					}}),
				},
			},
			wantErr: true,
			err:     errors.New("unknown error (status 400)"),
		},
		{
			name: "Delete succeeds",
			args: args{
				params: DeleteParams{
					UserName: "user bob",
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(""),
						StatusCode: 200,
					}}),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Delete(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.err != nil && !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Delete() actual error = '%v', want error '%v'", err, tt.err)
			}
		})
	}
}
