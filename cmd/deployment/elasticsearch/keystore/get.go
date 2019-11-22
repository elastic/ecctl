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

package cmdelasticsearchkeystore

import (
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment/elasticsearch"
	"github.com/elastic/ecctl/pkg/ecctl"
)

const keystoreShowExample = `$ ecctl deployment elasticsearch keystore show 4c052fb17f65467a9b3c36d060106377
{
  "secrets": {
    "s3.client.foobar.access_key": {
      "as_file": false
    },
    "s3.client.foobar.secret_key": {
      "as_file": false
    }
  }
}`

var getCmd = &cobra.Command{
	Use:     `show <cluster id>`,
	Short:   "Shows the current Elasticsearch keystore settings",
	Example: keystoreShowExample,
	PreRunE: cmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := elasticsearch.GetKeystore(elasticsearch.GetKeystoreParams{
			API:       ecctl.Get().API,
			ClusterID: args[0],
		})
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format("", res)
	},
}

func init() {
	Command.AddCommand(getCmd)
}
