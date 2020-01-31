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
)

func TestList(t *testing.T) {
	var runnerListSuccess = `
{
  "runners": [{
    "connected": true,
    "runner_id": "192.168.44.10"
  }]
}`
	type args struct {
		params Params
	}
	tests := []struct {
		name    string
		args    args
		want    *models.RunnerOverview
		wantErr error
	}{
		{
			name: "Runner list succeeds",
			args: args{params: Params{
				API: api.NewMock(mock.Response{Response: http.Response{
					Body:       mock.NewStringBody(runnerListSuccess),
					StatusCode: 200,
				}}),
			}},
			want: &models.RunnerOverview{
				Runners: []*models.RunnerInfo{
					{
						RunnerID:  ec.String("192.168.44.10"),
						Connected: ec.Bool(true),
					},
				},
			},
		},
		{
			name: "Runner list fails",
			args: args{params: Params{
				API: api.NewMock(mock.Response{Error: errors.New("error")}),
			}},
			want: nil,
			wantErr: &url.Error{
				Op:  "Get",
				URL: "https://mock-host/mock-path/platform/infrastructure/runners",
				Err: errors.New("error"),
			},
		},
		{
			name: "Runner list with an empty API",
			args: args{params: Params{
				API: nil,
			}},
			want:    nil,
			wantErr: errors.New("api reference is required for command"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := List(tt.args.params)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}
