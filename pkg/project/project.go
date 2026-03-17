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
	"bytes"
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

// Search is an alias for the Elasticsearch project type.
const Search ProjectType = "search"

// ValidateType checks whether the given string is a valid project type.
// It accepts "search" as an alias for "elasticsearch".
func ValidateType(t string) (ProjectType, error) {
	switch ProjectType(strings.ToLower(t)) {
	case Elasticsearch, Search:
		return Elasticsearch, nil
	case Observability:
		return Observability, nil
	case Security:
		return Security, nil
	default:
		return "", fmt.Errorf(
			"invalid project type %q, must be one of: elasticsearch (or search), observability, security", t,
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
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Alias     string            `json:"alias,omitempty"`
	Type      string            `json:"type"`
	RegionID  string            `json:"region_id"`
	CloudID   string            `json:"cloud_id,omitempty"`
	Endpoints map[string]string `json:"endpoints,omitempty"`
	Metadata  Metadata          `json:"metadata,omitempty"`
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

// Credentials holds the basic auth credentials returned when creating a project.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// CreateResult is the response from creating a project.
type CreateResult struct {
	Project
	Credentials Credentials `json:"credentials"`
}

// CreateParams are the parameters for creating a project.
type CreateParams struct {
	API      *api.API
	Host     string
	Type     string
	Name     string
	RegionID string
	Tier     string
	Client   *http.Client
}

// Validate ensures the parameters are usable.
func (p CreateParams) Validate() error {
	var merr = multierror.NewPrefixed("invalid project create params")
	if p.API == nil {
		merr = merr.Append(errors.New("api reference is required"))
	}
	if p.Host == "" {
		merr = merr.Append(errors.New("host is required"))
	}
	if p.Type == "" {
		merr = merr.Append(errors.New("type is required"))
	}
	if p.Name == "" {
		merr = merr.Append(errors.New("name is required"))
	}
	if p.RegionID == "" {
		merr = merr.Append(errors.New("region is required"))
	}
	if p.Tier != "" {
		pt, _ := ValidateType(p.Type)
		if pt == Elasticsearch {
			merr = merr.Append(errors.New("tier is not supported for elasticsearch projects"))
		}
	}
	return merr.ErrorOrNil()
}

func (p CreateParams) httpClient() *http.Client {
	if p.Client != nil {
		return p.Client
	}
	return &http.Client{}
}

// Create creates a new serverless project.
func Create(params CreateParams) (*CreateResult, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	pt, err := ValidateType(params.Type)
	if err != nil {
		return nil, err
	}

	host := strings.TrimRight(params.Host, "/")
	endpoint := fmt.Sprintf("%s/api/v1/serverless/projects/%s", host, pt)

	payload, err := buildCreatePayload(pt, params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := newRequest(http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	req = params.API.AuthWriter.AuthRequest(req)

	resp, err := params.httpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var result CreateResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	if result.Type == "" {
		result.Type = string(pt)
	}

	return &result, nil
}

// DeleteParams are the parameters for deleting a project.
type DeleteParams struct {
	API    *api.API
	Host   string
	ID     string
	Type   string
	Client *http.Client
}

// Validate ensures the parameters are usable.
func (p DeleteParams) Validate() error {
	var merr = multierror.NewPrefixed("invalid project delete params")
	if p.API == nil {
		merr = merr.Append(errors.New("api reference is required"))
	}
	if p.Host == "" {
		merr = merr.Append(errors.New("host is required"))
	}
	if p.ID == "" {
		merr = merr.Append(errors.New("project id is required"))
	}
	return merr.ErrorOrNil()
}

func (p DeleteParams) httpClient() *http.Client {
	if p.Client != nil {
		return p.Client
	}
	return &http.Client{}
}

// Delete deletes a serverless project. If Type is empty it auto-detects
// the project type by listing all projects first.
func Delete(params DeleteParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	pt, err := resolveProjectType(params.Type, params.ID, params.API, params.Host, params.Client)
	if err != nil {
		return err
	}

	host := strings.TrimRight(params.Host, "/")
	endpoint := fmt.Sprintf("%s/api/v1/serverless/projects/%s/%s", host, pt, params.ID)

	req, err := newRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	req = params.API.AuthWriter.AuthRequest(req)

	resp, err := params.httpClient().Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// ShowParams are the parameters for showing a project.
type ShowParams struct {
	API    *api.API
	Host   string
	ID     string
	Type   string
	Client *http.Client
}

// Validate ensures the parameters are usable.
func (p ShowParams) Validate() error {
	var merr = multierror.NewPrefixed("invalid project show params")
	if p.API == nil {
		merr = merr.Append(errors.New("api reference is required"))
	}
	if p.Host == "" {
		merr = merr.Append(errors.New("host is required"))
	}
	if p.ID == "" {
		merr = merr.Append(errors.New("project id is required"))
	}
	return merr.ErrorOrNil()
}

func (p ShowParams) httpClient() *http.Client {
	if p.Client != nil {
		return p.Client
	}
	return &http.Client{}
}

// Show retrieves a single serverless project by ID. If Type is empty it
// auto-detects the project type by listing all projects first.
func Show(params ShowParams) (*Project, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	pt, err := resolveProjectType(params.Type, params.ID, params.API, params.Host, params.Client)
	if err != nil {
		return nil, err
	}

	host := strings.TrimRight(params.Host, "/")
	endpoint := fmt.Sprintf("%s/api/v1/serverless/projects/%s/%s", host, pt, params.ID)

	req, err := newRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	req = params.API.AuthWriter.AuthRequest(req)

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

	var result Project
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	if result.Type == "" {
		result.Type = string(pt)
	}

	return &result, nil
}

func resolveProjectType(typeHint, id string, a *api.API, host string, client *http.Client) (ProjectType, error) {
	if typeHint != "" {
		return ValidateType(typeHint)
	}

	res, err := List(ListParams{
		API:    a,
		Host:   host,
		Client: client,
	})
	if err != nil {
		return "", fmt.Errorf("failed to auto-detect project type: %w", err)
	}

	for _, p := range res.Projects {
		if p.ID == id {
			return ValidateType(p.Type)
		}
	}

	return "", fmt.Errorf("project %q not found", id)
}

type securityProductType struct {
	ProductLine string `json:"product_line"`
	ProductTier string `json:"product_tier"`
}

func buildCreatePayload(pt ProjectType, params CreateParams) ([]byte, error) {
	base := map[string]interface{}{
		"name":      params.Name,
		"region_id": params.RegionID,
	}

	if params.Tier != "" {
		switch pt {
		case Observability:
			base["product_tier"] = params.Tier
		case Security:
			base["product_types"] = []securityProductType{
				{ProductLine: "security", ProductTier: params.Tier},
				{ProductLine: "cloud", ProductTier: params.Tier},
				{ProductLine: "endpoint", ProductTier: params.Tier},
			}
		}
	}

	return json.Marshal(base)
}

func newRequest(method, endpoint string, body io.Reader) (*http.Request, error) {
	reqURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint URL: %w", err)
	}

	req, err := http.NewRequest(method, reqURL.String(), body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func listByType(params ListParams, pt ProjectType) (*ListResponse, error) {
	const maxPages = 100

	host := strings.TrimRight(params.Host, "/")
	baseEndpoint := fmt.Sprintf("%s/api/v1/serverless/projects/%s", host, pt)

	var allItems []Project
	cursor := ""

	for page := 0; ; page++ {
		if page >= maxPages {
			return nil, fmt.Errorf("exceeded maximum number of pages (%d) while listing %s projects", maxPages, pt)
		}
		reqURL, err := url.Parse(baseEndpoint)
		if err != nil {
			return nil, fmt.Errorf("invalid base endpoint URL: %w", err)
		}

		if cursor != "" {
			q := reqURL.Query()
			q.Set("from", cursor)
			reqURL.RawQuery = q.Encode()
		}

		req, err := newRequest(http.MethodGet, reqURL.String(), nil)
		if err != nil {
			return nil, err
		}

		req = params.API.AuthWriter.AuthRequest(req)

		resp, err := params.httpClient().Do(req)
		if err != nil {
			return nil, fmt.Errorf("request failed: %w", err)
		}

		body, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
		}

		var page ListResponse
		if err := json.Unmarshal(body, &page); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}

		allItems = append(allItems, page.Items...)

		if page.NextPage == "" {
			break
		}
		cursor = page.NextPage
	}

	return &ListResponse{Items: allItems}, nil
}
