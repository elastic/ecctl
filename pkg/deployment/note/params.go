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
	"github.com/elastic/ecctl/pkg/deployment"
)

var (
	errEmptyNoteMessage = "note comment cannot be empty"
	errEmptyUserID      = "user id cannot be empty"
	errEmptyNoteID      = "note id cannot be empty"
)

// Params is used on Get and Update Notes
type Params struct {
	deployment.Params
}

// Use different resource types when this is supported by the API.
// For the time being, the notes endpoint only allows elasticsearch IDs.
func (params *Params) fillDefaults() error {
	esID, err := deployment.GetElasticsearchID(deployment.GetParams{
		API:          params.API,
		DeploymentID: params.ID,
	})
	if err != nil {
		return err
	}

	params.ID = esID
	return err
}
