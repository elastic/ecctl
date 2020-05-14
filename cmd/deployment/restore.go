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
	"github.com/elastic/cloud-sdk-go/pkg/api/deploymentapi"
	"github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/ecctl"
)

const restoreLong = `Use --restore-snapshot to automatically restore the latest available Elasticsearch snapshot (if applicable)`

var restoreExamples = `
$ ecctl deployment restore 5c17ad7c8df73206baa54b6e2829d9bc
{
  "id": "5c17ad7c8df73206baa54b6e2829d9bc"
}
`[1:]

var restoreCmd = &cobra.Command{
	Use:     "restore <deployment-id>",
	Short:   "Restores a previously shut down deployment and all of its associated sub-resources",
	Long:    restoreLong,
	Example: restoreExamples,
	PreRunE: cmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		restoreSnapshot, _ := cmd.Flags().GetBool("restore-snapshot")

		res, err := deploymentapi.Restore(deploymentapi.RestoreParams{
			API:             ecctl.Get().API,
			DeploymentID:    args[0],
			RestoreSnapshot: restoreSnapshot,
		})
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format("", res)
	},
}

func init() {
	Command.AddCommand(restoreCmd)
	restoreCmd.Flags().Bool("restore-snapshot", false, "Restores snapshots for those resources that allow it (Elasticsearch)")
}
