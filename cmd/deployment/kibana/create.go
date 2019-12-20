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
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment"
	"github.com/elastic/ecctl/pkg/deployment/depresource"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var createKibanaCmd = &cobra.Command{
	Use:     "create --id <deployment-id>",
	Short:   "Creates a Kibana instance",
	Long:    kibanaCreateLong,
	Example: kibanaCreateExample,
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

		var payload *models.KibanaPayload
		if err := cmdutil.FileOrStdin(cmd, "file"); err == nil {
			err := cmdutil.DecodeDefinition(cmd, "file", &payload)
			if err != nil && err != cmdutil.ErrNodefinitionLoaded {
				return err
			}
		}

		if payload == nil {
			p, err := depresource.NewKibana(depresource.NewStateless{
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

		// Returns the KibanaPayload skipping the creation of the resources.
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
					Kibana: []*models.KibanaPayload{payload},
				},
			},
		}

		res, err := deployment.Update(updateParams)
		if err != nil {
			return err
		}

		return cmdutil.Track(cmdutil.TrackParams{
			TrackResourcesParams: depresource.TrackResourcesParams{
				API:          ecctl.Get().API,
				Resources:    res.Resources,
				Orphaned:     res.ShutdownResources,
				OutputDevice: ecctl.Get().Config.OutputDevice,
			},
			Formatter: ecctl.Get().Formatter,
			Track:     track,
			Response:  res,
		})
	},
}

func init() {
	createKibanaCmd.Flags().BoolP("track", "t", false, cmdutil.TrackFlagMessage)
	createKibanaCmd.Flags().StringP("file", "f", "", "KibanaPayload file definition. See help for more information")
	createKibanaCmd.Flags().String("deployment-template", "", "Optional deployment template ID, automatically obtained from the current deployment")
	createKibanaCmd.Flags().String("version", "", "Version to use, if not specified, the deployment's stack version will be used")
	createKibanaCmd.Flags().String("ref-id", "kibana", "RefId for the Kibana deployment")
	createKibanaCmd.Flags().String("id", "", "Deployment ID where to create the Kibana deployment")
	createKibanaCmd.MarkFlagRequired("id")
	createKibanaCmd.Flags().String("elasticsearch-ref-id", "", "Optional Elasticsearch ref ID where the Kibana deployment will connect to")
	createKibanaCmd.Flags().String("name", "", "Optional name to set for the Kibana deployment (Overrides name if present)")
	createKibanaCmd.Flags().Int32("zones", 1, "Number of zones the deployment will span")
	createKibanaCmd.Flags().Int32("size", 1024, "Memory (RAM) in MB that each of the deployment nodes will have")
	createKibanaCmd.Flags().Bool("generate-payload", false, "Returns the KibanaPayload without actually creating the deployment resources")
}
