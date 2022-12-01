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

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var createCmd = &cobra.Command{
	Use:     "create <message> --resource-type <resource-type> --resource-id <resource-id>",
	Short:   cmdutil.AdminReqDescription("Creates a new resource comment"),
	PreRunE: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	RunE: func(cmd *cobra.Command, args []string) error {
		resourceType, _ := cmd.Flags().GetString("resource-type")
		resourceID, _ := cmd.Flags().GetString("resource-id")

		res, err := commentapi.Create(commentapi.CreateParams{
			API:          ecctl.Get().API,
			Region:       ecctl.Get().Config.Region,
			Message:      args[0],
			ResourceID:   resourceID,
			ResourceType: resourceType,
		})

		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format(filepath.Join("comment", "create"), res)
	},
}

func init() {
	initCreateFlags()
}

func initCreateFlags() {
	Command.AddCommand(createCmd)
	createCmd.Flags().String("resource-type", "", "The kind of resource that a comment belongs to. "+
		"Should be one of [elasticsearch, kibana, apm, appsearch, enterprise_search, allocator, constructor, runner, proxy].")
	createCmd.Flags().String("resource-id", "", "ID of the resource that a comment belongs to.")
	createCmd.MarkFlagRequired("resource-type")
	createCmd.MarkFlagRequired("resource-id")
}
