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

package cmdelasticsearch

import (
	"path/filepath"

	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment/elasticsearch"
	"github.com/elastic/ecctl/pkg/ecctl"
)

const elasticsearchSearchQueryExamples = `Read more about Query DSL in https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl.html`

var searchElasticsearchCmd = &cobra.Command{
	Use:     `search -f <query file>.json`,
	Short:   "Performs advanced clusters searching",
	Long:    elasticsearchSearchQueryExamples,
	PreRunE: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {

		sr, err := cmdutil.ParseQueryDSLFile(cmd.Flag("file").Value.String())
		if err != nil {
			return err
		}

		r, err := elasticsearch.SearchClusters(elasticsearch.SearchClusterParams{
			API:     ecctl.Get().API,
			Request: sr,
		})
		if err != nil {
			return err
		}
		return ecctl.Get().Formatter.Format(
			filepath.Join("elasticsearch", "list"),
			r,
		)
	},
}

func init() {
	Command.AddCommand(searchElasticsearchCmd)
	searchElasticsearchCmd.Flags().StringP("file", "f", "", "JSON file that contains JSON-style domain-specific language query")
	searchElasticsearchCmd.MarkFlagRequired("file")
}
