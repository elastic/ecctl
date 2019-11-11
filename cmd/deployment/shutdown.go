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

package cmddeployment

import (
	"os"

	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var shutdownCmd = &cobra.Command{
	Use:     "shutdown <deployment-id>",
	Short:   "Shutdown's a platform deployment",
	PreRunE: cmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")
		var msg = "This action will delete the specified deployment ID and its associated sub-resources, do you want to continue? [y/n]: "
		if !force && !cmdutil.ConfirmAction(msg, os.Stderr, os.Stdout) {
			return nil
		}

		skipSnapshot, _ := cmd.Flags().GetBool("skip-snapshot")
		hide, _ := cmd.Flags().GetBool("hide")

		res, err := deployment.Shutdown(deployment.ShutdownParams{
			API:          ecctl.Get().API,
			DeploymentID: args[0],
			SkipSnapshot: skipSnapshot,
			Hide:         hide,
		})
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format("deployment/shutdown", res)
	},
}

func init() {
	Command.AddCommand(shutdownCmd)
	shutdownCmd.Flags().Bool("skip-snapshot", false, "Skips taking an Elasticsearch snapshot prior to shutting down the deployment")
	shutdownCmd.Flags().Bool("hide", false, "Hides the deployment and its resources after it has been shut down")
}
