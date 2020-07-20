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
	"errors"
	"net/url"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"

	"github.com/elastic/ecctl/cmd/util/testutils"
)

func Test_deleteCmd(t *testing.T) {
	var reqAssertion = &mock.RequestAssertion{
		Header: api.DefaultWriteMockHeaders,
		Method: "DELETE",
		Path:   "/api/v1/deployments/templates/default",
		Host:   api.DefaultMockHost,
		Query: url.Values{
			"region": []string{"ece-region"},
		},
	}
	tests := []struct {
		name string
		args testutils.Args
		want testutils.Assertion
	}{
		{
			name: "fails due empty argument",
			args: testutils.Args{
				Cmd: deleteCmd,
				Args: []string{
					"delete",
				},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: errors.New(`required flag(s) "template-id" not set`),
			},
		},
		{
			name: "fails due to API error",
			args: testutils.Args{
				Cmd: deleteCmd,
				Args: []string{
					"delete", "--template-id=default",
				},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					mock.SampleInternalError(),
				}},
			},
			want: testutils.Assertion{
				Err: mock.MultierrorInternalError,
			},
		},
		{
			name: "succeeds with JSON format",
			args: testutils.Args{
				Cmd: deleteCmd,
				Args: []string{
					"delete", "--template-id=default",
				},
				Cfg: testutils.MockCfg{
					OutputFormat: "json",
					Responses: []mock.Response{
						mock.New200ResponseAssertion(
							reqAssertion, mock.NewStringBody("{}"),
						),
					},
				},
			},
			want: testutils.Assertion{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutils.RunCmdAssertion(t, tt.args, tt.want)
		})
	}
}
