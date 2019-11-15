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
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func TestQueryParams_Validate(t *testing.T) {
	type fields struct {
		AuthenticationParams util.AuthenticationParams
		ClusterParams        util.ClusterParams
		RestRequest          RestRequest
		Interactive          bool
		Verbose              bool
	}
	tests := []struct {
		name   string
		fields fields
		err    error
	}{
		{
			name: "validate succeeds on interactive",
			fields: fields{
				ClusterParams: util.ClusterParams{
					API:       new(api.API),
					ClusterID: util.ValidClusterID,
				},
				Interactive: true,
			},
		},
		{
			name: "validate succeeds using rest methods",
			fields: fields{
				ClusterParams: util.ClusterParams{
					API:       new(api.API),
					ClusterID: util.ValidClusterID,
				},
				RestRequest: RestRequest{
					Method: "POST",
				},
			},
		},
		{
			name: "validate fails on missing method (non interactive)",
			fields: fields{
				ClusterParams: util.ClusterParams{
					API:       new(api.API),
					ClusterID: util.ValidClusterID,
				},
				Interactive: false,
			},
			err: &multierror.Error{Errors: []error{
				errors.New("method needs to be specified"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := QueryParams{
				AuthenticationParams: tt.fields.AuthenticationParams,
				ClusterParams:        tt.fields.ClusterParams,
				RestRequest:          tt.fields.RestRequest,
				Interactive:          tt.fields.Interactive,
				Verbose:              tt.fields.Verbose,
			}
			if err := params.Validate(); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("QueryParams.Validate() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}

func TestQuery(t *testing.T) {
	type args struct {
		params QueryParams
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "fails on parameter validation failure",
			args: args{params: QueryParams{}},
			err: &multierror.Error{Errors: []error{
				errors.New("cluster id should have a length of 32 characters"),
				errors.New("method needs to be specified"),
			}},
		},
		{
			name: "Succeeds on non interactive and explicit credentials",
			args: args{params: QueryParams{
				RestRequest: RestRequest{Method: "GET"},
				Client: mock.NewClient(mock.Response{Response: http.Response{
					StatusCode: 200,
					Request:    &http.Request{Method: "GET"},
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body: mock.NewStringBody(`{
						"name": "My0iZew",
						"cluster_name": "some-host",
						"cluster_uuid": "URTP-Sp4TgywIL5BLqKd_Q",
						"version": {
						  "number": "6.5.0",
						  "build_flavor": "default",
						  "build_type": "tar",
						  "build_hash": "816e6f6",
						  "build_date": "2018-11-09T18:58:36.352602Z",
						  "build_snapshot": false,
						  "lucene_version": "7.5.0",
						  "minimum_wire_compatibility_version": "5.6.0",
						  "minimum_index_compatibility_version": "5.0.0"
						},
						"tagline": "You Know, for Search"
					  }`),
				}}),
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API: api.NewMock(mock.Response{Response: http.Response{
						StatusCode: http.StatusOK,
						Body: mock.NewStructBody(models.ElasticsearchClusterInfo{
							Metadata: &models.ClusterMetadataInfo{
								Endpoint: "some-host",
							},
						}),
					}}),
				},
				AuthenticationParams: util.AuthenticationParams{
					User: "some",
					Pass: "some",
				},
			}},
		},
		{
			name: "fails to obtain the elasticsearch cluster information",
			args: args{params: QueryParams{
				RestRequest: RestRequest{Method: "GET"},
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API:       api.NewMock(mock.New500Response(mock.NewStringBody(`{"error": "some error"}`))),
				},
			}},
			err: errors.New(`{"error": "some error"}`),
		},
		{
			name: "fails to create an elasticsearch-cli app",
			args: args{params: QueryParams{
				RestRequest: RestRequest{Method: "GET"},
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API: api.NewMock(mock.Response{Response: http.Response{
						StatusCode: http.StatusOK,
						Body: mock.NewStructBody(models.ElasticsearchClusterInfo{
							Metadata: &models.ClusterMetadataInfo{},
						}),
					}}),
				},
				AuthenticationParams: util.AuthenticationParams{
					User: "some",
					Pass: "some",
				},
			}},
			err: errors.New(`host "https://" is invalid`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Query(tt.args.params); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
