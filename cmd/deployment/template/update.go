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

package cmddeploymenttemplate

import (
	"errors"
	"fmt"
	"io"

	"github.com/elastic/cloud-sdk-go/pkg/api/deploymentapi/deptemplateapi"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/ecctl"
)

var updateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Updates an existing deployment template",
	PreRunE: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		var req *models.DeploymentTemplateRequestBody
		if err := sdkcmdutil.DecodeDefinition(cmd, "file", &req); err != nil {
			if errors.Is(err, io.EOF) {
				return fmt.Errorf("%s: %s", errReadingDefPrefix, errReadingDefMessage)
			}
			return fmt.Errorf("%s: %s", errReadingDefPrefix, err.Error())
		}

		templateID, _ := cmd.Flags().GetString("template-id")
		return deptemplateapi.Update(deptemplateapi.UpdateParams{
			API:        ecctl.Get().API,
			Region:     ecctl.Get().Config.Region,
			TemplateID: templateID,
			Request:    req,
		})
	},
}

func init() {
	Command.AddCommand(updateCmd)
	updateCmd.Flags().Bool("hide-instance-configurations", false, "Hides instance configurations - only visible when using the JSON output.")
	updateCmd.Flags().String("file", "", "It will cause the returned deployment templates which are valid for the specified stacks version.")
	updateCmd.Flags().String("template-id", "", "It will create the deployment template with the specified ID rather than auto-generating an ID.")
	cobra.MarkFlagFilename(updateCmd.Flags(), "file", "json")
	updateCmd.MarkFlagRequired("template-id")
}
