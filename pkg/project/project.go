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

package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/multierror"
)

// ProjectType represents the type of serverless project.
type ProjectType string

const (
	// Elasticsearch project type.
	Elasticsearch ProjectType = "elasticsearch"
	// Observability project type.
	Observability ProjectType = "observability"
	// Security project type.
	Security ProjectType = "security"
)

// AllTypes contains all valid project types.
var AllTypes = []ProjectType{Elasticsearch, Observability, Security}

// ValidateType checks whether the given string is a valid project type.
func ValidateType(t string) (ProjectType, error) {
	switch ProjectType(strings.ToLower(t)) {
	case Elasticsearch:
		return Elasticsearch, nil
	case Observability:
		return Observability, nil
	case Security:
		return Security, nil
	default:
		return "", fmt.Errorf(
			"invalid project type %q, must be one of: elasticsearch, observability, security", t,
		)
	}
}

// Metadata contains additional project metadata.
type Metadata struct {
	CreatedAt       string `json:"created_at,omitempty"`
	CreatedBy       string `json:"created_by,omitempty"`
	OrganizationID  string `json:"organization_id,omitempty"`
	SuspendedAt     string `json:"suspended_at,omitempty"`
	SuspendedReason string `json:"suspended_reason,omitempty"`
}

// Project represents a serverless project.
type Project struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Alias    string   `json:"alias,omitempty"`
	Type     string   `json:"type"`
	RegionID string   `json:"region_id"`
	CloudID  string   `json:"cloud_id,omitempty"`
	Metadata Metadata `json:"metadata,omitempty"`
}

// ListResponse is the response from listing projects of a given type.
type ListResponse struct {
	Items    []Project `json:"items"`
	NextPage string    `json:"next_page,omitempty"`
}

// ListResult is the aggregated result of listing projects, potentially
// across multiple project types.
type ListResult struct {
	Projects []Project `json:"projects"`
}

// ListParams are the parameters for listing projects.
type ListParams struct {
	API    *api.API
	Host   string
	Type   string
	Client *http.Client
}

// Validate ensures the parameters are usable.
func (p ListParams) Validate() error {
	var merr = multierror.NewPrefixed("invalid project list params")
	if p.API == nil {
		merr = merr.Append(errors.New("api reference is required"))
	}
	if p.Host == "" {
		merr = merr.Append(errors.New("host is required"))
	}
	return merr.ErrorOrNil()
}

func (p ListParams) httpClient() *http.Client {
	if p.Client != nil {
		return p.Client
	}
	return &http.Client{}
}

// List retrieves serverless projects. When Type is specified, only that
// project type is listed. Otherwise all types are queried.
func List(params ListParams) (*ListResult, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	types := AllTypes
	if params.Type != "" {
		pt, err := ValidateType(params.Type)
		if err != nil {
			return nil, err
		}
		types = []ProjectType{pt}
	}

	var result ListResult
	for _, pt := range types {
		resp, err := listByType(params, pt)
		if err != nil {
			return nil, fmt.Errorf("failed to list %s projects: %w", pt, err)
		}
		for i := range resp.Items {
			if resp.Items[i].Type == "" {
				resp.Items[i].Type = string(pt)
			}
		}
		result.Projects = append(result.Projects, resp.Items...)
	}

	return &result, nil
}

func listByType(params ListParams, pt ProjectType) (*ListResponse, error) {
	host := strings.TrimRight(params.Host, "/")
	endpoint := fmt.Sprintf("%s/api/v1/serverless/projects/%s", host, pt)

	reqURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint URL: %w", err)
	}

	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req = params.API.AuthWriter.AuthRequest(req)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := params.httpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var listResp ListResponse
	if err := json.Unmarshal(body, &listResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &listResp, nil
}
