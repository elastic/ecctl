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

	"github.com/elastic/cloud-sdk-go/pkg/api/deploymentapi/depresourceapi"
	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var startMaintCmd = &cobra.Command{
	Use:     "start-maintenance <deployment id> --kind <kind> [--all|--i <instance-id>,<instance-id>]",
	Short:   "Starts maintenance mode on a deployment resource",
	PreRunE: sdkcmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resKind, _ := cmd.Flags().GetString("kind")
		refID, _ := cmd.Flags().GetString("ref-id")
		instanceID, _ := cmd.Flags().GetStringSlice("instance-id")
		ignoreMissing, _ := cmd.Flags().GetBool("ignore-missing")
		all, _ := cmd.Flags().GetBool("all")
		force, _ := cmd.Flags().GetBool("force")

		if err := sdkcmdutil.IncompatibleFlags(cmd, "all", "instance-id"); err != nil {
			fmt.Fprintln(cmd.OutOrStderr(), err)
		}

		var msg = "This action will incur in downtime if used with the --all flag. Do you want to continue? [y/n]: "
		if all && !force && !sdkcmdutil.ConfirmAction(msg, os.Stdin, ecctl.Get().Config.OutputDevice) {
			return nil
		}

		_, err := depresourceapi.StartMaintenanceModeAllOrSpecified(depresourceapi.StartInstancesParams{
			StartParams: depresourceapi.StartParams{
				Params: depresourceapi.Params{
					API:          ecctl.Get().API,
					DeploymentID: args[0],
					Kind:         resKind,
					RefID:        refID,
				},
				All: all,
			},
			InstanceIDs:   instanceID,
			IgnoreMissing: ec.Bool(ignoreMissing),
		})
		if err != nil {
			return err
		}

		return nil

	},
}

func init() {
	Command.AddCommand(startMaintCmd)
	startMaintCmd.Flags().Bool("all", false, "Starts maintenance mode on all instances of a defined resource kind")
	startMaintCmd.Flags().Bool("ignore-missing", false, "If set and the specified instance does not exist, then quietly proceed to the next instance")
	cmdutil.AddKindFlag(startMaintCmd, "Required", true)
	startMaintCmd.MarkFlagRequired("kind")
	startMaintCmd.Flags().String("ref-id", "", "Optional deployment RefId, if not set, the RefId will be auto-discovered")
	startMaintCmd.Flags().StringSliceP("instance-id", "i", nil, "Deployment instance IDs to use (e.g. instance-0000000001)")
}
