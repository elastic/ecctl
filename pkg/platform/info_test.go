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

package platform

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

func TestGetInfo(t *testing.T) {
	type args struct {
		params GetInfoParams
	}
	tests := []struct {
		name string
		args args
		want *models.PlatformInfo
		err  error
	}{
		{
			name: "Succeeds",
			args: args{params: GetInfoParams{API: api.NewMock(mock.Response{
				Response: http.Response{
					Status:     http.StatusText(http.StatusOK),
					StatusCode: http.StatusOK,
					Body: mock.NewStructBody(models.PlatformInfo{
						EulaAccepted:     ec.Bool(false),
						PhoneHomeEnabled: ec.Bool(false),
					}),
				},
			})}},
			want: &models.PlatformInfo{
				EulaAccepted:     ec.Bool(false),
				PhoneHomeEnabled: ec.Bool(false),
			},
		},
		{
			name: "fails due to API error",
			args: args{params: GetInfoParams{
				API: api.NewMock(mock.New404Response(mock.NewStringBody(`{"error": "some error"}`))),
			}},
			err: errors.New(`{"error": "some error"}`),
		},
		{
			name: "fails due to parameter validation",
			err:  util.ErrAPIReq,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetInfo(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("GetInfo() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
