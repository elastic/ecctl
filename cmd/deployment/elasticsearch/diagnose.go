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
	"fmt"
	"os"
	"path/filepath"

	"github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/deployment/elasticsearch"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/util"
)

var generateElasticsearchDiagnosticsCmd = &cobra.Command{
	Use:     "diagnose <cluster id>",
	Short:   "Generates a diagnostics bundle for the cluster",
	Long:    "Generates a diagnostics bundle for the cluster, a timeout increase might be necessary",
	PreRunE: cmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		location := cmd.Flag("location").Value.String()
		filename := fmt.Sprint("diagnostic-", args[0], ".zip")
		path := filepath.Join(location, filename)

		f, err := os.Create(path)
		if err != nil {
			return err
		}
		defer f.Close()

		err = elasticsearch.Diagnose(elasticsearch.DiagnoseParams{
			ClusterParams: util.ClusterParams{
				ClusterID: args[0],
				API:       ecctl.Get().API,
			},
			Writer: f,
		})
		if err != nil {
			return err
		}

		fmt.Printf("Cluster [%s][Elasticsearch]: diagnostic bundle created: %s\n", args[0], path)
		return nil
	},
}

func init() {
	Command.AddCommand(generateElasticsearchDiagnosticsCmd)
	generateElasticsearchDiagnosticsCmd.Flags().StringP("location", "l", "./", "Directory location to store diagnostics file")
}
