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
	"os"

	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/deployment"
	"github.com/elastic/ecctl/pkg/deployment/depresource"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var deleteCmd = &cobra.Command{
	Use:     "delete <deployment id>",
	Short:   "Deletes a previously shut down AppSearch deployment resource",
	PreRunE: sdkcmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		refID, _ := cmd.Flags().GetString("ref-id")

		force, _ := cmd.Flags().GetBool("force")
		var msg = "This action will delete the AppSearch deployment resource and its configuration history. Do you want to continue? [y/n]: "
		if !force && !sdkcmdutil.ConfirmAction(msg, os.Stderr, os.Stdout) {
			return nil
		}

		return depresource.DeleteStateless(depresource.DeleteStatelessParams{
			ResourceParams: deployment.ResourceParams{
				API:          ecctl.Get().API,
				DeploymentID: args[0],
				Kind:         "appsearch",
				RefID:        refID,
			},
		})
	},
}

func init() {
	Command.AddCommand(deleteCmd)
	deleteCmd.Flags().String("ref-id", "", "Optional RefId, auto-discovered if not specified")
}
