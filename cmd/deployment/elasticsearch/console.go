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
	"errors"

	"github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/deployment/elasticsearch"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/util"
)

var (
	errConsoleQueryNeedsAtLeast3Arguments = errors.New("needs at least 3 arguments (<cluster id> <METHOD> <URL> [body])")
)

var openElasticsearchConsoleCmd = &cobra.Command{
	Use:   "console <cluster id>",
	Short: "Starts an interactive console with the cluster",
	Long: `Starts an interactive console that connects to the cluster. If no username or password is specified
it will use the system's credentials to connect to the cluster (needs admin privileges)`,
	Example: `
	If run without admin credentials, it will need a user and pass specification:
    ecctl elasticsearch console 18fc96c491b3d5e10e147463927a5349 --elasticsearch-user user --elasticsearch-pass pass
`,
	PreRunE: cmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query, _ := cmd.Flags().GetBool("query")
		if query && len(args) <= 2 {
			return errConsoleQueryNeedsAtLeast3Arguments
		}

		insecure, _ := cmd.Flags().GetBool("insecure")
		return elasticsearch.Query(elasticsearch.QueryParams{
			ClusterParams: util.ClusterParams{
				API:       ecctl.Get().API,
				ClusterID: args[0],
			},
			AuthenticationParams: util.AuthenticationParams{
				User:     cmd.Flag("elasticsearch-user").Value.String(),
				Pass:     cmd.Flag("elasticsearch-pass").Value.String(),
				Insecure: insecure,
			},
			Interactive: !query,
			RestRequest: restRequest(args),
		})
	},
}

func restRequest(args []string) elasticsearch.RestRequest {
	if len(args) == 1 {
		return elasticsearch.RestRequest{}
	}

	var body string
	if len(args) > 3 {
		body = args[3]
	}

	return elasticsearch.RestRequest{
		Method: args[1],
		Path:   args[2],
		Body:   body,
	}
}

func init() {
	Command.AddCommand(openElasticsearchConsoleCmd)
	openElasticsearchConsoleCmd.Flags().String("elasticsearch-user", "", "Set the elasticsearch user for the interactive console")
	openElasticsearchConsoleCmd.Flags().String("elasticsearch-pass", "", "Set the elasticsearch password for the interactive console")
	openElasticsearchConsoleCmd.Flags().Bool("insecure", false, "skips tls certificate validation (use at your own risk)")
	openElasticsearchConsoleCmd.Flags().Bool("query", false, "Instead of opening a console it will query the cluster with the arguments passed")
}
