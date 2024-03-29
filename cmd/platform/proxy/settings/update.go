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

package cmdproxysettings

import (
	"fmt"
	"path/filepath"

	proxysettingsapi "github.com/elastic/cloud-sdk-go/pkg/api/platformapi/proxyapi/settingsapi"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
)

const examples = `## Update only the defined proxy settings 
$ ecctl platform proxy settings update --file settings.json --region us-east-1

## Update by overriding all proxy settings 
$ ecctl platform proxy settings update --file settings.json --region us-east-1 --full
`

var platformProxySettingsUpdateCmd = &cobra.Command{
	Use:     "update --file settings.json",
	Short:   cmdutil.AdminReqDescription("Updates settings for all proxies"),
	Example: examples,
	PreRunE: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename, _ := cmd.Flags().GetString("file")
		version, _ := cmd.Flags().GetString("version")
		full, _ := cmd.Flags().GetBool("full")

		if filepath.Ext(filename) != ".json" {
			return fmt.Errorf("only files with json extension are supported")
		}

		var settings models.ProxiesSettings
		if err := sdkcmdutil.DecodeFile(filename, &settings); err != nil {
			return err
		}

		var updatedSettings *models.ProxiesSettings
		var updateErr error
		if full {
			updatedSettings, updateErr = proxysettingsapi.Set(proxysettingsapi.UpdateParams{
				API:             ecctl.Get().API,
				Region:          ecctl.Get().Config.Region,
				ProxiesSettings: &settings,
				Version:         version,
			})
		} else {
			updatedSettings, updateErr = proxysettingsapi.Patch(proxysettingsapi.UpdateParams{
				API:             ecctl.Get().API,
				Region:          ecctl.Get().Config.Region,
				ProxiesSettings: &settings,
				Version:         version,
			})
		}

		if updateErr != nil {
			return updateErr
		}
		return ecctl.Get().Formatter.Format("", updatedSettings)
	},
}

func init() {
	initUpdateFlags()
}

func initUpdateFlags() {
	Command.AddCommand(platformProxySettingsUpdateCmd)
	platformProxySettingsUpdateCmd.Flags().StringP("file", "f", "", "ProxiesSettings file definition. See https://www.elastic.co/guide/en/cloud-enterprise/current/ProxiesSettings.html for more information.")
	platformProxySettingsUpdateCmd.Flags().String("version", "", "If specified, checks for conflicts against the version of the repository configuration")
	platformProxySettingsUpdateCmd.Flags().Bool("full", false, "If set, a full update will be performed and all proxy settings will be overwritten. Any unspecified fields will be deleted.")
	cobra.MarkFlagFilename(platformProxySettingsUpdateCmd.Flags(), "file", "json")
	cobra.MarkFlagRequired(platformProxySettingsUpdateCmd.Flags(), "file")
}

func init() {
	Command.AddCommand(platformProxySettingsUpdateCmd)
}
