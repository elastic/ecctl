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
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/multierror"
	"github.com/elastic/cloud-sdk-go/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
)

const (
	homePrefix = "$HOME"
)

var (
	defaultClient = new(http.Client)
	defaultOutput = os.Stdout
	defaultInput  = os.Stdin
	defaultError  = os.Stderr
	defaultViper  = viper.New()

	ecctlHomePath = filepath.Join(homePrefix, ".ecctl")

	bashCompletionFunc = `__ecctl_valid_regions()
{
   COMPREPLY=($(echo ${EC_REGIONS}))
}
` + cmdutil.StatelessKindsCompFunc + "\n" +
		cmdutil.AllKindsCompFunc
)

var (
	versionInfo                 ecctl.VersionInfo
	excludedApplicationCommands = []string{
		"help", "version", "generate", "docs", "completions", "init",
	}
	messageErrHasNoPreRunCheck = "command %s/%s has no PreRunE check set"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:               "ecctl",
	Short:             "Elastic Cloud Control",
	SilenceErrors:     true,
	SilenceUsage:      true,
	DisableAutoGenTag: true,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	BashCompletionFunction: bashCompletionFunc,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		setupViper(defaultViper)
		if err := setupDebug(defaultViper.GetBool("trace"), defaultViper.GetBool("pprof")); err != nil {
			return err
		}

		return initApp(cmd, defaultClient, defaultViper)
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
// It returns the statuscode to be used by os.Exit.
func Execute(v ecctl.VersionInfo) int {
	defer stopDebug(defaultViper)

	populateValidArgs(RootCmd)
	versionInfo = v

	if err := RootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(RootCmd.OutOrStderr(), err)
		if ret, ok := err.(ecctl.ReturnCodeError); ok {
			return ret.ReturnCode()
		}
		return -1
	}
	return 0
}

func init() {
	RootCmd.PersistentFlags().String("config", "config", "Config name, used to have multiple configs in $HOME/.ecctl/<env>")
	RootCmd.PersistentFlags().String("host", "", "Base URL to use")
	RootCmd.PersistentFlags().String("user", "", "Username to use to authenticate (If empty will look for EC_USER environment variable)")
	RootCmd.PersistentFlags().String("pass", "", "Password to use to authenticate (If empty will look for EC_PASS environment variable)")
	RootCmd.PersistentFlags().String("api-key", "", "API key to use to authenticate (If empty will look for EC_API_KEY environment variable)")
	RootCmd.PersistentFlags().Bool("verbose", false, "Enable verbose mode")
	RootCmd.PersistentFlags().Bool("verbose-credentials", false, "When set, Authorization headers on the request/response trail will be displayed as plain text")
	RootCmd.PersistentFlags().String("verbose-file", "", "When set, the verbose request/response trail will be written to the defined file")
	RootCmd.PersistentFlags().String("output", "text", "Output format [text|json]")
	RootCmd.PersistentFlags().Bool("force", false, "Do not ask for confirmation")
	RootCmd.PersistentFlags().String("message", "", "A message to set on cluster operation")
	RootCmd.PersistentFlags().String("format", "", "Formats the output using a Go template")
	RootCmd.PersistentFlags().Bool("trace", false, "Enables tracing saves the trace to trace-20060102150405")
	RootCmd.PersistentFlags().Bool("pprof", false, "Enables pprofing and saves the profile to pprof-20060102150405")
	RootCmd.PersistentFlags().Bool("insecure", false, "Skips all TLS validation")
	RootCmd.PersistentFlags().BoolP("quiet", "q", false, "Suppresses the configuration file used for the run, if any")
	RootCmd.PersistentFlags().Duration("timeout", time.Second*30, "Timeout to use on all HTTP calls")
	RootCmd.PersistentFlags().String("region", "", "Elastic Cloud Hosted region")
	RootCmd.Flag("region").Annotations = map[string][]string{
		cobra.BashCompCustom: {"__ecctl_valid_regions"},
	}

	defaultViper.BindPFlags(RootCmd.PersistentFlags())
}

// setupViper configures a `*viper.Viper` instance for ecctl use.
func setupViper(v *viper.Viper) {
	v.SetEnvPrefix("EC")
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()
	v.AddConfigPath(ecctlHomePath)         // adding home directory as first search path
	v.SetConfigName(v.GetString("config")) // name of config file (without extension)

	// If a config file is found, read it in.
	if err := v.ReadInConfig(); err == nil && v.GetBool("verbose") {
		fmt.Fprintln(os.Stderr, "Using config file:", v.ConfigFileUsed())
	}
	// Register an alias value after the config file has been read.
	v.RegisterAlias("api_key", "api-key")
	v.RegisterAlias("verbose_file", "verbose-file")
	v.RegisterAlias("verbose_credentials", "verbose-credentials")
}

// populateValidArgs dynamically generates the validargs for all of the cobra
// commands and subcommands
func populateValidArgs(cmd *cobra.Command) {
	for _, c := range cmd.Commands() {
		var args = append(c.Aliases, c.Name())
		cmd.ValidArgs = append(cmd.ValidArgs, args...)
		if cmd.HasAvailableSubCommands() {
			populateValidArgs(c)
		}
	}
}

// GetCommand returns a child command from the command that is passed.
// If the command is not found, the parent is returned.
func GetCommand(command *cobra.Command, path ...string) *cobra.Command {
	for _, p := range path {
		for _, c := range command.Commands() {
			if c.Name() == p {
				return GetCommand(c, path[1:]...)
			}
		}
	}
	return command
}

func initApp(cmd *cobra.Command, client *http.Client, v *viper.Viper) error {
	for _, cmdName := range excludedApplicationCommands {
		if cmd.Name() == cmdName {
			return nil
		}
	}

	if client == nil {
		return errors.New("cmd: root http client cannot be nil")
	}

	var c = ecctl.Config{
		Client:       client,
		OutputDevice: output.NewDevice(defaultOutput),
		ErrorDevice:  defaultError,
		UserAgent:    strings.Join([]string{"ecctl", versionInfo.Version}, "/"),
	}
	if err := v.Unmarshal(&c); err != nil {
		return err
	}

	// Set the default region to `ece-region` when the endpoint is not the ESS endpoint.
	if c.Region == "" && c.Host != api.ESSEndpoint {
		c.Region = cmdutil.DefaultECERegion
	}

	err := api.ReturnErrOnly(ecctl.Instance(c))
	// When no config file has been read and initApp returns an error, tell
	// the user how to initialize the application.
	if err != nil && v.ConfigFileUsed() == "" {
		return multierror.NewPrefixed(
			`missing ecctl config file, please use the "ecctl init" command to initialize ecctl`, err,
		)
	}

	return err
}

func checkPreRunE(command *cobra.Command) error {
	var err = multierror.NewPrefixed("")
	for _, c := range command.Commands() {
		if command.HasSubCommands() {
			err = err.Append(checkPreRunE(c))
		}
		if c.PreRunE == nil {
			var message = fmt.Sprintf(messageErrHasNoPreRunCheck, command.Name(), c.Name())
			err = err.Append(errors.New(message))
		}
	}

	return err.ErrorOrNil()
}
