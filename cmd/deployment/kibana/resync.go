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

package cmdkibana

import (
	"fmt"

	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment/kibana"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var resyncKibanaCmd = &cobra.Command{
	Use:     "resync [<deployment id> | --all]",
	Short:   "Resynchronizes the search index and cache for the selected Kibana instance",
	PreRunE: cmdutil.CheckInputHas1ArgsOr0ArgAndAll,
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")

		if all {
			res, err := kibana.ResyncAll(kibana.ResyncAllParams{
				API: ecctl.Get().API,
			})
			if err != nil {
				return err
			}

			return ecctl.Get().Formatter.Format("", res)
		}

		fmt.Printf("Resynchronizing Kibana instance: %s\n", args[0])
		return kibana.Resync(kibana.DeploymentParams{
			API: ecctl.Get().API,
			ID:  args[0],
		})
	},
}

func init() {
	resyncKibanaCmd.Flags().Bool("all", false, "Resynchronizes the search index for all Kibana instances")
}
