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
	"net/url"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"

	"github.com/elastic/ecctl/pkg/util"
)

const errStr = `{
  "errors": null
}`

func TestResyncCluster(t *testing.T) {
	type args struct {
		params ResyncClusterParams
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "Fails due to parameter validation (Cluster ID)",
			args:    args{},
			wantErr: util.ErrClusterLength,
		},
		{
			name: "Fails due to parameter validation (API)",
			args: args{params: ResyncClusterParams{
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
				},
			}},
			wantErr: util.ErrAPIReq,
		},
		{
			name: "Fails due to unknown API response",
			args: args{params: ResyncClusterParams{
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API:       api.NewMock(mock.New500Response(mock.NewStringBody(`{}`))),
				},
			}},
			wantErr: errors.New(errStr),
		},
		{
			name: "Fails due to API error",
			args: args{params: ResyncClusterParams{
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API: api.NewMock(mock.Response{
						Error: errors.New("error with API"),
					}),
				},
			}},
			wantErr: &url.Error{
				Op:  "Post",
				URL: "https://mock-host/mock-path/clusters/elasticsearch/320b7b540dfc967a7a649c18e2fce4ed/_resync",
				Err: errors.New("error with API"),
			},
		},
		{
			name: "Succeeds to resynchronize Elasticsearch cluster without errors",
			args: args{params: ResyncClusterParams{
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API: api.NewMock(mock.Response{
						Response: http.Response{
							StatusCode: http.StatusOK,
							Body:       mock.NewStringBody(`{}`),
						},
					}),
				},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ResyncCluster(tt.args.params); !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("ResyncCluster() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
