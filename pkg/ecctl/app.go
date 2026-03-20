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
	"fmt"
	"os"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/auth"

	"github.com/elastic/ecctl/pkg/formatter"
)

// App stitches together all the parts to create the ecctl application.
type App struct {
	API       *api.API
	Formatter formatter.Formatter
	Config    Config
}

// NewApplication returns a fully initialized App, which will be called from the presentation layer (cmd)
func NewApplication(c Config) (*App, error) {
	cfg, err := newAPIConfig(c)
	if err != nil {
		return nil, err
	}

	apiInstance, err := api.NewAPI(cfg)
	if err != nil {
		return nil, err
	}

	// Sets any extra message that is passed to the commentator
	GetOperationInstance().Set(c.Message)

	var fmter formatter.Formatter = formatter.New(c.OutputDevice, c.Output)
	if c.Format != "" {
		fmter = formatter.NewText(&formatter.TextConfig{
			Output:   c.OutputDevice,
			Override: c.Format,
		})
	}

	return &App{
		API:       apiInstance,
		Formatter: fmter,
		Config:    c,
	}, nil
}

func newAPIConfig(cfg Config) (api.Config, error) {
	var empty api.Config
	if err := cfg.Validate(); err != nil {
		return empty, err
	}

	apiCfg := api.Config{
		Client:        cfg.Client,
		Host:          cfg.Host,
		SkipTLSVerify: cfg.Insecure,
		Timeout:       cfg.Timeout,
		VerboseSettings: api.VerboseSettings{
			Verbose:    cfg.Verbose,
			Device:     cfg.OutputDevice,
			RedactAuth: !cfg.VerboseCredentials,
		},
		SkipLogin:   cfg.SkipLogin,
		ErrorDevice: cfg.ErrorDevice,
		UserAgent:   cfg.UserAgent,
	}

	authWriter, err := auth.NewAuthWriter(auth.Config{
		APIKey: cfg.APIKey, Username: cfg.User, Password: cfg.Pass,
	})
	if err != nil {
		return empty, err
	}
	apiCfg.AuthWriter = authWriter

	if cfg.VerboseFile != "" {
		f, err := os.Create(cfg.VerboseFile)
		if err != nil {
			return empty, fmt.Errorf(
				`failed creating verbose file "%s": %w`, cfg.VerboseFile, err,
			)
		}
		apiCfg.Device = f
	}

	return apiCfg, nil
}
