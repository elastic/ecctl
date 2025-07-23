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
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/ecctl"
)

var listCmd = &cobra.Command{
	Use:     "list [--include-associations] [--single-region <region>]",
	Short:   "Lists the network security policies or traffic filter rulesets",
	PreRunE: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		region, _ := cmd.Flags().GetString("single-region")
		assoc, _ := cmd.Flags().GetBool("include-associations")

		res, err := trafficfilterapi.List(trafficfilterapi.ListParams{
			API:                 ecctl.Get().API,
			Region:              region,
			IncludeAssociations: assoc,
		})
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format("", res)
	},
}

func init() {
	initListFlags()
}

func initListFlags() {
	Command.AddCommand(listCmd)
	listCmd.Flags().String("single-region", "", "Optional flag to list traffic filters from a specific region only")
	listCmd.Flags().Bool("include-associations", false, "Optional flag to include all associated resources")
}
