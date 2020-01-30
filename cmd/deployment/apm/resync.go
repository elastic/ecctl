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
	"fmt"

	"github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/deployment/apm"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var resyncApmCmd = &cobra.Command{
	Use:     "resync {<deployment id> | --all}",
	Short:   "Resynchronizes the search index and cache for the selected APM deployment or all APM deployments",
	PreRunE: cmdutil.CheckInputHas1ArgsOr0ArgAndAll,
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")

		if all {
			fmt.Println("Resynchronizing all APM deployments")
			res, err := apm.ResyncAll(apm.ResyncAllParams{
				API: ecctl.Get().API,
			})
			if err != nil {
				return err
			}

			return ecctl.Get().Formatter.Format("", res)
		}

		fmt.Printf("Resynchronizing APM deployment: %s\n", args[0])
		return apm.Resync(apm.ResyncParams{
			API:       ecctl.Get().API,
			ClusterID: args[0],
		})
	},
}

func init() {
	resyncApmCmd.Flags().Bool("all", false, "Resynchronizes the search index for all APM instances")
}
