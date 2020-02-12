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

package cmdstack

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/platform/stack"
)

// Command is the top level stack command.
var Command = &cobra.Command{
	Use:     "stack",
	Short:   "Manages Elastic StackPacks",
	PreRunE: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var stackListCmd = &cobra.Command{
	Use:     "list",
	Short:   "Lists Elastic StackPacks",
	PreRunE: cobra.MaximumNArgs(0),
	RunE:    listStackPacks,
}

var stackShowCmd = &cobra.Command{
	Use:     "show",
	Short:   "Shows information about an Elastic StackPack",
	PreRunE: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := stack.Get(stack.GetParams{
			API:     ecctl.Get().API,
			Version: args[0],
		})
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format(filepath.Join("stack", "show"), s)
	},
}

func listStackPacks(cmd *cobra.Command, args []string) error {
	deleted, _ := cmd.Flags().GetBool("deleted")
	s, err := stack.List(stack.ListParams{
		API:     ecctl.Get().API,
		Deleted: deleted,
	})
	if err != nil {
		return err
	}

	return ecctl.Get().Formatter.Format(filepath.Join("stack", "list"), s)
}

var stackUploadCmd = &cobra.Command{
	Use:     "upload",
	Short:   cmdutil.AdminReqDescription("Uploads an Elastic StackPack"),
	PreRunE: cobra.MinimumNArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer f.Close()

		return stack.Upload(stack.UploadParams{
			API:       ecctl.Get().API,
			StackPack: f,
		})
	},
}

var stackDeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   cmdutil.AdminReqDescription("Deletes an Elastic StackPack"),
	PreRunE: cobra.MinimumNArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		return stack.Delete(stack.DeleteParams{
			API:     ecctl.Get().API,
			Version: args[0],
		})
	},
}

func init() {
	Command.AddCommand(
		stackListCmd,
		stackShowCmd,
		stackUploadCmd,
		stackDeleteCmd,
	)

	stackListCmd.Flags().BoolP("deleted", "d", false, "Shows deleted stackpacks")
}
