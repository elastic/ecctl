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
	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment/elasticsearch"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/util"
)

var shutdownElasticsearchCmd = &cobra.Command{
	Use:     "shutdown <cluster id>",
	Aliases: []string{"stop"},
	Short:   "Shuts down an Elasticsearch cluster, so that it no longer contains any running instances",
	PreRunE: sdkcmdutil.MinimumNArgsAndUUID(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		track, _ := cmd.Flags().GetBool("track")
		hide, _ := cmd.Flags().GetBool("hide")
		skipSnapshot, _ := cmd.Flags().GetBool("skip-snapshot")

		return elasticsearch.ShutdownCluster(elasticsearch.ShutdownClusterParams{
			ClusterParams: util.ClusterParams{
				ClusterID: args[0],
				API:       ecctl.Get().API,
			},
			Track: track,
			TrackChangeParams: cmdutil.NewTrackParams(cmdutil.TrackParamsConfig{
				App:        ecctl.Get(),
				ResourceID: args[0],
				Kind:       util.Elasticsearch,
				Track:      track,
			}).TrackChangeParams,
			Hide:         hide,
			SkipSnapshot: skipSnapshot,
		})
	},
}

func init() {
	Command.AddCommand(shutdownElasticsearchCmd)
	shutdownElasticsearchCmd.Flags().BoolP("track", "t", false, cmdutil.TrackFlagMessage)
	shutdownElasticsearchCmd.Flags().Bool("hide", false, "Hides the cluster from the listing")
	shutdownElasticsearchCmd.Flags().Bool("skip-snapshot", false, "Skips snapshotting before shutting down the cluster")
}
