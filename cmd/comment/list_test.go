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

func Test_listCmd(t *testing.T) {
	listJSONResponse := `{
  "values": [
    {
      "comment": {
        "id": "93420f71ca474b79aa4bc2aaa7f37e21",
        "message": "some message",
        "user_id": "root"
      },
      "metadata": {
        "created_time": "2021-03-17T13:05:06.958Z",
        "modified_time": "2021-03-17T13:05:06.958Z",
        "version": "8|14"
      }
    },
    {
      "comment": {
        "id": "e66e8336d4f3489193146a993779ce92",
        "message": "some message",
        "user_id": "root"
      },
      "metadata": {
        "created_time": "2021-03-17T12:21:58.202Z",
        "modified_time": "2021-03-17T12:21:58.202Z",
        "version": "9|12"
      }
    }
  ]
}
`
	tests := []struct {
		name string
		args testutils.Args
		want testutils.Assertion
	}{
		{
			name: "fails due to arguments passed",
			args: testutils.Args{
				Cmd:  listCmd,
				Args: []string{"list", "some-argument"},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: `accepts at most 0 arg(s), received 1`,
			},
		},
		{
			name: "fails due to missing resource-id",
			args: testutils.Args{
				Cmd:  listCmd,
				Args: []string{"list", "--resource-type", "allocator"},
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
				Cmd:  listCmd,
				Args: []string{"list", "--resource-id", "i-i123"},
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
				Cmd: listCmd,
				Args: []string{
					"list", "--resource-type", "allocator", "--resource-id", "i-123",
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
				Cmd: listCmd,
				Args: []string{
					"list", "--resource-type", "allocator", "--resource-id", "i-123",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "text",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultReadMockHeaders,
								Method: "GET",
								Path:   "/api/v1/regions/ece-region/comments/allocator/i-123",
								Host:   api.DefaultMockHost,
							},
							mock.NewStringBody(listJSONResponse),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: "COMMENT ID                         USER   MESSAGE        CREATED TIME               MODIFIED TIME              VERSION\n" +
					"93420f71ca474b79aa4bc2aaa7f37e21   root   some message   2021-03-17T13:05:06.958Z   2021-03-17T13:05:06.958Z   8|14\n" +
					"e66e8336d4f3489193146a993779ce92   root   some message   2021-03-17T12:21:58.202Z   2021-03-17T12:21:58.202Z   9|12\n",
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
