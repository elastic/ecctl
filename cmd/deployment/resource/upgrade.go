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
	"github.com/elastic/cloud-sdk-go/pkg/api/deploymentapi/depresourceapi"
	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
)

// upgradeCmd is the deployment subcommand
var upgradeCmd = &cobra.Command{
	Use:     "upgrade <deployment id> --kind <kind> --ref-id <ref-id>",
	Short:   "Upgrades a deployment resource",
	Long:    upgradeLong,
	PreRunE: sdkcmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resKind, _ := cmd.Flags().GetString("kind")
		refID, _ := cmd.Flags().GetString("ref-id")

		res, err := depresourceapi.UpgradeStateless(depresourceapi.Params{
			API:          ecctl.Get().API,
			DeploymentID: args[0],
			Kind:         resKind,
			RefID:        refID,
		})

		if err != nil {
			return err
		}

		track, _ := cmd.Flags().GetBool("track")
		return cmdutil.Track(cmdutil.NewTrackParams(cmdutil.TrackParamsConfig{
			App:          ecctl.Get(),
			DeploymentID: args[0],
			Track:        track,
			Response:     res,
		}))
	},
}

func init() {
	Command.AddCommand(upgradeCmd)
	upgradeCmd.Flags().BoolP("track", "t", false, cmdutil.TrackFlagMessage)
	cmdutil.AddKindFlag(upgradeCmd, "Required", true)
	upgradeCmd.MarkFlagRequired("kind")
	upgradeCmd.Flags().String("ref-id", "", "Optional deployment RefId, if not set, the RefId will be auto-discovered")
}
