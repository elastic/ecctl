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
	"errors"

	"github.com/elastic/cloud-sdk-go/pkg/client/deployments"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/util"
)

// AddParams is consumed by Add.
type AddParams struct {
	deployment.Params
	Message     string
	UserID      string
	Commentator ecctl.Commentator
}

// Validate ensures the parameters are valid
func (params AddParams) Validate() error {
	var merr = new(multierror.Error)

	if params.UserID == "" {
		merr = multierror.Append(merr, errors.New(errEmptyUserID))
	}

	if params.Message == "" {
		merr = multierror.Append(merr, errors.New(errEmptyNoteMessage))
	}

	merr = multierror.Append(merr, params.Params.Validate())

	return merr.ErrorOrNil()
}

// TODO: Use different resource types when this is supported by the API.
// For the time being, the notes endpoint only allows elasticsearch IDs.
func (params *AddParams) fillDefaults() error {
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

// Add posts a new message to the specified deployment
func Add(params AddParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	if err := params.fillDefaults(); err != nil {
		return err
	}

	var message = params.Message
	if params.Commentator != nil {
		message = params.Commentator.Message(message)
	}

	return util.ReturnErrOnly(params.V1API.Deployments.CreateDeploymentNote(
		deployments.NewCreateDeploymentNoteParams().
			WithDeploymentID(params.ID).
			WithBody(&models.Note{
				Message: ec.String(message),
				UserID:  params.UserID,
			}),
		params.AuthWriter,
	))
}
