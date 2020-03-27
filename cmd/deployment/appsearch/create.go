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

package cmdappsearch

import (
	"github.com/elastic/cloud-sdk-go/pkg/models"
	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment"
	"github.com/elastic/ecctl/pkg/deployment/depresource"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var createAppSearchCmd = &cobra.Command{
	Use:     "create --id <deployment-id>",
	Short:   "Creates an AppSearch instance",
	Long:    appsearchCreateLong,
	Example: appsearchCreateExample,
	PreRunE: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		var generatePayload, _ = cmd.Flags().GetBool("generate-payload")
		var zoneCount, _ = cmd.Flags().GetInt32("zones")
		var size, _ = cmd.Flags().GetInt32("size")
		var name, _ = cmd.Flags().GetString("name")
		var refID, _ = cmd.Flags().GetString("ref-id")
		var esRefID, _ = cmd.Flags().GetString("elasticsearch-ref-id")
		var id, _ = cmd.Flags().GetString("id")
		var version, _ = cmd.Flags().GetString("version")
		var dt, _ = cmd.Flags().GetString("deployment-template")
		var region = ecctl.Get().Config.Region
		if ecctl.Get().Config.Region == "" {
			region = cmdutil.DefaultECERegion
		}

		var payload *models.AppSearchPayload
		if err := sdkcmdutil.FileOrStdin(cmd, "file"); err == nil {
			err := sdkcmdutil.DecodeDefinition(cmd, "file", &payload)
			if err != nil && err != sdkcmdutil.ErrNodefinitionLoaded {
				return err
			}
		}

		if payload == nil {
			p, err := depresource.NewAppSearch(depresource.NewStateless{
				DeploymentID:       id,
				ElasticsearchRefID: esRefID,
				API:                ecctl.Get().API,
				RefID:              refID,
				Version:            version,
				Region:             region,
				TemplateID:         dt,
				Size:               size,
				ZoneCount:          zoneCount,
			})
			if err != nil {
				return err
			}
			payload = p
		}

		if payload.Region == nil || *payload.Region == "" {
			payload.Region = ec.String(region)
		}

		// Returns the AppSearchPayload skipping the creation of the resources.
		if generatePayload {
			return ecctl.Get().Formatter.Format("", payload)
		}

		var updateParams = deployment.UpdateParams{
			DeploymentID: id,
			API:          ecctl.Get().API,
			Request: &models.DeploymentUpdateRequest{
				// Setting PruneOrphans to false since we don't want any side
				// effects on the deployment when only a partial deployment
				// definition is sent.
				PruneOrphans: ec.Bool(false),
				Name:         name,
				Resources: &models.DeploymentUpdateResources{
					Appsearch: []*models.AppSearchPayload{payload},
				},
			},
		}

		res, err := deployment.Update(updateParams)
		if err != nil {
			return err
		}

		var track, _ = cmd.Flags().GetBool("track")
		return cmdutil.Track(cmdutil.NewTrackParams(cmdutil.TrackParamsConfig{
			App:          ecctl.Get(),
			DeploymentID: id,
			Track:        track,
			Response:     res,
		}))
	},
}

func init() {
	Command.AddCommand(createAppSearchCmd)
	createAppSearchCmd.Flags().BoolP("track", "t", false, cmdutil.TrackFlagMessage)
	createAppSearchCmd.Flags().StringP("file", "f", "", "AppSearchPayload file definition. See help for more information")
	createAppSearchCmd.Flags().String("deployment-template", "", "Optional deployment template ID, automatically obtained from the current deployment")
	createAppSearchCmd.Flags().String("version", "", "Optional version to use. If not specified, it will default to the deployment's stack version")
	createAppSearchCmd.Flags().String("ref-id", "main-appsearch", "RefId for the AppSearch deployment")
	createAppSearchCmd.Flags().String("id", "", "Deployment ID where to create the AppSearch deployment")
	createAppSearchCmd.MarkFlagRequired("id")
	createAppSearchCmd.Flags().String("elasticsearch-ref-id", "", "Optional Elasticsearch ref ID where the AppSearch deployment will connect to")
	createAppSearchCmd.Flags().String("name", "", "Optional name to set for the AppSearch deployment (Overrides name if present)")
	createAppSearchCmd.Flags().Int32("zones", 1, "Number of zones the deployment will span")
	createAppSearchCmd.Flags().Int32("size", 2048, "Memory (RAM) in MB that each of the deployment nodes will have")
	createAppSearchCmd.Flags().Bool("generate-payload", false, "Returns the AppSearchPayload without actually creating the deployment resources")
}
