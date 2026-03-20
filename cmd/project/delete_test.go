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

package cmdproject

import (
	"net/http"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api/mock"

	"github.com/elastic/ecctl/cmd/util/testutils"
	"github.com/elastic/ecctl/pkg/project"
)

func newDeleteOKResponse() mock.Response {
	return mock.Response{
		Response: http.Response{
			StatusCode: 200,
			Body:       mock.NewStringBody(`{}`),
		},
	}
}

func resetDeleteFlags() {
	deleteCmd.ResetFlags()
	initDeleteFlags()
}

func Test_deleteCmd(t *testing.T) {
	tests := []struct {
		name string
		args testutils.Args
		want testutils.Assertion
	}{
		{
			name: "fails due to missing project ID",
			args: testutils.Args{
				Cmd:  deleteCmd,
				Args: []string{"delete"},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: "requires at least 1 arg(s), only received 0",
			},
		},
		{
			name: "succeeds with explicit type",
			args: testutils.Args{
				Cmd:  deleteCmd,
				Args: []string{"delete", "abc123", "--type", "elasticsearch"},
				Cfg: testutils.MockCfg{
					Responses: []mock.Response{
						newDeleteOKResponse(),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: "Project \"abc123\" deletion scheduled.\n",
			},
		},
		{
			name: "succeeds with auto-detect type",
			args: testutils.Args{
				Cmd:  deleteCmd,
				Args: []string{"delete", "abc123"},
				Cfg: testutils.MockCfg{
					Responses: []mock.Response{
						newProjectListBody([]project.Project{
							{ID: "abc123", Name: "my-project", Type: "elasticsearch", RegionID: "aws-us-east-1"},
						}),
						newProjectListBody(nil),
						newProjectListBody(nil),
						newDeleteOKResponse(),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: "Project \"abc123\" deletion scheduled.\n",
			},
		},
		{
			name: "fails when project not found during auto-detect",
			args: testutils.Args{
				Cmd:  deleteCmd,
				Args: []string{"delete", "nonexistent"},
				Cfg: testutils.MockCfg{
					Responses: []mock.Response{
						newProjectListBody(nil),
						newProjectListBody(nil),
						newProjectListBody(nil),
					},
				},
			},
			want: testutils.Assertion{
				Err: `project "nonexistent" not found`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutils.RunCmdAssertion(t, tt.args, tt.want)
			resetDeleteFlags()
		})
	}
}
