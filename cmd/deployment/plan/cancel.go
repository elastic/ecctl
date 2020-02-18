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

package cmddeploymentplan

import (
	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment"
	"github.com/elastic/ecctl/pkg/deployment/depresource"
	"github.com/elastic/ecctl/pkg/ecctl"
)

// cancelPlan is the deployment subcommand
var cancelPlan = &cobra.Command{
	Use:     "cancel <deployment id> --type <type> [--ref-id <ref-id>]",
	Short:   "Cancels a resource's pending plan",
	PreRunE: sdkcmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")
		resType, _ := cmd.Flags().GetString("type")
		refID, _ := cmd.Flags().GetString("ref-id")

		_, err := depresource.CancelPlan(depresource.CancelPlanParams{
			ResourceParams: deployment.ResourceParams{
				API:          ecctl.Get().API,
				DeploymentID: args[0],
				Type:         resType,
				RefID:        refID,
			},
			ForceDelete: force,
		})

		return err
	},
}

func init() {
	Command.AddCommand(cancelPlan)
	cmdutil.AddTypeFlag(cancelPlan, "Required", true)
	cancelPlan.MarkFlagRequired("type")
	cancelPlan.Flags().String("ref-id", "", "Optional deployment RefId, if not set, the RefId will be auto-discovered")
}
