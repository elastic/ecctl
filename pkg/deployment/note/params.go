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

	"github.com/elastic/ecctl/pkg/deployment"
)

var (
	errEmptyNoteMessage = errors.New("note message cannot be empty")
	errEmptyUserID      = errors.New("user id cannot be empty")
	errEmptyNoteID      = errors.New("note id cannot be empty")
)

// ListParams is used by List
type ListParams struct {
	deployment.Params
}

// Params is used on Get and Update Notes
type Params struct {
	deployment.Params
	NoteID string
}

// Validate confirms the parmeters are valid
func (params Params) Validate() error {
	if params.NoteID == "" {
		return errEmptyNoteID
	}

	return params.Params.Validate()
}

// GetParams is used on ListNotes
type GetParams struct {
	Params
}

// UpdateParams is used on ListNotes
type UpdateParams struct {
	Params
	UserID  string
	Message string
}

// Validate confirms the parmeters are valid
func (params UpdateParams) Validate() error {
	if params.Message == "" {
		return errEmptyNoteMessage
	}

	if params.UserID == "" {
		return errEmptyUserID
	}

	if params.Params.NoteID == "" {
		return errEmptyNoteID
	}

	return params.Params.Validate()
}
