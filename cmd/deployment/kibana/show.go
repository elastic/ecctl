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

package cmdkibana

import (
	"path/filepath"

	"github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/deployment/deputil"
	"github.com/elastic/ecctl/pkg/deployment/kibana"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var showKibanaClusterCmd = showKibanaCluster()

func showKibanaCluster() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show <cluster id>",
		Short:   "Returns the cluster information",
		PreRunE: cmdutil.MinimumNArgsAndUUID(1),
		RunE:    runShowKibanaClusterCmd,
	}

	return cmd
}

func runShowKibanaClusterCmd(cmd *cobra.Command, args []string) error {
	logs, _ := cmd.Flags().GetBool("logs")
	plans, _ := cmd.Flags().GetBool("plans")
	metadata, _ := cmd.Flags().GetBool("metadata")
	settings, _ := cmd.Flags().GetBool("settings")
	defaults, _ := cmd.Flags().GetBool("defaults")
	if logs {
		plans = true
	}

	r, err := kibana.Get(kibana.ClusterParams{
		DeploymentParams: kibana.DeploymentParams{
			API: ecctl.Get().API,
			ID:  args[0]},
		QueryParams: deputil.QueryParams{
			ShowMetadata:     metadata,
			ShowPlanDefaults: defaults,
			ShowPlans:        plans,
			ShowPlanLogs:     logs,
			ShowSettings:     settings,
		}},
	)
	if err != nil {
		return err
	}

	return ecctl.Get().Formatter.Format(filepath.Join("kibana", "show"), r)
}

func init() {
	showKibanaClusterCmd.Flags().BoolP("metadata", "m", false, "View the cluster metadata")
	showKibanaClusterCmd.Flags().BoolP("defaults", "d", true, "View the cluster plan defaults")
	showKibanaClusterCmd.Flags().BoolP("plans", "p", false, "View the cluster plans")
	showKibanaClusterCmd.Flags().BoolP("logs", "l", false, "View the cluster plan logs")
	showKibanaClusterCmd.Flags().BoolP("settings", "s", false, "View the cluster settings")
}
