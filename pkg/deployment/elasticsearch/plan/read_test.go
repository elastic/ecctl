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

package plan

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"

	"github.com/elastic/ecctl/pkg/util"
)

func TestGet(t *testing.T) {
	type args struct {
		params GetParams
	}
	tests := []struct {
		name string
		args args
		want *models.ElasticsearchClusterPlanInfo
		err  error
	}{
		{
			name: "Fails due to missing API parameter validation",
			args: args{params: GetParams{ClusterParams: util.ClusterParams{
				ClusterID: "8b43ed5e277f7ea6f13606fcf4027f9c",
			}}},
			err: util.ErrAPIReq,
		},
		{
			name: "Fails due to missing cluster ID parameter validation",
			args: args{params: GetParams{ClusterParams: util.ClusterParams{
				API: api.NewMock(mock.Response{Response: http.Response{
					StatusCode: http.StatusNotFound,
					Status:     http.StatusText(http.StatusNotFound),
					Body:       mock.NewStringBody(`{}`),
				}}),
			}}},
			err: errors.New("cluster id should have a length of 32 characters"),
		},
		{
			name: "Succeeds with the current plan",
			args: args{params: GetParams{ClusterParams: util.ClusterParams{
				ClusterID: "8b43ed5e277f7ea6f13606fcf4027f9c",
				API: api.NewMock(mock.Response{Response: http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
					Body: mock.NewStructBody(models.ElasticsearchClusterPlansInfo{
						Current: &models.ElasticsearchClusterPlanInfo{PlanAttemptID: "someID"},
					}),
				}}),
			}}},
			want: &models.ElasticsearchClusterPlanInfo{PlanAttemptID: "someID"},
		},
		{
			name: "Succeeds with the pending plan",
			args: args{params: GetParams{Pending: true, ClusterParams: util.ClusterParams{
				ClusterID: "8b43ed5e277f7ea6f13606fcf4027f9c",
				API: api.NewMock(mock.Response{Response: http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
					Body: mock.NewStructBody(models.ElasticsearchClusterPlansInfo{
						Pending: &models.ElasticsearchClusterPlanInfo{PlanAttemptID: "pendingID"},
					}),
				}}),
			}}},
			want: &models.ElasticsearchClusterPlanInfo{PlanAttemptID: "pendingID"},
		},
		{
			name: "Fails when the pending plan is asked but there is none",
			args: args{params: GetParams{Pending: true, ClusterParams: util.ClusterParams{
				ClusterID: "8b43ed5e277f7ea6f13606fcf4027f9c",
				API: api.NewMock(mock.Response{Response: http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
					Body: mock.NewStructBody(models.ElasticsearchClusterPlansInfo{
						Pending: nil,
					}),
				}}),
			}}},
			err: errors.New("no pending plan"),
		},
		{
			name: "Fails when the API returns an error",
			args: args{params: GetParams{Pending: true, ClusterParams: util.ClusterParams{
				ClusterID: "8b43ed5e277f7ea6f13606fcf4027f9c",
				API: api.NewMock(mock.Response{Response: http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     http.StatusText(http.StatusInternalServerError),
					Body:       mock.NewStringBody(`{}`),
				}}),
			}}},
			err: errors.New("unknown error (status 500)"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHistory(t *testing.T) {
	type args struct {
		params GetHistoryParams
	}
	tests := []struct {
		name string
		args args
		want []*models.ElasticsearchClusterPlanInfo
		err  error
	}{
		{
			name: "Fails due to missing API parameter validation",
			args: args{params: GetHistoryParams{ClusterParams: util.ClusterParams{
				ClusterID: "8b43ed5e277f7ea6f13606fcf4027f9c",
			}}},
			err: util.ErrAPIReq,
		},
		{
			name: "Fails due to missing cluster ID parameter validation",
			args: args{params: GetHistoryParams{ClusterParams: util.ClusterParams{
				API: api.NewMock(mock.Response{Response: http.Response{
					StatusCode: http.StatusNotFound,
					Status:     http.StatusText(http.StatusNotFound),
					Body:       mock.NewStringBody(`{}`),
				}}),
			}}},
			err: errors.New("cluster id should have a length of 32 characters"),
		},
		{
			name: "Succeeds getting the plan history",
			args: args{params: GetHistoryParams{ClusterParams: util.ClusterParams{
				ClusterID: "8b43ed5e277f7ea6f13606fcf4027f9c",
				API: api.NewMock(mock.Response{Response: http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
					Body: mock.NewStructBody(models.ElasticsearchClusterPlansInfo{
						History: []*models.ElasticsearchClusterPlanInfo{
							{PlanAttemptID: "someID"},
						},
					}),
				}}),
			}}},
			want: []*models.ElasticsearchClusterPlanInfo{
				{PlanAttemptID: "someID"},
			},
		},
		{
			name: "Fails when the API returns an error",
			args: args{params: GetHistoryParams{ClusterParams: util.ClusterParams{
				ClusterID: "8b43ed5e277f7ea6f13606fcf4027f9c",
				API: api.NewMock(mock.Response{Response: http.Response{
					StatusCode: http.StatusInternalServerError,
					Status:     http.StatusText(http.StatusInternalServerError),
					Body:       mock.NewStringBody(`{}`),
				}}),
			}}},
			err: errors.New("unknown error (status 500)"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHistory(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("GetHistory() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetHistory() = %v, want %v", got, tt.want)
			}
		})
	}
}
