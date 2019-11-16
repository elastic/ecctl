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

package cmd

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/elastic/uptd"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/ecctl"
)

const (
	updateFmt = `
Your version of %s is out of date! The latest version is %s.
You can use the binaries from %s.`
)

var (
	errWrapError        = errors.New("\ncan't check newer version")
	primaryGithubToken  = "GITHUB_TOKEN"
	fallbackGithubToken = "HOMEBREW_GITHUB_API_TOKEN"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Shows ecctl version",
	PreRunE: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Fprint(cmd.OutOrStdout(), versionInfo)

		return checkUpdate(versionInfo, cmd.OutOrStderr())
	},
}

func githubToken() string {
	var envVars = []string{primaryGithubToken}
	if runtime.GOOS == "darwin" {
		envVars = append(envVars, fallbackGithubToken)
	}
	for _, key := range envVars {
		if token := os.Getenv(key); token != "" {
			return token
		}
	}
	return ""
}

func checkUpdate(version ecctl.VersionInfo, device io.Writer) error {
	githubUpdateProvider, err := uptd.NewGithubProvider(
		version.Organization, version.Repository, githubToken(),
	)
	if err != nil {
		fmt.Fprintln(device, errors.Wrap(err, errWrapError.Error()))
		return nil
	}

	uptodate, err := uptd.New(githubUpdateProvider, version.Version)
	if err != nil {
		fmt.Fprintln(device, errors.Wrap(err, errWrapError.Error()))
		return nil
	}

	res, err := uptodate.Check()
	if err != nil {
		fmt.Fprintln(device, errors.Wrap(err, errWrapError.Error()))
		return nil
	}

	if res.NeedsUpdate {
		var message = fmt.Sprintf(updateFmt, RootCmd.Name(),
			res.Latest.Version.String(), res.Latest.URL,
		)
		fmt.Fprintln(device, message)
	}
	return nil
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
