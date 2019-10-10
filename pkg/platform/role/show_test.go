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
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func TestShow(t *testing.T) {
	type args struct {
		params ShowParams
	}
	tests := []struct {
		name string
		args args
		want *models.RoleAggregate
		err  error
	}{
		{
			name: "fails on parameter validation",
			args: args{},
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				errors.New("role show: id cannot be empty"),
			}},
		},
		{
			name: "fails on api error",
			args: args{params: ShowParams{
				API: api.NewMock(mock.New500Response(mock.NewStringBody(
					"big failure",
				))),
				ID: "some",
			}},
			err: errors.New("unknown error (status 500)"),
		},
		{
			name: "succeeds",
			args: args{params: ShowParams{
				ID: "one",
				API: api.NewMock(mock.New200Response(mock.NewStructBody(
					models.RoleAggregate{
						ID: ec.String("one"),
					},
				))),
			}},
			want: &models.RoleAggregate{
				ID: ec.String("one"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Show(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Show() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Show() = %v, want %v", got, tt.want)
			}
		})
	}
}
