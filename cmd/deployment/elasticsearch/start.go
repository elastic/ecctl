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

	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment/elasticsearch"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/util"
)

var startElasticsearchCmd = &cobra.Command{
	Use:     "start <cluster id>",
	Short:   "Starts a stopped Elasticsearch cluster",
	PreRunE: sdkcmdutil.MinimumNArgsAndUUID(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := elasticsearch.GetCluster(elasticsearch.GetClusterParams{
			ClusterParams: util.ClusterParams{
				API:       ecctl.Get().API,
				ClusterID: args[0],
			},
		})
		if err != nil {
			return err
		}

		if c.Status == "started" {
			return fmt.Errorf("cluster [%s][Elasticsearch]: is already started", args[0])
		}

		track, _ := cmd.Flags().GetBool("track")
		if err := elasticsearch.RestartCluster(elasticsearch.RestartClusterParams{
			ClusterParams: util.ClusterParams{
				ClusterID: args[0],
				API:       ecctl.Get().API,
			},
			TrackParams: util.TrackParams{
				Track:  track,
				Output: ecctl.Get().Config.OutputDevice,
			},
		}); err != nil {
			return err
		}

		fmt.Printf("Cluster [%s][Elasticsearch]: started\n", args[0])
		return nil
	},
}

func init() {
	Command.AddCommand(startElasticsearchCmd)
	startElasticsearchCmd.Flags().BoolP("track", "t", false, cmdutil.TrackFlagMessage)
}
