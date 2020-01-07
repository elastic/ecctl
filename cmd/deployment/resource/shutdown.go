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

package cmddeploymentresource

import (
	"os"

	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment"
	"github.com/elastic/ecctl/pkg/deployment/depresource"
	"github.com/elastic/ecctl/pkg/ecctl"
)

// shutdownCmd is the deployment subcommand
var shutdownCmd = &cobra.Command{
	Use:     "shutdown <deployment id> --type <type> --ref-id <ref-id>",
	Short:   "Shuts down a deployment resource by its type and ref-id",
	Long:    shutdownLong,
	PreRunE: cmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resType, _ := cmd.Flags().GetString("type")
		refID, _ := cmd.Flags().GetString("ref-id")
		skipSnapshot, _ := cmd.Flags().GetBool("skip-snapshot")
		hide, _ := cmd.Flags().GetBool("hide")

		force, _ := cmd.Flags().GetBool("force")
		var msg = "This action will shut down a deployment's resource type. Do you want to continue? [y/n]: "
		if !force && !cmdutil.ConfirmAction(msg, os.Stderr, os.Stdout) {
			return nil
		}

		return depresource.Shutdown(depresource.ShutdownParams{
			ResourceParams: deployment.ResourceParams{
				API:          ecctl.Get().API,
				DeploymentID: args[0],
				Type:         resType,
				RefID:        refID,
			},
			SkipSnapshot: skipSnapshot,
			Hide:         hide,
		})

	},
}

func init() {
	Command.AddCommand(shutdownCmd)
	cmdutil.AddTypeFlag(shutdownCmd, "Required", true)
	shutdownCmd.MarkFlagRequired("type")
	shutdownCmd.Flags().String("ref-id", "", "Optional deployment RefId, auto-discovered if not specified")
	shutdownCmd.Flags().Bool("skip-snapshot", false, "Optional flag to toggle skipping the resource snapshot before shutting it down")
	shutdownCmd.Flags().Bool("hide", false, "Optionally hides the deployment resource from being listed by default")
}
