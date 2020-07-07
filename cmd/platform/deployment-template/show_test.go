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
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"

	"github.com/elastic/ecctl/cmd/util/testutils"
)

func Test_showCmd(t *testing.T) {
	var templateFormatDeployment = `
	{
  "name": "(Trial) Default Elasticsearch",
  "source": {
	"user_id": "1",
	"facilitator": "adminconsole",
	"date": "2018-04-19T18:16:57.297Z",
	"admin_id": "admin",
	"action": "deployments.create-template",
	"remote_addresses": ["52.205.1.231"]
  },
  "description": "Test default Elasticsearch trial template",
  "id": "default",
  "metadata": [{
	"key": "trial",
	"value": "true"
	}],
	"deployment_template": {
        "resources": {}
    },
	"system_owned": false
}`

	var templateFormatDeploymentResponse = `{
  "deployment_template": {
    "resources": {
      "apm": null,
      "appsearch": null,
      "elasticsearch": null,
      "enterprise_search": null,
      "kibana": null
    }
  },
  "description": "Test default Elasticsearch trial template",
  "id": "default",
  "metadata": [
    {
      "key": "trial",
      "value": "true"
    }
  ],
  "name": "(Trial) Default Elasticsearch",
  "source": {
    "action": "deployments.create-template",
    "admin_id": "admin",
    "date": "2018-04-19T18:16:57.297Z",
    "facilitator": "adminconsole",
    "remote_addresses": [
      "52.205.1.231"
    ],
    "user_id": "1"
  },
  "system_owned": false
}
`

	var templateFormatCluster = `
	{
  "name": "(Trial) Default Elasticsearch",
  "source": {
	"user_id": "1",
	"facilitator": "adminconsole",
	"date": "2018-04-19T18:16:57.297Z",
	"admin_id": "admin",
	"action": "deployments.create-template",
	"remote_addresses": ["52.205.1.231"]
  },
  "description": "Test default Elasticsearch trial template",
  "id": "default",
  "metadata": [{
	"key": "trial",
	"value": "true"
	}],
	"cluster_template": {
		"plan": {}
	},
	"system_owned": false
}`

	var templateFormatClusterResponse = `{
  "cluster_template": {
    "plan": {
      "cluster_topology": null,
      "elasticsearch": null
    }
  },
  "description": "Test default Elasticsearch trial template",
  "id": "default",
  "metadata": [
    {
      "key": "trial",
      "value": "true"
    }
  ],
  "name": "(Trial) Default Elasticsearch",
  "source": {
    "action": "deployments.create-template",
    "admin_id": "admin",
    "date": "2018-04-19T18:16:57.297Z",
    "facilitator": "adminconsole",
    "remote_addresses": [
      "52.205.1.231"
    ],
    "user_id": "1"
  },
  "system_owned": false
}
`

	tests := []struct {
		name string
		args testutils.Args
		want testutils.Assertion
	}{
		{
			name: "incorrect amount of arguments returns error",
			args: testutils.Args{
				Cmd: showCmd,
				Args: []string{
					"show",
				},
			},
			want: testutils.Assertion{
				Err: errors.New("accepts 1 arg(s), received 0"),
			},
		},
		{
			name: "succeeds showing deployment template",
			args: testutils.Args{
				Cmd: showCmd,
				Args: []string{
					"show", "default",
				},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					{
						Response: http.Response{
							Body:       mock.NewStringBody(templateFormatCluster),
							StatusCode: 200,
						},
						Assert: &mock.RequestAssertion{
							Header: api.DefaultReadMockHeaders,
							Method: "GET",
							Host:   api.DefaultMockHost,
							Query: url.Values{
								"format":                       {"cluster"},
								"show_instance_configurations": {"false"},
							},
							Path: "/api/v1/regions/ece-region/platform/configuration/templates/deployments/default",
						},
					},
				}},
			},
			want: testutils.Assertion{
				Stdout: templateFormatClusterResponse,
			},
		},
		{
			name: "succeeds showing deployment template",
			args: testutils.Args{
				Cmd: showCmd,
				Args: []string{
					"show", "default", "--template-format=deployment",
				},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					{
						Response: http.Response{
							Body:       mock.NewStringBody(templateFormatDeployment),
							StatusCode: 200,
						},
						Assert: &mock.RequestAssertion{
							Header: api.DefaultReadMockHeaders,
							Method: "GET",
							Host:   api.DefaultMockHost,
							Query: url.Values{
								"format":                       {"deployment"},
								"show_instance_configurations": {"false"},
							},
							Path: "/api/v1/regions/ece-region/platform/configuration/templates/deployments/default",
						},
					},
				}},
			},
			want: testutils.Assertion{
				Stdout: templateFormatDeploymentResponse,
			},
		},
		{
			name: "succeeds showing instance configurations",
			args: testutils.Args{
				Cmd: showCmd,
				Args: []string{
					"show", "default", "--show-instance-configurations",
				},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					{
						Response: http.Response{
							Body:       mock.NewStringBody(templateFormatDeployment),
							StatusCode: 200,
						},
						Assert: &mock.RequestAssertion{
							Header: api.DefaultReadMockHeaders,
							Method: "GET",
							Host:   api.DefaultMockHost,
							Query: url.Values{
								"format":                       {"deployment"},
								"show_instance_configurations": {"true"},
							},
							Path: "/api/v1/regions/ece-region/platform/configuration/templates/deployments/default",
						},
					},
				}},
			},
			want: testutils.Assertion{
				Stdout: templateFormatDeploymentResponse,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutils.RunCmdAssertion(t, tt.args, tt.want)
		})
	}
}
