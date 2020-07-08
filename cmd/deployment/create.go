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
	"errors"
	"fmt"

	"github.com/elastic/cloud-sdk-go/pkg/api/deploymentapi"
	"github.com/elastic/cloud-sdk-go/pkg/api/deploymentapi/depresourceapi"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/multierror"
	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var createCmd = &cobra.Command{
	Use:     "create {--file | --size <int> --zones <string> | --topology-element <obj>}",
	Short:   "Creates a deployment",
	PreRunE: cobra.NoArgs,
	// Switch back to non-temp constants when reads for deployment templates are available on ESS
	Long:    createLongTemp,
	Example: createExampleTemp,
	RunE: func(cmd *cobra.Command, args []string) error {
		var track, _ = cmd.Flags().GetBool("track")
		var generatePayload, _ = cmd.Flags().GetBool("generate-payload")
		var name, _ = cmd.Flags().GetString("name")
		var version, _ = cmd.Flags().GetString("version")
		var dt, _ = cmd.Flags().GetString("deployment-template")

		var esZoneCount, _ = cmd.Flags().GetInt32("es-zones")
		var esSize, _ = cmd.Flags().GetInt32("es-size")
		var esRefID, _ = cmd.Flags().GetString("es-ref-id")
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

		var skipFlagBased = cmd.Flag("deployment-template").Hidden

		var payload *models.DeploymentCreateRequest

		err := sdkcmdutil.DecodeDefinition(cmd, "file", &payload)
		if err := returnErrOnHidden(err, skipFlagBased); err != nil {
			merr := multierror.NewPrefixed("failed reading the file definition")
			return merr.Append(err,
				errors.New("could not read the specified file, please make sure it exists"),
			)
		}

		if payload == nil {
			var err error
			payload, err = depresourceapi.New(depresourceapi.NewParams{
				API:                  ecctl.Get().API,
				Name:                 name,
				DeploymentTemplateID: dt,
				Version:              version,
				Region:               ecctl.Get().Config.Region,
				Writer:               ecctl.Get().Config.ErrorDevice,
				Plugins:              plugin,
				TopologyElements:     te,
				ApmEnable:            apmEnable,
				AppsearchEnable:      appsearchEnable,
				ElasticsearchInstance: depresourceapi.InstanceParams{
					RefID:     esRefID,
					Size:      esSize,
					ZoneCount: esZoneCount,
				},
				KibanaInstance: depresourceapi.InstanceParams{
					RefID:     kibanaRefID,
					Size:      kibanaSize,
					ZoneCount: kibanaZoneCount,
				},
				ApmInstance: depresourceapi.InstanceParams{
					RefID:     apmRefID,
					Size:      apmSize,
					ZoneCount: apmZoneCount,
				},
				AppsearchInstance: depresourceapi.InstanceParams{
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
		reqID = deploymentapi.RequestID(reqID)

		var createParams = deploymentapi.CreateParams{
			API:       ecctl.Get().API,
			RequestID: reqID,
			Request:   payload,
		}

		res, err := deploymentapi.Create(createParams)
		if err != nil {
			fmt.Fprintln(cmd.ErrOrStderr(),
				"The deployment creation returned with an error. Use the displayed request ID to recreate the deployment resources",
			)
			fmt.Fprintln(cmd.ErrOrStderr(), "Request ID:", reqID)
			return err
		}

		return cmdutil.Track(cmdutil.NewTrackParams(cmdutil.TrackParamsConfig{
			App:          ecctl.Get(),
			DeploymentID: *res.ID,
			Track:        track,
			Response:     res,
		}))
	},
}

func returnErrOnHidden(err error, hidden bool) error {
	if hidden {
		return err
	}
	if err != nil && err != sdkcmdutil.ErrNodefinitionLoaded {
		return err
	}
	return nil
}

func init() {
	Command.AddCommand(createCmd)
	createCmd.Flags().StringP("file", "f", "", "DeploymentCreateRequest file definition. See help for more information")
	// Remove when reads for deployment templates are available on ESS
	createCmd.Flags().String("deployment-template", "default", "Deployment template ID on which to base the deployment from")
	createCmd.Flags().String("version", "", "Version to use, if not specified, the latest available stack version will be used")
	createCmd.Flags().String("name", "", "Optional name for the deployment")
	createCmd.Flags().BoolP("track", "t", false, cmdutil.TrackFlagMessage)
	createCmd.Flags().Bool("generate-payload", false, "Returns the deployment payload without actually creating the deployment resources")
	createCmd.Flags().String("request-id", "", "Optional request ID - Can be found in the Stderr device when a previous deployment creation failed. For more information see the examples in the help command page")

	createCmd.Flags().String("es-ref-id", "main-elasticsearch", "Optional RefId for the Elasticsearch deployment")
	createCmd.Flags().Int32("es-zones", 1, "Number of zones the Elasticsearch instances will span")
	createCmd.Flags().Int32("es-size", 4096, "Memory (RAM) in MB that each of the Elasticsearch instances will have")
	createCmd.Flags().StringArrayP("topology-element", "e", nil, "Optional Elasticsearch topology element definition. See help for more information")
	createCmd.Flags().StringSlice("plugin", nil, "Additional plugins to add to the Elasticsearch deployment")

	createCmd.Flags().String("kibana-ref-id", "main-kibana", "Optional RefId for the Kibana deployment")
	createCmd.Flags().Int32("kibana-zones", 1, "Number of zones the Kibana instances will span")
	createCmd.Flags().Int32("kibana-size", 1024, "Memory (RAM) in MB that each of the Kibana instances will have")

	createCmd.Flags().Bool("apm", false, "Enables APM for the deployment")
	createCmd.Flags().String("apm-ref-id", "main-apm", "Optional RefId for the APM deployment")
	createCmd.Flags().Int32("apm-zones", 1, "Number of zones the APM instances will span")
	createCmd.Flags().Int32("apm-size", 512, "Memory (RAM) in MB that each of the APM instances will have")

	createCmd.Flags().Bool("appsearch", false, "Enables App Search for the deployment")
	createCmd.Flags().String("appsearch-ref-id", "main-appsearch", "Optional RefId for the App Search deployment")
	createCmd.Flags().Int32("appsearch-zones", 1, "Number of zones the App Search instances will span")
	createCmd.Flags().Int32("appsearch-size", 2048, "Memory (RAM) in MB that each of the App Search instances will have")

	// The following flags will remain hidden until reads for deployment templates are available on ESS
	createCmd.Flags().MarkHidden("deployment-template")
	createCmd.Flags().MarkHidden("version")
	createCmd.Flags().MarkHidden("name")
	createCmd.Flags().MarkHidden("generate-payload")
	createCmd.Flags().MarkHidden("es-ref-id")
	createCmd.Flags().MarkHidden("es-zones")
	createCmd.Flags().MarkHidden("es-size")
	createCmd.Flags().MarkHidden("topology-element")
	createCmd.Flags().MarkHidden("plugin")
	createCmd.Flags().MarkHidden("kibana-ref-id")
	createCmd.Flags().MarkHidden("kibana-zones")
	createCmd.Flags().MarkHidden("kibana-size")
	createCmd.Flags().MarkHidden("apm")
	createCmd.Flags().MarkHidden("apm-ref-id")
	createCmd.Flags().MarkHidden("apm-zones")
	createCmd.Flags().MarkHidden("apm-size")
	createCmd.Flags().MarkHidden("appsearch")
	createCmd.Flags().MarkHidden("appsearch-ref-id")
	createCmd.Flags().MarkHidden("appsearch-zones")
	createCmd.Flags().MarkHidden("appsearch-size")
}
