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

//go:embed "testdata/list.json"
var listRawResp []byte

//go:embed "testdata/list-assoc.json"
var listAssocRawResp []byte

//go:embed "testdata/list-region.json"
var listRegionRawResp []byte

func Test_listCmd(t *testing.T) {

	var succeedResp = new(models.TrafficFilterRulesets)
	if err := succeedResp.UnmarshalBinary(listRawResp); err != nil {
		t.Fatal(err)
	}

	listJSONOutput, err := json.MarshalIndent(succeedResp, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	var succeedAssocResp = new(models.TrafficFilterRulesets)
	if err := succeedAssocResp.UnmarshalBinary(listAssocRawResp); err != nil {
		t.Fatal(err)
	}

	listAssocJSONOutput, err := json.MarshalIndent(succeedAssocResp, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	var succeedRegionResp = new(models.TrafficFilterRulesets)
	if err := succeedRegionResp.UnmarshalBinary(listRegionRawResp); err != nil {
		t.Fatal(err)
	}

	listRegionJSONOutput, err := json.MarshalIndent(succeedRegionResp, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		args testutils.Args
		want testutils.Assertion
	}{
		{
			name: "fails due to API error",
			args: testutils.Args{
				Cmd:  listCmd,
				Args: []string{"list"},
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
				Cmd:  listCmd,
				Args: []string{"list"},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultReadMockHeaders,
								Method: "GET",
								Path:   "/api/v1/deployments/traffic-filter/rulesets",
								Host:   api.DefaultMockHost,
								Query: url.Values{
									"include_associations": []string{"false"},
								},
							},
							mock.NewByteBody(listRawResp),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: string(listJSONOutput) + "\n",
			},
		},
		{
			name: "succeeds with include associations",
			args: testutils.Args{
				Cmd: listCmd,
				Args: []string{
					"list", "--include-associations",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultReadMockHeaders,
								Method: "GET",
								Path:   "/api/v1/deployments/traffic-filter/rulesets",
								Host:   api.DefaultMockHost,
								Query: url.Values{
									"include_associations": []string{"true"},
								},
							},
							mock.NewByteBody(listAssocRawResp),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: string(listAssocJSONOutput) + "\n",
			},
		},
		{
			name: "succeeds with region",
			args: testutils.Args{
				Cmd: listCmd,
				Args: []string{
					"list", "--single-region", "azure-eastus2",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultReadMockHeaders,
								Method: "GET",
								Path:   "/api/v1/deployments/traffic-filter/rulesets",
								Host:   api.DefaultMockHost,
								Query: url.Values{
									"include_associations": []string{"false"},
									"region":               []string{"azure-eastus2"},
								},
							},
							mock.NewByteBody(listRegionRawResp),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: string(listRegionJSONOutput) + "\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutils.RunCmdAssertion(t, tt.args, tt.want)
			tt.args.Cmd.ResetFlags()
			defer initListFlags()
		})
	}
}
