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

	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment"
	"github.com/elastic/ecctl/pkg/deployment/depresource"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var stopCmd = &cobra.Command{
	Use:     "stop <deployment id> --type <type> [--all|--i <instance-id>,<instance-id>]",
	Short:   "Stops a deployment resource",
	PreRunE: cmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resType, _ := cmd.Flags().GetString("type")
		refID, _ := cmd.Flags().GetString("ref-id")
		instanceID, _ := cmd.Flags().GetStringSlice("instance-id")
		ignoreMissing, _ := cmd.Flags().GetBool("ignore-missing")

		if refID == "" {
			var err error
			refID, err = getRefID(resType, args[0])
			if err != nil {
				return err
			}
		}

		if err := cmdutil.IncompatibleFlags(cmd, "all", "instance-id"); err != nil {
			fmt.Fprintln(cmd.OutOrStderr(), err)
		}

		if all, _ := cmd.Flags().GetBool("all"); all {
			_, err := depresource.Stop(depresource.StopParams{
				API:          ecctl.Get().API,
				DeploymentID: args[0],
				Type:         resType,
				RefID:        refID,
			})
			if err != nil {
				return err
			}

			return nil
		}

		_, err := depresource.StopInstances(depresource.StopInstancesParams{
			StopParams: depresource.StopParams{
				API:          ecctl.Get().API,
				DeploymentID: args[0],
				Type:         resType,
				RefID:        refID,
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

func getRefID(resType, depID string) (string, error) {
	getResourceParams := deployment.GetResourceParams{
		GetParams: deployment.GetParams{
			API:          ecctl.Get().API,
			DeploymentID: depID,
		},
		Type: resType,
	}

	refID, err := deployment.GetTypeRefID(getResourceParams)
	if err != nil {
		return "", err
	}

	return refID, nil
}

func init() {
	Command.AddCommand(stopCmd)
	stopCmd.Flags().Bool("all", false, "Stops all instances of a defined resource type")
	stopCmd.Flags().Bool("ignore-missing", false, "If set and the specified instance does not exist, then quietly proceed to the next instance")
	stopCmd.Flags().String("type", "", "Deployment resource type to stop (elasticsearch, kibana, apm, or appsearch)")
	stopCmd.MarkFlagRequired("type")
	stopCmd.Flags().String("ref-id", "", "Optional deployment RefId, if not set, the RefId will be auto-discovered")
	stopCmd.Flags().StringSliceP("instance-id", "i", nil, "Deployment instance IDs to stop (e.g. instance-0000000001)")
}
