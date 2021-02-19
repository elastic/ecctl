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
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/elastic/ecctl/pkg/ecctl"
)

var initCmd = &cobra.Command{
	Use:     "init",
	Short:   "Creates an initial configuration file.",
	PreRunE: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		fp := strings.Replace(ecctlHomePath, homePrefix, cmdutil.GetHomePath(runtime.GOOS), 1)
		if err := os.MkdirAll(fp, 0775); err != nil {
			return err
		}

		err := ecctl.InitConfig(ecctl.InitConfigParams{
			Client:           defaultClient,
			Viper:            defaultViper,
			Reader:           defaultInput,
			Writer:           defaultOutput,
			ErrWriter:        defaultError,
			PasswordReadFunc: term.ReadPassword,
			FilePath:         filepath.Join(fp, defaultViper.GetString("config")),
		})
		if err != nil {
			return err
		}

		// Only print this message when viper.ConfigFileUsed() since it means
		// the config hasn't been read via viper, thus just recently written
		// down in the path.
		if cfg := defaultViper.ConfigFileUsed(); cfg == "" {
			if err := defaultViper.ReadInConfig(); err != nil {
				return err
			}
			cfg = defaultViper.ConfigFileUsed()
			fmt.Fprintln(defaultOutput, "\nConfig written to", cfg)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
