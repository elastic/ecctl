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

	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment/elasticsearch/instances"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/util"
)

var stopElasticsearchInstancesRoutingCmd = &cobra.Command{
	Use:     "stop-routing <cluster id> [--all|--instances]",
	Aliases: []string{"pause-routing"},
	Short:   "Stops routing on specific Elasticsearch instances, use the --all flag to target all instances",
	PreRunE: cmdutil.MinimumNArgsAndUUID(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		var clusterParams = util.ClusterParams{
			ClusterID: args[0],
			API:       ecctl.Get().API,
		}

		if err := cmdutil.IncompatibleFlags(cmd, "all", "instance"); err != nil {
			fmt.Fprintln(cmd.OutOrStderr(), err)
		}

		instanceSlice, err := cmdutil.GetInstances(cmd, clusterParams, "instance")
		if err != nil {
			return err
		}

		if err := instances.StopRouting(instances.Params{
			ClusterParams: clusterParams,
			Instances:     instanceSlice,
		}); err != nil {
			return err
		}

		fmt.Printf("Cluster [%s][Elasticsearch]: routing stopped for %s\n", args[0], instanceSlice)
		return nil
	},
}

var startElasticsearchInstancesRoutingCmd = &cobra.Command{
	Use:     "start-routing <cluster id> [--all|--instances]",
	Aliases: []string{"resume-routing"},

	Short:   "Resumes routing on specific Elasticsearch instances, use the --all flag to target all instances",
	PreRunE: cmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var clusterParams = util.ClusterParams{
			ClusterID: args[0],
			API:       ecctl.Get().API,
		}

		if err := cmdutil.IncompatibleFlags(cmd, "all", "instance"); err != nil {
			fmt.Fprintln(cmd.OutOrStderr(), err)
		}

		instanceSlice, err := cmdutil.GetInstances(cmd, clusterParams, "instance")
		if err != nil {
			return err
		}

		if err := instances.StartRouting(instances.Params{
			ClusterParams: clusterParams,
			Instances:     instanceSlice,
		}); err != nil {
			return err
		}

		fmt.Printf("Cluster [%s][Elasticsearch]: routing started for %s\n", args[0], instanceSlice)
		return nil
	},
}

func init() {
	Command.AddCommand(stopElasticsearchInstancesRoutingCmd)
	Command.AddCommand(startElasticsearchInstancesRoutingCmd)

	stopElasticsearchInstancesRoutingCmd.Flags().StringSliceP("instance", "i", nil, "Instances to stop routing")
	stopElasticsearchInstancesRoutingCmd.Flags().Bool("all", false, "Stops routing to all instances")
	startElasticsearchInstancesRoutingCmd.Flags().StringSliceP("instance", "i", nil, "Instances to start routing")
	startElasticsearchInstancesRoutingCmd.Flags().Bool("all", false, "Starts routing to all instances")
}
