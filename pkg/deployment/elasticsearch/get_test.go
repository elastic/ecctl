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

	"github.com/elastic/ecctl/pkg/util"
)

func TestGetCluster(t *testing.T) {
	type args struct {
		params GetClusterParams
	}
	tests := []struct {
		name string
		args args
		want *models.ElasticsearchClusterInfo
		err  error
	}{
		{
			name: "Validate fails due to missing cluster ID",
			args: args{},
			err:  errors.New("cluster id should have a length of 32 characters"),
		},
		{
			name: "Validate fails due to missing API",
			args: args{params: GetClusterParams{
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
				},
			}},
			err: errors.New("api reference is required for command"),
		},
		{
			name: "Succeeds",
			args: args{params: GetClusterParams{
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API: api.NewMock(mock.Response{Response: http.Response{
						StatusCode: 200,
						Body: mock.NewStructBody(models.ElasticsearchClusterInfo{
							ClusterID: ec.String(util.ValidClusterID),
						}),
					}}),
				},
			}},
			want: &models.ElasticsearchClusterInfo{
				ClusterID: ec.String(util.ValidClusterID),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCluster(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("GetCluster() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCluster() = %v, want %v", got, tt.want)
			}
		})
	}
}
