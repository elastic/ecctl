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
	"github.com/elastic/cloud-sdk-go/pkg/models"
	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var updateCmd = &cobra.Command{
	Use:     `update -f <file definition.json>`,
	Short:   "Updates a deployment from a file definition, allowing certain flag overrides",
	Long:    updateLong,
	Example: updateExample,
	PreRunE: sdkcmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename, _ := cmd.Flags().GetString("file")
		var r models.DeploymentUpdateRequest
		if err := sdkcmdutil.DecodeFile(filename, &r); err != nil {
			return err
		}

		pruneOrphans, _ := cmd.Flags().GetBool("prune-orphans")
		r.PruneOrphans = &pruneOrphans

		force, _ := cmd.Flags().GetBool("force")
		var msg = `setting --prune-orphans to "true" will cause any resources not specified in the update request to be removed from the deployment, do you want to continue? [y/n]: `
		if pruneOrphans && !force && !sdkcmdutil.ConfirmAction(msg, os.Stdin, ecctl.Get().Config.OutputDevice) {
			return nil
		}

		skipSnapshot, _ := cmd.Flags().GetBool("skip-snapshot")
		hidePrunedOrphans, _ := cmd.Flags().GetBool("hide-pruned-orphans")

		var region = ecctl.Get().Config.Region
		if ecctl.Get().Config.Region == "" {
			region = cmdutil.DefaultECERegion
		}

		res, err := deploymentapi.Update(deploymentapi.UpdateParams{
			DeploymentID: args[0],
			API:          ecctl.Get().API,
			Overrides: deploymentapi.PayloadOverrides{
				Region: region,
			},
			Request:           &r,
			SkipSnapshot:      skipSnapshot,
			HidePrunedOrphans: hidePrunedOrphans,
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
		}))
	},
}

func init() {
	Command.AddCommand(updateCmd)
	updateCmd.Flags().BoolP("track", "t", false, cmdutil.TrackFlagMessage)
	updateCmd.Flags().Bool("prune-orphans", false, "When set to true, it will remove any resources not specified in the update request, treating the json file contents as the authoritative deployment definition")
	updateCmd.Flags().Bool("skip-snapshot", false, "Skips taking an Elasticsearch snapshot prior to shutting down the deployment")
	updateCmd.Flags().Bool("hide-pruned-orphans", false, "Hides orphaned resources that were shut down (only relevant if --prune-orphans=true)")
	updateCmd.Flags().StringP("file", "f", "", "Partial (default) or full JSON file deployment update payload")
	updateCmd.MarkFlagRequired("file")
	updateCmd.MarkFlagFilename("file", "json")
}
