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
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/platform/allocator"
)

const (
	allocatorFilterCmdMessage = `Post-query filter out allocators based on metadata tags, for instance 'instanceType:i3.large'`
	allocatorListMessage      = `Returns all allocators that have instances or are connected to the platform. Use --all flag or --output json to show all. Use --query to match any of the allocators properties.`
	allocatorQueryExample     = `

Query examples:

	* Allocators set to maintenance mode: --query status.maintenance_mode:true

	* Allocators with more than 10GB of capacity: --query capacity.memory.total:\>10240

  Read all the simple query string syntax in https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-query-string-query.html#query-string-syntax
	`
	queryFlagHelp          = "Searches allocators using an Elasticsearch search query string query"
	allocatorFilterExample = `

Filter examples:

	* Allocators with instance type i3.large : --filter instanceType:i3.large

	* Allocators with instance type i3.large AND instance family gcp.highcpu.1 : --filter instanceType:i3.large --filter instanceFamily:gcp.highcpu.1

Filter is a post-query action and doesn't support OR. All filter arguments are applied as AND.
  
Filter and query flags can be used in combination.
	`
)

func listAllocators(cmd *cobra.Command, args []string) error {
	var queryString = cmd.Flag("query").Value.String()
	if unhealthy, _ := cmd.Flags().GetBool("unhealthy"); unhealthy {
		queryString = allocator.UnhealthyQuery
	}

	allFlag, _ := strconv.ParseBool(cmd.Flag("all").Value.String())
	a, err := allocator.List(allocator.ListParams{
		API:        ecctl.Get().API,
		Query:      queryString,
		FilterTags: cmd.Flag("filter").Value.String(),
		ShowAll:    allFlag,
	})
	if err != nil {
		return err
	}

	templateName := "list"
	if metadata, _ := cmd.Flags().GetBool("metadata"); metadata {
		templateName = "listmetadata"
	}

	if cmd.Flag("output").Value.String() != "json" && !allFlag {
		fmt.Printf("Showing allocators that have instances or are connected in the platform. Use --all flag or --output json to show all\n")
	}

	return ecctl.Get().Formatter.Format(filepath.Join("allocator", templateName), a)
}

var listAllocatorsCmd = &cobra.Command{
	Use:     "list",
	Short:   allocatorListMessage,
	Long:    allocatorListMessage + allocatorQueryExample + allocatorFilterExample,
	PreRunE: cobra.MaximumNArgs(0),
	RunE:    listAllocators,
}

func init() {
	Command.AddCommand(listAllocatorsCmd)

	listAllocatorsCmd.Flags().StringArrayP("filter", "f", nil, allocatorFilterCmdMessage)
	listAllocatorsCmd.Flags().Bool("unhealthy", false, "Searches for unhealthy allocators")
	listAllocatorsCmd.Flags().String("query", "", queryFlagHelp)
	listAllocatorsCmd.Flags().Bool("metadata", false, "Shows allocators metadata")
	listAllocatorsCmd.Flags().Bool("all", false, "Shows all allocators (including those with no instances or not connected)")
}
