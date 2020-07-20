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

package cmdeskeystore

import (
	"github.com/spf13/cobra"

	"github.com/elastic/cloud-sdk-go/pkg/api/deploymentapi/eskeystoreapi"
	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"

	"github.com/elastic/ecctl/pkg/ecctl"
)

var showCmd = &cobra.Command{
	Use:     "show <deployment id> [--ref-id <ref-id>]",
	Short:   "Shows the settings from the Elasticsearch resource keystore",
	PreRunE: sdkcmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		refID, _ := cmd.Flags().GetString("ref-id")
		res, err := eskeystoreapi.Get(eskeystoreapi.GetParams{
			API:          ecctl.Get().API,
			DeploymentID: args[0],
			RefID:        refID,
		})
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format("deployment/eskeystore_show", res)
	},
}

func init() {
	Command.AddCommand(showCmd)
	showCmd.Flags().String("ref-id", "", "Optional ref_id to use for the Elasticsearch resource, auto-discovered if not specified.")
}
