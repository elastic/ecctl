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
	"sort"

	"github.com/blang/semver/v4"
	"github.com/elastic/cloud-sdk-go/pkg/api/deploymentapi"
	"github.com/elastic/cloud-sdk-go/pkg/api/deploymentapi/deptemplateapi"
	"github.com/elastic/cloud-sdk-go/pkg/api/stackapi"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/multierror"
	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var createCmd = &cobra.Command{
	Use:     "create {--file | --es-size <int> --es-zones <int> | --es-node-topology <obj>}",
	Short:   "Creates a deployment",
	PreRunE: cobra.NoArgs,
	Long:    createLong,
	Example: createExample,
	RunE: func(cmd *cobra.Command, args []string) error {
		track, _ := cmd.Flags().GetBool("track")
		generatePayload, _ := cmd.Flags().GetBool("generate-payload")
		name, _ := cmd.Flags().GetString("name")
		version, _ := cmd.Flags().GetString("version")
		file, _ := cmd.Flags().GetString("file")
		region := ecctl.Get().Config.Region

		if version == "" && file == "" {
			var err error
			version, err = getLatestStackVersion(region)
			if err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "--version not set, using latest: %s\n", version)
		}

		payload, err := newCreatePayload(cmd, version, region)
		if err != nil {
			return err
		}

		if err := deploymentapi.OverrideCreateOrUpdateRequest(payload,
			&deploymentapi.PayloadOverrides{
				Name:               name,
				Version:            version,
				Region:             region,
				ElasticsearchRefID: "main-elasticsearch",
				OverrideRefIDs:     true,
			},
		); err != nil {
			return err
		}

		// Returns the DeploymentCreateRequest skipping the creation of the resources.
		if generatePayload {
			return ecctl.Get().Formatter.Format("", payload)
		}

		reqID, _ := cmd.Flags().GetString("request-id")
		reqID = deploymentapi.RequestID(reqID)
		createParams := deploymentapi.CreateParams{
			API:       ecctl.Get().API,
			RequestID: reqID,
			Request:   payload,
		}

		res, err := deploymentapi.Create(createParams)
		if err != nil {
			_, _ = fmt.Fprintln(cmd.ErrOrStderr(),
				"The deployment creation returned with an error. Use the displayed request ID to recreate the deployment resources",
			)
			_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "Request ID:", reqID)
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

func init() {
	initFlags()
}

func initFlags() {
	Command.AddCommand(createCmd)
	createCmd.Flags().StringP("file", "f", "", "DeploymentCreateRequest file definition. See help for more information")
	createCmd.Flags().String("deployment-template", "", "Deployment template ID on which to base the deployment from")
	createCmd.Flags().String("version", "", "Version to use, if not specified, the latest available stack version will be used")
	createCmd.Flags().String("name", "", "Optional name for the deployment")
	createCmd.Flags().BoolP("track", "t", false, cmdutil.TrackFlagMessage)
	createCmd.Flags().Bool("generate-payload", false, "Returns the deployment payload without actually creating the deployment resources")
	createCmd.Flags().String("request-id", "", "Optional request ID - Can be found in the Stderr device when a previous deployment creation failed. For more information see the examples in the help command page")
	createCmd.Flags().Bool("minimum-size", false, "Shrink each Elasticsearch topology element to its minimum allowed size")
}

func getLatestStackVersion(region string) (string, error) {
	stacks, err := stackapi.List(stackapi.ListParams{
		API:    ecctl.Get().API,
		Region: region,
	})
	if err != nil {
		return "", fmt.Errorf("unable to fetch stack versions: %w", err)
	}
	if len(stacks.Stacks) == 0 {
		return "", errors.New("no stack versions available")
	}
	// List returns stacks sorted descending, so first is latest.
	return stacks.Stacks[0].Version, nil
}

func getDefaultTemplate(region, stackVersion string) (string, error) {
	showHidden := false
	templates, err := deptemplateapi.List(deptemplateapi.ListParams{
		API:                        ecctl.Get().API,
		Region:                     region,
		StackVersion:               stackVersion,
		ShowHidden:                 showHidden,
		HideInstanceConfigurations: true,
	})
	if err != nil {
		return "", err
	}

	var current []*models.DeploymentTemplateInfoV2
	for _, t := range templates {
		if !isLegacyTemplate(t) {
			current = append(current, t)
		}
	}
	if len(current) == 0 {
		return "", errors.New("no non-legacy deployment templates available for region")
	}

	sort.Slice(current, func(i, j int) bool {
		if current[i].Order == nil {
			return false
		}
		if current[j].Order == nil {
			return true
		}
		return *current[i].Order < *current[j].Order
	})

	return *current[0].ID, nil
}

func isLegacyTemplate(t *models.DeploymentTemplateInfoV2) bool {
	for _, m := range t.Metadata {
		if m.Key != nil && *m.Key == "legacy" && m.Value != nil && *m.Value == "true" {
			return true
		}
	}
	return false
}

func newCreatePayload(cmd *cobra.Command, version, region string) (*models.DeploymentCreateRequest, error) {
	file, _ := cmd.Flags().GetString("file")
	dt, _ := cmd.Flags().GetString("deployment-template")
	var payload models.DeploymentCreateRequest
	if file != "" {
		if err := sdkcmdutil.DecodeDefinition(cmd, "file", &payload); err != nil {
			merr := multierror.NewPrefixed("failed reading the file definition")
			return nil, merr.Append(err,
				errors.New("could not read the specified file, please make sure it exists"),
			)
		}
		return &payload, nil
	}

	if dt == "" {
		var err error
		dt, err = getDefaultTemplate(region, version)
		if err != nil {
			return nil, err
		}
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "--deployment-template not set, using %s\n", dt)
	}
	tpl, err := deptemplateapi.Get(deptemplateapi.GetParams{
		API:          ecctl.Get().API,
		TemplateID:   dt,
		Region:       region,
		StackVersion: version,
	})
	if err != nil {
		return nil, err
	}

	if es := tpl.DeploymentTemplate.Resources.Elasticsearch; len(es) > 0 {
		if es[0].Plan.DeploymentTemplate == nil {
			es[0].Plan.DeploymentTemplate = &models.DeploymentTemplateReference{}
		}

		es[0].Plan.DeploymentTemplate.ID = tpl.ID
	}

	result, err := removeUnsupportedResources(version, tpl.DeploymentTemplate)
	if err != nil {
		return nil, err
	}
	if minimumSize, _ := cmd.Flags().GetBool("minimum-size"); minimumSize {
		shrinkTopologySizesToMinimum(result)
	}
	return result, nil
}

func removeUnsupportedResources(version string, tpl *models.DeploymentCreateRequest) (*models.DeploymentCreateRequest, error) {
	vers, err := semver.Parse(version)
	if err != nil {
		return nil, fmt.Errorf("failed to parse version: %v", err)
	}
	if vers.Major >= 8 {
		tpl.Resources.Apm = nil
	}
	return tpl, nil
}

// shrinkTopologySizesToMinimum reduces the memory size and zone count of each
// Elasticsearch topology element to its minimum allowed values. The memory size
// is reduced to TopologyElementControl.Min, and the zone count is set to 1.
// Elements with a zero minimum (optional/disabled tiers) are left unchanged.
// Other resource types (Kibana, APM, Enterprise Search) are left untouched
// because they carry no minimum size constraint on the topology element itself.
func shrinkTopologySizesToMinimum(tpl *models.DeploymentCreateRequest) {
	if tpl.Resources == nil {
		return
	}
	for _, es := range tpl.Resources.Elasticsearch {
		if es.Plan == nil {
			continue
		}
		for _, node := range es.Plan.ClusterTopology {
			if node.Size == nil || node.Size.Value == nil {
				continue
			}
			if node.TopologyElementControl == nil || node.TopologyElementControl.Min == nil || node.TopologyElementControl.Min.Value == nil {
				continue
			}
			min := *node.TopologyElementControl.Min.Value
			if min > 0 {
				if *node.Size.Value > min {
					node.Size.Value = node.TopologyElementControl.Min.Value
				}
				node.ZoneCount = 1
			}
		}
	}
}
