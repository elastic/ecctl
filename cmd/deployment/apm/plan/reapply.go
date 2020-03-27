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

package cmdapmplan

import (
	"path/filepath"

	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment/apm"
	"github.com/elastic/ecctl/pkg/deployment/planutil"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/util"
)

var reapplyLatestPlanCmd = &cobra.Command{
	Use:     "reapply <cluster id>",
	Short:   "Reapplies the latest plan attempt, resetting all transient settings",
	PreRunE: sdkcmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var reparams = planutil.ReapplyParams{ID: args[0]}
		// Obtains the parameter's value and sets it in params
		for name := range util.FieldsOfStruct(reparams) {
			val, err := cmd.Flags().GetBool(name)
			if err != nil {
				return err
			}
			util.Set(&reparams, name, val)
		}

		track, _ := cmd.Flags().GetBool("track")
		p, err := apm.ReapplyLatestPlanAttempt(
			apm.PlanParams{
				Track: track,
				TrackChangeParams: cmdutil.NewTrackParams(cmdutil.TrackParamsConfig{
					App:        ecctl.Get(),
					ResourceID: args[0],
					Kind:       util.Apm,
					Track:      track,
				}).TrackChangeParams,
				API: ecctl.Get().API,
				ID:  args[0],
			},
			reparams,
		)
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format(
			filepath.Join(cmd.Parent().Name(), cmd.Name()),
			p,
		)
	},
}

func init() {
	reapplyLatestPlanCmd.Flags().Bool("track", false, cmdutil.TrackFlagMessage)

	// Dynamically adds the transient flags from the tagged struct fields
	for name, description := range util.FieldsOfStruct(planutil.ReapplyParams{}) {
		reapplyLatestPlanCmd.Flags().Bool(name, false, description)
	}
}
