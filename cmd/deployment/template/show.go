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

package cmddeploymenttemplate

import (
	"github.com/elastic/cloud-sdk-go/pkg/api/deploymentapi/deptemplateapi"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var showCmd = &cobra.Command{
	Use:     "show --template-id <template id>",
	Short:   cmdutil.AdminReqDescription("Displays a deployment template"),
	PreRunE: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		stackVersion, _ := cmd.Flags().GetString("stack-version")
		hideIC, _ := cmd.Flags().GetBool("hide-instance-configurations")
		templateID, _ := cmd.Flags().GetString("template-id")
		res, err := deptemplateapi.Get(deptemplateapi.GetParams{
			API:                        ecctl.Get().API,
			Region:                     ecctl.Get().Config.Region,
			StackVersion:               stackVersion,
			HideInstanceConfigurations: hideIC,
			TemplateID:                 templateID,
		})
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format("", res)
	},
}

func init() {
	Command.AddCommand(showCmd)
	showCmd.Flags().Bool("hide-instance-configurations", false, "Hides instance configurations - only visible when using the JSON output")
	showCmd.Flags().String("stack-version", "", "Optional filter to only return deployment templates which are valid for the specified stack version.")
	showCmd.Flags().String("template-id", "", "Required template ID to update.")
	showCmd.MarkFlagRequired("template-id")
}
