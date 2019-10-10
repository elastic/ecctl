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
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment/elasticsearch"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/util"
)

var deleteElasticsearchCmd = &cobra.Command{
	Use:     "delete <cluster id>",
	Short:   "Deletes an Elasticsearch cluster",
	PreRunE: cmdutil.MinimumNArgsAndUUID(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		if stop, _ := cmd.Flags().GetBool("stop"); stop {
			skipSnapshot, _ := cmd.Flags().GetBool("skip-snapshot")
			err := elasticsearch.ShutdownCluster(elasticsearch.ShutdownClusterParams{
				ClusterParams: util.ClusterParams{
					ClusterID: args[0],
					API:       ecctl.Get().API,
				},
				TrackParams: util.TrackParams{
					Track:  true,
					Output: ecctl.Get().Config.OutputDevice,
				},
				SkipSnapshot: skipSnapshot,
			})
			if err != nil {
				return err
			}
		}

		return elasticsearch.DeleteCluster(elasticsearch.DeleteClusterParams{
			ClusterParams: util.ClusterParams{
				ClusterID: args[0],
				API:       ecctl.Get().API,
			},
		})
	},
}

func init() {
	Command.AddCommand(deleteElasticsearchCmd)
	deleteElasticsearchCmd.Flags().Bool("stop", false, "Stops the deployment before deleting it")
	deleteElasticsearchCmd.Flags().Bool("skip-snapshot", false, "Skips snapshotting before shutting down the cluster (affects running clusters only)")
}
