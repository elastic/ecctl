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
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/elastic/cloud-sdk-go/pkg/input"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment/elasticsearch"
	"github.com/elastic/ecctl/pkg/deployment/elasticsearch/plan"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/util"
)

const (
	esCreateExample = `
  * Create a single node elasticsearch cluster
  ecctl deployment elasticsearch create --version 5.6.0 --zones 1 --capacity 2048
 
  * Create an elasticsearch cluster from a plan definition.
  ecctl deployment elasticsearch create --file definition.json
`
)

var createElasticsearchCmd = &cobra.Command{
	Use:     "create [--file] [--capacity|--version|--zones] [name]",
	Short:   "Creates an Elasticsearch cluster",
	PreRunE: cobra.MaximumNArgs(1),
	Example: esCreateExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		var name string
		if len(args) == 1 {
			name = args[0]
		}

		var track, _ = cmd.Flags().GetBool("track")
		var zoneCount, _ = cmd.Flags().GetInt32("zones")
		var capacity, _ = cmd.Flags().GetInt32("capacity")
		var plugin, _ = cmd.Flags().GetStringSlice("plugin")
		var def *models.ElasticsearchClusterPlan

		if err := cmdutil.FileOrStdin(cmd, "file"); err == nil {
			reader, _ := input.NewFileOrReader(os.Stdin, cmd.Flag("file").Value.String())
			if reader != nil {
				if err := json.NewDecoder(reader).Decode(&def); err != nil {
					return err
				}
			}
		}

		r, err := elasticsearch.Create(elasticsearch.CreateParams{
			API:            ecctl.Get().API,
			ClusterName:    name,
			PlanDefinition: def,
			LegacyParams: plan.LegacyParams{
				Version:   cmd.Flag("version").Value.String(),
				ZoneCount: zoneCount,
				Capacity:  capacity,
				Plugins:   plugin,
			},
			TrackParams: util.TrackParams{
				Track:  track,
				Output: ecctl.Get().Config.OutputDevice,
			},
		})
		if err != nil {
			return err
		}
		return ecctl.Get().Formatter.Format(
			filepath.Join("elasticsearch", "create"),
			r,
		)
	},
}

func init() {
	Command.AddCommand(createElasticsearchCmd)
	createElasticsearchCmd.Flags().String("file", "", "JSON plan definition file location")
	createElasticsearchCmd.Flags().StringP("version", "v", "", "Filter per version")
	createElasticsearchCmd.Flags().Int32P("zones", "z", 0, "Number of zones for the cluster")
	createElasticsearchCmd.Flags().Int32P("capacity", "c", 0, "Capacity per node")
	createElasticsearchCmd.Flags().BoolP("track", "t", false, cmdutil.TrackFlagMessage)
	createElasticsearchCmd.Flags().StringSlice("plugin", nil, "Additional plugins to add to the Elasticsearch cluster")
}
