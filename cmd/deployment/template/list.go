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

	"github.com/elastic/ecctl/pkg/ecctl"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "Lists deployment templates",
	PreRunE: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		stackVersion, _ := cmd.Flags().GetString("stack-version")
		filter, _ := cmd.Flags().GetString("filter")
		hideIC, _ := cmd.Flags().GetBool("hide-instance-configurations")

		res, err := deptemplateapi.List(deptemplateapi.ListParams{
			API:                        ecctl.Get().API,
			Region:                     ecctl.Get().Config.Region,
			MetadataFilter:             filter,
			StackVersion:               stackVersion,
			HideInstanceConfigurations: hideIC,
		})
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format("deployment-template/list", res)
	},
}

func init() {
	Command.AddCommand(listCmd)
	listCmd.Flags().Bool("hide-instance-configurations", false, "Hides instance configurations - only visible when using the JSON output")
	listCmd.Flags().String("stack-version", "", "It will cause the returned deployment templates which are valid for the specified stacks version.")
	listCmd.Flags().String("filter", "", "Optional key:value pair that acts as a filter and excludes any template that does not have a matching metadata item associated.")
}
