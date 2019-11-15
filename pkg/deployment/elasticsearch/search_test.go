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

package elasticsearch

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

func TestSearchClusters(t *testing.T) {
	type args struct {
		params SearchClusterParams
	}
	tests := []struct {
		name    string
		args    args
		want    *models.ElasticsearchClustersInfo
		wantErr bool
		error   error
	}{
		{
			name: "Search clusters fails if search Request is invalid",
			args: args{params: SearchClusterParams{
				Request: models.SearchRequest{Query: &models.QueryContainer{Exists: &models.ExistsQuery{Field: nil}}},
				API: api.NewMock(mock.Response{Response: http.Response{
					Body:       mock.NewStringBody("{}"),
					StatusCode: 200,
				}}),
			}},
			wantErr: true,
		},
		{
			name:    "Search clusters fails if api reference is empty",
			args:    args{params: SearchClusterParams{}},
			want:    nil,
			wantErr: true,
			error:   errors.New("api reference is required for command"),
		},
		{
			name: "Search clusters fails if search api call fails",
			args: args{params: SearchClusterParams{
				Request: models.SearchRequest{Query: &models.QueryContainer{}},
				API:     api.NewMock(mock.New404Response(mock.NewStringBody(`{}`))),
			}},
			wantErr: true,
			error:   errors.New("{}"),
		},
		{
			name: "Search clusters succeeds if search api call succeeds",
			args: args{params: SearchClusterParams{
				Request: models.SearchRequest{Query: &models.QueryContainer{}},
				API: api.NewMock(mock.Response{Response: http.Response{
					Body: mock.NewStringBody(`{
						"match_count": 0,
						"return_count": 0,
						"elasticsearch_clusters": []
					  }`),
					StatusCode: 200,
				}}),
			}},
			want: &models.ElasticsearchClustersInfo{
				MatchCount:            0,
				ReturnCount:           ec.Int32(0),
				ElasticsearchClusters: []*models.ElasticsearchClusterInfo{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SearchClusters(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchClusters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.error != nil && !reflect.DeepEqual(err, tt.error) {
				t.Errorf("SearchClusters() actual error = '%v', want error '%v'", err, tt.error)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchClusters() = %v, want %v", got, tt.want)
			}
		})
	}
}
