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
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment/depresource"
	"github.com/elastic/ecctl/pkg/ecctl"
)

// upgradeCmd is the deployment subcommand
var upgradeCmd = &cobra.Command{
	Use:     "upgrade <deployment id> --type <type> --ref-id <ref-id>",
	Short:   "Upgrades a deployment resource",
	Long:    upgradeLong,
	PreRunE: cmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resType, _ := cmd.Flags().GetString("type")
		refID, _ := cmd.Flags().GetString("ref-id")

		res, err := depresource.UpgradeStateless(depresource.UpgradeStatelessParams{
			API:          ecctl.Get().API,
			DeploymentID: args[0],
			Type:         resType,
			RefID:        refID,
		})

		if err != nil {
			return err
		}

		if track, _ := cmd.Flags().GetBool("track"); !track {
			return nil
		}

		return depresource.TrackResources(depresource.TrackResourcesParams{
			API:          ecctl.Get().API,
			OutputDevice: ecctl.Get().Config.OutputDevice,
			Resources: []*models.DeploymentResource{{
				ID:    ec.String(res.ResourceID),
				Kind:  ec.String(resType),
				RefID: ec.String(refID),
			}},
		})
	},
}

func init() {
	Command.AddCommand(upgradeCmd)
	upgradeCmd.Flags().BoolP("track", "t", false, cmdutil.TrackFlagMessage)
	upgradeCmd.Flags().String("type", "", "Optional stateless deployment type to upgrade (kibana, apm, or appsearch)")
	upgradeCmd.MarkFlagRequired("type")
	upgradeCmd.Flags().String("ref-id", "", "Optional deployment RefId, if not set, the RefId will be auto-discovered")
	upgradeCmd.MarkFlagRequired("ref-id")
}
