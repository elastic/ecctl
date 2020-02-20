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

package cmdrunner

import (
	"encoding/json"
	"path/filepath"
	"strconv"

	"github.com/elastic/cloud-sdk-go/pkg/models"
	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/platform/runner"
)

// Command represents the runner search command.
var searchCmd = &cobra.Command{
	Use:     "search",
	Short:   "Performs advanced runner searching",
	Long:    "Read more about Query DSL in https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl.html",
	PreRunE: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")
		query, _ := cmd.Flags().GetString("query")

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

		r, err := runner.Search(
			runner.SearchParams{
				Params: runner.Params{
					API: ecctl.Get().API,
				},
				Request: sr,
			})
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format(filepath.Join("runner", "list"), r)
	},
}

func init() {
	Command.AddCommand(searchCmd)
	searchCmd.Flags().StringP("file", "f", "", "JSON file that contains JSON-style domain-specific language query")
	searchCmd.Flags().String("query", "", "Optional argument that contains a JSON-style domain-specific language query")
}
