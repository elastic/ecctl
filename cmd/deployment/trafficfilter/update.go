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
	"encoding/json"
	"errors"

	"github.com/elastic/cloud-sdk-go/pkg/api/deploymentapi/trafficfilterapi"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/ecctl"
)

const updateExample = `
* Return the current traffic filter state as a valid update payload.
  ecctl deployment traffic-filter update <traffic-filter id> --generate-payload > update.json

* After editing the file with your new values pass it as an argument to the --file flag.
  ecctl deployment traffic-filter update <traffic-filter id> --file update.json`

var updateCmd = &cobra.Command{
	Use:     "update <traffic-filter id> {--file <file-path> | --generate-payload}",
	Short:   "Updates a traffic-filter",
	Example: updateExample,
	PreRunE: cmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		genPayload, _ := cmd.Flags().GetBool("generate-payload")
		file, _ := cmd.Flags().GetString("file")

		if err := flagRequirements(genPayload, file); err != nil {
			return err
		}

		if genPayload {
			res, err := trafficfilterapi.Get(trafficfilterapi.GetParams{
				API: ecctl.Get().API,
				ID:  args[0],
			})
			if err != nil {
				return err
			}

			enc := json.NewEncoder(ecctl.Get().Config.OutputDevice)
			enc.SetIndent("", "  ")
			return enc.Encode(
				trafficfilterapi.NewUpdateRequestFromGet(res),
			)
		}

		var req models.TrafficFilterRulesetRequest
		if err := sdkcmdutil.DecodeFile(file, &req); err != nil {
			return err
		}
		res, err := trafficfilterapi.Update(trafficfilterapi.UpdateParams{
			API: ecctl.Get().API,
			ID:  args[0],
			Req: &req,
		})
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format("", res)
	},
}

func init() {
	initUpdateFlags()
}

func initUpdateFlags() {
	Command.AddCommand(updateCmd)
	updateCmd.Flags().Bool("generate-payload", false, "Outputs JSON which can be used as an argument for the --file flag.")
	updateCmd.Flags().String("file", "", "Path to the file containing the update JSON definition.")
	updateCmd.MarkFlagFilename("file", "json")
}

func flagRequirements(genPayload bool, file string) error {
	if genPayload == true && file != "" {
		return errors.New("both --file and --generate-payload are set. Only one may be used")
	}

	if genPayload == false && file == "" {
		return errors.New("one of --file or --generate-payload must be set")
	}

	return nil
}
