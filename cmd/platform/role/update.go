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

package cmdrole

import (
	"fmt"

	"github.com/elastic/cloud-sdk-go/pkg/api/platformapi/roleapi"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/ecctl"
)

var updateCmd = &cobra.Command{
	Use:     "update <role>",
	Short:   "Updates an existing platform role",
	PreRunE: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename, _ := cmd.Flags().GetString("file")
		var r models.Role
		if err := cmdutil.DecodeFile(filename, &r); err != nil {
			return err
		}
		if *r.ID != args[0] {
			return fmt.Errorf("role id [%s] cannot be found in the role file", args[0])
		}
		return roleapi.Update(roleapi.UpdateParams{
			API:    ecctl.Get().API,
			ID:     args[0],
			Role:   &r,
			Region: ecctl.Get().Config.Region,
		})
	},
}

func init() {
	initFlags()
}

func initFlags() {
	Command.AddCommand(updateCmd)
	updateCmd.Flags().String("file", "", "File name of the role to update")
	cobra.MarkFlagFilename(updateCmd.Flags(), "file", "json")
	cobra.MarkFlagRequired(updateCmd.Flags(), "file")
}
