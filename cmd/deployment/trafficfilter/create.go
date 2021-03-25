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

package cmddeploymenttrafficfilter

import (
	"github.com/elastic/cloud-sdk-go/pkg/api/deploymentapi/trafficfilterapi"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/ecctl"
)

var createCmd = &cobra.Command{
	Use:     "create --region <region> --name <filter name> --type <filter type> --source <filter source>,<filter source> ",
	Short:   "Creates traffic filter rulesets",
	PreRunE: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		region := ecctl.Get().Config.Region
		name, _ := cmd.Flags().GetString("name")
		ftype, _ := cmd.Flags().GetString("type")
		description, _ := cmd.Flags().GetString("description")
		source, _ := cmd.Flags().GetStringSlice("source")
		include, _ := cmd.Flags().GetBool("include-by-default")

		var rules []*models.TrafficFilterRule
		for _, rule := range source {
			rules = append(rules, &models.TrafficFilterRule{
				Source: rule,
			})
		}

		res, err := trafficfilterapi.Create(trafficfilterapi.CreateParams{
			API: ecctl.Get().API,
			Req: &models.TrafficFilterRulesetRequest{
				Description:      description,
				IncludeByDefault: ec.Bool(include),
				Name:             ec.String(name),
				Region:           ec.String(region),
				Rules:            rules,
				Type:             ec.String(ftype),
			},
		})
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format("", res)
	},
}

func init() {
	initCreateFlags()
}

func initCreateFlags() {
	Command.AddCommand(createCmd)
	createCmd.Flags().String("name", "", "Name for the traffic filter.")
	createCmd.Flags().String("description", "", "Optional description for the traffic filter.")
	createCmd.Flags().String("type", "", "Type of traffic filter. Can be one of [ip, vpce])")
	createCmd.Flags().StringSlice("source", nil, "List of allowed traffic filter sources. Can be IP addresses, CIDR masks, or VPC endpoint IDs")
	createCmd.Flags().Bool("include-by-default", false, "If set, any future eligible deployments will have this filter applied automatically.")
}
