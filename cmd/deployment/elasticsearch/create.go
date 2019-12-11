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
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment"
	"github.com/elastic/ecctl/pkg/deployment/depresource"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var createElasticsearchCmd = &cobra.Command{
	Use:     "create {--file|--size <int> --version <string> --zones <string>|--topology-element <obj>}",
	Short:   "Creates a deployment with (only) an Elasticsearch resource in it",
	PreRunE: cobra.MaximumNArgs(0),
	Long:    esCreateLong,
	Example: esCreateExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		var track, _ = cmd.Flags().GetBool("track")
		var generatePayload, _ = cmd.Flags().GetBool("generate-payload")
		var zoneCount, _ = cmd.Flags().GetInt32("zones")
		var size, _ = cmd.Flags().GetInt32("size")
		var plugin, _ = cmd.Flags().GetStringSlice("plugin")
		var name, _ = cmd.Flags().GetString("name")
		var refID, _ = cmd.Flags().GetString("ref-id")
		var version, _ = cmd.Flags().GetString("version")
		var dt, _ = cmd.Flags().GetString("deployment-template")
		var te, _ = cmd.Flags().GetStringArray("topology-element")
		var region string
		if ecctl.Get().Config.Region == "" {
			region = cmdutil.DefaultECERegion
		}

		var p *models.ElasticsearchPayload
		if err := cmdutil.FileOrStdin(cmd, "file"); err == nil {
			err := cmdutil.DecodeDefinition(cmd, "file", &p)
			if err != nil && err != cmdutil.ErrNodefinitionLoaded {
				return err
			}
		}

		payload, err := depresource.ParseElasticsearchInput(depresource.ParseElasticsearchInputParams{
			NewElasticsearchParams: depresource.NewElasticsearchParams{
				API:        ecctl.Get().API,
				RefID:      refID,
				Version:    version,
				Plugins:    plugin,
				Region:     region,
				TemplateID: dt,
			},
			Size:             size,
			ZoneCount:        zoneCount,
			Payload:          p,
			Writer:           ecctl.Get().Config.ErrorDevice,
			TopologyElements: te,
		})
		if err != nil {
			return err
		}

		// Returns the ElasticsearchPayload skipping the creation of the resources.
		if generatePayload {
			return ecctl.Get().Formatter.Format("", payload)
		}

		var createParams = deployment.CreateParams{
			API: ecctl.Get().API,
			Request: &models.DeploymentCreateRequest{
				Name: name,
				Resources: &models.DeploymentCreateResources{
					Elasticsearch: []*models.ElasticsearchPayload{payload},
				},
			},
		}

		res, err := deployment.Create(createParams)
		if err != nil {
			return err
		}

		return cmdutil.Track(cmdutil.TrackParams{
			TrackResourcesParams: depresource.TrackResourcesParams{
				API:          ecctl.Get().API,
				Resources:    res.Resources,
				OutputDevice: ecctl.Get().Config.OutputDevice,
			},
			Formatter: ecctl.Get().Formatter,
			Track:     track,
			Response:  res,
		})
	},
}

func init() {
	Command.AddCommand(createElasticsearchCmd)
	createElasticsearchCmd.Flags().String("file", "", "ElasticsearchPayload file definition. See help for more information")
	createElasticsearchCmd.Flags().String("deployment-template", "default", "Deployment template ID on which to base the deployment from")
	createElasticsearchCmd.Flags().StringArrayP("topology-element", "e", nil, "Topology element definition. See help for more information")
	createElasticsearchCmd.Flags().String("version", "", "Version to use, if not specified, the latest available stack version will be used")
	createElasticsearchCmd.Flags().String("name", "", "Optional name for the Elasticsearch deployment")
	createElasticsearchCmd.Flags().String("ref-id", "elasticsearch", "RefId for the Elasticsearch deployment")
	createElasticsearchCmd.Flags().Int32("zones", 1, "Number of zones the deployment will span")
	createElasticsearchCmd.Flags().Int32("size", 4096, "Memory (RAM) in MB that each of the deployment nodes will have")
	createElasticsearchCmd.Flags().BoolP("track", "t", false, cmdutil.TrackFlagMessage)
	createElasticsearchCmd.Flags().Bool("generate-payload", false, "Returns the ElasticsearchPayload without actually creating the deployment resources")
	createElasticsearchCmd.Flags().StringSlice("plugin", nil, "Additional plugins to add to the Elasticsearch deployment")
}
