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

package cmdelasticsearchplan

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/elastic/cloud-sdk-go/pkg/input"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment/elasticsearch/plan"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/util"
)

var updateElasticsearchPlanCmd = &cobra.Command{
	Use:   "update <cluster id>",
	Short: "Applies or validates the provided plan and tracks the resulting change attempt",
	Long: `Applies the provided plan and tracks the resulting change attempt.
If --validate is set, the plan will simply be validated by the API failing if the plan is invalid. 
The plan can either be provided via the --file parameter or piped via stdin. 
The plan must be a valid ElasticsearchClusterPlan and will be validated against the specified cluster by the API.`,
	PreRunE: cobra.MinimumNArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		if err := sdkcmdutil.FileOrStdin(cmd, "file"); err != nil {
			return err
		}

		f, err := input.NewFileOrReader(os.Stdin, cmd.Flag("file").Value.String())
		if err != nil {
			return err
		}
		defer f.Close()

		var planModel models.ElasticsearchClusterPlan
		if err := json.NewDecoder(f).Decode(&planModel); err != nil {
			return err
		}

		track, _ := cmd.Flags().GetBool("track")
		validate, _ := cmd.Flags().GetBool("validate")
		p, err := plan.Update(plan.UpdateParams{
			API:          ecctl.Get().API,
			ID:           args[0],
			ValidateOnly: validate,
			Plan:         planModel,
			TrackParams: util.TrackParams{
				Track:  track,
				Output: ecctl.Get().Config.OutputDevice,
			},
		})
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format(
			filepath.Join("elasticsearch", "update"),
			p,
		)
	},
}

func init() {
	updateElasticsearchPlanCmd.Flags().StringP("file", "f", "", "Provide the location of a file containing the plan JSON")
	updateElasticsearchPlanCmd.Flags().BoolP("validate", "v", false, "Only validate the plan")
	updateElasticsearchPlanCmd.Flags().Bool("track", true, cmdutil.TrackFlagMessage)
}
