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

package cmdinstanceconfig

import (
	"path/filepath"
	"strconv"

	"github.com/elastic/cloud-sdk-go/pkg/api/platformapi/instanceconfigapi"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var showCmd = &cobra.Command{
	Use:     "show <config id>",
	Short:   cmdutil.AdminReqDescription("Shows an instance configuration"),
	PreRunE: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		configVersion, _ := cmd.Flags().GetString("config-version")
		showDeleted, _ := cmd.Flags().GetBool("show-deleted")

		var version *int64 = nil

		if configVersion != "" {
			parsedVersion, err := strconv.ParseInt(configVersion, 10, 64)

			version = ec.Int64(parsedVersion)

			if err != nil {
				return err
			}
		}

		res, err := instanceconfigapi.Get(instanceconfigapi.GetParams{
			API:           ecctl.Get().API,
			Region:        ecctl.Get().Config.Region,
			ID:            args[0],
			ConfigVersion: version,
			ShowDeleted:   showDeleted,
		})

		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format(filepath.Join("instance-configuration", "show"), res)
	},
}

func init() {
	Command.AddCommand(showCmd)
	showCmd.Flags().StringP("config-version", "v", "", "Instance configuration version")
	showCmd.Flags().Bool("show-deleted", false, "If set to true, allows to show deleted instance configurations")
}
