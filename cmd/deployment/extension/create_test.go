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

package cmddeploymentextension

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"

	"github.com/elastic/ecctl/cmd/util/testutils"
)

//go:embed "testdata/create.json"
var createRawResp []byte

func Test_createCmd(t *testing.T) {
	var succeedResp = new(models.Extension)
	if err := succeedResp.UnmarshalBinary(createRawResp); err != nil {
		t.Fatal(err)
	}

	createJSONOutput, err := json.MarshalIndent(succeedResp, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		args testutils.Args
		want testutils.Assertion
	}{
		{
			name: "fails due to empty argument",
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
				Err: `accepts 1 arg(s), received 0`,
			},
		},
		{
			name: "fails due to missing flags",
			args: testutils.Args{
				Cmd: createCmd,
				Args: []string{
					"create", "mybundle",
				},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: `required flag(s) "type", "version" not set`,
			},
		},
		{
			name: "fails due to mismatching flags",
			args: testutils.Args{
				Cmd: createCmd,
				Args: []string{
					"create", "mybundle", "--version", "*", "--type", "bundle",
					"--download-url", "example.com", "--file", "hi.zip",
				},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: "both --file and --download-url are set. Only one may be used",
			},
		},
		{
			name: "fails due to API error",
			args: testutils.Args{
				Cmd: createCmd,
				Args: []string{
					"create", "mybundle", "--version", "*", "--type", "bundle",
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
			name: "succeeds",
			args: testutils.Args{
				Cmd: createCmd,
				Args: []string{
					"create", "mybundle", "--version", "*", "--type", "bundle", "--description", "Why hello there",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultWriteMockHeaders,
								Method: "POST",
								Path:   "/api/v1/deployments/extensions",
								Host:   api.DefaultMockHost,
								Body:   mock.NewStringBody(`{"description":"Why hello there","extension_type":"bundle","name":"mybundle","version":"*"}` + "\n"),
							},
							mock.NewByteBody(createRawResp),
						),
					},
				},
			},
			want: testutils.Assertion{
				Err: string(createJSONOutput) + "\n",
			},
		},
		{
			name: "succeeds with upload",
			args: testutils.Args{
				Cmd: createCmd,
				Args: []string{
					"create", "mybundle", "--version", "*", "--type", "bundle", "--file",
					"./testdata/extension.zip", "--description", "Why hello there",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						mock.New200Response(
							mock.NewByteBody(createRawResp),
						),
						mock.New200Response(
							mock.NewByteBody(createRawResp),
						),
					},
				},
			},
			want: testutils.Assertion{
				Err: string(createJSONOutput) + "\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutils.RunCmdAssertion(t, tt.args, tt.want)
			tt.args.Cmd.ResetFlags()
			defer initCreateFlags()
		})
	}
}
