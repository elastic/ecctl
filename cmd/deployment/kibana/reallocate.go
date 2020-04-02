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
	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment/kibana"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/util"
)

var reallocateKibanaClusterCmd = &cobra.Command{
	Use:     "reallocate <cluster id>",
	Short:   "Reallocates Kibana instances",
	Long:    "Reallocates Kibana instances. If no \"--instances\" are specified all of the nodes will be restarted",
	PreRunE: sdkcmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		instances, _ := cmd.Flags().GetStringSlice("instances")
		return kibana.Reallocate(kibana.ReallocateParams{
			DeploymentParams: kibana.DeploymentParams{
				ID:    args[0],
				API:   ecctl.Get().API,
				Track: true,
				TrackChangeParams: cmdutil.NewTrackParams(cmdutil.TrackParamsConfig{
					App:        ecctl.Get(),
					ResourceID: args[0],
					Kind:       util.Kibana,
					Track:      true,
				}).TrackChangeParams,
			},
			Instances: instances,
		})
	},
}

func init() {
	reallocateKibanaClusterCmd.Flags().StringSliceP("instances", "i", nil, "Reallocate only specific instances")
}
