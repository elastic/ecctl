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

package role

import (
	"errors"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func TestSetBlessings(t *testing.T) {
	type args struct {
		params SetBlessingsParams
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "fails on parameter validation",
			args: args{},
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				errors.New("role set blessing: blessing definitions cannot be empty"),
				errors.New("role set blessing: id cannot be empty"),
			}},
		},
		{
			name: "fails updating the role",
			args: args{params: SetBlessingsParams{
				API: api.NewMock(mock.New500Response(mock.NewStringBody(
					`{"error": "failed updating role"}`,
				))),
				Blessings: &models.Blessings{},
				ID:        "one",
			}},
			err: errors.New(`{"error": "failed updating role"}`),
		},
		{
			name: "succeeds",
			args: args{params: SetBlessingsParams{
				API:       api.NewMock(mock.New200Response(mock.NewStringBody(""))),
				Blessings: &models.Blessings{},
				ID:        "one",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetBlessings(tt.args.params); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("SetBlessings() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
