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
	"github.com/elastic/cloud-sdk-go/pkg/api/commentapi"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var deleteCmd = &cobra.Command{
	Use:     "delete <comment id> --resource-type <resource-type> --resource-id <resource-id>",
	Short:   cmdutil.AdminReqDescription("Shows information about a resource comment"),
	PreRunE: cobra.ExactValidArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resourceType, _ := cmd.Flags().GetString("resource-type")
		resourceID, _ := cmd.Flags().GetString("resource-id")
		version, _ := cmd.Flags().GetString("version")

		err := commentapi.Delete(commentapi.DeleteParams{
			API:          ecctl.Get().API,
			Region:       ecctl.Get().Config.Region,
			CommentID:    args[0],
			ResourceID:   resourceID,
			ResourceType: resourceType,
			Version:      version,
		})

		if err != nil {
			return err
		}
		_, err = ecctl.Get().Config.OutputDevice.Write([]byte("comment deleted successfully"))
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	initDeleteFlags()
}

func initDeleteFlags() {
	Command.AddCommand(deleteCmd)

	deleteCmd.Flags().String("resource-type", "", "The kind of resource that a comment belongs to. "+
		"Should be one of [elasticsearch, kibana, apm, appsearch, enterprise_search, allocator, constructor, runner, proxy].")
	deleteCmd.Flags().String("resource-id", "", "ID of the resource that the comment belongs to.")
	deleteCmd.Flags().String("version", "", "If specified then checks for conflicts against the version stored in the persistent store.")

	deleteCmd.MarkFlagRequired("resource-type")
	deleteCmd.MarkFlagRequired("resource-id")
}
