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
	"path/filepath"

	"github.com/elastic/cloud-sdk-go/pkg/api/platformapi/enrollmenttokenapi"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/ecctl"
)

const (
	roleFlag     = "role"
	validityFlag = "validity"
)

const tokenCreateExamples = `  ecctl [globalFlags] enrollment-token create --role coordinator
  ecctl [globalFlags] enrollment-token create --role coordinator --role proxy
  ecctl [globalFlags] enrollment-token create --role allocator --validity 120s
  ecctl [globalFlags] enrollment-token create --role allocator --validity 2h`

var createTokenCmd = &cobra.Command{
	Use:     "create --role <ROLE>",
	Short:   "Creates an enrollment token for role(s)",
	Example: tokenCreateExamples,
	PreRunE: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		roles, _ := cmd.Flags().GetStringArray(roleFlag)
		validity, _ := cmd.Flags().GetDuration(validityFlag)

		res, err := enrollmenttokenapi.Create(enrollmenttokenapi.CreateParams{
			API:      ecctl.Get().API,
			Region:   ecctl.Get().Config.Region,
			Roles:    roles,
			Duration: validity,
		})
		if err != nil {
			return err
		}
		return ecctl.Get().Formatter.Format(filepath.Join("token", "create"), res)
	},
}

func init() {
	Command.AddCommand(createTokenCmd)

	createTokenCmd.Flags().StringArrayP(roleFlag, "r", nil, "Role(s) to associate the tokens with.")
	createTokenCmd.Flags().DurationP(validityFlag, "v", 0, "Time in seconds for which this token is valid. Currently this will make the token ephemeral (persistent: false)")
	cobra.MarkFlagRequired(createTokenCmd.Flags(), "role")
}
