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
	"errors"
	"os"

	"github.com/elastic/cloud-sdk-go/pkg/api/deploymentapi/extensionapi"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/ecctl"
)

var createCmd = &cobra.Command{
	Use:     "create <extension name> --version <version> --type <extension type> {--file <file-path> | --download-url <url>} [--description <description>]",
	Short:   "Creates an extension",
	PreRunE: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		version, _ := cmd.Flags().GetString("version")
		extType, _ := cmd.Flags().GetString("type")
		url, _ := cmd.Flags().GetString("download-url")
		description, _ := cmd.Flags().GetString("description")
		file, _ := cmd.Flags().GetString("file")

		if url != "" && file != "" {
			return errors.New("both --file and --download-url are set. Only one may be used")
		}

		res, err := extensionapi.Create(extensionapi.CreateParams{
			API:         ecctl.Get().API,
			Name:        args[0],
			Version:     version,
			Type:        extType,
			DownloadURL: url,
			Description: description,
		})
		if err != nil {
			return err
		}

		if file != "" {
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer func() { _ = f.Close() }()

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
	initCreateFlags()
}

func initCreateFlags() {
	Command.AddCommand(createCmd)
	createCmd.Flags().String("version", "", "Elastic stack version. Numeric version for plugins, e.g. 7.10.0. Major version e.g. 7.*, or wildcards e.g. * for bundles.")
	createCmd.Flags().String("type", "", "Extension type. Can be one of [bundle, plugin].")
	createCmd.Flags().String("download-url", "", "Optional flag to define the URL to download the extension archive.")
	createCmd.Flags().String("description", "", "Optional flag to add a description to the extension.")
	createCmd.Flags().String("file", "", "Optional flag to upload an extension from a local file path.")
	cobra.MarkFlagRequired(createCmd.Flags(), "version")
	cobra.MarkFlagRequired(createCmd.Flags(), "type")
}
