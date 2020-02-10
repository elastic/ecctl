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

package runner

import (
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"
)

func TestShow(t *testing.T) {
	var runnerShow = `
{
  "connected": true,
  "runner_id": "192.168.44.10"
}`

	type args struct {
		params ShowParams
	}
	tests := []struct {
		name    string
		args    args
		want    *models.RunnerInfo
		wantErr error
	}{
		{
			name: "Show runner succeeds",
			args: args{
				params: ShowParams{
					ID: "192.168.44.10",
					Params: Params{
						API: api.NewMock(mock.Response{Response: http.Response{
							Body:       mock.NewStringBody(runnerShow),
							StatusCode: 200,
						}}),
					},
				},
			},
			want: &models.RunnerInfo{
				RunnerID:  ec.String("192.168.44.10"),
				Connected: ec.Bool(true),
			},
		},
		{
			name: "Show runner fails",
			args: args{
				params: ShowParams{
					ID: "192.168.44.10",
					Params: Params{
						API: api.NewMock(mock.Response{Error: errors.New("error")}),
					},
				},
			},
			want: nil,
			wantErr: &url.Error{
				Op:  "Get",
				URL: "https://mock-host/mock-path/platform/infrastructure/runners/192.168.44.10",
				Err: errors.New("error"),
			},
		},
		{
			name: "Show runner fails due to validation",
			args: args{
				params: ShowParams{},
			},
			want: nil,
			wantErr: &multierror.Error{
				Errors: []error{
					errors.New("id field cannot be empty"),
					errors.New("api reference is required for command"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Show(tt.args.params)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Show() = %v, want %v", got, tt.want)
			}
		})
	}
}
