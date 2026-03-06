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

func init() {
	// Mirror the persistent --region flag from cmd/root.go so that
	// createCmd can inherit it when tested in isolation.
	if Command.PersistentFlags().Lookup("region") == nil {
		Command.PersistentFlags().String("region", "", "Elastic Cloud Hosted or Serverless region")
	}
}

func newProjectCreateBody(p project.CreateResult) mock.Response {
	b, _ := json.Marshal(p)
	return mock.Response{
		Response: http.Response{
			StatusCode: 201,
			Body:       mock.NewByteBody(b),
		},
	}
}

func resetCreateFlags() {
	createCmd.ResetFlags()
	initCreateFlags()
}

func Test_createCmd(t *testing.T) {
	var created = project.CreateResult{
		Project: project.Project{
			ID:       "abc123",
			Name:     "my-project",
			Type:     "elasticsearch",
			RegionID: "aws-us-east-1",
			Alias:    "my-project-abc123",
		},
		Credentials: project.Credentials{
			Username: "admin",
			Password: "secret-password",
		},
	}

	createdJSON, _ := json.MarshalIndent(created, "", "  ")

	tests := []struct {
		name string
		args testutils.Args
		want testutils.Assertion
	}{
		{
			name: "fails due to missing required flags",
			args: testutils.Args{
				Cmd:  createCmd,
				Args: []string{"create"},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: `required flag(s) "name", "type" not set`,
			},
		},
		{
			name: "fails due to invalid type",
			args: testutils.Args{
				Cmd:  createCmd,
				Args: []string{"create", "--type", "invalid", "--name", "test", "--region", "aws-us-east-1"},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: `invalid project type "invalid", must be one of: elasticsearch (or search), observability, security`,
			},
		},
		{
			name: "succeeds with JSON format",
			args: testutils.Args{
				Cmd:  createCmd,
				Args: []string{"create", "--type", "elasticsearch", "--name", "my-project", "--region", "aws-us-east-1"},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						newProjectCreateBody(created),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: string(createdJSON) + "\n",
			},
		},
		{
			name: "succeeds with text format",
			args: testutils.Args{
				Cmd:  createCmd,
				Args: []string{"create", "--type", "elasticsearch", "--name", "my-project", "--region", "aws-us-east-1"},
				Cfg: testutils.MockCfg{
					OutputFormat: "text",
					Responses: []mock.Response{
						newProjectCreateBody(created),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: "ID       NAME         TYPE            REGION          ALIAS\n" +
					"abc123   my-project   elasticsearch   aws-us-east-1   my-project-abc123",
			},
		},
		{
			name: "succeeds using search as type alias",
			args: testutils.Args{
				Cmd:  createCmd,
				Args: []string{"create", "--type", "search", "--name", "my-project", "--region", "aws-us-east-1"},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						newProjectCreateBody(created),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: string(createdJSON) + "\n",
			},
		},
		{
			name: "fails when tier is used with elasticsearch type",
			args: testutils.Args{
				Cmd:  createCmd,
				Args: []string{"create", "--type", "elasticsearch", "--name", "test", "--region", "aws-us-east-1", "--tier", "complete"},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: "invalid project create params: 1 error occurred:\n\t* tier is not supported for elasticsearch projects\n\n",
			},
		},
		{
			name: "succeeds with tier for observability",
			args: testutils.Args{
				Cmd:  createCmd,
				Args: []string{"create", "--type", "observability", "--name", "my-obs-project", "--region", "aws-us-east-1", "--tier", "complete"},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						newProjectCreateBody(project.CreateResult{
							Project: project.Project{
								ID:       "obs123",
								Name:     "my-obs-project",
								Type:     "observability",
								RegionID: "aws-us-east-1",
							},
							Credentials: project.Credentials{
								Username: "admin",
								Password: "secret",
							},
						}),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: `{
  "id": "obs123",
  "name": "my-obs-project",
  "type": "observability",
  "region_id": "aws-us-east-1",
  "metadata": {},
  "credentials": {
    "username": "admin",
    "password": "secret"
  }
}
`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutils.RunCmdAssertion(t, tt.args, tt.want)
			resetCreateFlags()
		})
	}
}
