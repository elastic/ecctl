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
	"github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/deployment"
	"github.com/elastic/ecctl/pkg/ecctl"
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

		res, err := deployment.Search(deployment.SearchParams{
			API:     ecctl.Get().API,
			Request: &sr,
		})

		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format("deployment/search", res)
	},
}

func init() {
	Command.AddCommand(searchCmd)
	searchCmd.Flags().StringP("file", "f", "", "JSON file that contains JSON-style domain-specific language query")
	searchCmd.MarkFlagRequired("file")
}
