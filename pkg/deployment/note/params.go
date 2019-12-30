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

	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment"
)

var (
	errEmptyNoteMessage = "note message cannot be empty"
	errEmptyUserID      = "user id cannot be empty"
	errEmptyNoteID      = "note id cannot be empty"
)

// Params is used on Get and Update Notes
type Params struct {
	deployment.Params
	NoteID string
}

// Validate confirms the parmeters are valid
func (params Params) Validate() error {
	var merr = new(multierror.Error)

	if params.NoteID == "" {
		merr = multierror.Append(merr, errors.New(errEmptyNoteID))
	}

	merr = multierror.Append(merr, params.Params.Validate())

	return merr.ErrorOrNil()
}

func getElasticsearchID(params deployment.GetParams) (string, error) {
	res, err := deployment.Get(params)
	if err != nil {
		return "", err
	}

	deploymentID := *res.Resources.Elasticsearch[0].ID
	return deploymentID, nil
}
