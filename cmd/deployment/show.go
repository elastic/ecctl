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

package cmddeployment

import (
	"errors"

	"github.com/elastic/cloud-sdk-go/pkg/util/booleans"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment"
	"github.com/elastic/ecctl/pkg/deployment/deputil"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var showCmd = &cobra.Command{
	Use:     "show <deployment-id>",
	Short:   "Shows the specified deployment resources",
	PreRunE: cmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		apm, _ := cmd.Flags().GetBool("apm")
		appsearch, _ := cmd.Flags().GetBool("app-search")
		elasticsearch, _ := cmd.Flags().GetBool("elasticsearch")
		kibana, _ := cmd.Flags().GetBool("kibana")
		planLogs, _ := cmd.Flags().GetBool("plan-logs")
		planDefaults, _ := cmd.Flags().GetBool("plan-defaults")
		metadata, _ := cmd.Flags().GetBool("metadata")
		settings, _ := cmd.Flags().GetBool("settings")
		plans, _ := cmd.Flags().GetBool("plans")
		var showPlans bool
		if planLogs || planDefaults || plans {
			showPlans = true
		}

		resourceFlags := []bool{apm, appsearch, elasticsearch, kibana}
		if !booleans.CheckNoneOrOneIsTrue(resourceFlags) {
			return errors.New("deployment: only one of --apm, --app-search, --elasticsearch, --kibana flags are allowed")
		}

		getParams := deployment.GetParams{
			API:          ecctl.Get().API,
			DeploymentID: args[0],
			QueryParams: deputil.QueryParams{
				ShowPlans:        showPlans,
				ShowPlanLogs:     planLogs,
				ShowPlanDefaults: planDefaults,
				ShowMetadata:     metadata,
				ShowSettings:     settings,
			},
		}

		if apm {
			res, err := deployment.GetApm(getParams)
			if err != nil {
				return err
			}
			return ecctl.Get().Formatter.Format("deployment/show", res)
		}

		if appsearch {
			res, err := deployment.GetAppSearch(getParams)
			if err != nil {
				return err
			}
			return ecctl.Get().Formatter.Format("deployment/show", res)
		}

		if elasticsearch {
			res, err := deployment.GetElasticsearch(getParams)
			if err != nil {
				return err
			}
			return ecctl.Get().Formatter.Format("deployment/show", res)
		}

		if kibana {
			res, err := deployment.GetKibana(getParams)
			if err != nil {
				return err
			}
			return ecctl.Get().Formatter.Format("deployment/show", res)
		}

		res, err := deployment.Get(getParams)
		if err != nil {
			return err
		}
		return ecctl.Get().Formatter.Format("deployment/show", res)
	},
}

func init() {
	showCmd.Flags().Bool("apm", false, "Shows APM resource information if any")
	showCmd.Flags().Bool("app-search", false, "Shows App Search resource information if any")
	showCmd.Flags().Bool("kibana", false, "Shows Kibana resource information if any")
	showCmd.Flags().Bool("elasticsearch", false, "Shows Elasticsearch resource information")
	showCmd.Flags().Bool("plans", false, "Shows the deployment plans")
	showCmd.Flags().Bool("plan-logs", false, "Shows the deployment plan logs")
	showCmd.Flags().Bool("plan-defaults", false, "Shows the deployment plan defaults")
	showCmd.Flags().BoolP("metadata", "m", false, "Shows the deployment metadata")
	showCmd.Flags().BoolP("settings", "s", false, "Shows the deployment settings")
}
