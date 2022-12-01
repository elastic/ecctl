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
	"net/url"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"

	"github.com/elastic/ecctl/cmd/util/testutils"
)

//go:embed "testdata/update-payload.json"
var updatePayloadRawResp []byte

//go:embed "testdata/update.json"
var updateRawResp []byte

func Test_updateCmd(t *testing.T) {
	var succeedPayloadResp = new(models.UpdateExtensionRequest)
	if err := succeedPayloadResp.UnmarshalBinary(updatePayloadRawResp); err != nil {
		t.Fatal(err)
	}

	updatePayloadJSONOutput, err := json.MarshalIndent(succeedPayloadResp, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	var succeedResp = new(models.Extension)
	if err := succeedResp.UnmarshalBinary(updateRawResp); err != nil {
		t.Fatal(err)
	}

	updateJSONOutput, err := json.MarshalIndent(succeedResp, "", "  ")
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
				Cmd: updateCmd,
				Args: []string{
					"update",
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
				Cmd: updateCmd,
				Args: []string{
					"update", "2649887448",
				},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: `one of --file or --generate-payload must be set`,
			},
		},
		{
			name: "fails due to mismatching flags",
			args: testutils.Args{
				Cmd: updateCmd,
				Args: []string{
					"update", "2649887448", "--generate-payload", "--file", "testdata/update.json",
				},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: "both --file and --generate-payload are set. Only one may be used",
			},
		},
		{
			name: "fails due to API error",
			args: testutils.Args{
				Cmd: updateCmd,
				Args: []string{
					"update", "2649887448", "--generate-payload",
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
			name: "succeeds to generate payload",
			args: testutils.Args{
				Cmd: updateCmd,
				Args: []string{
					"update", "2649887448", "--generate-payload",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultReadMockHeaders,
								Method: "GET",
								Path:   "/api/v1/deployments/extensions/2649887448",
								Host:   api.DefaultMockHost,
								Query: url.Values{
									"include_deployments": []string{"false"},
								},
							},
							mock.NewByteBody(updatePayloadRawResp),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: string(updatePayloadJSONOutput) + "\n",
			},
		},
		{
			name: "succeeds to update",
			args: testutils.Args{
				Cmd: updateCmd,
				Args: []string{
					"update", "2649887448", "--file", "./testdata/update-payload.json",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultWriteMockHeaders,
								Method: "POST",
								Path:   "/api/v1/deployments/extensions/2649887448",
								Body:   mock.NewStringBody(`{"description":"hello","extension_type":"bundle","name":"mybundle","version":"*"}` + "\n"),
								Host:   api.DefaultMockHost,
							},
							mock.NewByteBody(updateRawResp),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: string(updateJSONOutput) + "\n",
			},
		},
		{
			name: "succeeds with upload",
			args: testutils.Args{
				Cmd: updateCmd,
				Args: []string{
					"update", "2649887448", "--file", "./testdata/update-payload.json",
					"--extension-file", "./testdata/extension.zip",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						mock.New200Response(
							mock.NewByteBody(updateRawResp),
						),
						mock.New200Response(
							mock.NewByteBody(updateRawResp),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: string(updateJSONOutput) + "\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutils.RunCmdAssertion(t, tt.args, tt.want)
			tt.args.Cmd.ResetFlags()
			defer initUpdateFlags()
		})
	}
}
