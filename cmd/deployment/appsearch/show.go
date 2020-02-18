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

package cmdappsearch

import (
	"path/filepath"

	"github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/deployment"
	"github.com/elastic/ecctl/pkg/deployment/deputil"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var showAppSearchCmd = &cobra.Command{
	Use:     "show <deployment id>",
	Short:   "Shows the specified AppSearch deployment",
	PreRunE: cmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		refID, _ := cmd.Flags().GetString("ref-id")
		planLogs, _ := cmd.Flags().GetBool("plan-logs")
		planDefaults, _ := cmd.Flags().GetBool("plan-defaults")
		metadata, _ := cmd.Flags().GetBool("metadata")
		settings, _ := cmd.Flags().GetBool("settings")
		plans, _ := cmd.Flags().GetBool("plans")
		var showPlans bool
		if planLogs || planDefaults || plans {
			showPlans = true
		}
		res, err := deployment.GetResource(deployment.GetResourceParams{
			Type: "appsearch",
			GetParams: deployment.GetParams{
			API: ecctl.Get().API,
			DeploymentID:  args[0],
			RefID: refID,
			QueryParams: deputil.QueryParams{
				ShowPlans:        showPlans,
				ShowPlanLogs:     planLogs,
				ShowPlanDefaults: planDefaults,
				ShowMetadata:     metadata,
				ShowSettings:     settings,
			},
		},
		})
		if err != nil {
			return err
		}
		return ecctl.Get().Formatter.Format(filepath.Join("appsearch", "show"), res)
	},
}

func init() {
	showAppSearchCmd.Flags().String("ref-id", "", "Optional RefId, auto-discovered if not specified")
	showAppSearchCmd.Flags().Bool("plans", false, "Shows the deployment plans")
	showAppSearchCmd.Flags().Bool("plan-logs", false, "Shows the deployment plan logs")
	showAppSearchCmd.Flags().Bool("plan-defaults", false, "Shows the deployment plan defaults")
	showAppSearchCmd.Flags().BoolP("metadata", "m", false, "Shows the deployment metadata")
	showAppSearchCmd.Flags().BoolP("settings", "s", false, "Shows the deployment settings")
}
