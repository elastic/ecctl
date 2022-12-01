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

package cmdrole

import (
	_ "embed"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api/mock"

	"github.com/elastic/ecctl/cmd/util/testutils"
)

//go:embed "testdata/update-response.json"
var roleRawResponseData []byte

func Test_updateCmd(t *testing.T) {
	tests := []struct {
		name string
		args testutils.Args
		want testutils.Assertion
	}{
		{
			name: "fails due empty argument",
			args: testutils.Args{
				Cmd: updateCmd,
				Args: []string{
					"update",
				},
			},
			want: testutils.Assertion{
				Err: "accepts 1 arg(s), received 0",
			},
		},
		{
			name: "fails due to API error",
			args: testutils.Args{
				Cmd: updateCmd,
				Args: []string{
					"update", "proxy", "--file=./testdata/update-request.json",
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
			name: "fails due to missing file option",
			args: testutils.Args{
				Cmd: updateCmd,
				Args: []string{
					"update", "proxy",
				},
			},
			want: testutils.Assertion{
				Err: `required flag(s) "file" not set`,
			},
		},
		{
			name: "fails due to inconsistent role ID",
			args: testutils.Args{
				Cmd: updateCmd,
				Args: []string{
					"update", "my-role", "--file=./testdata/update-request.json",
				},
			},
			want: testutils.Assertion{
				Err: "role id [my-role] cannot be found in the role file",
			},
		},
		{
			name: "succeeds to update role",
			args: testutils.Args{
				Cmd: updateCmd,
				Args: []string{
					"update", "proxy", "--file=./testdata/update-request.json",
				},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.New200Response(mock.NewStringBody(string(roleRawResponseData))),
				}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutils.RunCmdAssertion(t, tt.args, tt.want)
			tt.args.Cmd.ResetFlags()
			defer initFlags()
		})
	}
}
