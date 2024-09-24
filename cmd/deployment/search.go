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

package cmddeployment

import (
	"fmt"
	"github.com/elastic/cloud-sdk-go/pkg/api/deploymentapi"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/spf13/cobra"
)

const searchQueryLong = `Read more about Query DSL in https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl.html`

var searchExamples = `
$ cat query_string_query.json
{
    "query": {
        "query_string": {
            "query": "name: admin"
        }
    }
}
$ ecctl deployment search -f query_string_query.json
[...]`[1:]

var searchCmd = &cobra.Command{
	Use:     `search -f <query file.json>`,
	Short:   "Performs advanced deployment search using the Elasticsearch Query DSL",
	Long:    searchQueryLong,
	Example: searchExamples,
	PreRunE: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		sr, err := cmdutil.ParseQueryDSLFile(cmd.Flag("file").Value.String())
		if err != nil {
			return err
		}

		returnAllMatches, _ := cmd.Flags().GetBool("all-matches")
		if returnAllMatches && sr.Sort == nil {
			return fmt.Errorf("The query must include a sort-field when using --all-matches. Example: \"sort\": [\"id\"]")
		}

		batchSize, _ := cmd.Flags().GetInt32("size")

		var result *models.DeploymentsSearchResponse
		var cursor string
		for i := 0; i < 100; i++ {
			sr.Cursor = cursor
			if returnAllMatches {
				// Custom batch-size to override any size already set in the input query
				sr.Size = batchSize
			}

			res, err := deploymentapi.Search(deploymentapi.SearchParams{
				API:     ecctl.Get().API,
				Request: &sr,
			})

			if err != nil {
				return err
			}

			cursor = res.Cursor

			if result == nil {
				result = res
				result.Cursor = "" // Hide cursor in output
			} else {
				result.Deployments = append(result.Deployments, res.Deployments...)
				newReturnCount := *result.ReturnCount + *res.ReturnCount
				result.ReturnCount = &newReturnCount
				result.MatchCount = newReturnCount
			}

			if *res.ReturnCount == 0 || !returnAllMatches {
				break
			}
		}

		return ecctl.Get().Formatter.Format("deployment/search", result)
	},
}

func init() {
	Command.AddCommand(searchCmd)
	searchCmd.Flags().StringP("file", "f", "", "JSON file that contains JSON-style domain-specific language query")
	searchCmd.MarkFlagRequired("file")
	searchCmd.Flags().BoolP("all-matches", "a", false,
		"Uses a cursor to return all matches of the query (ignoring the size in the query). This can be used to query more than 10k results.")
	searchCmd.Flags().Int32("size", 500, "Defines the size per request when using the --all-matches option.")
}
