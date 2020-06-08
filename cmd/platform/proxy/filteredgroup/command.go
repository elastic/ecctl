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

package cmdfilteredgroup

import (
	"path/filepath"

	"github.com/elastic/cloud-sdk-go/pkg/api/platformapi/proxy/filteredgroup"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/ecctl"
)

// Command represents the top level filtered-group command.
var Command = &cobra.Command{
	Use:     "filtered-group",
	Short:   "Manages proxies filtered group",
	PreRunE: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var platformProxyFilteredGroupShowCmd = &cobra.Command{
	Use:     "show <filtered group id>",
	Short:   "Shows details for proxies filtered group",
	PreRunE: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		response, err := filteredgroup.Get(filteredgroup.CommonParams{
			API: ecctl.Get().API,
			ID:  args[0],
		})

		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format(filepath.Join("filtered-group", "show"), response)
	},
}

var platformProxyFilteredGroupCreateCmd = &cobra.Command{
	Use:     "create <filtered group id> --filters <key1=value1,key2=value2> --expected-proxies-count <int>",
	Short:   "Creates proxies filtered group",
	PreRunE: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		filters, _ := cmd.Flags().GetStringToString("filters")
		expectedProxiesCount, _ := cmd.Flags().GetInt32("expected-proxies-count")

		response, err := filteredgroup.Create(filteredgroup.CreateParams{
			CommonParams: filteredgroup.CommonParams{
				API: ecctl.Get().API,
				ID:  args[0],
			},
			Filters:              filters,
			ExpectedProxiesCount: expectedProxiesCount,
		})

		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format(filepath.Join("filtered-group", "create"), response)
	},
}

var platformProxyFilteredGroupDeleteCmd = &cobra.Command{
	Use:     "delete <filtered group id>",
	Short:   "Deletes proxies filtered group",
	PreRunE: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return filteredgroup.Delete(filteredgroup.CommonParams{
			API: ecctl.Get().API,
			ID:  args[0],
		})
	},
}

var platformProxyFilteredGroupUpdateCmd = &cobra.Command{
	Use:     "update <filtered group id> --filters <key1=value1,key2=value2> --expected-proxies-count <int> --version <int>",
	Short:   "Updates proxies filtered group",
	PreRunE: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		filters, _ := cmd.Flags().GetStringToString("filters")
		expectedProxiesCount, _ := cmd.Flags().GetInt32("expected-proxies-count")
		version, _ := cmd.Flags().GetInt64("version")

		response, err := filteredgroup.Update(filteredgroup.UpdateParams{
			CreateParams: filteredgroup.CreateParams{
				CommonParams: filteredgroup.CommonParams{
					API: ecctl.Get().API,
					ID:  args[0],
				},
				Filters:              filters,
				ExpectedProxiesCount: expectedProxiesCount,
			},
			Version: version,
		})

		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format(filepath.Join("filtered-group", "update"), response)
	},
}

func init() {
	Command.AddCommand(
		platformProxyFilteredGroupShowCmd,
		platformProxyFilteredGroupCreateCmd,
		platformProxyFilteredGroupDeleteCmd,
		platformProxyFilteredGroupUpdateCmd,
	)
	platformProxyFilteredGroupCreateCmd.Flags().StringToString("filters", make(map[string]string), "Filters for proxies group")
	platformProxyFilteredGroupCreateCmd.Flags().Int32("expected-proxies-count", 1, "Expected proxies count in filtered group")
	platformProxyFilteredGroupUpdateCmd.Flags().StringToString("filters", make(map[string]string), "fFlters for proxies group")
	platformProxyFilteredGroupUpdateCmd.Flags().Int32("expected-proxies-count", 1, "Expected proxies count in filtered group")
	platformProxyFilteredGroupUpdateCmd.Flags().Int64("version", 0, "Version update for filtered group")
}
