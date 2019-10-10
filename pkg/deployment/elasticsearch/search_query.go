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

package elasticsearch

import (
	"errors"

	"github.com/elastic/cloud-sdk-go/pkg/models"
)

// NewQueryParams specifies user to perform the query on
type NewQueryParams struct {
	UserID string
}

// Validate ensures the structure is usable
func (params NewQueryParams) Validate() error {
	if params.UserID == "" {
		return errors.New("query user id cannot be empty")
	}
	return nil
}

// NewQueryByUserID creates query to match on userID
func NewQueryByUserID(params NewQueryParams) ([]*models.QueryContainer, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	var queries []*models.QueryContainer
	queries = append(queries,
		&models.QueryContainer{
			Term: map[string]models.TermQuery{
				"settings.metadata.owner_id": {
					Value: params.UserID,
				},
			},
		})

	return queries, nil
}

// NewStartedClusterQuery matches clusters currently running (started)
func NewStartedClusterQuery(includeInitializing bool) []*models.QueryContainer {
	var queries []*models.QueryContainer

	queries = append(queries,
		&models.QueryContainer{
			Term: map[string]models.TermQuery{
				"status": {
					Value: "started",
				},
			},
		})
	if includeInitializing {
		queries = append(queries,
			&models.QueryContainer{
				Term: map[string]models.TermQuery{
					"status": {
						Value: "initializing",
					},
				},
			})
	}

	return queries
}
