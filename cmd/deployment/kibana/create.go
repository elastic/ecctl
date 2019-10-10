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

package cmdkibana

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/elastic/cloud-sdk-go/pkg/input"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment/kibana"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/util"
)

var createKibanaCmd = &cobra.Command{
	Use:   "create -f <definition.json>",
	Short: "Creates a Kibana Cluster",

	PreRunE: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmdutil.FileOrStdin(cmd, "file-template"); err != nil {
			return err
		}

		def, err := parseKibanaDefinitionFile(cmd.Flag("file-template").Value.String())
		if err != nil {
			return err
		}

		if id := cmd.Flag("id").Value.String(); id != "" {
			def.ElasticsearchClusterID = &id
		}

		if name := cmd.Flag("name").Value.String(); name != "" {
			def.ClusterName = name
		}

		track, _ := cmd.Flags().GetBool("track")
		r, err := kibana.Create(kibana.CreateParams{
			CreateKibanaRequest: &def,
			DeploymentParams: kibana.DeploymentParams{
				API: ecctl.Get().API,
				ID:  cmd.Flag("id").Value.String(),
				TrackParams: util.TrackParams{
					Track:  track,
					Output: ecctl.Get().Config.OutputDevice,
				},
			},
		})
		if err != nil {
			return err
		}
		return ecctl.Get().Formatter.Format(
			filepath.Join("kibana", "create"),
			r,
		)
	},
}

func parseKibanaDefinitionFile(fp string) (models.CreateKibanaRequest, error) {
	var kinabaDef models.CreateKibanaRequest

	defFile, err := input.NewFileOrReader(os.Stdin, fp)
	if err != nil {
		return kinabaDef, err
	}
	defer defFile.Close()

	if err := json.NewDecoder(defFile).Decode(&kinabaDef); err != nil {
		return kinabaDef, err
	}

	return kinabaDef, nil
}

func init() {
	createKibanaCmd.Flags().BoolP("track", "t", false, cmdutil.TrackFlagMessage)
	createKibanaCmd.Flags().StringP("file-template", "f", "", "JSON file that contains the Kibana cluster definition")
	createKibanaCmd.Flags().String("id", "", "Optional ID to set for the Elasticsearch cluster (Overrides ID if present).")
	createKibanaCmd.Flags().String("name", "", "Optional name to set for the Kibana cluster (Overrides name if present).")
}
