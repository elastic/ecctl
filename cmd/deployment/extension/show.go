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

package cmddeploymentextension

import (
	"github.com/elastic/cloud-sdk-go/pkg/api/deploymentapi/extensionapi"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/ecctl"
)

var showCmd = &cobra.Command{
	Use:     "show <extension id> [--include-deployments]",
	Short:   "Shows information about a deployment extension",
	PreRunE: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		deps, _ := cmd.Flags().GetBool("include-deployments")

		res, err := extensionapi.Get(extensionapi.GetParams{
			API:                ecctl.Get().API,
			ExtensionID:        args[0],
			IncludeDeployments: deps,
		})
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format("", res)
	},
}

func init() {
	initShowFlags()
}

func initShowFlags() {
	Command.AddCommand(showCmd)
	showCmd.Flags().Bool("include-deployments", false, "Include deployments referencing this extension. Up to only 10000 deployments will be included.")
}
