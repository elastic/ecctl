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

var platformProxySettingsUpdateCmd = &cobra.Command{
	Use:     "update --file settings.json",
	Short:   cmdutil.AdminReqDescription("Updates proxies settings"),
	PreRunE: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename, _ := cmd.Flags().GetString("file")
		version, _ := cmd.Flags().GetString("version")
		patch, _ := cmd.Flags().GetBool("patch")

		if filepath.Ext(filename) != ".json" {
			return fmt.Errorf("only files with json extension are supported")
		}

		var settings models.ProxiesSettings
		if err := sdkcmdutil.DecodeFile(filename, &settings); err != nil {
			return err
		}

		var updatedSettings *models.ProxiesSettings
		var updateErr error
		if patch {
			updatedSettings, updateErr = proxysettingsapi.Patch(proxysettingsapi.UpdateParams{
				API:             ecctl.Get().API,
				Region:          ecctl.Get().Config.Region,
				ProxiesSettings: &settings,
				Version:         version,
			})
		} else {
			updatedSettings, updateErr = proxysettingsapi.Set(proxysettingsapi.UpdateParams{
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
	platformProxySettingsUpdateCmd.Flags().String("file", "", "File name containing the proxy settings")
	platformProxySettingsUpdateCmd.Flags().String("version", "", "If specified, checks for conflicts against the version of the settings")
	platformProxySettingsUpdateCmd.Flags().Bool("patch", true, "Set this to false to perform a full update and "+
		"overwrite the proxy settings - all unspecified fields are deleted.\n"+
		"When it's set to true, the operation does a partial update and only the fields that are referenced in the given file are changed")
	cobra.MarkFlagFilename(platformProxySettingsUpdateCmd.Flags(), "file", "json")
	cobra.MarkFlagRequired(platformProxySettingsUpdateCmd.Flags(), "file")
}

func init() {
	Command.AddCommand(platformProxySettingsUpdateCmd)
}
