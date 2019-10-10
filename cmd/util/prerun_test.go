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

package cmdutil

import (
	"testing"

	"github.com/spf13/cobra"
)

const validClusterID = "8b43ed5e277f7ea6f13606fcf4027f9c"

func TestMinimumNArgsAndUUID(t *testing.T) {
	type args struct {
		argsCount int
	}
	type funcArgs struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name    string
		args    args
		fargs   funcArgs
		wantErr bool
	}{
		{
			name: "Test 1 arg",
			args: args{
				argsCount: 1,
			},
			fargs: funcArgs{
				cmd:  new(cobra.Command),
				args: []string{validClusterID},
			},
			wantErr: false,
		},
		{
			name: "Test 1 arg Invalid ID",
			args: args{
				argsCount: 1,
			},
			fargs: funcArgs{
				cmd:  new(cobra.Command),
				args: []string{"invalid ID"},
			},
			wantErr: true,
		},
		{
			name: "Test 0 arg fails",
			args: args{
				argsCount: 1,
			},
			fargs: funcArgs{
				cmd: new(cobra.Command),
			},
			wantErr: true,
		},
		{
			name: "Test 5 arg succeeds",
			args: args{
				argsCount: 5,
			},
			fargs: funcArgs{
				cmd:  new(cobra.Command),
				args: []string{validClusterID, "b", "c", "d", "e"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MinimumNArgsAndUUID(tt.args.argsCount)(tt.fargs.cmd, tt.fargs.args); (err != nil) != tt.wantErr {
				t.Errorf("MinimumNArgsAndUUID() = returned %v, when want %v", err, tt.wantErr)
			}
		})
	}
}
