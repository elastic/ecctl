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

package cmdeskeystore

import (
	_ "embed"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"

	"github.com/elastic/ecctl/cmd/util/testutils"
	"github.com/elastic/ecctl/pkg/util"
)

//go:embed "testdata/update.json"
var updateRawResp []byte

func Test_updateCmd(t *testing.T) {
	var succeedResp = new(models.KeystoreContents)
	if err := succeedResp.UnmarshalBinary(updateRawResp); err != nil {
		t.Fatal(err)
	}

	updateJSONOutput, err := json.MarshalIndent(succeedResp, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	rawBody, err := succeedResp.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	rawBody = append(rawBody, []byte("\n")...)

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
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: "requires at least 1 arg(s), only received 0",
			},
		},
		{
			name: "fails due empty definition",
			args: testutils.Args{
				Cmd: updateCmd,
				Args: []string{
					"update", util.ValidClusterID,
				},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: "failed reading keystore secret definition: provide a valid keystore secret definition using the --file flag",
			},
		},
		{
			name: "fails due to API error",
			args: testutils.Args{
				Cmd: updateCmd,
				Args: []string{
					"update", util.ValidClusterID, "--file=./testdata/update.json",
					"--ref-id=main-elasticsearch",
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
				Cmd: updateCmd,
				Args: []string{
					"update", util.ValidClusterID, "--file=./testdata/update.json",
					"--ref-id=main-elasticsearch",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultWriteMockHeaders,
								Method: "PATCH",
								Path:   "/api/v1/deployments/320b7b540dfc967a7a649c18e2fce4ed/elasticsearch/main-elasticsearch/keystore",
								Host:   api.DefaultMockHost,
								Body:   mock.NewByteBody(rawBody),
							},
							mock.NewByteBody(updateRawResp),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: string(updateJSONOutput) + "\n",
			},
		},
		{
			name: "succeeds with text format",
			args: testutils.Args{
				Cmd: updateCmd,
				Args: []string{
					"update", util.ValidClusterID, "--file=./testdata/update.json",
					"--ref-id=main-elasticsearch",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "text",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultWriteMockHeaders,
								Method: "PATCH",
								Path:   "/api/v1/deployments/320b7b540dfc967a7a649c18e2fce4ed/elasticsearch/main-elasticsearch/keystore",
								Host:   api.DefaultMockHost,
								Body:   mock.NewByteBody(rawBody),
							},
							mock.NewByteBody(updateRawResp),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: "SECRET        VALUE              AS FILE\n" +
					"secret-name   the secret value   <nil>\n",
			},
		},
		{
			name: "succeeds with text format and and ref-id auto-discovery",
			args: testutils.Args{
				Cmd: updateCmd,
				Args: []string{
					"update", util.ValidClusterID, "--file=./testdata/update.json",
					"--ref-id=",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "text",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultReadMockHeaders,
								Method: "GET",
								Host:   api.DefaultMockHost,
								Path:   "/api/v1/deployments/320b7b540dfc967a7a649c18e2fce4ed",
								Query: url.Values{
									"convert_legacy_plans": {"false"},
									"show_metadata":        {"false"},
									"show_plan_defaults":   {"false"},
									"show_plan_history":    {"false"},
									"show_plan_logs":       {"false"},
									"show_plans":           {"false"},
									"show_settings":        {"false"},
									"show_system_alerts":   {"5"},
								},
							},
							mock.NewStructBody(models.DeploymentGetResponse{
								Healthy: ec.Bool(true),
								ID:      ec.String(mock.ValidClusterID),
								Resources: &models.DeploymentResources{
									Elasticsearch: []*models.ElasticsearchResourceInfo{{
										ID:    ec.String(mock.ValidClusterID),
										RefID: ec.String("elasticsearch"),
									}},
								},
							}),
						),
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultWriteMockHeaders,
								Method: "PATCH",
								Path:   "/api/v1/deployments/320b7b540dfc967a7a649c18e2fce4ed/elasticsearch/elasticsearch/keystore",
								Host:   api.DefaultMockHost,
								Body:   mock.NewByteBody(rawBody),
							},
							mock.NewByteBody(updateRawResp),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: "SECRET        VALUE              AS FILE\n" +
					"secret-name   the secret value   <nil>\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutils.RunCmdAssertion(t, tt.args, tt.want)
		})
	}
}
