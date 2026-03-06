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
	"encoding/json"
	"net/http"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api/mock"

	"github.com/elastic/ecctl/cmd/util/testutils"
	"github.com/elastic/ecctl/pkg/project"
)

func newProjectShowBody(p project.Project) mock.Response {
	b, _ := json.Marshal(p)
	return mock.Response{
		Response: http.Response{
			StatusCode: 200,
			Body:       mock.NewByteBody(b),
		},
	}
}

func resetShowFlags() {
	showCmd.ResetFlags()
	initShowFlags()
}

func Test_showCmd(t *testing.T) {
	var proj = project.Project{
		ID:       "abc123",
		Name:     "my-project",
		Type:     "elasticsearch",
		RegionID: "aws-us-east-1",
		Alias:    "my-project-abc123",
		CloudID:  "my-project:abc==",
		Endpoints: map[string]string{
			"elasticsearch": "https://abc123.es.us-east-1.aws.elastic.cloud",
			"kibana":        "https://abc123.kb.us-east-1.aws.elastic.cloud",
		},
	}

	projJSON, _ := json.MarshalIndent(proj, "", "  ")

	tests := []struct {
		name string
		args testutils.Args
		want testutils.Assertion
	}{
		{
			name: "fails due to missing project ID",
			args: testutils.Args{
				Cmd:  showCmd,
				Args: []string{"show"},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: "requires at least 1 arg(s), only received 0",
			},
		},
		{
			name: "succeeds with explicit type and JSON format",
			args: testutils.Args{
				Cmd:  showCmd,
				Args: []string{"show", "abc123", "--type", "elasticsearch"},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						newProjectShowBody(proj),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: string(projJSON) + "\n",
			},
		},
		{
			name: "succeeds with explicit type and text format",
			args: testutils.Args{
				Cmd:  showCmd,
				Args: []string{"show", "abc123", "--type", "elasticsearch"},
				Cfg: testutils.MockCfg{
					OutputFormat: "text",
					Responses: []mock.Response{
						newProjectShowBody(proj),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: "ID       NAME         TYPE            REGION          ALIAS\n" +
					"abc123   my-project   elasticsearch   aws-us-east-1   my-project-abc123\n" +
					"\n" +
					"ENDPOINTS:\n" +
					"  elasticsearch:   https://abc123.es.us-east-1.aws.elastic.cloud\n" +
					"  kibana:          https://abc123.kb.us-east-1.aws.elastic.cloud\n",
			},
		},
		{
			name: "succeeds with auto-detect type",
			args: testutils.Args{
				Cmd:  showCmd,
				Args: []string{"show", "abc123"},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						newProjectListBody([]project.Project{
							{ID: "abc123", Name: "my-project", Type: "elasticsearch", RegionID: "aws-us-east-1"},
						}),
						newProjectListBody(nil),
						newProjectListBody(nil),
						newProjectShowBody(proj),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: string(projJSON) + "\n",
			},
		},
		{
			name: "fails when project not found during auto-detect",
			args: testutils.Args{
				Cmd:  showCmd,
				Args: []string{"show", "nonexistent"},
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
			resetShowFlags()
		})
	}
}
