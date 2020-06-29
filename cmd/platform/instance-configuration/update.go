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

package cmdinstanceconfig

import (
	"os"

	"github.com/elastic/cloud-sdk-go/pkg/api/platformapi/instanceconfigapi"
	"github.com/elastic/cloud-sdk-go/pkg/input"
	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/ecctl"
)

var updateCmd = &cobra.Command{
	Use:     "update <config id> -f <config.json>",
	Short:   "Overwrites an instance configuration",
	PreRunE: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := sdkcmdutil.FileOrStdin(cmd, "file"); err != nil {
			return err
		}

		file, err := input.NewFileOrReader(os.Stdin, cmd.Flag("file").Value.String())
		if err != nil {
			return err
		}
		defer file.Close()

		config, err := instanceconfigapi.NewConfig(file)
		if err != nil {
			return err
		}

		return instanceconfigapi.Update(instanceconfigapi.UpdateParams{
			API:    ecctl.Get().API,
			Region: ecctl.Get().Config.Region,
			ID:     args[0],
			Config: config,
		})
	},
}

func init() {
	Command.AddCommand(updateCmd)

	updateCmd.Flags().StringP("file", "f", "", "Instance configuration JSON file definition")
	updateCmd.MarkFlagRequired("file")
}
