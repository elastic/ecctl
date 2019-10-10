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
	"fmt"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_apm"
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_kibana"
	"github.com/elastic/cloud-sdk-go/pkg/client/deployments"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/elastic/cloud-sdk-go/pkg/util/slice"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/util"
)

// AddParams is consumed by Add.
type AddParams struct {
	*api.API

	ID          string
	Type        string
	Message     string
	UserID      string
	Commentator ecctl.Commentator
}

// Validate ensures the parameters are valid
func (params AddParams) Validate() error {
	var err = new(multierror.Error)
	if params.API == nil {
		err = multierror.Append(err, util.ErrAPIReq)
	}

	if len(params.ID) != 32 {
		err = multierror.Append(err, errors.New("invalid id"))
	}

	if params.UserID == "" {
		err = multierror.Append(err, errors.New("user id cannot be empty"))
	}

	if !slice.HasString(util.ValidTypes, params.Type) {
		err = multierror.Append(err,
			fmt.Errorf("invalid type %s: valid types are %v", params.Type, util.ValidTypes),
		)
	}

	if params.Message == "" {
		err = multierror.Append(err, errors.New("message cannot be empty"))
	}

	return err.ErrorOrNil()
}

// Add posts a new message to the specified deployment
func Add(params AddParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	var deploymentID string
	switch params.Type {
	case "elasticsearch":
		deploymentID = params.ID
	case "kibana":
		res, err := params.V1API.ClustersKibana.GetKibanaCluster(
			clusters_kibana.NewGetKibanaClusterParams().
				WithClusterID(params.ID),
			params.AuthWriter,
		)
		if err != nil {
			return api.UnwrapError(err)
		}
		deploymentID = *res.Payload.ElasticsearchCluster.ElasticsearchID
	case "apm":
		res, err := params.V1API.ClustersApm.GetApmCluster(
			clusters_apm.NewGetApmClusterParams().
				WithClusterID(params.ID),
			params.AuthWriter,
		)
		if err != nil {
			return api.UnwrapError(err)
		}
		deploymentID = *res.Payload.ElasticsearchCluster.ElasticsearchID
	}

	var message = params.Message
	if params.Commentator != nil {
		message = params.Commentator.Message(message)
	}

	return util.ReturnErrOnly(params.V1API.Deployments.CreateDeploymentNote(
		deployments.NewCreateDeploymentNoteParams().
			WithDeploymentID(deploymentID).
			WithBody(&models.Note{
				Message: ec.String(message),
				UserID:  params.UserID,
			}),
		params.AuthWriter,
	))
}

// Update updates a note from its deployment and note ID
func Update(params UpdateParams) (*models.Note, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.API.V1API.Deployments.UpdateDeploymentNote(
		deployments.NewUpdateDeploymentNoteParams().
			WithDeploymentID(params.ID).
			WithNoteID(params.NoteID).
			WithBody(&models.Note{
				Message: ec.String(params.Message),
				UserID:  params.UserID,
			}),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}
