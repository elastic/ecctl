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
	"github.com/elastic/cloud-sdk-go/pkg/util/slice"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment"
	"github.com/elastic/ecctl/pkg/deployment/deputil"
	"github.com/elastic/ecctl/pkg/ecctl"
)

const showExample = `
* Shows kibana resource information from a given deployment:
  ecctl deployment show <deployment-id> --type kibana

* Shows apm resource information from a given deployment with a specified ref-id.
  ecctl deployment show <deployment-id> --type apm --ref-id apm-server`

var acceptedTypes = []string{"apm", "appsearch", "elasticsearch", "kibana"}

var showCmd = &cobra.Command{
	Use:     "show <deployment-id>",
	Short:   "Shows the specified deployment resources",
	Example: showExample,
	PreRunE: cmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resourceType, _ := cmd.Flags().GetString("type")
		if resourceType != "" && !slice.HasString(acceptedTypes, resourceType) {
			return errors.Errorf(`"%v" is not a valid resource type. Accepted resource types are: %v`, resourceType, acceptedTypes)
		}

		planLogs, _ := cmd.Flags().GetBool("plan-logs")
		planDefaults, _ := cmd.Flags().GetBool("plan-defaults")
		metadata, _ := cmd.Flags().GetBool("metadata")
		settings, _ := cmd.Flags().GetBool("settings")
		plans, _ := cmd.Flags().GetBool("plans")
		showPlans := planLogs || planDefaults || plans

		refID, _ := cmd.Flags().GetString("ref-id")
		getParams := deployment.GetParams{
			API:          ecctl.Get().API,
			DeploymentID: args[0],
			RefID:        refID,
			QueryParams: deputil.QueryParams{
				ShowPlans:        showPlans,
				ShowPlanLogs:     planLogs,
				ShowPlanDefaults: planDefaults,
				ShowMetadata:     metadata,
				ShowSettings:     settings,
			},
		}

		res, err := deployment.GetResource(deployment.GetResourceParams{
			GetParams: getParams,
			Type:      resourceType,
		})

		if err != nil {
			return err
		}
		return ecctl.Get().Formatter.Format("deployment/show", res)
	},
}

func init() {
	Command.AddCommand(showCmd)
	cmdutil.AddTypeFlag(showCmd, "Optional", true)
	showCmd.Flags().String("ref-id", "", "Optional deployment type RefId, if not set, the RefId will be auto-discovered")
	showCmd.Flags().Bool("plans", false, "Shows the deployment plans")
	showCmd.Flags().Bool("plan-logs", false, "Shows the deployment plan logs")
	showCmd.Flags().Bool("plan-defaults", false, "Shows the deployment plan defaults")
	showCmd.Flags().BoolP("metadata", "m", false, "Shows the deployment metadata")
	showCmd.Flags().BoolP("settings", "s", false, "Shows the deployment settings")
}
