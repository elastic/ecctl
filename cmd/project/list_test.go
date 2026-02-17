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
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api/mock"

	"github.com/elastic/ecctl/cmd/util/testutils"
	"github.com/elastic/ecctl/pkg/project"
)

func newProjectListBody(projects []project.Project) mock.Response {
	body := project.ListResponse{Items: projects}
	b, _ := json.Marshal(body)
	return mock.New200Response(mock.NewByteBody(b))
}

func initListFlags() {
	listCmd.ResetFlags()
	listCmd.Flags().String("type", "", "Filters by project type (elasticsearch, observability, security)")
}

func Test_listCmd(t *testing.T) {
	var esProjects = []project.Project{
		{
			ID:       "abc123",
			Name:     "my-es-project",
			Type:     "elasticsearch",
			RegionID: "aws-us-east-1",
			Alias:    "my-es-project-abc123",
		},
	}
	var obsProjects = []project.Project{
		{
			ID:       "def456",
			Name:     "my-obs-project",
			Type:     "observability",
			RegionID: "gcp-us-central1",
			Alias:    "my-obs-project-def456",
		},
	}
	var secProjects = []project.Project{
		{
			ID:       "ghi789",
			Name:     "my-sec-project",
			Type:     "security",
			RegionID: "azure-eastus2",
			Alias:    "",
		},
	}

	var allProjects project.ListResult
	allProjects.Projects = append(allProjects.Projects, esProjects...)
	allProjects.Projects = append(allProjects.Projects, obsProjects...)
	allProjects.Projects = append(allProjects.Projects, secProjects...)

	allProjectsJSON, _ := json.MarshalIndent(allProjects, "", "  ")

	var esResult = project.ListResult{Projects: esProjects}
	esResultJSON, _ := json.MarshalIndent(esResult, "", "  ")

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
				Err: `unknown command "some-argument" for "project list"`,
			},
		},
		{
			name: "fails due to invalid type",
			args: testutils.Args{
				Cmd:  listCmd,
				Args: []string{"list", "--type", "invalid"},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: `invalid project type "invalid", must be one of: elasticsearch, observability, security`,
			},
		},
		{
			name: "succeeds filtering by elasticsearch type with JSON format",
			args: testutils.Args{
				Cmd:  listCmd,
				Args: []string{"list", "--type", "elasticsearch"},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						newProjectListBody(esProjects),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: string(esResultJSON) + "\n",
			},
		},
		{
			name: "succeeds filtering by elasticsearch type with text format",
			args: testutils.Args{
				Cmd:  listCmd,
				Args: []string{"list", "--type", "elasticsearch"},
				Cfg: testutils.MockCfg{
					OutputFormat: "text",
					Responses: []mock.Response{
						newProjectListBody(esProjects),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: "ID       NAME            TYPE            REGION          ALIAS\n" +
					"abc123   my-es-project   elasticsearch   aws-us-east-1   my-es-project-abc123\n",
			},
		},
		{
			name: "succeeds listing all project types with JSON format",
			args: testutils.Args{
				Cmd:  listCmd,
				Args: []string{"list"},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						newProjectListBody(esProjects),
						newProjectListBody(obsProjects),
						newProjectListBody(secProjects),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: string(allProjectsJSON) + "\n",
			},
		},
		{
			name: "succeeds listing all project types with text format",
			args: testutils.Args{
				Cmd:  listCmd,
				Args: []string{"list"},
				Cfg: testutils.MockCfg{
					OutputFormat: "text",
					Responses: []mock.Response{
						newProjectListBody(esProjects),
						newProjectListBody(obsProjects),
						newProjectListBody(secProjects),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: "ID       NAME             TYPE            REGION            ALIAS\n" +
					"abc123   my-es-project    elasticsearch   aws-us-east-1     my-es-project-abc123\n" +
					"def456   my-obs-project   observability   gcp-us-central1   my-obs-project-def456\n" +
					"ghi789   my-sec-project   security        azure-eastus2     -\n",
			},
		},
		{
			name: "succeeds with empty project list",
			args: testutils.Args{
				Cmd:  listCmd,
				Args: []string{"list", "--type", "security"},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						newProjectListBody(nil),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: `{
  "projects": null
}
`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutils.RunCmdAssertion(t, tt.args, tt.want)
			initListFlags()
		})
	}
}
