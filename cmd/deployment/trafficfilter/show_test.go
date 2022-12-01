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

package cmddeploymenttrafficfilter

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

//go:embed "testdata/show-assoc.json"
var showAssocRawResp []byte

func Test_showCmd(t *testing.T) {
	var succeedResp = new(models.TrafficFilterRulesetInfo)
	if err := succeedResp.UnmarshalBinary(showRawResp); err != nil {
		t.Fatal(err)
	}

	showJSONOutput, err := json.MarshalIndent(succeedResp, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	var succeedAssocResp = new(models.TrafficFilterRulesetInfo)
	if err := succeedAssocResp.UnmarshalBinary(showAssocRawResp); err != nil {
		t.Fatal(err)
	}

	showAssocJSONOutput, err := json.MarshalIndent(succeedAssocResp, "", "  ")
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
					"show", "11111111111111111111111111111111",
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
					"show", "4e974d9476534d35b12fbdcfd0acee0a",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultReadMockHeaders,
								Method: "GET",
								Path:   "/api/v1/deployments/traffic-filter/rulesets/4e974d9476534d35b12fbdcfd0acee0a",
								Host:   api.DefaultMockHost,
								Query: url.Values{
									"include_associations": []string{"false"},
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
			name: "succeeds with include associations",
			args: testutils.Args{
				Cmd: showCmd,
				Args: []string{
					"show", "4e974d9476534d35b12fbdcfd0acee0a", "--include-associations",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultReadMockHeaders,
								Method: "GET",
								Path:   "/api/v1/deployments/traffic-filter/rulesets/4e974d9476534d35b12fbdcfd0acee0a",
								Host:   api.DefaultMockHost,
								Query: url.Values{
									"include_associations": []string{"true"},
								},
							},
							mock.NewByteBody(showAssocRawResp),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: string(showAssocJSONOutput) + "\n",
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
