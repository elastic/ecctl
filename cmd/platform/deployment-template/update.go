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
	"fmt"

	"github.com/elastic/cloud-sdk-go/pkg/api/platformapi/configurationtemplateapi"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var updateCmd = &cobra.Command{
	Use:     "update <template id> -f <template file>.json",
	Short:   cmdutil.DeprecatedDescription("Updates a platform deployment template"),
	PreRunE: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tc, err := parseTemplateFile(cmd.Flag("file-template").Value.String())
		if err != nil {
			return err
		}

		if err := configurationtemplateapi.UpdateTemplate(
			configurationtemplateapi.UpdateTemplateParams{
				API:                    ecctl.Get().API,
				Region:                 ecctl.Get().Config.Region,
				ID:                     args[0],
				DeploymentTemplateInfo: tc,
			},
		); err != nil {
			return err
		}

		fmt.Printf("Successfully updated deployment template %v \n", args[0])
		return nil

	},
}

func init() {
	Command.AddCommand(updateCmd)
	updateCmd.Flags().StringP("file-template", "f", "", "YAML or JSON file that contains the deployment template configuration")
}
