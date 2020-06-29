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

package testutils

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/output"

	"github.com/elastic/ecctl/pkg/ecctl"
)

const (
	defaultOutputFormat = "json"
	defaultAPIKey       = "dummy"
	defaultRegion       = "ece-region"
)

// MockCfg represents a small and targeted amount of `ecctl.Config` options
// aimed at making mocking convenient and easy.
type MockCfg struct {
	Responses    []mock.Response
	Out          io.Writer
	Err          io.Writer
	OutputFormat string
	Format       string
	Region       string

	Force   bool
	Verbose bool
}

func fillDefaults(cfg MockCfg) MockCfg {
	if cfg.OutputFormat == "" {
		cfg.OutputFormat = defaultOutputFormat
	}

	if cfg.Region == "" {
		cfg.Region = defaultRegion
	}

	if cfg.Err == nil {
		cfg.Err = new(bytes.Buffer)
	}

	if cfg.Out == nil {
		cfg.Out = new(bytes.Buffer)
	}

	return cfg
}

func newConfig(cfg MockCfg) ecctl.Config {
	cfg = fillDefaults(cfg)
	return ecctl.Config{
		Client:       mock.NewClient(cfg.Responses...),
		OutputDevice: output.NewDevice(cfg.Out),
		ErrorDevice:  cfg.Err,
		Output:       cfg.OutputFormat,
		Format:       cfg.Format,
		Host:         fmt.Sprintf("https://%s", api.DefaultMockHost),
		APIKey:       defaultAPIKey,
		Force:        cfg.Force,
		Verbose:      cfg.Verbose,
	}
}

// MockApp initiates a mocked app from a MockCfg.
func MockApp(t *testing.T, cfg MockCfg) func() {
	if _, err := ecctl.Instance(newConfig(cfg)); err != nil {
		t.Error(err)
	}

	return ecctl.Cleanup
}
