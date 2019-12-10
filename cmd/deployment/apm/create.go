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
	"fmt"

	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment"
	"github.com/elastic/ecctl/pkg/deployment/depresource"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var createApmCmd = &cobra.Command{
	Use:     "create -f <definition.json>",
	Short:   "Creates a Apm instance",
	Long:    apmCreateLong,
	Example: apmCreateExample,
	PreRunE: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		var track, _ = cmd.Flags().GetBool("track")
		var generatePayload, _ = cmd.Flags().GetBool("generate-payload")
		var zoneCount, _ = cmd.Flags().GetInt32("zones")
		var size, _ = cmd.Flags().GetInt32("size")
		var name, _ = cmd.Flags().GetString("name")
		var refID, _ = cmd.Flags().GetString("ref-id")
		var esRefID, _ = cmd.Flags().GetString("elasticsearch-ref-id")
		var id, _ = cmd.Flags().GetString("id")
		var version, _ = cmd.Flags().GetString("version")
		var dt, _ = cmd.Flags().GetString("deployment-template")
		var region string
		if ecctl.Get().Config.Region == "" {
			region = cmdutil.DefaultECERegion
		}

		var payload *models.ApmPayload
		if err := cmdutil.FileOrStdin(cmd, "file"); err == nil {
			err := cmdutil.DecodeDefinition(cmd, "file", &payload)
			if err != nil && err != cmdutil.ErrNodefinitionLoaded {
				return err
			}
		}

		if payload == nil {
			p, err := depresource.NewApm(depresource.NewStateless{
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

		// Returns the ApmPayload skipping the creation of the resources.
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
					Apm: []*models.ApmPayload{payload},
				},
			},
		}

		res, err := deployment.Update(updateParams)
		if err != nil {
			return err
		}

		if err := ecctl.Get().Formatter.Format("", res); err != nil {
			if !track {
				return err
			}
			fmt.Fprintln(ecctl.Get().Config.OutputDevice, err)
		}

		if !track {
			return nil
		}

		return depresource.TrackResources(depresource.TrackResourcesParams{
			API:          ecctl.Get().API,
			Resources:    res.Resources,
			Orphaned:     res.ShutdownResources,
			OutputDevice: ecctl.Get().Config.OutputDevice,
		})
	},
}

func init() {
	Command.AddCommand(createApmCmd)
	createApmCmd.Flags().BoolP("track", "t", false, cmdutil.TrackFlagMessage)
	createApmCmd.Flags().StringP("file", "f", "", "ApmPayload file definition. See help for more information")
	createApmCmd.Flags().String("deployment-template", "", "Optional deployment template ID, automatically obtained from the current deployment")
	createApmCmd.Flags().String("version", "", "Version to use, if not specified, the deployment's stack version will be used")
	createApmCmd.Flags().String("ref-id", "apm", "RefId for the Apm deployment")
	createApmCmd.Flags().String("id", "", "Deployment ID where to create the Apm deployment")
	createApmCmd.MarkFlagRequired("id")
	createApmCmd.Flags().String("elasticsearch-ref-id", "", "Optional Elasticsearch ref ID where the Apm deployment will connect to")
	createApmCmd.Flags().String("name", "", "Optional name to set for the Apm deployment (Overrides name if present)")
	createApmCmd.Flags().Int32("zones", 1, "Number of zones the deployment will span")
	createApmCmd.Flags().Int32("size", 512, "Memory (RAM) in MB that each of the deployment nodes will have")
	createApmCmd.Flags().Bool("generate-payload", false, "Returns the ApmPayload without actually creating the deployment resources")
}
