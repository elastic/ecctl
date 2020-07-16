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

var (
	errReadingDefPrefix  = "failed reading deployment template definition"
	errReadingDefMessage = "provide a valid deployment template definition using the --file flag"
)

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Creates a new deployment template",
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
		res, err := deptemplateapi.Create(deptemplateapi.CreateParams{
			API:        ecctl.Get().API,
			Region:     ecctl.Get().Config.Region,
			TemplateID: templateID,
			Request:    req,
		})
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format("id", models.IDResponse{ID: &res})
	},
}

func init() {
	Command.AddCommand(createCmd)
	createCmd.Flags().Bool("hide-instance-configurations", false, "Hides instance configurations - only visible when using the JSON output.")
	createCmd.Flags().String("file", "", "It will cause the returned deployment templates which are valid for the specified stacks version.")
	createCmd.Flags().String("template-id", "", "It will create the deployment template with the specified ID rather than auto-generating an ID.")
	cobra.MarkFlagFilename(createCmd.Flags(), "file", "json")
}
