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
	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/ecctl"
)

var createCmd = &cobra.Command{
	Use:     "create -f <template file>.json",
	Short:   "Creates a platform deployment template",
	PreRunE: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := sdkcmdutil.FileOrStdin(cmd, "file-template"); err != nil {
			return err
		}

		tc, err := parseTemplateFile(cmd.Flag("file-template").Value.String())
		if err != nil {
			return err
		}

		if id := cmd.Flag("id").Value.String(); id != "" {
			tc.ID = id
		}

		tid, err := configurationtemplateapi.CreateTemplate(configurationtemplateapi.CreateTemplateParams{
			Region:                 ecctl.Get().Config.Region,
			API:                    ecctl.Get().API,
			ID:                     tc.ID,
			DeploymentTemplateInfo: tc,
		})

		if err != nil {
			return err
		}

		fmt.Printf("Successfully created deployment template %v \n", tid)
		return nil

	},
}

func init() {
	Command.AddCommand(createCmd)
	createCmd.Flags().StringP("file-template", "f", "", "YAML or JSON file that contains the deployment template configuration")
	createCmd.Flags().String("id", "", "Optional ID to set for the deployment template (Overrides ID if present)")
}
