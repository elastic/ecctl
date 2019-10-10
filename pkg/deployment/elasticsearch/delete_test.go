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

	"github.com/elastic/ecctl/pkg/util"
)

func TestDeleteCluster(t *testing.T) {
	type args struct {
		params DeleteClusterParams
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "succeeds",
			args: args{params: DeleteClusterParams{
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API: api.NewMock(mock.Response{
						Response: http.Response{
							StatusCode: 200,
							Body: mock.NewStructBody(models.ElasticsearchClusterInfo{
								Status: stoppedState,
							}),
						},
					}, mock.Response{
						Response: http.Response{
							StatusCode: 200,
							Body:       mock.NewStringBody(`{}`),
						},
					}),
				},
			}},
		},
		{
			name: "fails to delete the cluster",
			args: args{params: DeleteClusterParams{
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API: api.NewMock(mock.Response{
						Response: http.Response{
							StatusCode: 200,
							Body: mock.NewStructBody(models.ElasticsearchClusterInfo{
								Status: stoppedState,
							}),
						},
					}, mock.Response{
						Response: http.Response{
							StatusCode: 500,
							Body:       mock.NewStringBody(`{}`),
						},
					}),
				},
			}},
			err: errors.New("unknown error (status 500)"),
		},
		{
			name: "fails due to cluster not being stopped",
			args: args{params: DeleteClusterParams{
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API: api.NewMock(mock.Response{
						Response: http.Response{
							StatusCode: 200,
							Body: mock.NewStructBody(models.ElasticsearchClusterInfo{
								Status: "started",
							}),
						},
					}),
				},
			}},
			err: errors.New("elasticsearch delete: deployment must be stopped"),
		},
		{
			name: "fails due to api failure",
			args: args{params: DeleteClusterParams{
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API: api.NewMock(mock.Response{
						Response: http.Response{
							StatusCode: 500,
							Body:       mock.NewStringBody(`{}`),
						},
					}),
				},
			}},
			err: errors.New("unknown error (status 500)"),
		},
		{
			name: "fails due to missing cluster ID",
			args: args{params: DeleteClusterParams{}},
			err:  errors.New("cluster id should have a length of 32 characters"),
		},
		{
			name: "fails due to missing API",
			args: args{params: DeleteClusterParams{
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
				},
			}},
			err: errors.New("api reference is required for command"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteCluster(tt.args.params); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("DeleteCluster() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
