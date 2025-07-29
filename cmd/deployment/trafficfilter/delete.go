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
	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/ecctl"
)

var deleteCmd = &cobra.Command{
	Use:     "delete <ruleset id> [--ignore-associations]",
	Short:   "Deletes a network security policy or traffic filter ruleset",
	PreRunE: sdkcmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		assoc, _ := cmd.Flags().GetBool("ignore-associations")

		return trafficfilterapi.Delete(trafficfilterapi.DeleteParams{
			API:                ecctl.Get().API,
			ID:                 args[0],
			IgnoreAssociations: assoc,
		})
	},
}

func init() {
	initDeleteFlags()
}

func initDeleteFlags() {
	Command.AddCommand(deleteCmd)
	deleteCmd.Flags().Bool("ignore-associations", false, "Optional flag to delete the ruleset even if it has associated rules")
}
