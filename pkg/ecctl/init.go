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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/elastic/ecctl/pkg/user"

	"github.com/elastic/cloud-sdk-go/pkg/input"
	"github.com/elastic/cloud-sdk-go/pkg/output"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/spf13/viper"
)

const (
	_ = iota
	apiKeyChoice
	userPassChoice
)

const (
	_ = iota
	textFormatChoice
	jsonFormatChoice
)

const (
	disclaimer      = "Welcome to the Elastic Cloud CLI! This command will guide you through authenticating and setting some default values.\n\n"
	redacted        = "[REDACTED]"
	settingsPathMsg = "Found existing settings in %s. Here's a JSON representation of what they look like:\n"

	missingConfigMsg  = `Missing configuration file, would you like to initialise it? [y/n]: `
	existingConfigMsg = `Would you like to change your current settings? [y/n]: `

	hostMsg = "Enter the URL of the Elastic Cloud API or your ECE installation: "

	apiKeyMsg = "Paste your API Key and press enter: "
	userMsg   = "Type in your username: "
	//nolint
	passMsg = "Type in your password: "

	validCredentialsMsg = "Your credentials seem to be valid, and show you're authenticated as \"%s\".\n\n"
)

var (
	authChoiceMsg = `
Which authentication mechanism would you like to use?
  [1] API Keys (Recommended).
  [2] Username and Password login.

Please enter your choice: `

	formatChoiceMsg = `
What default output format would you like?
  [1] text - Human-readable output format, commands with no output templates defined will fall back to JSON.
  [2] json - JSON formatted output API responses.

Please enter a choice: `

	finalMsg = `
You're all set! Here are some commands to try:
  $ ecctl auth user key list
  $ ecctl deployment elasticsearch list`[1:]
)

// PassFunc represents the function used to consume a password.
type PassFunc func(fd int) ([]byte, error)

// InitConfigParams is consumed by InitConfig
type InitConfigParams struct {
	// Viper instance
	Viper *viper.Viper

	// Reader is Input reader
	Reader io.Reader

	// Writer is where the output will be written
	Writer io.Writer

	// ErrWriter is where any Errors will be written
	ErrWriter io.Writer

	// PasswordReadFunc is the function used to read a password from the file descriptor.
	PasswordReadFunc PassFunc

	// Client used to perform authentication validations.
	Client *http.Client

	// FilePath of the configuration
	FilePath string
}

// Validate ensures the parameters are usable.
func (params InitConfigParams) Validate() error {
	var merr = new(multierror.Error)
	if params.Viper == nil {
		merr = multierror.Append(merr, errors.New("init: viper instance cannot be nil"))
	}

	if params.Reader == nil {
		merr = multierror.Append(merr, errors.New("init: input reader cannot be nil"))
	}

	if params.Writer == nil {
		merr = multierror.Append(merr, errors.New("init: output writer cannot be nil"))
	}

	if params.ErrWriter == nil {
		merr = multierror.Append(merr, errors.New("init: error writer cannot be nil"))
	}

	if params.PasswordReadFunc == nil {
		merr = multierror.Append(merr, errors.New("init: password read function cannot be nil"))
	}

	if params.Client == nil {
		merr = multierror.Append(merr, errors.New("init: http client cannot be nil"))
	}

	return merr.ErrorOrNil()
}

// InitConfig initialises a configuration file or changes an existing one
// it thought of as a mechanism for users to onboard easily and have a guided
// and interactive configuration bootstrap.
func InitConfig(params InitConfigParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	fmt.Fprint(params.Writer, disclaimer)

	configRead := params.Viper.ReadInConfig() == nil
	var confirmationMessage = missingConfigMsg
	if configRead {
		confirmationMessage = existingConfigMsg
		params.FilePath = params.Viper.ConfigFileUsed()

		if err := printConfig(params.Writer, params.Viper); err != nil {
			return err
		}
	}

	outputDevice := output.NewDevice(params.Writer)
	scanner := input.NewScanner(params.Reader, outputDevice)

	shouldExit := !strings.EqualFold(scanner.Scan(confirmationMessage), "y")
	if shouldExit {
		return nil
	}

	// Populate config
	var cfg = Config{
		Client:       params.Client,
		OutputDevice: outputDevice,
		ErrorDevice:  params.ErrWriter,
	}
	if err := params.Viper.Unmarshal(&cfg); err != nil {
		return err
	}

	cfg.Host = scanner.Scan(hostMsg)

	if err := askOutputFormat(&cfg, scanner, params.Writer, params.ErrWriter); err != nil {
		return err
	}

	if err := askAuthMechanism(&cfg, scanner, params.Writer, params.PasswordReadFunc); err != nil {
		return err
	}

	if err := validateAuth(cfg, params.Writer); err != nil {
		return err
	}

	fmt.Fprintln(params.Writer, finalMsg)

	if err := setViperConfig(cfg, params.Viper); err != nil {
		return err
	}

	// It's better to write the config as is since it ommits defaults and
	// empties vs viper's behaviour in `WriteConfig`.
	return writeConfig(cfg, params.FilePath, ".json")
}

func writeConfig(cfg Config, filePath, ext string) error {
	configBytes, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	if !strings.HasSuffix(filePath, ext) {
		filePath += ext
	}

	return ioutil.WriteFile(filePath, configBytes, 0666)
}

func printConfig(writer io.Writer, v *viper.Viper) error {
	fmt.Fprintf(writer, settingsPathMsg, v.ConfigFileUsed())
	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return err
	}

	if c.Pass != "" {
		c.Pass = redacted
	}
	if c.APIKey != "" {
		c.APIKey = redacted
	}

	enc := json.NewEncoder(writer)
	enc.SetIndent("", "  ")
	return enc.Encode(c)
}

func askOutputFormat(cfg *Config, scanner *input.Scanner, writer, errWriter io.Writer) error {
	formatChoiceRaw := scanner.Scan(formatChoiceMsg)
	fmt.Fprintln(writer)
	formatChoice, err := strconv.Atoi(formatChoiceRaw)
	if err != nil {
		return err
	}

	cfg.Output = "text"
	switch formatChoice {
	case textFormatChoice:
	case jsonFormatChoice:
		cfg.Output = "json"
	default:
		fmt.Fprintln(errWriter, "invalid choice, defaulting to \"text\"")
	}

	return nil
}

func askAuthMechanism(cfg *Config, scanner *input.Scanner, writer io.Writer, passFunc PassFunc) error {
	authChoiceRaw := scanner.Scan(authChoiceMsg)
	fmt.Fprintln(writer)
	authChoice, err := strconv.Atoi(authChoiceRaw)
	if err != nil {
		return err
	}

	switch authChoice {
	default:
		return errors.New("invalid authentication choice")
	case apiKeyChoice:
		apikey, err := ReadSecret(writer, passFunc, apiKeyMsg)
		if err != nil {
			return err
		}
		cfg.APIKey = string(apikey)
		cfg.User, cfg.Pass = "", ""
	case userPassChoice:
		cfg.User = scanner.Scan(userMsg)

		pass, err := ReadSecret(writer, passFunc, passMsg)
		if err != nil {
			return err
		}
		cfg.Pass = string(pass)
		cfg.APIKey = ""
	}

	fmt.Fprintln(writer)
	return nil
}

func validateAuth(cfg Config, writer io.Writer) error {
	a, err := NewApplication(cfg)
	if err != nil {
		return err
	}

	u, err := user.GetCurrent(user.GetCurrentParams{API: a.API})
	if err != nil {
		return err
	}

	fmt.Fprintf(writer, validCredentialsMsg, *u.UserName)

	return nil
}

func setViperConfig(cfg Config, v *viper.Viper) error {
	b, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	v.SetConfigType("json")
	return v.ReadConfig(bytes.NewReader(b))
}
