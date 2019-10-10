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

package cmduser

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
)

func TestCheckInputHas1ArgOr0ArgAndCurrent(t *testing.T) {
	type funcArgs struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name     string
		funcArgs funcArgs
		wantErr  bool
		err      error
	}{
		{
			name: "Succeeds if command has one argument",
			funcArgs: funcArgs{
				cmd:  new(cobra.Command),
				args: []string{"Tomasa"},
			},
			wantErr: false,
		},
		{
			name: "Returns an error if command has more than one argument",
			funcArgs: funcArgs{
				cmd:  new(cobra.Command),
				args: []string{"Tomasa", "Juana"},
			},
			wantErr: true,
			err:     errors.New(" needs 1 argument or the --current flag"),
		},
		{
			name: "Returns an error if command has no argument or --current flag",
			funcArgs: funcArgs{
				cmd:  new(cobra.Command),
				args: []string{},
			},
			wantErr: true,
			err:     errors.New(" needs 1 argument or the --current flag"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkInputHas1ArgOr0ArgAndCurrent(tt.funcArgs.cmd, tt.funcArgs.args)

			if (err != nil) != tt.wantErr {
				t.Errorf("checkInputHas1ArgOr0ArgAndCurrent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.err == nil {
				t.Errorf("checkInputHas1ArgOr0ArgAndCurrent() expected errors = '%v' but no errors returned", tt.err)
			}

			if tt.wantErr && err.Error() != tt.err.Error() {
				t.Errorf("checkInputHas1ArgOr0ArgAndCurrent() expected errors = '%v' but got %v", tt.err, err)
			}
		})
	}
}

func TestCheckInputHas1ArgOr0ArgAndCurrent_Flag(t *testing.T) {
	type funcArgs struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name     string
		funcArgs funcArgs
		wantErr  bool
		err      error
	}{
		{
			name: "Succeeds if command has no arguments and --current flag",
			funcArgs: funcArgs{
				cmd:  new(cobra.Command),
				args: []string{},
			},
			wantErr: false,
		},
		{
			name: "Returns an error if command has and argument and --current flag",
			funcArgs: funcArgs{
				cmd:  new(cobra.Command),
				args: []string{"Tomasa"},
			},
			wantErr: true,
			err:     errors.New(" needs 1 argument or the --current flag"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.funcArgs.cmd.Flags().Bool("current", true, "shows details of the current user")

			err := checkInputHas1ArgOr0ArgAndCurrent(tt.funcArgs.cmd, tt.funcArgs.args)

			if (err != nil) != tt.wantErr {
				t.Errorf("checkInputHas1ArgOr0ArgAndCurrent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.err == nil {
				t.Errorf("checkInputHas1ArgOr0ArgAndCurrent() expected errors = '%v' but no errors returned", tt.err)
			}

			if tt.wantErr && err.Error() != tt.err.Error() {
				t.Errorf("checkInputHas1ArgOr0ArgAndCurrent() expected errors = '%v' but got %v", tt.err, err)
			}
		})
	}
}
