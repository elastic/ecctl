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

const updateExample = `
  * Update a platform user
	ecctl user update --username xochitl --role ece_platform_viewer --email xo@example.com
	
  * Update the current platform user
    ecctl user update --current --email xo@example.com --password
`

var updateCmd = &cobra.Command{
	Use:     "update <username> --role <role>",
	Short:   cmdutil.AdminReqDescription("Updates a platform user"),
	PreRunE: checkInputHas1ArgOr0ArgAndCurrent,
	Example: updateExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		var password []byte
		var err error

		role, _ := cmd.Flags().GetStringSlice("role")
		insecure, _ := cmd.Flags().GetString("insecure-password")
		pwFlag, _ := cmd.Flags().GetBool("password")
		name, _ := cmd.Flags().GetString("fullname")
		email, _ := cmd.Flags().GetString("email")

		if pwFlag || insecure != "" {
			message := "enter new password for user: "
			password, err = cmdutil.InsecureOrSecurePassword(insecure, message, true)
			if err != nil {
				return err
			}
		}

		if current, _ := cmd.Flags().GetBool("current"); current {
			userInfo, err := userapi.GetCurrent(userapi.GetCurrentParams{
				API: ecctl.Get().API,
			})
			if err != nil {
				return err
			}

			res, err := userapi.UpdateCurrent(userapi.UpdateParams{
				UserName: *userInfo.UserName,
				API:      ecctl.Get().API,
				Password: password,
				FullName: name,
				Email:    email,
				Roles:    role,
			})
			if err != nil {
				return err
			}

			return ecctl.Get().Formatter.Format("user/user-details", res)
		}

		res, err := userapi.Update(userapi.UpdateParams{
			UserName: args[0],
			API:      ecctl.Get().API,
			Password: password,
			FullName: name,
			Email:    email,
			Roles:    role,
		})
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format("user/user-details", res)
	},
}

func init() {
	updateCmd.Flags().Bool("current", false, "updates details of the current user")
	updateCmd.Flags().Bool("password", false, "if set, updates the user's password securely \n(must use a minimum of 8 characters )")
	updateCmd.Flags().String("fullname", "", "user's full name")
	updateCmd.Flags().String("email", "", "user's email address (must be in a valid email format)")
	updateCmd.Flags().StringSlice("role", nil, "role or roles assigned to the user. Available roles: \nece_platform_admin, ece_platform_viewer, ece_deployment_manager, ece_deployment_viewer")
	updateCmd.Flags().String("insecure-password", "", "[INSECURE] user plaintext password")
}
