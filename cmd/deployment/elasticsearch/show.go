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
	"path/filepath"

	"github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/deployment/elasticsearch"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/util"
)

var showElasticsearchmd = &cobra.Command{
	Use:     "show <cluster id>",
	Short:   "Displays information about the specified cluster",
	PreRunE: cmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var logs, _ = cmd.Flags().GetBool("logs")
		var plans, _ = cmd.Flags().GetBool("plans")
		var metadata, _ = cmd.Flags().GetBool("metadata")
		var defaults, _ = cmd.Flags().GetBool("defaults")
		var settings, _ = cmd.Flags().GetBool("settings")

		if logs {
			plans = true
		}
		c, err := elasticsearch.GetCluster(elasticsearch.GetClusterParams{
			ClusterParams: util.ClusterParams{
				API:       ecctl.Get().API,
				ClusterID: args[0],
			},
			Metadata:     metadata,
			PlanDefaults: defaults,
			Settings:     settings,
			Plans:        plans,
			Logs:         logs,
		})
		if err != nil {
			return err
		}
		return ecctl.Get().Formatter.Format(
			filepath.Join("elasticsearch", "show"),
			c,
		)
	},
}

func init() {
	Command.AddCommand(showElasticsearchmd)
	showElasticsearchmd.Flags().BoolP("metadata", "m", false, "View the cluster metadata")
	showElasticsearchmd.Flags().BoolP("defaults", "d", true, "View the cluster plan defaults")
	showElasticsearchmd.Flags().BoolP("plans", "p", false, "View the cluster plans")
	showElasticsearchmd.Flags().BoolP("logs", "l", false, "View the cluster plan logs")
	showElasticsearchmd.Flags().BoolP("settings", "s", false, "View the cluster settings")
}
