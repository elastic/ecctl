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

package cmddeployment

import (
	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/elastic/ecctl/cmd/util/testutils"
	"testing"
)

func Test_searchCmd(t *testing.T) {
	tests := []struct {
		name string
		args testutils.Args
		want testutils.Assertion
	}{
		{
			name: "all-matches collects deployments using multiple requests",
			args: testutils.Args{
				Cmd: searchCmd,
				Args: []string{
					"search",
					"-f",
					"testdata/search_query.json",
					"--all-matches",
					"--size",
					"250",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultWriteMockHeaders,
								Method: "POST",
								Path:   "/api/v1/deployments/_search",
								Host:   api.DefaultMockHost,
								Body: mock.NewStructBody(models.SearchRequest{
									Query: &models.QueryContainer{
										MatchAll: struct{}{},
									},
									Size: 250,
									Sort: []interface{}{"id"},
								}),
							},
							mock.NewStructBody(models.DeploymentsSearchResponse{
								Cursor: "cursor1",
								Deployments: []*models.DeploymentSearchResponse{
									{ID: ec.String("d1")},
									{ID: ec.String("d2")},
								},
								MatchCount:      3,
								MinimalMetadata: nil,
								ReturnCount:     ec.Int32(2),
							}),
						),
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultWriteMockHeaders,
								Method: "POST",
								Path:   "/api/v1/deployments/_search",
								Host:   api.DefaultMockHost,
								Body: mock.NewStructBody(models.SearchRequest{
									Cursor: "cursor1",
									Query: &models.QueryContainer{
										MatchAll: struct{}{},
									},
									Size: 250,
									Sort: []interface{}{"id"},
								}),
							},
							mock.NewStructBody(models.DeploymentsSearchResponse{
								Cursor: "cursor2",
								Deployments: []*models.DeploymentSearchResponse{
									{ID: ec.String("d3")},
								},
								MatchCount:      3,
								MinimalMetadata: nil,
								ReturnCount:     ec.Int32(1),
							}),
						),
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultWriteMockHeaders,
								Method: "POST",
								Path:   "/api/v1/deployments/_search",
								Host:   api.DefaultMockHost,
								Body: mock.NewStructBody(models.SearchRequest{
									Cursor: "cursor2",
									Query: &models.QueryContainer{
										MatchAll: struct{}{},
									},
									Size: 250,
									Sort: []interface{}{"id"},
								}),
							},
							mock.NewStructBody(models.DeploymentsSearchResponse{
								Cursor:          "cursor3",
								Deployments:     []*models.DeploymentSearchResponse{},
								MatchCount:      3,
								MinimalMetadata: nil,
								ReturnCount:     ec.Int32(0),
							}),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: string(expectedOutput) + "\n",
			},
		},
		{
			name: "all-matches requires a query with a sort",
			args: testutils.Args{
				Cmd: searchCmd,
				Args: []string{
					"search",
					"-f",
					"testdata/search_query_no_sort.json",
					"--all-matches",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses:    []mock.Response{},
				},
			},
			want: testutils.Assertion{
				Err: "The query must include a sort-field when using --all-matches. Example: \"sort\": [\"id\"]",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutils.RunCmdAssertion(t, tt.args, tt.want)
		})
	}
}

var expectedOutput = `{
  "deployments": [
    {
      "healthy": null,
      "id": "d1",
      "name": null,
      "resources": null
    },
    {
      "healthy": null,
      "id": "d2",
      "name": null,
      "resources": null
    },
    {
      "healthy": null,
      "id": "d3",
      "name": null,
      "resources": null
    }
  ],
  "match_count": 3,
  "minimal_metadata": null,
  "return_count": 3
}`
