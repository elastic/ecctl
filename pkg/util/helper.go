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

package util

import (
	"errors"
)

const (
	// DefaultESPort is the default port on ES clusters running on cloud
	DefaultESPort = 9243
	// DefaultClientTimeout is the default client http timeout for console requests

	// DefaultIndexPollerRate is the polling interval in seconds for tab-completion on index endpoints
	DefaultIndexPollerRate = 15
)

var (
	// ErrAPIReq is the message returned when API reference is required for a command
	ErrAPIReq = errors.New("api reference is required for command")
	// ErrDeploymentID is the message returned when a provided cluster id is not of the expected length (32 chars)
	ErrDeploymentID = errors.New("deployment id should have a length of 32 characters")
	// ErrIDCannotBeEmpty is the message returned when an ID field is empty
	ErrIDCannotBeEmpty = errors.New("id field cannot be empty")

	// SkipMaintenanceHeaders tells the EC proxy layer to still send requests to the
	// underlying cluster instances even if they are in maintenance mode
	SkipMaintenanceHeaders = map[string]string{
		"X-Found-Bypass-Maintenance": "true",
	}
)
