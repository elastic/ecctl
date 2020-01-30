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

package cmdelasticsearch

import (
	"fmt"
	"time"

	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment/elasticsearch"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/util"
)

var restartElasticsearchCmd = &cobra.Command{
	Use:     "restart <cluster id>",
	Short:   "Restarts an Elasticsearch cluster",
	PreRunE: sdkcmdutil.MinimumNArgsAndUUID(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		track, _ := cmd.Flags().GetBool("track")
		rollingByName, _ := cmd.Flags().GetBool("rolling-by-name")
		rollingByZone, _ := cmd.Flags().GetBool("rolling-by-zone")
		skipSnapshot, _ := cmd.Flags().GetBool("skip-snapshot")
		shardInitWaitTime, _ := cmd.Flags().GetDuration("shard-init-wait-time")
		err := elasticsearch.RestartCluster(elasticsearch.RestartClusterParams{
			ClusterParams: util.ClusterParams{
				ClusterID: args[0],
				API:       ecctl.Get().API,
			},
			TrackParams: util.TrackParams{
				Track:  track,
				Output: ecctl.Get().Config.OutputDevice,
			},
			RollingByName:     rollingByName,
			RollingByZone:     rollingByZone,
			SkipSnapshot:      skipSnapshot,
			ShardInitWaitTime: shardInitWaitTime,
		})
		if err != nil {
			return err
		}

		message := "restarted"
		if !track {
			message = "restart triggered"
		}

		fmt.Printf("Cluster [%s][Elasticsearch]: %s\n", args[0], message)
		return nil
	},
}

func init() {
	Command.AddCommand(restartElasticsearchCmd)
	restartElasticsearchCmd.Flags().BoolP("track", "t", false, cmdutil.TrackFlagMessage)
	restartElasticsearchCmd.Flags().Bool("rolling-by-name", false, "Performs the restart in a rolling fashion (one instance at a time)")
	restartElasticsearchCmd.Flags().Bool("rolling-by-zone", false, "Performs the restart in a rolling fashion (one logical zone at a time)")
	restartElasticsearchCmd.Flags().Bool("skip-snapshot", false, "Prevents snapshotting prior to restart")
	restartElasticsearchCmd.Flags().Duration("shard-init-wait-time", 10*time.Minute, "Time to wait for shards that show no progress of initializing, before rolling the next group")
}
