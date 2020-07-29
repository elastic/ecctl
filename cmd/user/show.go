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
	"fmt"

	"github.com/elastic/cloud-sdk-go/pkg/api/userapi"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var showCmd = &cobra.Command{
	Use:     "show <user name>",
	Short:   cmdutil.AdminReqDescription("Shows details of a specified user"),
	PreRunE: checkInputHas1ArgOr0ArgAndCurrent,
	RunE: func(cmd *cobra.Command, args []string) error {
		if current, _ := cmd.Flags().GetBool("current"); current {
			res, err := userapi.GetCurrent(userapi.GetCurrentParams{
				API: ecctl.Get().API,
			})
			if err != nil {
				return err
			}
			return ecctl.Get().Formatter.Format("user/user-details", res)
		}

		res, err := userapi.Get(userapi.GetParams{
			API:      ecctl.Get().API,
			UserName: args[0],
		})
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format("user/user-details", res)
	},
}

func checkInputHas1ArgOr0ArgAndCurrent(cmd *cobra.Command, args []string) error {
	var currentFlag, _ = cmd.Flags().GetBool("current")
	if len(args) == 1 && !currentFlag {
		return nil
	}

	if len(args) == 0 && currentFlag {
		return nil
	}

	return fmt.Errorf("%s needs 1 argument or the --current flag", cmd.Name())
}

func init() {
	showCmd.Flags().Bool("current", false, "Shows details of the current user")
}
