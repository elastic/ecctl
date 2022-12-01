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
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"

	"github.com/elastic/ecctl/cmd/util/testutils"
)

//go:embed "testdata/create.json"
var createRawResp []byte

func Test_createCmd(t *testing.T) {
	var succeedResp = new(models.TrafficFilterRulesets)
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
			name: "fails due to API error",
			args: testutils.Args{
				Cmd:  createCmd,
				Args: []string{"create"},
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
				Args: []string{"create", "--name", "hi", "--type", "ip",
					"--include-by-default", "--description", "bob", "--source", "0.0.0.0/0,0.0.0.0/1"},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultWriteMockHeaders,
								Method: "POST",
								Path:   "/api/v1/deployments/traffic-filter/rulesets",
								Host:   api.DefaultMockHost,
								Body:   mock.NewStringBody(`{"description":"bob","include_by_default":true,"name":"hi","region":"ece-region","rules":[{"source":"0.0.0.0/0"},{"source":"0.0.0.0/1"}],"type":"ip"}` + "\n"),
							},
							mock.NewByteBody(createRawResp),
						),
					},
				},
			},
			want: testutils.Assertion{
				Err: string(createJSONOutput),
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
