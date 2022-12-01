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

package cmddeployment

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

//go:embed "testdata/show.json"
var showRawResp []byte

//go:embed "testdata/show_apm.json"
var showApmResp []byte

//go:embed "testdata/want_generate-payload.json"
var wantGeneratePayload []byte

//go:embed "testdata/show-resource.json"
var showResourceRawResp []byte

func Test_showCmd(t *testing.T) {
	var succeedResp = new(models.DeploymentGetResponse)
	if err := succeedResp.UnmarshalBinary(showRawResp); err != nil {
		t.Fatal(err)
	}

	showJSONOutput, err := json.MarshalIndent(succeedResp, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	var succeedESResp = new(models.ElasticsearchResourceInfo)
	if err := succeedESResp.UnmarshalBinary(showResourceRawResp); err != nil {
		t.Fatal(err)
	}

	showESJSONOutput, err := json.MarshalIndent(succeedESResp, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	var succeedKibanaResp = new(models.KibanaResourceInfo)
	if err := succeedKibanaResp.UnmarshalBinary(showResourceRawResp); err != nil {
		t.Fatal(err)
	}

	showKibanaJSONOutput, err := json.MarshalIndent(succeedKibanaResp, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	var succeedApmResp = new(models.ApmResourceInfo)
	if err := succeedApmResp.UnmarshalBinary(showResourceRawResp); err != nil {
		t.Fatal(err)
	}

	showApmJSONOutput, err := json.MarshalIndent(succeedApmResp, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	var succeedAppSearchResp = new(models.AppSearchResourceInfo)
	if err := succeedAppSearchResp.UnmarshalBinary(showResourceRawResp); err != nil {
		t.Fatal(err)
	}

	showAppSearchJSONOutput, err := json.MarshalIndent(succeedAppSearchResp, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	var succeedEnterpriseSearchResp = new(models.EnterpriseSearchResourceInfo)
	if err := succeedEnterpriseSearchResp.UnmarshalBinary(showResourceRawResp); err != nil {
		t.Fatal(err)
	}

	showEnterpriseSearchJSONOutput, err := json.MarshalIndent(succeedEnterpriseSearchResp, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		args testutils.Args
		want testutils.Assertion
	}{
		{
			name: "fails due empty argument",
			args: testutils.Args{
				Cmd: showCmd,
				Args: []string{
					"show",
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
			name: "fails due to API error",
			args: testutils.Args{
				Cmd: showCmd,
				Args: []string{
					"show", "29337f77410e23ab30e15c280060facf",
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
					"show", "29337f77410e23ab30e15c280060facf",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultReadMockHeaders,
								Method: "GET",
								Path:   "/api/v1/deployments/29337f77410e23ab30e15c280060facf",
								Host:   api.DefaultMockHost,
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
							mock.NewByteBody(showRawResp),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: string(showJSONOutput) + "\n",
			},
		},
		{
			name: "succeeds with `--generate-update-payload`",
			args: testutils.Args{
				Cmd: showCmd,
				Args: []string{
					"show", "29337f77410e23ab30e15c280060facf",
					"--generate-update-payload",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultReadMockHeaders,
								Method: "GET",
								Path:   "/api/v1/deployments/29337f77410e23ab30e15c280060facf",
								Host:   api.DefaultMockHost,
								Query: url.Values{
									"convert_legacy_plans": {"false"},
									"show_metadata":        {"false"},
									"show_plan_defaults":   {"false"},
									"show_plan_history":    {"false"},
									"show_plan_logs":       {"false"},
									"show_plans":           {"true"},
									"show_settings":        {"true"},
									"show_system_alerts":   {"5"},
								},
							},
							mock.NewByteBody(showApmResp),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: string(wantGeneratePayload),
			},
		},
		{
			name: "succeeds with elasticsearch kind",
			args: testutils.Args{
				Cmd: showCmd,
				Args: []string{
					"show", "29337f77410e23ab30e15c280060facf",
					"--kind=elasticsearch", "--ref-id=main-elasticsearch",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultReadMockHeaders,
								Method: "GET",
								Path:   "/api/v1/deployments/29337f77410e23ab30e15c280060facf/elasticsearch/main-elasticsearch",
								Host:   api.DefaultMockHost,
								Query: url.Values{
									"show_metadata":      {"false"},
									"show_plan_defaults": {"false"},
									"show_plan_logs":     {"false"},
									"show_plans":         {"false"},
									"show_settings":      {"false"},
									"show_system_alerts": {"5"},
								},
							},
							mock.NewByteBody(showRawResp),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: string(showESJSONOutput) + "\n",
			},
		},
		{
			name: "succeeds with kibana kind",
			args: testutils.Args{
				Cmd: showCmd,
				Args: []string{
					"show", "29337f77410e23ab30e15c280060facf",
					"--kind=kibana", "--ref-id=main-kibana",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultReadMockHeaders,
								Method: "GET",
								Path:   "/api/v1/deployments/29337f77410e23ab30e15c280060facf/kibana/main-kibana",
								Host:   api.DefaultMockHost,
								Query: url.Values{
									"show_metadata":      {"false"},
									"show_plan_defaults": {"false"},
									"show_plan_logs":     {"false"},
									"show_plans":         {"false"},
									"show_settings":      {"false"},
								},
							},
							mock.NewByteBody(showResourceRawResp),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: string(showKibanaJSONOutput) + "\n",
			},
		},
		{
			name: "succeeds with apm kind",
			args: testutils.Args{
				Cmd: showCmd,
				Args: []string{
					"show", "29337f77410e23ab30e15c280060facf",
					"--kind=apm", "--ref-id=main-apm",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultReadMockHeaders,
								Method: "GET",
								Path:   "/api/v1/deployments/29337f77410e23ab30e15c280060facf/apm/main-apm",
								Host:   api.DefaultMockHost,
								Query: url.Values{
									"show_metadata":      {"false"},
									"show_plan_defaults": {"false"},
									"show_plan_logs":     {"false"},
									"show_plans":         {"false"},
									"show_settings":      {"false"},
								},
							},
							mock.NewByteBody(showResourceRawResp),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: string(showApmJSONOutput) + "\n",
			},
		},
		{
			name: "succeeds with appsearch kind",
			args: testutils.Args{
				Cmd: showCmd,
				Args: []string{
					"show", "29337f77410e23ab30e15c280060facf",
					"--kind=appsearch", "--ref-id=main-appsearch",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultReadMockHeaders,
								Method: "GET",
								Path:   "/api/v1/deployments/29337f77410e23ab30e15c280060facf/appsearch/main-appsearch",
								Host:   api.DefaultMockHost,
								Query: url.Values{
									"show_metadata":      {"false"},
									"show_plan_defaults": {"false"},
									"show_plan_logs":     {"false"},
									"show_plans":         {"false"},
									"show_settings":      {"false"},
								},
							},
							mock.NewByteBody(showResourceRawResp),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: string(showAppSearchJSONOutput) + "\n",
			},
		},
		{
			name: "succeeds with enterprise_search kind",
			args: testutils.Args{
				Cmd: showCmd,
				Args: []string{
					"show", "29337f77410e23ab30e15c280060facf",
					"--kind=enterprise_search", "--ref-id=main-enterprise_search",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultReadMockHeaders,
								Method: "GET",
								Path:   "/api/v1/deployments/29337f77410e23ab30e15c280060facf/enterprise_search/main-enterprise_search",
								Host:   api.DefaultMockHost,
								Query: url.Values{
									"show_metadata":      {"false"},
									"show_plan_defaults": {"false"},
									"show_plan_logs":     {"false"},
									"show_plans":         {"false"},
									"show_settings":      {"false"},
								},
							},
							mock.NewByteBody(showResourceRawResp),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: string(showEnterpriseSearchJSONOutput) + "\n",
			},
		},
		{
			name: "succeeds with show flags set",
			args: testutils.Args{
				Cmd: showCmd,
				Args: []string{
					"show", "29337f77410e23ab30e15c280060facf", "--plans",
					"--plan-logs", "--plan-defaults", "--plan-history",
					"--metadata", "--settings",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultReadMockHeaders,
								Method: "GET",
								Path:   "/api/v1/deployments/29337f77410e23ab30e15c280060facf",
								Host:   api.DefaultMockHost,
								Query: url.Values{
									"convert_legacy_plans": {"false"},
									"show_metadata":        {"true"},
									"show_plan_defaults":   {"true"},
									"show_plan_history":    {"true"},
									"show_plan_logs":       {"true"},
									"show_plans":           {"true"},
									"show_settings":        {"true"},
									"show_system_alerts":   {"5"},
								},
							},
							mock.NewByteBody(showResourceRawResp),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: string(showJSONOutput) + "\n",
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
