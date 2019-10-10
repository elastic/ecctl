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

package kibana

import (
	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_kibana"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment/deputil"
	"github.com/elastic/ecctl/pkg/util"
)

// ClusterParams contains parameters used to fetch cluster's data
type ClusterParams struct {
	DeploymentParams

	// Optional parameters
	deputil.QueryParams
}

// ListParams contains parameters used to fetch cluster's data
type ListParams struct {
	*api.API
	Version string

	// Optional parameters
	deputil.QueryParams
}

// Validate is the implementation for the ecctl.Validator interface
func (params ListParams) Validate() error {
	var merr = new(multierror.Error)

	if params.API == nil {
		merr = multierror.Append(merr, util.ErrAPIReq)
	}

	return merr.ErrorOrNil()
}

// Get returns the kibana cluster
func Get(params ClusterParams) (*models.KibanaClusterInfo, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.DeploymentParams.API.V1API.ClustersKibana.GetKibanaCluster(
		clusters_kibana.NewGetKibanaClusterParams().
			WithClusterID(params.ID).
			WithShowMetadata(&params.ShowMetadata).
			WithShowPlanDefaults(&params.ShowPlanDefaults).
			WithShowPlans(&params.ShowPlans).
			WithShowPlanLogs(&params.ShowPlanLogs).
			WithShowSettings(&params.ShowSettings),
		params.API.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}

// List lists all the kibana clusters matching the filters
func List(params ListParams) (*models.KibanaClustersInfo, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.API.V1API.ClustersKibana.GetKibanaClusters(
		clusters_kibana.NewGetKibanaClustersParams().WithSize(&params.Size).
			WithShowPlans(ec.Bool(true)).
			WithTimeout(util.GetTimeoutFromSize(params.Size)).
			WithShowMetadata(ec.Bool(params.ShowMetadata)),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return filterByVersion(res.Payload, params.Version), nil
}

func filterByVersion(p *models.KibanaClustersInfo, v string) *models.KibanaClustersInfo {
	if v == "" {
		return p
	}

	result := models.KibanaClustersInfo{}

	for _, c := range p.KibanaClusters {
		if c.PlanInfo == nil || c.PlanInfo.Current == nil || c.PlanInfo.Current.Plan == nil {
			continue
		}

		if c.PlanInfo.Current.Plan.Kibana.Version == v {
			result.KibanaClusters = append(result.KibanaClusters, c)
		}
	}

	count := int32(len(result.KibanaClusters))
	result.ReturnCount = &count

	return &result
}
