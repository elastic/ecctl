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

package instances

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

func TestList(t *testing.T) {
	type args struct {
		params util.ClusterParams
	}
	tests := []struct {
		name string
		args args
		want []string
		err  error
	}{
		{
			name: "returns a single instance",
			args: args{params: util.ClusterParams{
				ClusterID: util.ValidClusterID,
				API: api.NewMock(mock.Response{Response: http.Response{
					StatusCode: 200,
					Body: mock.NewStructBody(models.ElasticsearchClusterInfo{
						Topology: &models.ClusterTopologyInfo{
							Instances: []*models.ClusterInstanceInfo{
								{InstanceName: ec.String("instance-000000")},
							},
						},
					}),
				}}),
			}},
			want: []string{"instance-000000"},
		},
		{
			name: "returns multiple instances",
			args: args{params: util.ClusterParams{
				ClusterID: util.ValidClusterID,
				API: api.NewMock(mock.Response{Response: http.Response{
					StatusCode: 200,
					Body: mock.NewStructBody(models.ElasticsearchClusterInfo{
						Topology: &models.ClusterTopologyInfo{
							Instances: []*models.ClusterInstanceInfo{
								{InstanceName: ec.String("instance-000000")},
								{InstanceName: ec.String("instance-000001")},
								{InstanceName: ec.String("instance-000002")},
							},
						},
					}),
				}}),
			}},
			want: []string{"instance-000000", "instance-000001", "instance-000002"},
		},
		{
			name: "fails due to API error",
			args: args{params: util.ClusterParams{
				ClusterID: util.ValidClusterID,
				API:       api.NewMock(mock.New500Response(mock.NewStringBody(`{"error": "some error"}`))),
			}},
			err: errors.New(`{"error": "some error"}`),
		},
		{
			name: "fails due parameter validation",
			args: args{params: util.ClusterParams{}},
			err:  util.ErrClusterLength,
		},
		{
			name: "fails due parameter validation",
			args: args{params: util.ClusterParams{
				ClusterID: util.ValidClusterID,
			}},
			err: util.ErrAPIReq,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := List(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("List() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}
