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

package cmddeploymenttemplate

import (
	"io/ioutil"
	"net/url"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"

	"github.com/elastic/ecctl/cmd/util/testutils"
)

func Test_createCmd(t *testing.T) {
	createRawResp, err := ioutil.ReadFile("./testdata/create.json")
	if err != nil {
		t.Fatal(err)
	}

	var succeedResp = new(models.DeploymentTemplateRequestBody)
	if err := succeedResp.UnmarshalBinary(createRawResp); err != nil {
		t.Fatal(err)
	}

	createBody, err := succeedResp.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		args testutils.Args
		want testutils.Assertion
	}{
		{
			name: "fails due empty argument",
			args: testutils.Args{
				Cmd: createCmd,
				Args: []string{
					"create",
				},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: "failed reading deployment template definition: provide a valid deployment template definition using the --file flag",
			},
		},
		{
			name: "fails due to API error",
			args: testutils.Args{
				Cmd: createCmd,
				Args: []string{
					"create", "--file=testdata/create.json",
				},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: mock.MultierrorInternalError.Error(),
			},
		},
		{
			name: "succeeds with JSON format",
			args: testutils.Args{
				Cmd: createCmd,
				Args: []string{
					"create", "--file=testdata/create.json",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						mock.New201ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultWriteMockHeaders,
								Method: "POST",
								Path:   "/api/v1/deployments/templates",
								Host:   api.DefaultMockHost,
								Query: url.Values{
									"region": []string{"ece-region"},
								},
								Body: mock.NewStringBody(string(createBody) + "\n"),
							},
							mock.NewStructBody(models.IDResponse{
								ID: ec.String("some-randomly-generated-id"),
							}),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: "{\n  \"id\": \"some-randomly-generated-id\"\n}\n",
			},
		},
		{
			name: "succeeds with text format",
			args: testutils.Args{
				Cmd: createCmd,
				Args: []string{
					"create", "--file=testdata/create.json",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "text",
					Responses: []mock.Response{
						mock.New201ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultWriteMockHeaders,
								Method: "POST",
								Path:   "/api/v1/deployments/templates",
								Host:   api.DefaultMockHost,
								Query: url.Values{
									"region": []string{"ece-region"},
								},
								Body: mock.NewStringBody(string(createBody) + "\n"),
							},
							mock.NewStructBody(models.IDResponse{
								ID: ec.String("some-randomly-generated-id"),
							}),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: "ID\nsome-randomly-generated-id\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutils.RunCmdAssertion(t, tt.args, tt.want)
		})
	}
}
