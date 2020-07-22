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
	"encoding/json"
	"path/filepath"
	"strconv"

	"github.com/elastic/cloud-sdk-go/pkg/api/platformapi/allocatorapi"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
)

const (
	fileArg       = "file"
	queryArg      = "query"
	queryExamples = `Read more about Query DSL in https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl.html`
)

// Command represents the allocator search command.
var searchAllocatorCmd = &cobra.Command{
	Use:     `search`,
	Short:   cmdutil.AdminReqDescription("Performs advanced allocator searching"),
	Long:    queryExamples,
	PreRunE: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString(fileArg)
		query, _ := cmd.Flags().GetString(queryArg)

		if err := cmdutil.ConflictingFlags(cmd, "file", "query"); err != nil {
			return err
		}

		if err := cmdutil.MustUseAFlag(cmd, "file", "query"); err != nil {
			return err
		}

		var sr models.SearchRequest
		var err error

		if len(file) > 0 {
			sr, err = sdkcmdutil.ParseQueryDSLFile(cmd.Flag("file").Value.String())
		} else {
			err = json.Unmarshal([]byte(query), &sr)
			if err != nil {
				uq, _ := strconv.Unquote(query)
				err = json.Unmarshal([]byte(uq), &sr)
			}
		}

		if err != nil {
			return err
		}

		r, err := allocatorapi.Search(allocatorapi.SearchParams{
			API:     ecctl.Get().API,
			Request: sr,
			Region:  ecctl.Get().Config.Region,
		})
		if err != nil {
			return err
		}
		return ecctl.Get().Formatter.Format(filepath.Join("allocator", "list"), r)
	},
}

func init() {
	Command.AddCommand(searchAllocatorCmd)
	searchAllocatorCmd.Flags().StringP(fileArg, "f", "", "JSON file that contains JSON-style domain-specific language query")
	searchAllocatorCmd.Flags().String(queryArg, "", "Optional argument that contains a JSON-style domain-specific language query")
}
