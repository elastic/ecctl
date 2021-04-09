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

package cmdproxysettings

import (
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"

	"github.com/elastic/ecctl/cmd/util/testutils"
)

func Test_showCmd(t *testing.T) {
	const expectedProxiesSettings = `{
  "expected_proxies_count": 5,
  "http_settings": {
    "cookie_secret": "some-secret",
    "dashboards_base_url": "some-url",
    "disconnected_cutoff": 2,
    "minimum_proxy_services": 1,
    "sso_settings": {
      "cookie_name": "some-cookie",
      "default_redirect_path": "some-path",
      "dont_log_requests": true,
      "maintenance_bypass_cookie_name": "some-other-cookie",
      "max_age": 10,
      "sso_secret": "sso-secret"
    },
    "user_cookie_key": "user-cookie"
  },
  "signature_secret": "signature-secret",
  "signature_valid_for_millis": 60
}`

	tests := []struct {
		name string
		args testutils.Args
		want testutils.Assertion
	}{
		{
			name: "fails due to provided args",
			args: testutils.Args{
				Cmd:  platformProxySettingsShowCmd,
				Args: []string{"show", "some-arg"},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: `accepts at most 0 arg(s), received 1`,
			},
		},
		{
			name: "fails due to API error",
			args: testutils.Args{
				Cmd: platformProxySettingsShowCmd,
				Args: []string{
					"show",
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
				Cmd: platformProxySettingsShowCmd,
				Args: []string{
					"show",
				},
				Cfg: testutils.MockCfg{
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultReadMockHeaders,
								Method: "GET",
								Path:   "/api/v1/regions/ece-region/platform/infrastructure/proxies/settings",
								Host:   api.DefaultMockHost,
							},
							mock.NewStringBody(expectedProxiesSettings),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: expectedProxiesSettings + "\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutils.RunCmdAssertion(t, tt.args, tt.want)
			tt.args.Cmd.ResetFlags()
		})
	}
}
