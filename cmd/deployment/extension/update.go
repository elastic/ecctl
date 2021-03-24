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

package cmddeploymentextension

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/elastic/cloud-sdk-go/pkg/api/deploymentapi/extensionapi"
	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/ecctl"
)

const updateExample = `
* Return the current extension state as a valid update payload.
  ecctl deployment extension update <extension id> --generate-payload > update.json

* After editing the file with your new values pass it as an argument to the --file flag.
  ecctl deployment extension update <extension id> --file update.json

* The extensions uploaded from a local file will remain unchanged unless the --extension-file flag is used.
  ecctl deployment extension update <extension id> --file update.json --extension-file extension.zip`

var updateCmd = &cobra.Command{
	Use:     "update <extension id> {--file <file-path> | --generate-payload} [--extension-file <file path>]",
	Short:   "Updates an extension",
	Example: updateExample,
	PreRunE: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		genPayload, _ := cmd.Flags().GetBool("generate-payload")
		file, _ := cmd.Flags().GetString("file")
		extFile, _ := cmd.Flags().GetString("extension-file")

		if err := flagRequirements(genPayload, file); err != nil {
			return err
		}

		if genPayload {
			res, err := extensionapi.Get(extensionapi.GetParams{
				API:         ecctl.Get().API,
				ExtensionID: args[0],
			})
			if err != nil {
				return err
			}

			enc := json.NewEncoder(ecctl.Get().Config.OutputDevice)
			enc.SetIndent("", "  ")
			return enc.Encode(
				extensionapi.NewUpdateRequestFromGet(res),
			)
		}

		var req extensionapi.UpdateParams
		if err := sdkcmdutil.DecodeFile(file, &req); err != nil {
			return err
		}
		res, err := extensionapi.Update(extensionapi.UpdateParams{
			API:         ecctl.Get().API,
			ExtensionID: args[0],
			Name:        req.Name,
			Version:     req.Version,
			Type:        req.Type,
			DownloadURL: req.DownloadURL,
			Description: req.Description,
		})
		if err != nil {
			return err
		}

		if extFile != "" {
			f, err := os.Open(extFile)
			if err != nil {
				return err
			}
			defer f.Close()

			res2, err := extensionapi.Upload(extensionapi.UploadParams{
				API:         ecctl.Get().API,
				ExtensionID: *res.ID,
				File:        f,
			})
			if err != nil {
				return err
			}

			return ecctl.Get().Formatter.Format("", res2)
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
	updateCmd.Flags().String("extension-file", "", "Optional flag to upload an extension from a local file path.")
	updateCmd.MarkFlagFilename("file", "json")
}

func flagRequirements(genPayload bool, file string) error {
	if genPayload && file != "" {
		return errors.New("both --file and --generate-payload are set. Only one may be used")
	}

	if !genPayload && file == "" {
		return errors.New("one of --file or --generate-payload must be set")
	}

	return nil
}
