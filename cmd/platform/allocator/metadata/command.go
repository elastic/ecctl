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

package cmdallocatormetadata

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/platform/allocator"
)

// Command represents the allocator metadata command.
var Command = &cobra.Command{
	Use:     "metadata",
	Short:   "Manages an allocator's metadata",
	PreRunE: cobra.MaximumNArgs(0),
	Run:     func(cmd *cobra.Command, args []string) { cmd.Help() },
}

var allocatorMetadataSetCmd = &cobra.Command{
	Use:     "set <allocator id> <key> <value>",
	Short:   "Sets or updates a single metadata item to a given allocators metadata",
	PreRunE: cobra.MinimumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {

		var params = &allocator.MetadataSetParams{
			API:   ecctl.Get().API,
			ID:    args[0],
			Key:   args[1],
			Value: args[2],
		}
		err := allocator.SetAllocatorMetadataItem(*params)

		if err != nil {
			return err
		}
		return nil
	},
}

var allocatorMetadataDeleteCmd = &cobra.Command{
	Use:     "delete <allocator id> <key>",
	Short:   "Deletes a single metadata item from a given allocators metadata",
	PreRunE: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {

		var params = &allocator.MetadataDeleteParams{
			API: ecctl.Get().API,
			ID:  args[0],
			Key: args[1],
		}
		err := allocator.DeleteAllocatorMetadataItem(*params)

		if err != nil {
			return err
		}
		return nil
	},
}

var allocatorMetadataShowCmd = &cobra.Command{
	Use:     "show <allocator id>",
	Short:   "Shows allocator metadata",
	PreRunE: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var params = allocator.MetadataGetParams{
			API: ecctl.Get().API,
			ID:  args[0],
		}
		res, err := allocator.GetAllocatorMetadata(params)

		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format(filepath.Join("allocator", "show"), res)
	},
}

func init() {
	Command.AddCommand(
		allocatorMetadataShowCmd,
		allocatorMetadataSetCmd,
		allocatorMetadataDeleteCmd,
	)
}
