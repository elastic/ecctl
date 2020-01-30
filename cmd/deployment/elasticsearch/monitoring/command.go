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

package cmdelasticsearchmonitoring

import (
	"github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/deployment/elasticsearch/monitoring"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/util"
)

// Command represents the elasticsearch monitoring subcommand
var Command = &cobra.Command{
	Use:     "monitoring",
	Short:   "Manages monitoring for an Elasticsearch cluster",
	PreRunE: cobra.MaximumNArgs(0),

	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var enableMonitoringElasticsearchClusterCmd = &cobra.Command{
	Use:     "enable <monitored cluster id> <monitoring cluster id>",
	Short:   "Enables monitoring for the cluster by sending data to a monitoring cluster you specify",
	PreRunE: cmdutil.MinimumNArgsAndUUID(2),

	RunE: func(cmd *cobra.Command, args []string) error {
		return monitoring.Enable(monitoring.EnableParams{
			ClusterParams: util.ClusterParams{
				API:       ecctl.Get().API,
				ClusterID: args[0],
			},
			TargetID: args[1],
		})
	},
}

var disableMonitoringElasticsearchClusterCmd = &cobra.Command{
	Use:   "disable <monitored cluster id>",
	Short: "Disables monitoring in the specified cluster",

	PreRunE: cmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return monitoring.Disable(monitoring.DisableParams{
			ClusterParams: util.ClusterParams{
				API:       ecctl.Get().API,
				ClusterID: args[0],
			},
		})
	},
}

func init() {
	Command.AddCommand(
		enableMonitoringElasticsearchClusterCmd,
		disableMonitoringElasticsearchClusterCmd,
	)
}
