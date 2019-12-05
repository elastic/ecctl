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

package depresource

import (
	"errors"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/platform_configuration_templates"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

const (
	// DefaultKibanahRefID is used when the RefID is not specified.
	DefaultKibanahRefID = "kibana"
)

// NewKibanaParams is consumed by NewKibana.
type NewKibanaParams struct {
	*api.API

	// Required: DeploymentID.
	DeploymentID string

	// Required region name.
	Region string

	// Optional ElasticsearchRefID.
	ElasticsearchRefID string

	// Optional name. If not specified it defaults to the autogeneratd ID.
	Name string

	// Optional RefID for the kibana deployment resource.
	RefID string

	// Optional: Version is the Kibana Version. If not set it'll automatically
	// be set to the Elasticsearch RefID version.
	Version string

	// Kibana deployment size.
	Size      int32
	ZoneCount int32

	// Optional Deployment Template ID.
	TemplateID string
}

func (params *NewKibanaParams) fillDefaults() {
	if params.RefID == "" {
		params.RefID = DefaultKibanahRefID
	}
}

// Validate ensures the parameters are usable by the consuming function.
func (params *NewKibanaParams) Validate() error {
	var merr = new(multierror.Error)

	if params.API == nil {
		merr = multierror.Append(merr, util.ErrAPIReq)
	}

	if len(params.DeploymentID) != 32 {
		merr = multierror.Append(merr, util.ErrDeploymentID)
	}

	if params.Region == "" {
		merr = multierror.Append(merr, errors.New("deployment topology: region cannot be empty"))
	}

	return merr.ErrorOrNil()
}

// NewKibana creates a *models.KibanaPayload from the parameters.
// It relies on a simplified single dimension memory size and zone count to
// construct the Kibana's topology.
func NewKibana(params NewKibanaParams) (*models.KibanaPayload, error) {
	params.fillDefaults()
	if err := params.Validate(); err != nil {
		return nil, err
	}

	// When either not set, we obtain the values from the running deployment.
	// Overriding either or both is done at the end of the if.
	if err := getTemplateAndRefID(&params); err != nil {
		return nil, err
	}

	// Obtain the deployment template so we can create the kibana topology from
	// the specified sizes. The sizing overrides are done in newKibanaPayload.
	res, err := params.V1API.PlatformConfigurationTemplates.GetDeploymentTemplate(
		platform_configuration_templates.NewGetDeploymentTemplateParams().
			WithTemplateID(params.TemplateID).
			WithShowInstanceConfigurations(ec.Bool(true)),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	var clusterTopology = res.Payload.ClusterTemplate.Kibana.Plan.ClusterTopology
	var topology = models.KibanaClusterTopologyElement{Size: new(models.TopologySize)}
	if len(clusterTopology) > 0 {
		topology = *clusterTopology[0]
	}
	var payload = newKibanaPayload(params, topology)

	return &payload, nil
}

func newKibanaPayload(params NewKibanaParams, topology models.KibanaClusterTopologyElement) models.KibanaPayload {
	if params.Size > 0 {
		topology.Size.Value = ec.Int32(params.Size)
	}
	if params.ZoneCount > 0 {
		topology.ZoneCount = params.ZoneCount
	}

	return models.KibanaPayload{
		ElasticsearchClusterRefID: ec.String(params.ElasticsearchRefID),
		DisplayName:               params.Name,
		Region:                    ec.String(params.Region),
		RefID:                     ec.String(params.RefID),
		Plan: &models.KibanaClusterPlan{
			Kibana:          &models.KibanaConfiguration{Version: params.Version},
			ClusterTopology: []*models.KibanaClusterTopologyElement{&topology},
		},
	}
}

// When either of TemplateID or ElasticsearchRefID are not set, obtain the
// values from the running deployment. Overriding either or both when not set.
// Receives a pointer to the structure, meaning it modifies its values.
func getTemplateAndRefID(params *NewKibanaParams) error {
	if params.TemplateID == "" || params.ElasticsearchRefID == "" {
		res, err := GetDeploymentInfo(GetDeploymentInfoParams{
			API:          params.API,
			DeploymentID: params.DeploymentID,
		})
		if err != nil {
			return err
		}
		if params.TemplateID == "" {
			params.TemplateID = res.DeploymentTemplate
		}
		if params.ElasticsearchRefID == "" {
			params.ElasticsearchRefID = res.RefID
		}
		return nil
	}
	return nil
}
