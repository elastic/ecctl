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

package cmdenrollmenttoken

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
	enrollmenttoken "github.com/elastic/ecctl/pkg/platform/enrollment-token"
)

const (
	roleFlag     = "role"
	validityFlag = "validity"
)

// Command represents the enrollment-token subcomand.
var Command = &cobra.Command{
	Use:     "enrollment-token",
	Short:   fmt.Sprintf("Manages tokens %v", cmdutil.PlatformAdminRequired),
	PreRunE: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

const tokenCreateExamples = `  ecctl [globalFlags] enrollment-token create --role coordinator
  ecctl [globalFlags] enrollment-token create --role coordinator --role proxy
  ecctl [globalFlags] enrollment-token create --role allocator --validity 120s
  ecctl [globalFlags] enrollment-token create --role allocator --validity 2h`

var createTokenCmd = &cobra.Command{
	Use:     "create --role <ROLE>",
	Short:   "Creates an enrollment token for role(s)",
	Example: tokenCreateExamples,
	PreRunE: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		roles, _ := cmd.Flags().GetStringArray(roleFlag)
		validity, _ := cmd.Flags().GetDuration(validityFlag)

		res, err := enrollmenttoken.Create(enrollmenttoken.CreateParams{
			API:      ecctl.Get().API,
			Roles:    roles,
			Duration: validity,
		})
		if err != nil {
			return err
		}
		return ecctl.Get().Formatter.Format(filepath.Join("token", "create"), res)
	},
}

var deleteTokenCmd = &cobra.Command{
	Use:     "delete <enrollment-token>",
	Short:   "Deletes an enrollment token",
	PreRunE: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := enrollmenttoken.Delete(enrollmenttoken.DeleteParams{
			API:   ecctl.Get().API,
			Token: args[0],
		})
		if err != nil {
			return err
		}
		fmt.Printf("Token %s deleted.\n", args[0])
		return nil
	},
}

func listTokens(cmd *cobra.Command, args []string) error {
	res, err := enrollmenttoken.List(enrollmenttoken.ListParams{
		API: ecctl.Get().API,
	})
	if err != nil {
		return err
	}
	return ecctl.Get().Formatter.Format(filepath.Join("token", "list"), res)
}

var listTokensCmd = &cobra.Command{
	Use:     "list",
	Short:   "Retrieves a list of persistent enrollment tokens",
	PreRunE: cobra.MaximumNArgs(0),
	RunE:    listTokens,
}

func init() {
	Command.AddCommand(
		createTokenCmd,
		deleteTokenCmd,
		listTokensCmd,
	)

	createTokenCmd.Flags().StringArrayP(roleFlag, "r", nil, "Role(s) to associate the tokens with.")
	cobra.MarkFlagRequired(createTokenCmd.Flags(), "role")

	createTokenCmd.Flags().DurationP(validityFlag, "v", 0, "Time in seconds for which this token is valid. Currently this will make the token ephemeral (persistent: false)")
}
