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
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/deployment/elasticsearch/plan"
	"github.com/elastic/ecctl/pkg/deployment/planutil"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/util"
)

var reapplyLatestElasticsearchPlanCmd = &cobra.Command{
	Use:     "reapply <cluster id>",
	Short:   "Reapplies the latest plan attempt, resetting all transient settings",
	PreRunE: cobra.MinimumNArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		commonParams := &planutil.ReapplyParams{ID: args[0]}
		for name := range util.FieldsOfStruct(commonParams) {
			val, err := cmd.Flags().GetBool(name)
			if err != nil {
				return err
			}
			util.Set(commonParams, name, val)
		}

		var params = &plan.ReapplyParams{
			ReapplyParams: *commonParams,
			API:           ecctl.Get().API,
			TrackParams: util.TrackParams{
				Track:  true,
				Output: ecctl.Get().Config.OutputDevice,
			},
		}
		for name := range util.FieldsOfStruct(params) {
			val, err := cmd.Flags().GetBool(name)
			if err != nil {
				return err
			}
			util.Set(params, name, val)
		}

		p, err := plan.Reapply(*params)
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format(
			filepath.Join("elasticsearch", "reapply"),
			p,
		)
	},
}

func init() {
	// Dynamically adds the transient flags from the tagged struct fields
	for name, description := range util.FieldsOfStruct(plan.ReapplyParams{}) {
		reapplyLatestElasticsearchPlanCmd.Flags().Bool(name, false, description)
	}
	// Dynamically adds the transient flags from the tagged struct fields
	for name, description := range util.FieldsOfStruct(plan.ReapplyParams{}.ReapplyParams) {
		reapplyLatestElasticsearchPlanCmd.Flags().Bool(name, false, description)
	}
}
