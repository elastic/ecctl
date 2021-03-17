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

package cmdcomment

import (
	"path/filepath"

	"github.com/elastic/cloud-sdk-go/pkg/api/commentapi"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/ecctl"
)

var showCmd = &cobra.Command{
	Use:     "show <id> --resource-type <resource-type> --resource-id <resource-id>",
	Short:   "Shows a resource comment",
	PreRunE: cobra.ExactValidArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resourceType, _ := cmd.Flags().GetString("resource-type")
		resourceID, _ := cmd.Flags().GetString("resource-id")

		res, err := commentapi.Get(commentapi.GetParams{
			API:          ecctl.Get().API,
			Region:       ecctl.Get().Config.Region,
			CommentID:    args[0],
			ResourceID:   resourceID,
			ResourceType: resourceType,
		})

		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format(filepath.Join("comment", "show"), res)
	},
}

func init() {
	initShowFlags()
}

func initShowFlags() {
	Command.AddCommand(showCmd)

	showCmd.Flags().String("resource-type", "", "The kind of Resource that a Comment belongs to. Should be one of [elasticsearch, kibana, apm, appsearch, enterprise_search, allocator, constructor, runner, proxy].")
	showCmd.Flags().String("resource-id", "", "Id of the Resource that a Comment belongs to.")

	showCmd.MarkFlagRequired("resource-type")
	showCmd.MarkFlagRequired("resource-id")
}
