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

package cmddeploymentdemplate

import (
	"github.com/elastic/cloud-sdk-go/pkg/api/platformapi/configurationtemplateapi"

	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/ecctl"
)

var pullCmd = &cobra.Command{
	Use:     "pull --path <path>",
	Short:   "Downloads deployment template into a local folder",
	PreRunE: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return configurationtemplateapi.PullToFolder(configurationtemplateapi.PullToFolderParams{
			API:    ecctl.Get().API,
			Region: ecctl.Get().Config.Region,
			Folder: cmd.Flag("path").Value.String(),
		})
	},
}

func init() {
	Command.AddCommand(pullCmd)
	pullCmd.Flags().StringP("path", "p", "", "Local path where to store deployment templates")
	pullCmd.MarkFlagRequired("path")
}
