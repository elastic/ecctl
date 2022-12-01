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

package cmddeploymenttemplate

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

func Test_listCmd(t *testing.T) {
	var reqAssertion = &mock.RequestAssertion{
		Header: api.DefaultReadMockHeaders,
		Method: "GET",
		Path:   "/api/v1/deployments/templates",
		Host:   api.DefaultMockHost,
		Query: url.Values{
			"region":                       []string{"ece-region"},
			"show_hidden":                  []string{"false"},
			"show_instance_configurations": []string{"true"},
		},
	}

	var succeedResp []*models.DeploymentTemplateInfoV2
	if err := json.Unmarshal(listRawResp, &succeedResp); err != nil {
		t.Fatal(err)
	}

	listJSONOutput, err := json.MarshalIndent(succeedResp, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	textOutput := `ID        NAME      SYSTEM   DESCRIPTION` + "\n"
	textOutput += `default   Default   true     Default deployment template for clusters` + "\n"

	tests := []struct {
		name string
		args testutils.Args
		want testutils.Assertion
	}{
		{
			name: "fails due to API error",
			args: testutils.Args{
				Cmd: listCmd,
				Args: []string{
					"list",
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
			name: "succeeds with JSON format",
			args: testutils.Args{
				Cmd: listCmd,
				Args: []string{
					"list",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							reqAssertion, mock.NewByteBody(listRawResp),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: string(listJSONOutput) + "\n",
			},
		},
		{
			name: "succeeds with text format",
			args: testutils.Args{
				Cmd: listCmd,
				Args: []string{
					"list",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "text",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							reqAssertion, mock.NewByteBody(listRawResp),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: textOutput,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutils.RunCmdAssertion(t, tt.args, tt.want)
		})
	}
}
