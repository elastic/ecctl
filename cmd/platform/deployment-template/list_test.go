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

package cmddeploymentdemplate

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"

	"github.com/elastic/ecctl/cmd/util/testutils"
)

func Test_listCmd(t *testing.T) {
	var templateList = `[
	{
	  "description": "Test default Elasticsearch trial template",
	  "id": "default",
	  "name": "(Trial) Default Elasticsearch",
	  "system_owned": false
	},
	{
	  "description": "Test default Elasticsearch template",
	  "id": "default-appsearch",
	  "name": "Default Elasticsearch",
	  "system_owned": true
	}
  ]
`
	var templateListResponse = `[
  {
    "description": "Test default Elasticsearch trial template",
    "id": "default",
    "metadata": null,
    "name": "(Trial) Default Elasticsearch",
    "system_owned": false
  },
  {
    "description": "Test default Elasticsearch template",
    "id": "default-appsearch",
    "metadata": null,
    "name": "Default Elasticsearch",
    "system_owned": true
  }
]
`

	tests := []struct {
		name string
		args testutils.Args
		want testutils.Assertion
	}{
		{
			name: "succeeds listing deployment template",
			args: testutils.Args{
				Cmd: showCmd,
				Args: []string{
					"list",
				},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					{
						Response: http.Response{
							Body:       mock.NewStringBody(templateList),
							StatusCode: 200,
						},
						Assert: &mock.RequestAssertion{
							Header: api.DefaultReadMockHeaders,
							Method: "GET",
							Host:   api.DefaultMockHost,
							Query: url.Values{
								"format":                       {"cluster"},
								"show_hidden":                  {"false"},
								"show_instance_configurations": {"false"},
							},
							Path: "/api/v1/regions/ece-region/platform/configuration/templates/deployments",
						},
					},
				}},
			},
			want: testutils.Assertion{
				Stdout: templateListResponse,
			},
		},
		{
			name: "succeeds showing deployment template",
			args: testutils.Args{
				Cmd: showCmd,
				Args: []string{
					"list", "--template-format=deployment",
				},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					{
						Response: http.Response{
							Body:       mock.NewStringBody(templateList),
							StatusCode: 200,
						},
						Assert: &mock.RequestAssertion{
							Header: api.DefaultReadMockHeaders,
							Method: "GET",
							Host:   api.DefaultMockHost,
							Query: url.Values{
								"format":                       {"deployment"},
								"show_hidden":                  {"false"},
								"show_instance_configurations": {"false"},
							},
							Path: "/api/v1/regions/ece-region/platform/configuration/templates/deployments",
						},
					},
				}},
			},
			want: testutils.Assertion{
				Stdout: templateListResponse,
			},
		},
		{
			name: "succeeds showing instance configurations",
			args: testutils.Args{
				Cmd: showCmd,
				Args: []string{
					"list", "--show-instance-configurations",
				},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					{
						Response: http.Response{
							Body:       mock.NewStringBody(templateList),
							StatusCode: 200,
						},
						Assert: &mock.RequestAssertion{
							Header: api.DefaultReadMockHeaders,
							Method: "GET",
							Host:   api.DefaultMockHost,
							Query: url.Values{
								"format":                       {"deployment"},
								"show_hidden":                  {"false"},
								"show_instance_configurations": {"true"},
							},
							Path: "/api/v1/regions/ece-region/platform/configuration/templates/deployments",
						},
					},
				}},
			},
			want: testutils.Assertion{
				Stdout: templateListResponse,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutils.RunCmdAssertion(t, tt.args, tt.want)
		})
	}
}
