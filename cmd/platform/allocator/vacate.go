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

package cmdallocator

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/platform/allocator"
)

const vacateExamples = `  ecctl [globalFlags] allocator vacate i-05e245252362f7f1d
  # Move everything from multiple allocators
  ecctl [globalFlags] allocator vacate i-05e245252362f7f1d i-2362f7f1d252362f7

  # filter by a cluster kind
  ecctl [globalFlags] allocator vacate -k kibana i-05e245252362f7f1d

  # Only move specific cluster IDs
  ecctl [globalFlags] allocator vacate -c f521dedb07194c478fbbc6624f3bbf8f -c f404eea372cc4ea5bd553d47a09705cd i-05e245252362f7f1d

  # Specify multiple allocator targets
  ecctl [globalFlags] allocator vacate -t i-05e245252362f7f2d -t i-2362f7f1d252362f7 i-05e245252362f7f1d
  ecctl [globalFlags] allocator vacate --target i-05e245252362f7f2d --target i-2362f7f1d252362f7 --kind kibana i-05e245252362f7f1d

  # Set the allocators to maintenance mode before vacating them
  ecctl [globalFlags] allocator vacate --maintenance -t i-05e245252362f7f2d -t i-2362f7f1d252362f7 i-05e245252362f7f1d

  # Set the amount of maximum moves to happen at any time
  ecctl [globalFlags] allocator vacate --concurrency 10 i-05e245252362f7f1d

  # Override the allocator health auto discovery
  ecctl [globalFlags] allocator vacate --allocator-down=true i-05e245252362f7f1d

  # Override the skip_snapshot setting
  ecctl [globalFlags] allocator vacate --skip-snapshot=true i-05e245252362f7f1d -c f521dedb07194c478fbbc6624f3bbf8f

  # Override the skip_data_migration setting
  ecctl [globalFlags] allocator vacate --skip-data-migration=true i-05e245252362f7f1d -c f521dedb07194c478fbbc6624f3bbf8f
  `

var vacateAllocatorCmd = &cobra.Command{
	Use:     "vacate <source>",
	Short:   "Moves all the clusters from the specified allocator",
	Example: vacateExamples,
	PreRunE: cobra.MinimumNArgs(1),
	Aliases: []string{"move-nodes"},

	RunE: func(cmd *cobra.Command, args []string) error {
		concurrency, err := strconv.ParseUint(cmd.Flag("concurrency").Value.String(), 10, 64)
		if err != nil {
			return err
		}

		var allocatorDownRaw = cmd.Flag("allocator-down").Value.String()
		if len(args) > 1 && allocatorDownRaw != "" {
			return errors.New("cannot specify multiple allocators with --allocator-down")
		}

		allocatorDown, err := strconv.ParseBool(allocatorDownRaw)
		if err != nil && allocatorDownRaw != "" {
			return err
		}

		clusters, err := cmd.Flags().GetStringArray("cluster")
		if err != nil {
			return err
		}

		moveOnly, _ := cmd.Flags().GetBool("move-only")

		overrideFailsafeRaw, _ := cmd.Flags().GetBool("override-failsafe")
		overrideFailsafe, err := cmdutil.ActionConfirm(strconv.FormatBool(overrideFailsafeRaw), "--override-failsafe flag specified. Are you sure you want to proceed? [y/N]: ")
		if err != nil {
			return err
		}

		target, _ := cmd.Flags().GetStringSlice("target")
		kind, _ := cmd.Flags().GetString("kind")

		skipSnapshotRaw, err := cmd.Flags().GetString("skip-snapshot")
		if err != nil {
			return err
		}

		skipSnapshot, err := cmdutil.ActionConfirm(skipSnapshotRaw, "--skip-snapshot flag specified. Are you sure you want to proceed? [y/N]: ")
		if err != nil {
			return err
		}

		skipDataMigrationRaw, err := cmd.Flags().GetString("skip-data-migration")
		if err != nil {
			return err
		}

		err = validateSkipDataMigration(clusters, moveOnly)
		if err != nil && skipDataMigrationRaw != "" {
			return err
		}

		skipDataMigration, err := cmdutil.ActionConfirm(skipDataMigrationRaw, "--skip-data-migration flag specified. Are you sure you want to proceed? [y/N]: ")
		if err != nil {
			return err
		}

		setAllocatorMaintenance, _ := cmd.Flags().GetBool("maintenance")

		var merr error
		// Only sets the allocator to maintenance mode when the flag is specified
		if setAllocatorMaintenance {
			for _, id := range args {
				var params = allocator.MaintenanceParams{
					API: ecctl.Get().API,
					ID:  id,
				}
				if err := allocator.StartMaintenance(params); err != nil {
					merr = multierror.Append(merr, err)
				}
			}
		}

		if merr != nil {
			fmt.Fprint(ecctl.Get().Config.OutputDevice, merr)
		}

		var params = &allocator.VacateParams{
			API:                 ecctl.Get().API,
			Allocators:          args,
			PreferredAllocators: target,
			ClusterFilter:       clusters,
			KindFilter:          kind,
			Concurrency:         uint16(concurrency),
			Output:              ecctl.Get().Config.OutputDevice,
			MoveOnly:            ec.Bool(moveOnly),
			PlanOverrides: allocator.PlanOverrides{
				SkipSnapshot:      skipSnapshot,
				SkipDataMigration: skipDataMigration,
				OverrideFailsafe:  overrideFailsafe,
			},
		}
		if len(args) == 1 && allocatorDownRaw != "" {
			params.AllocatorDown = &allocatorDown
		}

		return allocator.Vacate(params)
	},
}

func validateSkipDataMigration(clusters []string, moveOnly bool) error {
	if len(clusters) < 1 || !moveOnly {
		return errors.New("skip data migration is not available if there are no cluster IDs specified or move-only is set to false")
	}

	return nil
}

func init() {
	Command.AddCommand(vacateAllocatorCmd)

	vacateAllocatorCmd.Flags().StringP("kind", "k", "", "Kind of workload to vacate (elasticsearch|kibana)")
	vacateAllocatorCmd.Flags().StringArrayP("cluster", "c", nil, "Cluster IDs to include in the vacate")
	vacateAllocatorCmd.Flags().StringArrayP("target", "t", nil, "Target allocator(s) on which to place the vacated workload")
	vacateAllocatorCmd.Flags().BoolP("maintenance", "m", false, "Whether to set the allocator(s) in maintenance before performing the vacate")
	vacateAllocatorCmd.Flags().Uint("concurrency", 8, "Maximum number of concurrent moves to perform at any time")
	vacateAllocatorCmd.Flags().String("allocator-down", "", "Disables the allocator health auto-discovery, setting the allocator-down to either [true|false]")
	vacateAllocatorCmd.Flags().Bool("move-only", true, "Keeps the cluster in its current -possibly broken- state and just does the bare minimum to move the requested instances across to another allocator. [true|false]")
	vacateAllocatorCmd.Flags().Bool("override-failsafe", false, "If false (the default) then the plan will fail out if it believes the requested sequence of operations can result in data loss - this flag will override some of these restraints. [true|false]")
	vacateAllocatorCmd.Flags().String("skip-snapshot", "", "Skips the snapshot operation on the specified cluster IDs. ONLY available when the cluster IDs are specified. [true|false]")
	vacateAllocatorCmd.Flags().String("skip-data-migration", "", "Skips the data-migration operation on the specified cluster IDs. ONLY available when the cluster IDs are specified and --move-only is true. [true|false]")
}
