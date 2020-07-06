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

package cmddeploymentdemplate

import (
	"errors"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/multierror"

	"github.com/elastic/ecctl/cmd/util/testutils"
)

func Test_showCmd(t *testing.T) {
	tests := []struct {
		name string
		args testutils.Args
		want testutils.Assertion
	}{
		{
			name: "incorrect amount of arguments returns error",
			args: testutils.Args{
				Cmd: showCmd,
				Args: []string{
					"show",
				},
			},
			want: testutils.Assertion{
				Err: errors.New("accepts 1 arg(s), received 0"),
			},
		},
		{
			name: "unspecified region returns error",
			args: testutils.Args{
				Cmd: showCmd,
				Args: []string{
					"show", "default", "--template-format=hi",
				},
			},
			want: testutils.Assertion{
				Err: multierror.NewPrefixed("invalid deployment template get params",
					errors.New("region not specified and is required for this operation"),
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutils.RunCmdAssertion(t, tt.args, tt.want)
		})
	}
}
