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

package ecctl

import (
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/elastic/cloud-sdk-go/pkg/output"
	"github.com/elastic/cloud-sdk-go/pkg/util/slice"
	multierror "github.com/hashicorp/go-multierror"
)

const (
	// JSONOutput is the json output format
	JSONOutput = "json"
	// TextOutput is the text (templated) output format
	TextOutput = "text"
)

var (
	errCannotSpecifyJSONOutputAndCustomFormat = errors.New("cannot specify json output with format flag")
	errInvalidOutputFormat                    = errors.New("output must be one either json or text")
	errInvalidOutputDevice                    = errors.New("output device must not be nil")
	errInvalidErrorDevice                     = errors.New("error device must not be nil")
	errInvalidEmptyAuthenticaitonSettings     = errors.New("apikey or user and pass must be specified")
	errInvalidBothAuthenticaitonSettings      = errors.New("cannot specify both apikey and user / pass")
)

// Config contains the application configuration
type Config struct {
	User    string `json:"user,omitempty"`
	Pass    string `json:"pass,omitempty"`
	Host    string `json:"host,omitempty"`
	APIKey  string `json:"apikey,omitempty"`
	Region  string `json:"region,omitempty"`
	Output  string `json:"output,omitempty"`
	Message string `json:"message,omitempty"`
	Format  string `json:"format,omitempty"`

	OutputDevice *output.Device `json:"-"`
	ErrorDevice  io.Writer      `json:"-"`
	Client       *http.Client   `json:"-"`

	Timeout time.Duration `json:"timeout,omitempty"`

	Verbose  bool `json:"verbose,omitempty"`
	Force    bool `json:"force,omitempty"`
	Insecure bool `json:"insecure,omitempty"`
	// SkipLogin skips loging in when user and pass are set.
	SkipLogin bool `json:"-"`

	// SkipLogin skips loging in when user and pass are set.
	UserAgent string `json:"-"`
}

// Validate checks that the application config is a valid one
func (c *Config) Validate() error {
	var err = new(multierror.Error)
	if !slice.HasString([]string{JSONOutput, TextOutput}, c.Output) {
		err = multierror.Append(err, errInvalidOutputFormat)
	}

	var allCreds = c.APIKey != "" && (c.User != "" || c.Pass != "")
	if allCreds {
		err = multierror.Append(err, errInvalidBothAuthenticaitonSettings)
	}

	var emptyCreds = c.APIKey == "" && (c.User == "" || c.Pass == "")
	if emptyCreds {
		err = multierror.Append(err, errInvalidEmptyAuthenticaitonSettings)
	}

	if c.Output == JSONOutput && c.Format != "" {
		err = multierror.Append(err, errCannotSpecifyJSONOutputAndCustomFormat)
	}

	if c.OutputDevice == nil {
		err = multierror.Append(err, errInvalidOutputDevice)
	}

	if c.ErrorDevice == nil {
		err = multierror.Append(err, errInvalidErrorDevice)
	}

	return err.ErrorOrNil()
}
