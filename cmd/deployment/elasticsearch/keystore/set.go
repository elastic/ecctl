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
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment/elasticsearch"
	"github.com/elastic/ecctl/pkg/ecctl"
)

const setKeystoreLong = `Manages the keystore settings of an Elasticsearch cluster.
Note that the underlying API call uses the PATCH method, meaning that each operation is add/modify only.`

var setKeystoreExamples = `
$ cat keystore_example.json
{
    "secrets": {
        "s3.client.foobar.access_key": {
            "value": "AKIXAIQFKXPHIFXSILUWPA",
            "as_file": false
        },
        "s3.client.foobar.secret_key": {
            "value": "18qXOpY2zGlApay1237dLXh+LG1X5LUNWjTHq5X1SWjf++m+p0"
        }
    }
}
$ ecctl deployment elasticsearch keystore set 4c052fb17f65467a9b3c36d060106377 --file keystore_example.json
{
  "secrets": {
    "s3.client.foobar.access_key": {
      "as_file": false
    },
    "s3.client.foobar.secret_key": {
      "as_file": false
    }
  }
}`[1:]

var setCmd = &cobra.Command{
	Use:     `set <cluster id> -f <file definition.json>`,
	Short:   "Updates an Elasticsearch cluster keystore with the contents of a file",
	Long:    setKeystoreLong,
	Example: setKeystoreExamples,
	PreRunE: cmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename, _ := cmd.Flags().GetString("file")
		var req models.KeystoreContents
		if err := cmdutil.DecodeFile(filename, &req); err != nil {
			return err
		}

		res, err := elasticsearch.SetKeystore(elasticsearch.SetKeystoreParams{
			API:       ecctl.Get().API,
			ClusterID: args[0],
			Request:   &req,
		})
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format("", res)
	},
}

func init() {
	Command.AddCommand(setCmd)
	setCmd.Flags().StringP("file", "f", "", "JSON file that contains JSON-style domain-specific keystore definition")
	setCmd.MarkFlagRequired("file")
	setCmd.MarkFlagFilename("file", "*.json")
}
