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

//go:embed "testdata/show.json"
var showRawResp []byte

//go:embed "testdata/show-dep.json"
var showDepRawResp []byte

func Test_showCmd(t *testing.T) {
	var succeedResp = new(models.Extension)
	if err := succeedResp.UnmarshalBinary(showRawResp); err != nil {
		t.Fatal(err)
	}

	showJSONOutput, err := json.MarshalIndent(succeedResp, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	var succeedDepResp = new(models.Extension)
	if err := succeedDepResp.UnmarshalBinary(showDepRawResp); err != nil {
		t.Fatal(err)
	}

	showDepJSONOutput, err := json.MarshalIndent(succeedDepResp, "", "  ")
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
				Cmd: showCmd,
				Args: []string{
					"show",
				},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: `requires at least 1 arg(s), only received 0`,
			},
		},
		{
			name: "fails due to API error",
			args: testutils.Args{
				Cmd: showCmd,
				Args: []string{
					"show", "2536123203",
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
				Cmd: showCmd,
				Args: []string{
					"show", "2536123203",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultReadMockHeaders,
								Method: "GET",
								Path:   "/api/v1/deployments/extensions/2536123203",
								Host:   api.DefaultMockHost,
								Query: url.Values{
									"include_deployments": []string{"false"},
								},
							},
							mock.NewByteBody(showRawResp),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: string(showJSONOutput) + "\n",
			},
		},
		{
			name: "succeeds with include deployments",
			args: testutils.Args{
				Cmd: showCmd,
				Args: []string{
					"show", "2536123203", "--include-deployments",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultReadMockHeaders,
								Method: "GET",
								Path:   "/api/v1/deployments/extensions/2536123203",
								Host:   api.DefaultMockHost,
								Query: url.Values{
									"include_deployments": []string{"true"},
								},
							},
							mock.NewByteBody(showDepRawResp),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: string(showDepJSONOutput) + "\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutils.RunCmdAssertion(t, tt.args, tt.want)
			tt.args.Cmd.ResetFlags()
			defer initShowFlags()
		})
	}
}
