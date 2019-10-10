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
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	errInvalidUUID = errors.New("ID is invalid")
)

// MinimumNArgsAndUUID ensures that the command has at least N number of
// arguments and the first argument is 32 characters long.
func MinimumNArgsAndUUID(argsCount int) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(argsCount)(cmd, args); err != nil {
			return err
		}
		if len(args[0]) != 32 {
			cmd.Help()
			return errInvalidUUID
		}
		return nil
	}
}

// CheckInputHas2ArgsOr1ArgAndAll checks that the input has either:
// * 2 arguments
// * 1 argument and the --all flag.
func CheckInputHas2ArgsOr1ArgAndAll(cmd *cobra.Command, args []string) error {
	var allFlag, _ = cmd.Flags().GetBool("all")
	if len(args) == 2 || (len(args) == 1 && allFlag) {
		return nil
	}
	return fmt.Errorf("%s needs 2 arguments or 1 argument and the --all flag", cmd.Name())
}

// CheckInputHas1ArgsOr0ArgAndAll checks that the input has either:
// * 1 argument
// * 0 arguments and the --all flag.
func CheckInputHas1ArgsOr0ArgAndAll(cmd *cobra.Command, args []string) error {
	var allFlag, _ = cmd.Flags().GetBool("all")
	if len(args) == 1 || (allFlag && len(args) == 0) {
		return nil
	}
	return fmt.Errorf("%s needs 1 arguments or --all flag", cmd.Name())
}
