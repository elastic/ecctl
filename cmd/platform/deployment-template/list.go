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

package cmddeploymentdemplate

import (
	"path/filepath"

	"github.com/elastic/cloud-sdk-go/pkg/api/platformapi/configurationtemplateapi"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/ecctl"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "Lists the platform deployment templates",
	PreRunE: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		showInstanceConfig, _ := cmd.Flags().GetBool(showInstanceConfigurations)
		stackVersion, _ := cmd.Flags().GetString(stackVersion)
		metadataFilter, _ := cmd.Flags().GetString(filter)
		format, _ := cmd.Flags().GetString(format)

		res, err := configurationtemplateapi.ListTemplates(configurationtemplateapi.ListTemplateParams{
			API:                ecctl.Get().API,
			Region:             ecctl.Get().Config.Region,
			ShowInstanceConfig: showInstanceConfig,
			StackVersion:       stackVersion,
			Metadata:           metadataFilter,
			Format:             format,
		})

		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format(filepath.Join("deployment-template", "list"), res)
	},
}

func init() {
	Command.AddCommand(listCmd)
	listCmd.Flags().BoolP(showInstanceConfigurations, "", false, "Shows instance configurations - only visible when using the JSON output")
	listCmd.Flags().String(stackVersion, "", "If present, it will cause the returned deployment templates to be adapted to return only the elements allowed in that version.")
	listCmd.Flags().String(filter, "", "Optional key:value pair that acts as a filter and excludes any template that does not have a matching metadata item associated.")
	listCmd.Flags().String(format, "cluster", "If 'deployment' is specified, the deployment_template is populated in the response. If 'cluster' is specified, the cluster_template is populated in the response.")
}
