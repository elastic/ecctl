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

package cmdinstanceconfig

import (
	"os"
	"path/filepath"

	"github.com/elastic/cloud-sdk-go/pkg/input"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/platform/instanceconfig"
)

// Command is the top instance-config subcommand.
var Command = &cobra.Command{
	Use:     "instance-configuration",
	Short:   "Manages instance configurations",
	PreRunE: cobra.MaximumNArgs(0),
	Run:     func(cmd *cobra.Command, args []string) { cmd.Help() },
}

var platformInstanceConfigurationListCmd = &cobra.Command{
	Use:     "list",
	Short:   "Lists the instance configurations",
	PreRunE: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := instanceconfig.List(instanceconfig.ListParams{
			API: ecctl.Get().API,
		})

		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format(filepath.Join("instance-configuration", "list"), res)
	},
}

var platformInstanceConfigurationCreateCmd = &cobra.Command{
	Use:     "create -f <config.json>",
	Short:   "Creates a new instance configuration",
	PreRunE: cobra.MaximumNArgs(0),

	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmdutil.FileOrStdin(cmd, "file"); err != nil {
			return err
		}

		file, err := input.NewFileOrReader(os.Stdin, cmd.Flag("file").Value.String())
		if err != nil {
			return err
		}
		defer file.Close()

		config, err := instanceconfig.NewConfig(file)
		if err != nil {
			return err
		}

		if id := cmd.Flag("id").Value.String(); id != "" {
			config.ID = id
		}

		res, err := instanceconfig.Create(instanceconfig.CreateParams{
			API:    ecctl.Get().API,
			Config: config,
		})

		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format(filepath.Join("instance-configuration", "create"), res)
	},
}

var platformInstanceConfigurationShowCmd = &cobra.Command{
	Use:     "show <config id>",
	Short:   "Shows an instance configuration",
	PreRunE: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := instanceconfig.Get(instanceconfig.GetParams{
			API: ecctl.Get().API,
			ID:  args[0],
		})

		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format(filepath.Join("instance-configuration", "show"), res)
	},
}

var platformInstanceConfigurationDeleteCmd = &cobra.Command{
	Use:     "delete <config id>",
	Short:   "Deletes an instance configuration",
	PreRunE: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return instanceconfig.Delete(instanceconfig.DeleteParams{
			API: ecctl.Get().API,
			ID:  args[0],
		})
	},
}

var platformInstanceConfigurationUpdateCmd = &cobra.Command{
	Use:     "update <config id> -f <config.json>",
	Short:   "Overwrites an instance configuration",
	PreRunE: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmdutil.FileOrStdin(cmd, "file"); err != nil {
			return err
		}

		file, err := input.NewFileOrReader(os.Stdin, cmd.Flag("file").Value.String())
		if err != nil {
			return err
		}
		defer file.Close()

		config, err := instanceconfig.NewConfig(file)
		if err != nil {
			return err
		}

		return instanceconfig.Update(instanceconfig.UpdateParams{
			API:    ecctl.Get().API,
			ID:     args[0],
			Config: config,
		})
	},
}

var platformInstanceConfigurationPullCmd = &cobra.Command{
	Use:     "pull --path <path>",
	Short:   "Downloads instance configuration into a local folder",
	PreRunE: cobra.MaximumNArgs(0),

	RunE: func(cmd *cobra.Command, args []string) error {

		return instanceconfig.PullToFolder(instanceconfig.PullToFolderParams{
			API:    ecctl.Get().API,
			Folder: cmd.Flag("path").Value.String(),
		})
	},
}

func init() {
	Command.AddCommand(
		platformInstanceConfigurationListCmd,
		platformInstanceConfigurationCreateCmd,
		platformInstanceConfigurationShowCmd,
		platformInstanceConfigurationDeleteCmd,
		platformInstanceConfigurationUpdateCmd,
		platformInstanceConfigurationPullCmd,
	)

	platformInstanceConfigurationCreateCmd.Flags().StringP("file", "f", "", "Instance configuration JSON file definition")
	platformInstanceConfigurationCreateCmd.Flags().String("id", "", "Optional ID to set for the instance configuration (Overrides id if present)")
	platformInstanceConfigurationUpdateCmd.Flags().StringP("file", "f", "", "Instance configuration JSON file definition")

	platformInstanceConfigurationPullCmd.Flags().StringP("path", "p", "", "Local path with instance configuration.")
	platformInstanceConfigurationPullCmd.MarkFlagRequired("path")
}
