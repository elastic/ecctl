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

	"github.com/elastic/ecctl/pkg/deployment/deputil"
	"github.com/elastic/ecctl/pkg/deployment/elasticsearch"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var elasticSearchQueryExamples = `
* Returns ES Clusters that contain the string "logging" in the cluster name:
	ecctl elasticsearch list --query cluster_name:logging

* Returns ES Cluster with cluster ID bedd14a60f0871b53636496b5b4b0014:
	ecctl elasticsearch list --query cluster_id:bedd14a60f0871b53636496b5b4b0014

* Returns ES Clusters that contain the string "prod cluster" in the cluster name:
	ecctl elasticsearch list --query cluster_name:'prod cluster'

* To search in both cluster name and ID, ignore the query field:
	ecctl elasticsearch list --query logging

Read all the simple query string syntax in https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-query-string-query.html#query-string-syntax
`[1:]

var listElasticsearchCmd = &cobra.Command{
	Use:     "list",
	Short:   "Returns the list of Elasticsearch clusters",
	Example: elasticSearchQueryExamples,
	PreRunE: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		size, _ := cmd.Flags().GetInt64("size")
		metadata, _ := cmd.Flags().GetBool("metadata")

		r, err := elasticsearch.List(elasticsearch.ListParams{
			API:     ecctl.Get().API,
			Version: cmd.Flag("version").Value.String(),
			QueryParams: deputil.QueryParams{
				ShowMetadata: metadata,
				Size:         size,
				Query:        cmd.Flag("query").Value.String(),
			},
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
	Command.AddCommand(listElasticsearchCmd)
	listElasticsearchCmd.Flags().BoolP("metadata", "m", false, "Shows deployment metadata")
	listElasticsearchCmd.Flags().Int64P("size", "s", 100, "Sets the upper limit of ES clusters to return")
	listElasticsearchCmd.Flags().StringP("version", "v", "", "Filters per version")
	listElasticsearchCmd.Flags().String("query", "", "Searches clusters using an Elasticsearch search query string query")
}
