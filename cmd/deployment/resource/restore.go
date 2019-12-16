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
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment/depresource"
	"github.com/elastic/ecctl/pkg/ecctl"
)

// restoreCmd is the deployment subcommand
var restoreCmd = &cobra.Command{
	Use:     "restore <deployment id> --type <type> --ref-id <ref-id>",
	Short:   "Restores a previously shut down deployment resource",
	PreRunE: cmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resType, _ := cmd.Flags().GetString("type")
		refID, _ := cmd.Flags().GetString("ref-id")
		restoreSnapshot, _ := cmd.Flags().GetBool("restore-snapshot")

		force, _ := cmd.Flags().GetBool("force")
		var esType = resType == "elasticsearch"
		var dataLoss = esType && !restoreSnapshot
		var msg = "This action restores an Elasticsearch resource without its snapshot which might incur data loss, do you want to continue? [y/n]: "
		if dataLoss && !force && !cmdutil.ConfirmAction(msg, os.Stderr, os.Stdout) {
			return nil
		}

		if restoreSnapshot && !esType {
			fmt.Fprintf(ecctl.Get().Config.ErrorDevice,
				"Using --restore-snapshot on resource type \"%s\" will have no effect\n", resType,
			)
		}

		return depresource.Restore(depresource.RestoreParams{
			API:             ecctl.Get().API,
			DeploymentID:    args[0],
			Type:            resType,
			RefID:           refID,
			RestoreSnapshot: restoreSnapshot,
		})
	},
}

func init() {
	Command.AddCommand(restoreCmd)
	restoreCmd.Flags().String("type", "", "Required deployment type to restore (elasticsearch, kibana, apm, or appsearch)")
	restoreCmd.MarkFlagRequired("type")
	restoreCmd.Flags().String("ref-id", "", "Required deployment RefId")
	restoreCmd.MarkFlagRequired("ref-id")
	restoreCmd.Flags().Bool("restore-snapshot", false, "Optional flag to toggle restoring a snapshot for an Elasticsearch resource. It has no effect for other resources")
}
