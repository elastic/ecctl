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

	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"

	"github.com/elastic/ecctl/cmd/util/testutils"
)

func Test_updateCmd(t *testing.T) {
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
			name: "fails due to provided arguments",
			args: testutils.Args{
				Cmd:  platformProxySettingsUpdateCmd,
				Args: []string{"update", "some-arg"},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: `accepts 0 arg(s), received 1`,
			},
		},
		{
			name: "fails due to missing file option",
			args: testutils.Args{
				Cmd:  platformProxySettingsUpdateCmd,
				Args: []string{"update"},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: `required flag(s) "file" not set`,
			},
		},
		{
			name: "fails due to wrong filename provided",
			args: testutils.Args{
				Cmd:  platformProxySettingsUpdateCmd,
				Args: []string{"update", "--file", "./testfiles/not-a-json-file.txt"},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: `only files with json extension are supported`,
			},
		},
		{
			name: "fails due to invalid json contents",
			args: testutils.Args{
				Cmd:  platformProxySettingsUpdateCmd,
				Args: []string{"update", "--file", "./testfiles/settings-invalid.json"},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: `invalid character 'h' in literal true (expecting 'r')`,
			},
		},
		{
			name: "fails due to API error",
			args: testutils.Args{
				Cmd: platformProxySettingsUpdateCmd,
				Args: []string{
					"update", "--file", "./testfiles/settings.json",
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
			name: "succeeds (patch)",
			args: testutils.Args{
				Cmd: platformProxySettingsUpdateCmd,
				Args: []string{
					"update", "--file", "./testfiles/settings.json",
				},
				Cfg: testutils.MockCfg{
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultWriteMockHeaders,
								Method: "PATCH",
								Path:   "/api/v1/regions/ece-region/platform/infrastructure/proxies/settings",
								Host:   api.DefaultMockHost,
								Body:   mock.NewStructBody(proxySettings()),
							},
							mock.NewStructBody(proxySettings()),
						),
					},
				},
			},
			want: testutils.Assertion{
				Stdout: expectedProxiesSettings + "\n",
			},
		},
		{
			name: "succeeds (put)",
			args: testutils.Args{
				Cmd: platformProxySettingsUpdateCmd,
				Args: []string{
					"update", "--file", "./testfiles/settings.json", "--patch=false",
				},
				Cfg: testutils.MockCfg{
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							&mock.RequestAssertion{
								Header: api.DefaultWriteMockHeaders,
								Method: "PUT",
								Path:   "/api/v1/regions/ece-region/platform/infrastructure/proxies/settings",
								Host:   api.DefaultMockHost,
								Body:   mock.NewStructBody(proxySettings()),
							},
							mock.NewStructBody(proxySettings()),
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
			defer initUpdateFlags()
		})
	}
}

func proxySettings() *models.ProxiesSettings {
	return &models.ProxiesSettings{
		ExpectedProxiesCount: ec.Int32(5),
		HTTPSettings: &models.ProxiesHTTPSettings{
			CookieSecret:         ec.String("some-secret"),
			DashboardsBaseURL:    ec.String("some-url"),
			DisconnectedCutoff:   ec.Int64(2),
			MinimumProxyServices: ec.Int32(1),
			SsoSettings: &models.ProxiesSSOSettings{
				CookieName:                  ec.String("some-cookie"),
				DefaultRedirectPath:         ec.String("some-path"),
				DontLogRequests:             ec.Bool(true),
				MaintenanceBypassCookieName: ec.String("some-other-cookie"),
				MaxAge:                      ec.Int64(10),
				SsoSecret:                   ec.String("sso-secret"),
			},
			UserCookieKey: ec.String("user-cookie"),
		},
		SignatureSecret:         ec.String("signature-secret"),
		SignatureValidForMillis: ec.Int64(60),
	}
}
