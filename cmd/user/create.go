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
	"github.com/elastic/cloud-sdk-go/pkg/api/userapi"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
)

const createExample = `
  * Create a platform user who has two roles assigned
    ecctl user create --username sam89 --role ece_platform_viewer --role ece_deployment_viewer
`

var createCmd = &cobra.Command{
	Use:     "create --username <username> --role <role>",
	Short:   cmdutil.AdminReqDescription("Creates a new platform user"),
	PreRunE: cobra.MinimumNArgs(0),
	Example: createExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		insecure := cmd.Flag("insecure-password").Value.String()
		message := "enter new password for user: "
		password, err := cmdutil.InsecureOrSecurePassword(insecure, message, true)
		if err != nil {
			return err
		}

		role, _ := cmd.Flags().GetStringSlice("role")

		res, err := userapi.Create(userapi.CreateParams{
			API:      ecctl.Get().API,
			Password: password,
			UserName: cmd.Flag("username").Value.String(),
			FullName: cmd.Flag("fullname").Value.String(),
			Email:    cmd.Flag("email").Value.String(),
			Roles:    role,
		})
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format("user/user-details", res)
	},
}

func init() {
	createCmd.Flags().String("username", "", "Unique username for the platform user")
	createCmd.Flags().String("fullname", "", "User's full name")
	createCmd.Flags().String("email", "", "User's email address (must be in a valid email format)")
	createCmd.Flags().StringSlice("role", nil, "Role or roles assigned to the user. Available roles: \nece_platform_admin, ece_platform_viewer, ece_deployment_manager, ece_deployment_viewer")
	createCmd.Flags().String("insecure-password", "", "[INSECURE] User's plaintext password")
	cobra.MarkFlagRequired(createCmd.Flags(), "username")
	cobra.MarkFlagRequired(createCmd.Flags(), "role")
}
