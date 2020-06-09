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

package cmdallocator

import (
	"path/filepath"

	"github.com/elastic/cloud-sdk-go/pkg/api/platformapi/allocatorapi"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/ecctl"
)

func showAllocator(cmd *cobra.Command, args []string) error {
	a, err := allocatorapi.Get(allocatorapi.GetParams{
		API: ecctl.Get().API,
		ID:  args[0],
	})
	if err != nil {
		return err
	}

	templateName := "show"
	if metadata, _ := cmd.Flags().GetBool("metadata"); metadata {
		templateName = "showmetadata"
	}

	return ecctl.Get().Formatter.Format(filepath.Join("allocator", templateName), a)
}

var showAllocatorCmd = &cobra.Command{
	Use:     "show <allocator id>",
	Short:   "Returns information about the allocator",
	PreRunE: cobra.MinimumNArgs(1),
	RunE:    showAllocator,
}

func init() {
	Command.AddCommand(showAllocatorCmd)
	showAllocatorCmd.Flags().Bool("metadata", false, "Show allocator metadata")
}
