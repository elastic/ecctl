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

var showCmd = &cobra.Command{
	Use:     "show <template id>",
	Short:   "Shows information about a specific platform deployment template",
	PreRunE: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		show, _ := cmd.Flags().GetBool("show-instance-configurations")
		format, _ := cmd.Flags().GetString(format)

		res, err := configurationtemplateapi.GetTemplate(configurationtemplateapi.GetTemplateParams{
			API:                ecctl.Get().API,
			Region:             ecctl.Get().Config.Region,
			ID:                 args[0],
			ShowInstanceConfig: show,
			Format:             format,
		})

		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format(filepath.Join("deployment-template", "show"), res)
	},
}

func init() {
	Command.AddCommand(showCmd)
	showCmd.Flags().BoolP(showInstanceConfigurations, "", false, "Shows instance configurations")
	showCmd.Flags().String(format, "deployment", "If deployment is specified deployment_template is populated in the response, If cluster is specified cluster_template is populated in the response. (Defaults to deployment)")
}
