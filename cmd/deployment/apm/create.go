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

package cmdapm

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/elastic/cloud-sdk-go/pkg/input"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment/apm"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/util"
)

var createApmCmd = &cobra.Command{
	Use:     "create -f <file definition> --id <deployment id>",
	Short:   "Creates an APM instance in the selected deployment",
	PreRunE: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmdutil.FileOrStdin(cmd, "file"); err != nil {
			return err
		}

		f, err := input.NewFileOrReader(os.Stdin, cmd.Flag("file").Value.String())
		if err != nil {
			return err
		}
		defer f.Close()

		var definition models.CreateApmRequest
		if err := json.NewDecoder(f).Decode(&definition); err != nil {
			return err
		}

		track, _ := cmd.Flags().GetBool("track")
		res, err := apm.Create(apm.CreateParams{
			API:              ecctl.Get().API,
			ID:               cmd.Flag("id").Value.String(),
			CreateApmRequest: definition,
			TrackParams: util.TrackParams{
				Track:  track,
				Output: ecctl.Get().Config.OutputDevice,
			},
		})
		if err != nil {
			return err
		}
		return ecctl.Get().Formatter.Format(filepath.Join("apm", "create"), res)
	},
}

func init() {
	createApmCmd.Flags().StringP("file", "f", "", "Deployment definition")
	createApmCmd.Flags().String("id", "", "Overrides the deployment ID on which it's tied")
	createApmCmd.Flags().Bool("track", false, cmdutil.TrackFlagMessage)
}
