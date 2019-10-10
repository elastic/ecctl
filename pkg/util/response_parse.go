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

package util

import (
	"errors"
	"reflect"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/plan"
	multierror "github.com/hashicorp/go-multierror"
)

// ParseCUResponseParams is used by ParseCUResponse.
type ParseCUResponseParams struct {
	*api.API
	CreateResponse interface{}
	UpdateResponse interface{}
	Err            error
	TrackParams
}

// Validate ensures that the parameters are usable by the consuming function.
func (params ParseCUResponseParams) Validate() error {
	if params.Err != nil {
		return api.UnwrapError(params.Err)
	}

	var merr = new(multierror.Error)
	if params.API == nil {
		merr = multierror.Append(merr, errors.New("parse response: API cannot be empty"))
	}

	if params.CreateResponse == nil && params.UpdateResponse == nil {
		merr = multierror.Append(merr,
			errors.New("parse response: One of Create or Update response must be populated"),
		)
	}

	return multierror.Append(merr, params.TrackParams.Validate()).ErrorOrNil()
}

// ParseCUResponse parses the create / update response
func ParseCUResponse(params ParseCUResponseParams) (*models.ClusterCrudResponse, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	var res = params.CreateResponse
	if res == nil || reflect.ValueOf(res).IsNil() {
		res = params.UpdateResponse
	}

	field := reflect.ValueOf(res).Elem().FieldByName("Payload")
	if field.IsNil() || !field.IsValid() {
		return nil, errors.New("failed to obtain Payload field from Response")
	}

	response, ok := field.Interface().(*models.ClusterCrudResponse)
	if !ok {
		return nil, errors.New("failed casting Payload to *models.ClusterCrudResponse")
	}

	if !params.Track {
		return response, nil
	}

	var id, kind = response.ElasticsearchClusterID, "elasticsearch"
	if id == "" {
		id = response.KibanaClusterID
		kind = "kibana"
	}

	return response, TrackCluster(TrackClusterParams{
		Output: params.Output,
		TrackParams: plan.TrackParams{
			API:           params.API,
			PollFrequency: params.PollFrequency,
			MaxRetries:    params.MaxRetries,
			ID:            id,
			Kind:          kind,
		},
	})
}
