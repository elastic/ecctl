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

import (
	"fmt"
)

// GetResourceParams is consumed by GetResource.
type GetResourceParams struct {
	GetParams

	Kind string
}

// GetResource is a high level function which either returns the top level
// deployment information when no params.Kind is specified, or it returns a
// specific deployment resource information by RefID. If no RefID is defined,
// It will perform an additional API call to obtain the top level
func GetResource(params GetResourceParams) (interface{}, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	var noRefIDAndKind = params.GetParams.RefID == "" && params.Kind != ""
	if noRefIDAndKind {
		refID, err := GetKindRefID(params)
		if err != nil {
			return nil, err
		}
		params.GetParams.RefID = refID
	}

	switch params.Kind {
	case "apm":
		return GetApm(params.GetParams)
	case "kibana":
		return GetKibana(params.GetParams)
	case "elasticsearch":
		return GetElasticsearch(params.GetParams)
	case "appsearch":
		return GetAppSearch(params.GetParams)
	default:
		// If the is specified but not supported, return an error.
		if params.Kind != "" {
			return nil, fmt.Errorf(
				"deployment get: resource kind %s is not valid", params.Kind,
			)
		}
		return Get(params.GetParams)
	}
}

// GetKindRefID obtains a resource kind RefID. If the kind is not supported
// an error is returned.
func GetKindRefID(params GetResourceParams) (string, error) {
	res, err := Get(params.GetParams)
	if err != nil {
		return "", err
	}

	var refID string
	switch params.Kind {
	case "apm":
		for _, resource := range res.Resources.Apm {
			refID = *resource.RefID
		}
	case "kibana":
		for _, resource := range res.Resources.Kibana {
			refID = *resource.RefID
		}
	case "elasticsearch":
		for _, resource := range res.Resources.Elasticsearch {
			refID = *resource.RefID
		}
	case "appsearch":
		for _, resource := range res.Resources.Appsearch {
			refID = *resource.RefID
		}
	}

	if refID == "" {
		return "", fmt.Errorf("deployment get: resource kind %s is not available", params.Kind)
	}

	return refID, nil
}
