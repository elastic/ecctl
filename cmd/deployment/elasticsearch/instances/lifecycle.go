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

	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment/elasticsearch/instances"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/util"
)

// Command represents the instances subcommand
var Command = &cobra.Command{
	Use:     "instances",
	Short:   "Manages elasticsearch at the instance level",
	PreRunE: cobra.MaximumNArgs(0),

	Run: func(cmd *cobra.Command, args []string) { cmd.Help() },
}

var pauseCmd = &cobra.Command{
	Use:     "pause <cluster id> [--all|--instances]",
	Short:   "Pauses (stops) specific Elasticsearch instances, use the --all flag to target all instances",
	PreRunE: sdkcmdutil.MinimumNArgsAndUUID(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		var clusterParams = util.ClusterParams{
			API:       ecctl.Get().API,
			ClusterID: args[0],
		}

		if err := sdkcmdutil.IncompatibleFlags(cmd, "all", "instance"); err != nil {
			fmt.Fprintln(cmd.OutOrStderr(), err)
		}

		instanceSlice, err := cmdutil.GetInstances(cmd, clusterParams, "instance")
		if err != nil {
			return err
		}

		if err := instances.Pause(instances.Params{
			Instances:     instanceSlice,
			ClusterParams: clusterParams,
		}); err != nil {
			return err
		}

		fmt.Printf("Cluster [%s][Elasticsearch]: paused %s\n", args[0], instanceSlice)
		return nil
	},
}

var resumeCmd = &cobra.Command{
	Use:     "resume <cluster id> [--all|--instances]",
	Short:   "Resumes (starts) specific Elasticsearch instances, use the --all flag to target all instances",
	PreRunE: sdkcmdutil.MinimumNArgsAndUUID(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		var clusterParams = util.ClusterParams{
			API:       ecctl.Get().API,
			ClusterID: args[0],
		}

		if err := sdkcmdutil.IncompatibleFlags(cmd, "all", "instance"); err != nil {
			fmt.Fprintln(cmd.OutOrStderr(), err)
		}

		instanceSlice, err := cmdutil.GetInstances(cmd, clusterParams, "instance")
		if err != nil {
			return err
		}

		if err := instances.Resume(instances.Params{
			Instances:     instanceSlice,
			ClusterParams: clusterParams,
		}); err != nil {
			return err
		}

		fmt.Printf("Cluster [%s][Elasticsearch]: resumed %s\n", args[0], instanceSlice)
		return nil
	},
}

func init() {
	Command.AddCommand(
		pauseCmd,
		resumeCmd,
	)
	pauseCmd.Flags().StringSliceP("instance", "i", []string{}, "Instances to pause")
	pauseCmd.Flags().Bool("all", false, "Targets all the cluster instances")
	resumeCmd.Flags().StringSliceP("instance", "i", []string{}, "Instances to resume")
	resumeCmd.Flags().Bool("all", false, "Targets all the cluster instances")
}
