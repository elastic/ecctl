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
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/multierror"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
)

func TestSearch(t *testing.T) {
	var runnerSearchSuccess = `
{
  "runners": [{
    "connected": true,
	"runner_id": "192.168.44.10"
  },{
	"connected": true,
	"runner_id": "192.168.44.11" 
  }]
}`
	var searchReqErr = `validation failure list:
validation failure list:
validation failure list:
field in body is required`

	type args struct {
		params SearchParams
	}
	tests := []struct {
		name    string
		args    args
		want    *models.RunnerOverview
		wantErr bool
		err     error
	}{
		{
			name: "fails validation",
			args: args{params: SearchParams{
				Request: models.SearchRequest{Query: &models.QueryContainer{Exists: &models.ExistsQuery{Field: nil}}},
			}},
			wantErr: true,
			err: multierror.NewPrefixed("runner search",
				errors.New("api reference is required for command"),
				errors.New(searchReqErr),
			),
		},
		{
			name: "fails if search api call fails",
			args: args{params: SearchParams{
				Request: models.SearchRequest{Query: &models.QueryContainer{}},
				Params: Params{
					API: api.NewMock(mock.New404Response(mock.NewStringBody(`{"error": "some error"}`))),
				},
			}},
			wantErr: true,
			err:     errors.New(`{"error": "some error"}`),
		},
		{
			name: "succeeds if search api call succeeds",
			args: args{params: SearchParams{
				Request: models.SearchRequest{Query: &models.QueryContainer{}},
				Params: Params{
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(runnerSearchSuccess),
						StatusCode: 200,
					}}),
				},
			}},
			want: &models.RunnerOverview{
				Runners: []*models.RunnerInfo{
					{
						RunnerID:  ec.String("192.168.44.10"),
						Connected: ec.Bool(true),
					}, {
						RunnerID:  ec.String("192.168.44.11"),
						Connected: ec.Bool(true),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Search(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.err != nil && (err.Error() != tt.err.Error()) {
				t.Errorf("Search() actual error = '%v', want error '%v'", err, tt.err)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Search() = %v, want %v", got, tt.want)
			}
		})
	}
}
