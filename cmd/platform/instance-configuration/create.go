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
	"path/filepath"

	"github.com/elastic/cloud-sdk-go/pkg/api/platformapi/instanceconfigapi"
	"github.com/elastic/cloud-sdk-go/pkg/input"
	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/ecctl"
)

var createCmd = &cobra.Command{
	Use:     "create -f <config.json>",
	Short:   "Creates a new instance configuration",
	PreRunE: cobra.NoArgs,
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

		if id := cmd.Flag("id").Value.String(); id != "" {
			config.ID = id
		}

		res, err := instanceconfigapi.Create(instanceconfigapi.CreateParams{
			API:    ecctl.Get().API,
			Region: ecctl.Get().Config.Region,
			Config: config,
		})

		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format(filepath.Join("instance-configuration", "create"), res)
	},
}

func init() {
	Command.AddCommand(createCmd)

	createCmd.Flags().StringP("file", "f", "", "Instance configuration JSON file definition")
	createCmd.Flags().String("id", "", "Optional ID to set for the instance configuration (Overrides id if present)")
	createCmd.MarkFlagRequired("file")
}
