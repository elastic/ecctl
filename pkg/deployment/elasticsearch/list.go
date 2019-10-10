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

package elasticsearch

import (
	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_elasticsearch"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment/deputil"
	"github.com/elastic/ecctl/pkg/util"
)

const (
	defaultSearchSize = 100

	// defaultSystemAlertsNumber controls the number of system alerts that are
	// shown in the elasticsearch cluster list by default.
	defaultSystemAlertsNumber = int64(5)
)

// ListParams is used to obtain a number of clusters
type ListParams struct {
	*api.API

	// Version post processes the list results filtering out any versions which
	// don't match the specified one. Using this parameter makes the resulting
	// list equal or smaller to the specified size.
	Version string

	// SystemAlerts is number of system alerts to include in the response.
	// For example, the number of forced restarts caused from a limited amount
	// of memory. Only numbers greater than zero return a field. If 0, defaults
	// to 5.
	SystemAlerts int64

	// Optional parameters
	deputil.QueryParams
}

// Validate ensures the parameters are usable.
func (params ListParams) Validate() error {
	var err = new(multierror.Error)
	if params.API == nil {
		err = multierror.Append(err, util.ErrAPIReq)
	}

	return err.ErrorOrNil()
}

func (params ListParams) fillValues() {
	if params.Size <= 0 {
		params.Size = defaultSearchSize
	}

	if params.SystemAlerts <= 0 {
		params.SystemAlerts = defaultSystemAlertsNumber
	}
}

// List lists all the elasticsearch clusters matching the filters
func List(params ListParams) (*models.ElasticsearchClustersInfo, error) {
	params.fillValues()
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.V1API.ClustersElasticsearch.GetEsClusters(
		clusters_elasticsearch.NewGetEsClustersParams().
			WithSize(ec.Int64(params.Size)).
			WithShowPlans(ec.Bool(true)).
			WithQ(ec.String(params.Query)).
			WithShowSystemAlerts(ec.Int64(params.SystemAlerts)).
			WithTimeout(util.GetTimeoutFromSize(params.Size)).
			WithShowMetadata(ec.Bool(params.ShowMetadata)),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return filterVersion(res.Payload, params.Version), nil
}

// filterVersion filters clusters based on the given version
func filterVersion(clusters *models.ElasticsearchClustersInfo, version string) *models.ElasticsearchClustersInfo {
	if version == "" {
		return clusters
	}

	result := new(models.ElasticsearchClustersInfo)
	for _, c := range clusters.ElasticsearchClusters {
		if c.PlanInfo == nil || c.PlanInfo.Current == nil || c.PlanInfo.Current.Plan == nil {
			continue
		}
		if c.PlanInfo.Current.Plan.Elasticsearch.Version == version {
			result.ElasticsearchClusters = append(result.ElasticsearchClusters, c)
		}
	}
	count := int32(len(result.ElasticsearchClusters))
	result.ReturnCount = &count
	return result
}
