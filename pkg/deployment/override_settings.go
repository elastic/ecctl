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

package deployment

import "github.com/elastic/cloud-sdk-go/pkg/models"

// PayloadOverrides represent the override settings to
type PayloadOverrides struct {
	// If set, it will override the deployment name.
	Name string

	// If set, it will override the region when not present in the
	// DeploymentCreateRequest.
	// Note this behaviour is different from the rest of overrides
	// since this field tends to be populated by the global region
	// field which is implicit (by config) rather than explicit by flag.
	Region string

	// If set, it'll override all versions to match this one.
	Version string
}

// setOverrides sets a series of overrides to either the Create or Update
// deployment request. See PayloadOverrides to understand which how each
// field of that struct affects the behavior of this function.
func setOverrides(req interface{}, overrides *PayloadOverrides) {
	if req == nil || overrides == nil {
		return
	}

	var apm []*models.ApmPayload
	var appsearch []*models.AppSearchPayload
	var elasticsearch []*models.ElasticsearchPayload
	var kibana []*models.KibanaPayload
	switch t := req.(type) {
	case *models.DeploymentUpdateRequest:
		if t.Resources == nil {
			return
		}
		apm, appsearch = t.Resources.Apm, t.Resources.Appsearch
		elasticsearch, kibana = t.Resources.Elasticsearch, t.Resources.Kibana
	case *models.DeploymentCreateRequest:
		if overrides.Name != "" {
			t.Name = overrides.Name
		}
		if t.Resources == nil {
			return
		}
		apm, appsearch = t.Resources.Apm, t.Resources.Appsearch
		elasticsearch, kibana = t.Resources.Elasticsearch, t.Resources.Kibana
	}

	overrideByPayload(
		apm, appsearch, elasticsearch, kibana,
		overrides.Region, overrides.Version,
	)
}

// nolint
func overrideByPayload(apm []*models.ApmPayload, appsearch []*models.AppSearchPayload,
	elasticsearch []*models.ElasticsearchPayload, kibana []*models.KibanaPayload, region, version string) {
	for _, resource := range apm {
		if resource.Region == nil && region != "" {
			resource.Region = &region
		}

		if version != "" {
			if resource.Plan != nil && resource.Plan.Apm != nil {
				resource.Plan.Apm.Version = version
			}
		}
	}

	for _, resource := range appsearch {
		if resource.Region == nil && region != "" {
			resource.Region = &region
		}

		if version != "" {
			if resource.Plan != nil && resource.Plan.Appsearch != nil {
				resource.Plan.Appsearch.Version = version
			}
		}
	}

	for _, resource := range elasticsearch {
		if resource.Region == nil && region != "" {
			resource.Region = &region
		}

		if version != "" {
			if resource.Plan != nil && resource.Plan.Elasticsearch != nil {
				resource.Plan.Elasticsearch.Version = version
			}
		}
	}

	for _, resource := range kibana {
		if resource.Region == nil && region != "" {
			resource.Region = &region
		}

		if version != "" {
			if resource.Plan != nil && resource.Plan.Kibana != nil {
				resource.Plan.Kibana.Version = version
			}
		}
	}
}
