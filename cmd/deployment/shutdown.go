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

	"github.com/elastic/cloud-sdk-go/pkg/api/deploymentapi"
	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var shutdownCmd = &cobra.Command{
	Use:     "shutdown <deployment-id>",
	Short:   "Shuts down a deployment and all of its associated sub-resources",
	PreRunE: sdkcmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")
		var msg = "This action will delete the specified deployment ID and its associated sub-resources, do you want to continue? [y/n]: "
		if !force && !sdkcmdutil.ConfirmAction(msg, os.Stdin, ecctl.Get().Config.OutputDevice) {
			return nil
		}

		skipSnapshot, _ := cmd.Flags().GetBool("skip-snapshot")

		res, err := deploymentapi.Shutdown(deploymentapi.ShutdownParams{
			API:          ecctl.Get().API,
			DeploymentID: args[0],
			SkipSnapshot: skipSnapshot,
		})
		if err != nil {
			return err
		}

		var track, _ = cmd.Flags().GetBool("track")
		return cmdutil.Track(cmdutil.NewTrackParams(cmdutil.TrackParamsConfig{
			App:          ecctl.Get(),
			DeploymentID: args[0],
			Track:        track,
			Response:     res,
			Template:     "deployment/shutdown",
		}))
	},
}

func init() {
	Command.AddCommand(shutdownCmd)
	shutdownCmd.Flags().BoolP("track", "t", false, cmdutil.TrackFlagMessage)
	shutdownCmd.Flags().Bool("skip-snapshot", false, "Skips taking an Elasticsearch snapshot prior to shutting down the deployment")
}
