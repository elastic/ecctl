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

package cmdkibana

import (
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/cmd/util/testutils"
)

func TestRunShowKibanaClusterCmd(t *testing.T) {
	if err := testutils.MockInitApp(); err != nil {
		fmt.Println(err.Error())
		return
	}

	type args struct {
		cmd  *cobra.Command
		args []string
	}

	tests := []struct {
		name     string
		args     args
		logsFlag bool
		err      string
	}{
		{
			name: "Returns an error when input is invalid",
			args: args{
				cmd:  showKibanaClusterCmd,
				args: []string{""},
			},
			err: `id "" is invalid`,
		},
		{
			name: "Returns the expected valid api call when command is run without specifying flags",
			args: args{
				cmd:  showKibanaClusterCmd,
				args: []string{"d558cdf210dc4737960906a83b1db52c"},
			},
			err: "Get http://localhost/api/v1/clusters/kibana/d558cdf210dc4737960906a83b1db52c?convert_legacy_plans=false&show_metadata=false&show_plan_defaults=true&show_plan_logs=false&show_plans=false&show_settings=false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := runShowKibanaClusterCmd(tt.args.cmd, tt.args.args)
			if !strings.Contains(err.Error(), tt.err) {
				t.Errorf("runShowKibanaClusterCmd() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}

func TestRunShowKibanaClusterCmd_Flags(t *testing.T) {
	if err := testutils.MockInitApp(); err != nil {
		fmt.Println(err.Error())
		return
	}

	type args struct {
		cmd  *cobra.Command
		args []string
	}

	tests := []struct {
		name     string
		args     args
		logsFlag bool
		err      string
	}{
		{
			name: "Sets both show_plan_defaults and show_plan_logs as true when --logs flag is set",
			args: args{
				cmd:  showKibanaClusterCmd,
				args: []string{"d558cdf210dc4737960906a83b1db52c"},
			},
			logsFlag: true,
			err:      "show_plan_logs=true&show_plans=true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.logsFlag {
				tt.args.cmd.ResetFlags()
				tt.args.cmd.Flags().BoolP("logs", "l", true, "View the cluster plan logs")
			}

			err := runShowKibanaClusterCmd(tt.args.cmd, tt.args.args)
			if !strings.Contains(err.Error(), tt.err) {
				t.Errorf("runShowKibanaClusterCmd() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
