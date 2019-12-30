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

package note

import (
	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/deployments"
	"github.com/elastic/cloud-sdk-go/pkg/models"

	"github.com/elastic/ecctl/pkg/deployment"
)

// ListParams is used by List
type ListParams struct {
	deployment.Params
}

// TODO: Use different resource types when this is supported by the API.
// For the time being, the notes endpoint only allows elasticsearch IDs.
func (params *ListParams) fillDefaults() error {
	esID, err := getElasticsearchID(deployment.GetParams{
		API:          params.API,
		DeploymentID: params.ID,
	})
	if err != nil {
		return err
	}

	params.ID = esID
	return err
}

// List lists all of the notes for the deployment
func List(params ListParams) (*models.Notes, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	if err := params.fillDefaults(); err != nil {
		return nil, err
	}

	res, err := params.API.V1API.Deployments.GetDeploymentNotes(
		deployments.NewGetDeploymentNotesParams().
			WithDeploymentID(params.ID),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}
