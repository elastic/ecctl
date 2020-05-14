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

func TestGetKey(t *testing.T) {
	type args struct {
		params GetKeyParams
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
				errors.New("get key requires a key id"),
			),
		},
		{
			name: "fails due to API error",
			args: args{params: GetKeyParams{
				API: api.NewMock(
					mock.New404Response(mock.NewStructBody(
						&models.BasicFailedReply{Errors: []*models.BasicFailedReplyElement{
							{
								Code:    ec.String("key.not_found"),
								Message: ec.String("key not found"),
							},
						}},
					)),
				),
				ID: "somekey",
			}},
			err: multierror.NewPrefixed("api error",
				errors.New("key.not_found: key not found"),
			),
		},
		{
			name: "succeeds",
			args: args{params: GetKeyParams{
				API: api.NewMock(
					mock.New200Response(mock.NewStructBody(models.APIKeyResponse{
						Key: "somekeyvalue",
						ID:  ec.String("somekey"),
					})),
				),
				ID: "somekey",
			}},
			want: &models.APIKeyResponse{
				Key: "somekeyvalue",
				ID:  ec.String("somekey"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetKey(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("GetKey() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
