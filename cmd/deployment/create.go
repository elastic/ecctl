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

package cmddeployment

import (
	"fmt"
	"os"

	"github.com/elastic/cloud-sdk-go/pkg/models"
	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment"
	"github.com/elastic/ecctl/pkg/deployment/depresource"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/util"
)

var createCmd = &cobra.Command{
	Use:     "create {--file | --size <int> --zones <string> | --topology-element <obj>}",
	Short:   "Creates a deployment",
	PreRunE: cobra.MaximumNArgs(0),
	Long:    createLong,
	Example: createExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		var track, _ = cmd.Flags().GetBool("track")
		var generatePayload, _ = cmd.Flags().GetBool("generate-payload")
		var name, _ = cmd.Flags().GetString("name")
		var version, _ = cmd.Flags().GetString("version")
		var dt, _ = cmd.Flags().GetString("deployment-template")

		var zoneCount, _ = cmd.Flags().GetInt32("zones")
		var size, _ = cmd.Flags().GetInt32("size")
		var esRefID, _ = cmd.Flags().GetString("ref-id")
		var te, _ = cmd.Flags().GetStringArray("topology-element")
		var plugin, _ = cmd.Flags().GetStringSlice("plugin")

		var kibanaZoneCount, _ = cmd.Flags().GetInt32("kibana-zones")
		var kibanaSize, _ = cmd.Flags().GetInt32("kibana-size")
		var kibanaRefID, _ = cmd.Flags().GetString("kibana-ref-id")

		var apmEnable, _ = cmd.Flags().GetBool("apm")
		var apmZoneCount, _ = cmd.Flags().GetInt32("apm-zones")
		var apmSize, _ = cmd.Flags().GetInt32("apm-size")
		var apmRefID, _ = cmd.Flags().GetString("apm-ref-id")

		var appsearchEnable, _ = cmd.Flags().GetBool("appsearch")
		var appsearchZoneCount, _ = cmd.Flags().GetInt32("appsearch-zones")
		var appsearchSize, _ = cmd.Flags().GetInt32("appsearch-size")
		var appsearchRefID, _ = cmd.Flags().GetString("appsearch-ref-id")

		region := ecctl.Get().Config.Region
		if region == "" {
			region = cmdutil.DefaultECERegion
		}

		var payload *models.DeploymentCreateRequest
		if err := sdkcmdutil.FileOrStdin(cmd, "file"); err == nil {
			err := sdkcmdutil.DecodeDefinition(cmd, "file", &payload)
			if err != nil && err != sdkcmdutil.ErrNodefinitionLoaded {
				return err
			}
		}

		if payload == nil {
			var err error
			payload, err = depresource.New(depresource.NewParams{
				API:                  ecctl.Get().API,
				Name:                 name,
				DeploymentTemplateID: dt,
				Version:              version,
				Region:               region,
				Writer:               ecctl.Get().Config.ErrorDevice,
				Plugins:              plugin,
				TopologyElements:     te,
				ApmEnable:            apmEnable,
				AppsearchEnable:      appsearchEnable,
				ElasticsearchInstance: depresource.InstanceParams{
					RefID:     esRefID,
					Size:      size,
					ZoneCount: zoneCount,
				},
				KibanaInstance: depresource.InstanceParams{
					RefID:     kibanaRefID,
					Size:      kibanaSize,
					ZoneCount: kibanaZoneCount,
				},
				ApmInstance: depresource.InstanceParams{
					RefID:     apmRefID,
					Size:      apmSize,
					ZoneCount: apmZoneCount,
				},
				AppsearchInstance: depresource.InstanceParams{
					RefID:     appsearchRefID,
					Size:      appsearchSize,
					ZoneCount: appsearchZoneCount,
				},
			})
			if err != nil {
				return err
			}
		}

		// Returns the DeploymentCreateRequest skipping the creation of the resources.
		if generatePayload {
			return ecctl.Get().Formatter.Format("", payload)
		}

		reqID, _ := cmd.Flags().GetString("request-id")
		if reqID == "" {
			reqID = util.RandomString(64)
		}

		var createParams = deployment.CreateParams{
			API:       ecctl.Get().API,
			RequestID: reqID,
			Request:   payload,
		}

		res, err := deployment.Create(createParams)
		if err != nil {
			fmt.Fprintln(os.Stderr,
				"The deployment creation returned with an error, please use the displayed idempotency token to recreate the deployment resources",
			)
			fmt.Fprintln(os.Stderr, "Idempotency token:", reqID)
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
			OutputDevice: ecctl.Get().Config.OutputDevice,
		})
	},
}

func init() {
	Command.AddCommand(createCmd)
	createCmd.Flags().StringP("file", "f", "", "DeploymentCreateRequest file definition. See help for more information")
	createCmd.Flags().String("deployment-template", "default", "Deployment template ID on which to base the deployment from")
	createCmd.Flags().String("version", "", "Version to use, if not specified, the latest available stack version will be used")
	createCmd.Flags().String("name", "", "Optional name for the deployment")
	createCmd.Flags().BoolP("track", "t", false, cmdutil.TrackFlagMessage)
	createCmd.Flags().Bool("generate-payload", false, "Returns the deployment payload without actually creating the deployment resources")
	createCmd.Flags().String("request-id", "", "Optional idempotency token - Can be found in the Stderr device when a previous deployment creation failed, for more information see the examples in the help command page")

	createCmd.Flags().String("ref-id", "main-elasticsearch", "Optional RefId for the Elasticsearch deployment")
	createCmd.Flags().Int32("zones", 1, "Number of zones the Elasticsearch instances will span")
	createCmd.Flags().Int32("size", 4096, "Memory (RAM) in MB that each of the Elasticsearch instances will have")
	createCmd.Flags().StringArrayP("topology-element", "e", nil, "Optional Elasticsearch topology element definition. See help for more information")
	createCmd.Flags().StringSlice("plugin", nil, "Additional plugins to add to the Elasticsearch deployment")

	createCmd.Flags().String("kibana-ref-id", "main-kibana", "Optional RefId for the Kibana deployment")
	createCmd.Flags().Int32("kibana-zones", 1, "Number of zones the Kibana instances will span")
	createCmd.Flags().Int32("kibana-size", 1024, "Memory (RAM) in MB that each of the Kibana instances will have")

	createCmd.Flags().Bool("apm", false, "Enables APM for the deployment")
	createCmd.Flags().String("apm-ref-id", "main-apm", "Optional RefId for the APM deployment")
	createCmd.Flags().Int32("apm-zones", 1, "Number of zones the APM instances will span")
	createCmd.Flags().Int32("apm-size", 512, "Memory (RAM) in MB that each of the APM instances will have")

	createCmd.Flags().Bool("appsearch", false, "Enables AppSearch for the deployment")
	createCmd.Flags().String("appsearch-ref-id", "main-appsearch", "Optional RefId for the AppSearch deployment")
	createCmd.Flags().Int32("appsearch-zones", 1, "Number of zones the AppSearch instances will span")
	createCmd.Flags().Int32("appsearch-size", 2048, "Memory (RAM) in MB that each of the AppSearch instances will have")
}
