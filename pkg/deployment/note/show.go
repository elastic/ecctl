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

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/deployments_notes"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	multierror "github.com/hashicorp/go-multierror"
)

// GetParams is used on List
type GetParams struct {
	Params
	NoteID string
}

// Validate confirms the parmeters are valid
func (params GetParams) Validate() error {
	var merr = new(multierror.Error)

	if params.NoteID == "" {
		merr = multierror.Append(merr, errors.New(errEmptyNoteID))
	}

	merr = multierror.Append(merr, params.Params.Validate())

	return merr.ErrorOrNil()
}

// Get obtains a note from a deployment and note ID
func Get(params GetParams) (*models.Note, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	if err := params.fillDefaults(); err != nil {
		return nil, err
	}

	res, err := params.V1API.DeploymentsNotes.GetDeploymentNote(
		deployments_notes.NewGetDeploymentNoteParams().
			WithDeploymentID(params.ID).
			WithNoteID(params.NoteID),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}
