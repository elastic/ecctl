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

package cmdtrafficfilterassoc

import (
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"

	"github.com/elastic/ecctl/cmd/util/testutils"
)

func Test_createCmd(t *testing.T) {
	tests := []struct {
		name string
		args testutils.Args
		want testutils.Assertion
	}{
		{
			name: "fails due to empty argument",
			args: testutils.Args{
				Cmd:  createCmd,
				Args: []string{"create"},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: `requires at least 1 arg(s), only received 0`,
			},
		},
		{
			name: "fails due to missing flag",
			args: testutils.Args{
				Cmd: createCmd,
				Args: []string{
					"create", "11111111111111111111111111111111",
				},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: `required flag(s) "deployment-id" not set`,
			},
		},
		{
			name: "fails due to API error",
			args: testutils.Args{
				Cmd: createCmd,
				Args: []string{
					"create", "11111111111111111111111111111111", "--deployment-id", "11",
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
					"create", "a8b649146d4149479b03c2076886e4de", "--deployment-id", "0a1214498c68421ca09908545a0e25ca",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultWriteMockHeaders,
								Method: "POST",
								Path:   "/api/v1/deployments/traffic-filter/rulesets/a8b649146d4149479b03c2076886e4de/associations",
								Host:   api.DefaultMockHost,
								Body:   mock.NewStringBody(`{"entity_type":"deployment","id":"0a1214498c68421ca09908545a0e25ca"}` + "\n"),
							},
							mock.NewByteBody(nil),
						),
					},
				},
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
