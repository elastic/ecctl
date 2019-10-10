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
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func TestListKeys(t *testing.T) {
	var listKeysSomeIDUser = models.APIKeysResponse{
		Keys: []*models.APIKeyResponse{
			{ID: ec.String("10"), UserID: "someid"},
			{ID: ec.String("11"), UserID: "someid"},
			{ID: ec.String("12"), UserID: "someid"},
		},
	}
	type args struct {
		params ListKeysParams
	}
	tests := []struct {
		name string
		args args
		want *models.APIKeysResponse
		err  error
	}{
		{
			name: "fails due to parameter validation",
			args: args{},
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
			}},
		},
		{
			name: "fails due to API error on all call",
			args: args{params: ListKeysParams{
				API: api.NewMock(
					mock.New500Response(mock.NewStringBody("")),
				),
			}},
			err: errors.New(`unknown error (status 500)`),
		},
		{
			name: "fails due to API error",
			args: args{params: ListKeysParams{
				API: api.NewMock(
					mock.New500Response(mock.NewStringBody("")),
				),
			}},
			err: errors.New(`unknown error (status 500)`),
		},
		{
			name: "succeeds listing keys",
			args: args{params: ListKeysParams{
				API: api.NewMock(
					mock.New200Response(mock.NewStructBody(listKeysSomeIDUser)),
				),
			}},
			want: &listKeysSomeIDUser,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListKeys(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("ListKeys() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}
