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
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/elastic/cloud-sdk-go/pkg/input"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/platform/deployment"
)

const (
	showInstanceConfigurations = "show-instance-configurations"
	stackVersion               = "stack-version"
	filter                     = "filter"
)

// Command represents the top level deployment-template command.
var Command = &cobra.Command{
	Use:     "deployment-template",
	Short:   "Manages deployment templates",
	PreRunE: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var platformDeploymentTemplateListCmd = &cobra.Command{
	Use:     "list",
	Short:   "Lists the platform deployment templates",
	PreRunE: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		showInstanceConfig, _ := cmd.Flags().GetBool(showInstanceConfigurations)
		stackVersion, _ := cmd.Flags().GetString(stackVersion)
		metadataFilter, _ := cmd.Flags().GetString(filter)

		res, err := deployment.ListTemplates(deployment.ListTemplateParams{
			API:                ecctl.Get().API,
			ShowInstanceConfig: showInstanceConfig,
			StackVersion:       stackVersion,
			Metadata:           metadataFilter,
		})

		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format(filepath.Join("deployment-template", "list"), res)
	},
}

var platformDeploymentTemplateShowCmd = &cobra.Command{
	Use:     "show <template id>",
	Short:   "Shows information about a specific platform deployment template",
	PreRunE: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		show, _ := cmd.Flags().GetBool("show-instance-configurations")
		res, err := deployment.GetTemplate(deployment.GetTemplateParams{
			TemplateParams: deployment.TemplateParams{
				API: ecctl.Get().API,
				ID:  args[0],
			},
			ShowInstanceConfig: show,
		})

		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format(filepath.Join("deployment-template", "show"), res)
	},
}

var platformDeploymentTemplateDeleteCmd = &cobra.Command{
	Use:     "delete <template id>",
	Short:   "Deletes a specific platform deployment template",
	PreRunE: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := deployment.DeleteTemplate(deployment.GetTemplateParams{
			TemplateParams: deployment.TemplateParams{
				API: ecctl.Get().API,
				ID:  args[0],
			},
		})

		if err != nil {
			return err
		}

		fmt.Printf("Successfully deleted deployment template %v \n", args[0])
		return nil
	},
}

var platformDeploymentTemplateCreateCmd = &cobra.Command{
	Use:     "create -f <template file>.json",
	Short:   "Creates a platform deployment template",
	PreRunE: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmdutil.FileOrStdin(cmd, "file-template"); err != nil {
			return err
		}

		tc, err := parseTemplateFile(cmd.Flag("file-template").Value.String())
		if err != nil {
			return err
		}

		if id := cmd.Flag("id").Value.String(); id != "" {
			tc.ID = id
		}

		tid, err := deployment.CreateTemplate(deployment.CreateTemplateParams{
			DeploymentTemplateInfo: tc,
			API:                    ecctl.Get().API,
		})

		if err != nil {
			return err
		}

		fmt.Printf("Successfully created deployment template %v \n", tid)
		return nil

	},
}

var platformDeploymentTemplateUpdateCmd = &cobra.Command{
	Use:     "update <template id> -f <template file>.json",
	Short:   "Updates a platform deployment template",
	PreRunE: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmdutil.FileOrStdin(cmd, "file-template"); err != nil {
			return err
		}
		tc, err := parseTemplateFile(cmd.Flag("file-template").Value.String())
		if err != nil {
			return err
		}
		if err := deployment.UpdateTemplate(
			deployment.UpdateTemplateParams{
				TemplateParams: deployment.TemplateParams{
					API: ecctl.Get().API,
					ID:  args[0],
				},

				DeploymentTemplateInfo: tc,
			},
		); err != nil {
			return err
		}

		fmt.Printf("Successfully updated deployment template %v \n", args[0])
		return nil

	},
}

var platformDeploymentTemplatePullCmd = &cobra.Command{
	Use:     "pull --path <path>",
	Short:   "Downloads deployment template into a local folder",
	PreRunE: cobra.MaximumNArgs(0),

	RunE: func(cmd *cobra.Command, args []string) error {
		return deployment.PullToFolder(deployment.PullTemplateToFolderParams{
			TemplateToFolderParams: deployment.TemplateToFolderParams{
				API:    ecctl.Get().API,
				Folder: cmd.Flag("path").Value.String(),
			},
		})
	},
}

func parseTemplateFile(fp string) (*models.DeploymentTemplateInfo, error) {
	if fp != "" {
		if ext := path.Ext(fp); ext != ".json" {
			return nil, errors.New("unsupported file type, only json template files are currently supported")
		}
	}

	templateFile, err := input.NewFileOrReader(os.Stdin, fp)
	if err != nil {
		return nil, err
	}
	defer templateFile.Close()

	var templateConfiguration models.DeploymentTemplateInfo
	if err := json.NewDecoder(templateFile).Decode(&templateConfiguration); err != nil {
		return nil, err
	}

	return &templateConfiguration, nil
}

func init() {
	Command.AddCommand(
		platformDeploymentTemplateListCmd,
		platformDeploymentTemplateShowCmd,
		platformDeploymentTemplateDeleteCmd,
		platformDeploymentTemplateCreateCmd,
		platformDeploymentTemplateUpdateCmd,
		platformDeploymentTemplatePullCmd,
	)

	platformDeploymentTemplateShowCmd.Flags().BoolP(showInstanceConfigurations, "", false, "Shows instance configurations")

	platformDeploymentTemplateListCmd.Flags().BoolP(showInstanceConfigurations, "", false, "Shows instance configurations - only visible when using the JSON output")
	platformDeploymentTemplateListCmd.Flags().String(stackVersion, "", "If present, it will cause the returned deployment templates to be adapted to return only the elements allowed in that version.")
	platformDeploymentTemplateListCmd.Flags().String(filter, "", "Optional key/value pair in the form of key:value that will act as a filter and exclude any templates that do not have a matching metadata item associated")

	platformDeploymentTemplateCreateCmd.Flags().StringP("file-template", "f", "", "YAML or JSON file that contains the deployment template configuration")
	platformDeploymentTemplateCreateCmd.Flags().String("id", "", "Optional ID to set for the deployment template (Overrides ID if present)")

	platformDeploymentTemplateUpdateCmd.Flags().StringP("file-template", "f", "", "YAML or JSON file that contains the deployment template configuration")

	platformDeploymentTemplatePullCmd.Flags().StringP("path", "p", "", "Local path where to store deployment templates")
	platformDeploymentTemplatePullCmd.MarkFlagRequired("path")
}
