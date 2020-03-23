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
	"io"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/models"
)

// NewParams is consumed by New()
type NewParams struct {
	*api.API

	Name                   string
	Version                string
	DeploymentTemplateID   string
	Region                 string
	ElasticsearchRefID     string
	KibanaRefID            string
	ApmRefID               string
	AppsearchRefID         string
	ElasticsearchSize      int32
	ElasticsearchZoneCount int32
	KibanaSize             int32
	KibanaZoneCount        int32
	ApmSize                int32
	ApmZoneCount           int32
	AppsearchSize          int32
	AppsearchZoneCount     int32
	ApmEnable              bool
	AppsearchEnable        bool
	Writer                 io.Writer
	Plugins                []string
	TopologyElements       []string
}

// New creates the payload for a deployment
func New(params NewParams) (*models.DeploymentCreateRequest, error) {
	esPayload, err := ParseElasticsearchInput(ParseElasticsearchInputParams{
		NewElasticsearchParams: NewElasticsearchParams{
			API:        params.API,
			RefID:      params.ElasticsearchRefID,
			Version:    params.Version,
			Plugins:    params.Plugins,
			Region:     params.Region,
			TemplateID: params.DeploymentTemplateID,
		},
		Size:             params.ElasticsearchSize,
		ZoneCount:        params.ElasticsearchZoneCount,
		Writer:           params.Writer,
		TopologyElements: params.TopologyElements,
	})
	if err != nil {
		return nil, err
	}

	kibanaPayload, err := NewKibana(NewStateless{
		ElasticsearchRefID: params.ElasticsearchRefID,
		API:                params.API,
		RefID:              params.KibanaRefID,
		Version:            params.Version,
		Region:             params.Region,
		TemplateID:         params.DeploymentTemplateID,
		Size:               params.KibanaSize,
		ZoneCount:          params.KibanaZoneCount,
	})
	if err != nil {
		return nil, err
	}

	resources := models.DeploymentCreateResources{
		Elasticsearch: []*models.ElasticsearchPayload{esPayload},
		Kibana:        []*models.KibanaPayload{kibanaPayload},
	}

	if params.ApmEnable {
		apmPayload, err := NewApm(NewStateless{
			ElasticsearchRefID: params.ElasticsearchRefID,
			API:                params.API,
			RefID:              params.ApmRefID,
			Version:            params.Version,
			Region:             params.Region,
			TemplateID:         params.DeploymentTemplateID,
			Size:               params.ApmSize,
			ZoneCount:          params.ApmZoneCount,
		})
		if err != nil {
			return nil, err
		}

		resources.Apm = []*models.ApmPayload{apmPayload}
	}

	if params.AppsearchEnable {
		appsearchPayload, err := NewAppSearch(NewStateless{
			ElasticsearchRefID: params.ElasticsearchRefID,
			API:                params.API,
			RefID:              params.AppsearchRefID,
			Version:            params.Version,
			Region:             params.Region,
			TemplateID:         params.DeploymentTemplateID,
			Size:               params.AppsearchSize,
			ZoneCount:          params.AppsearchZoneCount,
		})
		if err != nil {
			return nil, err
		}

		resources.Appsearch = []*models.AppSearchPayload{appsearchPayload}
	}

	payload := &models.DeploymentCreateRequest{
		Name:      params.Name,
		Resources: &resources,
	}

	return payload, nil
}
