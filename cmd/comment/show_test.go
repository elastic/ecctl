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

package cmdcomment

import (
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"

	"github.com/elastic/ecctl/cmd/util/testutils"
)

func Test_showCmd(t *testing.T) {
	tests := []struct {
		name string
		args testutils.Args
		want testutils.Assertion
	}{
		{
			name: "fails due to missing message",
			args: testutils.Args{
				Cmd:  showCmd,
				Args: []string{"show"},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: `accepts 1 arg(s), received 0`,
			},
		},
		{
			name: "fails due to missing resource-id",
			args: testutils.Args{
				Cmd:  showCmd,
				Args: []string{"show", "some-message-id", "--resource-type", "allocator"},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: `required flag(s) "resource-id" not set`,
			},
		},
		{
			name: "fails due to missing resource-type",
			args: testutils.Args{
				Cmd:  showCmd,
				Args: []string{"show", "some-message-id", "--resource-id", "i-i123"},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: `required flag(s) "resource-type" not set`,
			},
		},
		{
			name: "fails due to API error",
			args: testutils.Args{
				Cmd: showCmd,
				Args: []string{
					"show", "some-message", "--resource-type", "allocator", "--resource-id", "i-123",
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
					"show", "some-message-id", "--resource-type", "allocator", "--resource-id", "i-123",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "text",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultReadMockHeaders,
								Method: "GET",
								Path:   "/api/v1/regions/ece-region/comments/allocator/i-123/some-message-id",
								Host:   api.DefaultMockHost,
							},
							mock.NewStringBody(`{"id":"random-generated-id", "user_id":"some-user", "message":"some-message"}`),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: "COMMENT ID            USER        MESSAGE\n" +
					"random-generated-id   some-user   some-message\n",
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
