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

package cmdapm

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/deployment/apm"
	"github.com/elastic/ecctl/pkg/deployment/deputil"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var listApmCmd = &cobra.Command{
	Use:     "list",
	Short:   "Lists the APM deployments",
	PreRunE: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		planLogs, _ := cmd.Flags().GetBool("plan-logs")
		planDefaults, _ := cmd.Flags().GetBool("plan-defaults")
		metadata, _ := cmd.Flags().GetBool("metadata")
		plans, _ := cmd.Flags().GetBool("plans")
		hidden, _ := cmd.Flags().GetBool("hidden")
		size, _ := cmd.Flags().GetInt64("size")
		res, err := apm.List(apm.ListParams{
			API: ecctl.Get().API,
			QueryParams: deputil.QueryParams{
				ShowPlans:        plans,
				ShowPlanLogs:     planLogs,
				ShowPlanDefaults: planDefaults,
				ShowMetadata:     metadata,
				Size:             size,
				Query:            cmd.Flag("query").Value.String(),
				ShowHidden:       hidden,
			},
		})
		if err != nil {
			return err
		}
		return ecctl.Get().Formatter.Format(filepath.Join("apm", "list"), res)
	},
}

func init() {
	listApmCmd.Flags().Bool("plans", false, "Shows deployment plans")
	listApmCmd.Flags().Bool("plan-logs", false, "Shows deployment plan logs")
	listApmCmd.Flags().Bool("plan-defaults", false, "Shows deployment plan defaults")
	listApmCmd.Flags().BoolP("metadata", "m", false, "Shows deployment metadata")
	listApmCmd.Flags().Bool("hidden", false, "Shows hidden deployments")
	listApmCmd.Flags().String("query", "", "Custom Elasticsearch query to filter deployment")
	listApmCmd.Flags().Int64P("size", "s", 100, "Limits the number of deployments to the size value")
}
