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

package allocator

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
)

func TestSearch(t *testing.T) {
	type args struct {
		params SearchParams
	}
	tests := []struct {
		name    string
		args    args
		want    *models.AllocatorOverview
		wantErr bool
		error   error
	}{
		{
			name: "fails if search request is invalid",
			args: args{params: SearchParams{
				Request: models.SearchRequest{Query: &models.QueryContainer{Exists: &models.ExistsQuery{Field: nil}}},
				API: api.NewMock(mock.Response{Response: http.Response{
					Body:       mock.NewStringBody(""),
					StatusCode: 200,
				}}),
			}},
			want:    nil,
			wantErr: true,
			error:   nil,
		},
		{
			name:    "fails if api reference is empty",
			args:    args{params: SearchParams{}},
			want:    nil,
			wantErr: true,
			error:   errors.New("api reference is required for command"),
		},
		{
			name: "fails if search api call fails",
			args: args{params: SearchParams{
				Request: models.SearchRequest{Query: &models.QueryContainer{}},
				API: api.NewMock(mock.Response{Response: http.Response{
					Body:       mock.NewStringBody(`{"error":"invalid query"}`),
					StatusCode: 404,
				}}),
			}},
			want:    nil,
			wantErr: true,
			error:   errors.New("unknown error (status 404)"),
		},
		{
			name: "succeeds if search api call succeeds",
			args: args{params: SearchParams{
				Request: models.SearchRequest{Query: &models.QueryContainer{}},
				API: api.NewMock(mock.Response{Response: http.Response{
					Body: mock.NewStringBody(`{
						"zones": [
						  {
							"allocators": null,
							"zone_id": "us-east-1a"
						  },
						  {
							"allocators": null,
							"zone_id": "us-east-1c"
						  },
						  {
							"allocators": null,
							"zone_id": "us-east-1e"
						  }
						]
					  }`),
					StatusCode: 200,
				}}),
			}},
			want: &models.AllocatorOverview{
				Zones: []*models.AllocatorZoneInfo{
					{ZoneID: ec.String("us-east-1a")},
					{ZoneID: ec.String("us-east-1c")},
					{ZoneID: ec.String("us-east-1e")},
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

			if tt.wantErr && tt.error != nil && !reflect.DeepEqual(err, tt.error) {
				t.Errorf("Search() actual error = '%v', want error '%v'", err, tt.error)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Search() = %v, want %v", got, tt.want)
			}
		})
	}
}
