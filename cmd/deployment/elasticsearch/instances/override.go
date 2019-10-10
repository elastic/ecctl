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

package cmdelasticsearchinstances

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment/elasticsearch/instances"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/util"
)

var setESInstanceOverrideCmd = &cobra.Command{
	Use:   "override-capacity <cluster id>",
	Short: "Sets overrides at the instance level, use the --all flag to target all instances",
	Long: `Sets overrides at the instance level, use the --all flag to target all instances.
If the multiplier flag is not set, the override will be set using the current configured capacity * 2.
Only works for overrides <= 65535`,
	Example: `
Set all the instances in the cluster to the original (plan) capacity:
    ecctl deployment elasticsearch instances override-capacity 18fc96c491b3d5e10e147463927a5349 --reset

Set all the instances in the cluster to 4096 capacity:
    ecctl deployment elasticsearch instances override-capacity 18fc96c491b3d5e10e147463927a5349 --all --value 4096

Set all the instances in the cluster to 2x multiplier:
    ecctl deployment elasticsearch instances override-capacity 18fc96c491b3d5e10e147463927a5349 --all --multiplier 2

Set the cluster instance to 3x of its current capacity:
    ecctl deployment elasticsearch instances override-capacity 18fc96c491b3d5e10e147463927a5349 --instance instance-0000000003 --multiplier 3`,
	PreRunE: cmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		clusterParams := util.ClusterParams{
			API:       ecctl.Get().API,
			ClusterID: args[0],
		}

		if err := cmdutil.IncompatibleFlags(cmd, "all", "instance"); err != nil {
			fmt.Fprintln(cmd.OutOrStderr(), err)
		}

		instanceSlice, err := cmdutil.GetInstances(cmd, clusterParams, "instance")
		if err != nil {
			return err
		}

		multiplier, _ := cmd.Flags().GetUint8("multiplier")
		value, _ := cmd.Flags().GetUint16("value")
		reset, _ := cmd.Flags().GetBool("reset")

		res, err := instances.OverrideCapacity(instances.OverrideCapacityParams{
			ClusterParams: clusterParams,
			Instances:     instanceSlice,
			BoostFactor:   multiplier,
			Value:         uint64(value),
			Reset:         reset,
		})
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format(
			filepath.Join("elasticsearch", "overridecapacity"),
			res,
		)
	},
}

func init() {
	Command.AddCommand(setESInstanceOverrideCmd)
	setESInstanceOverrideCmd.Flags().StringSliceP("instance", "i", nil, "Instances on which to apply the override")
	setESInstanceOverrideCmd.Flags().Bool("all", false, "Applies the override to all of the instances")
	setESInstanceOverrideCmd.Flags().Uint8("multiplier", 0, "Capacity multiplier")
	setESInstanceOverrideCmd.Flags().Uint16P("value", "", 0, "Absolute value of instance override memory (in MBs)")
	setESInstanceOverrideCmd.Flags().BoolP("reset", "", false, "Resets the instance(s) memory to the original value found in the current plan")
}
