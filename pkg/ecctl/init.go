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
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/elastic/cloud-sdk-go/pkg/api/deploymentapi"
	"github.com/elastic/cloud-sdk-go/pkg/api/userapi"
	"github.com/elastic/cloud-sdk-go/pkg/input"
	"github.com/elastic/cloud-sdk-go/pkg/multierror"
	"github.com/elastic/cloud-sdk-go/pkg/output"
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
	_ = iota
	essInfraChoice
	eceInfraChoice
	esspInfraChoice
)

const (
	_ = iota
	gcpUsCentral1Choice
	gcpUsEast1Choice
	gcpUsEast4Choice
	gcpUsWest1Choice
	gcpNorthamericaNortheast1Choice
	gcpSouthamericaEast1Choice
	gcpAustraliaSoutheast1Choice
	gcpEuropeWest1Choice
	gcpEuropeWest2Choice
	gcpEuropeWest3Choice
	gcpAsiaNortheast1Choice
	gcpAsiaSouth1Choice
	gcpAsiaSoutheast1Choice

	awsUsEast1Choice
	awsUsWest1Choice
	awsUsWest2Choice
	awsEuCentral1Choice
	awsEuWest2Choice
	awsEuWest1Choice
	awsApNortheast1Choice
	awsApSoutheast1Choice
	awsApSoutheast2Choice
	awsSaEast1Choice

	azureEastUs2Choice
	azureWestUs2Choice
	azureWestEuropeChoice
	azureUkSouthChoice
	azureJapanEastChoice
	azureSouthEastAsiaChoice
)

const (
	disclaimer      = "Welcome to Elastic Cloud Control (ecctl)! This command will guide you through authenticating and setting some default values.\n\n"
	redacted        = "[REDACTED]"
	settingsPathMsg = "Found existing settings in %s. Here's a JSON representation of what they look like:\n"

	missingConfigMsg  = `Missing configuration file, would you like to initialise it? [y/n]: `
	existingConfigMsg = `Would you like to change your current settings? [y/n]: `

	essHostAddress = "https://api.elastic-cloud.com"
	eceHostMsg     = "Enter the URL of your ECE installation: "
	esspHostMsg    = "Enter the URL of your ESSP installation: "
	essChoiceMsg   = "Using \"%s\" as the API endpoint.\n"

	essAPIKeyCreateMsg = "Create a new Elasticsearch Service API key (https://cloud.elastic.co/deployment-features/keys) and/or"
	apiKeyMsg          = "Paste your API Key and press enter: "
	userMsg            = "Type in your username: "
	//nolint
	passMsg = "Type in your password: "

	validCredentialsMsg            = "Your credentials seem to be valid, and show you're authenticated as \"%s\".\n\n"
	validCredentialsAlternativeMsg = "Your credentials seem to be valid.\n\n"
	invalidCredentialsMsg          = "Your credentials couldn't be validated. Make sure they're correct and try again"
)

var (
	hostChoiceMsg = `
Select which type of Elastic Cloud offering you will be working with:
  [1] Elasticsearch Service (default).
  [2] Elastic Cloud Enterprise (ECE).
  [3] Elasticsearch Service Private (ESSP).

Please enter your choice: `

	regionChoiceMsg = `
Select a region you would like to have as default:

  GCP
  [1] us-central1 (Iowa)
  [2] us-east1 (S. Carolina)
  [3] us-east4 (N. Virginia)
  [4] us-west1 (Oregon)
  [5] northamerica-northeast1 (Montreal)
  [6] southamerica-east1 (São Paulo)
  [7] australia-southeast1 (Sydney)
  [8] europe-west1 (Belgium)
  [9] europe-west2 (London)
  [10] europe-west3 (Frankfurt)
  [11] asia-northeast1 (Tokyo)
  [12] asia-south1 (Mumbai)
  [13] asia-southeast1 (Singapore)

  AWS
  [14] us-east-1 (N. Virginia)
  [15] us-west-1 (N. California)
  [16] us-west-2 (Oregon)
  [17] eu-central-1 (Frankfurt)
  [18] eu-west-2 (London)
  [19] eu-west-1 (Ireland)
  [20] ap-northeast-1 (Tokyo)
  [21] ap-southeast-1 (Singapore)
  [22] ap-southeast-2 (Sydney)
  [23] sa-east-1 (São Paulo)

  Azure
  [24] eastus2 (Virginia)
  [25] westus2 (Washington)
  [26] westeurope (Netherlands)
  [27] uksouth (London)
  [28] japaneast (Tokyo)
  [29] southeastasia (Singapore)

Please enter your choice: `

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
  $ ecctl deployment list`[1:]

	// Remove once we have an endpoint available to list regions.
	essRegions = map[int]string{
		gcpUsCentral1Choice:             "gcp-us-central1",
		gcpUsEast1Choice:                "gcp-us-east1",
		gcpUsEast4Choice:                "gcp-us-east4",
		gcpUsWest1Choice:                "gcp-us-west1",
		gcpNorthamericaNortheast1Choice: "gcp-northamerica-northeast1",
		gcpSouthamericaEast1Choice:      "gcp-southamerica-east1",
		gcpAustraliaSoutheast1Choice:    "gcp-australia-southeast1",
		gcpEuropeWest1Choice:            "gcp-europe-west1",
		gcpEuropeWest2Choice:            "gcp-europe-west2",
		gcpEuropeWest3Choice:            "gcp-europe-west3",
		gcpAsiaNortheast1Choice:         "gcp-asia-northeast1",
		gcpAsiaSouth1Choice:             "gcp-asia-south1",
		gcpAsiaSoutheast1Choice:         "gcp-asia-southeast1",

		awsUsEast1Choice:      "us-east-1",
		awsUsWest1Choice:      "us-west-1",
		awsUsWest2Choice:      "us-west-2",
		awsEuCentral1Choice:   "aws-eu-central-1",
		awsEuWest2Choice:      "aws-eu-west-2",
		awsEuWest1Choice:      "eu-west-1",
		awsApNortheast1Choice: "ap-northeast-1",
		awsApSoutheast1Choice: "ap-southeast-1",
		awsApSoutheast2Choice: "ap-southeast-2",
		awsSaEast1Choice:      "sa-east-1",

		azureEastUs2Choice:       "azure-eastus2",
		azureWestUs2Choice:       "azure-westus2",
		azureWestEuropeChoice:    "azure-westeurope",
		azureUkSouthChoice:       "azure-uksouth",
		azureJapanEastChoice:     "azure-japaneast",
		azureSouthEastAsiaChoice: "azure-southeastasia",
	}
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
	var merr = multierror.NewPrefixed("invalid init configuration")
	if params.Viper == nil {
		merr = merr.Append(errors.New("viper instance cannot be nil"))
	}

	if params.Reader == nil {
		merr = merr.Append(errors.New("input reader cannot be nil"))
	}

	if params.Writer == nil {
		merr = merr.Append(errors.New("output writer cannot be nil"))
	}

	if params.ErrWriter == nil {
		merr = merr.Append(errors.New("error writer cannot be nil"))
	}

	if params.PasswordReadFunc == nil {
		merr = merr.Append(errors.New("password read function cannot be nil"))
	}

	if params.Client == nil {
		merr = merr.Append(errors.New("http client cannot be nil"))
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

	// Insecure is set to true by default to allow API calls against HTTPS
	// endpoints with self-signed certificates.
	cfg.Insecure = true

	if err := askInfraSelection(&cfg, scanner, params.Writer, params.ErrWriter, params.PasswordReadFunc); err != nil {
		return err
	}

	if err := askOutputFormat(&cfg, scanner, params.Writer, params.ErrWriter); err != nil {
		return err
	}

	if err := validateAuth(cfg, params.Writer); err != nil {
		return err
	}

	fmt.Fprintln(params.Writer, finalMsg)

	if err := setViperConfig(cfg, params.Viper); err != nil {
		return err
	}

	// It's better to write the config as is since it omits defaults and
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

	return os.WriteFile(filePath, configBytes, 0666)
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

func askInfraSelection(cfg *Config, scanner *input.Scanner, writer, errWriter io.Writer, passFunc PassFunc) error {
	infraChoiceRaw := scanner.Scan(hostChoiceMsg)
	fmt.Fprintln(writer)
	infraChoice, err := strconv.Atoi(infraChoiceRaw)
	if err != nil {
		return err
	}

	cfg.Host = essHostAddress
	switch infraChoice {
	case essInfraChoice:
		fmt.Fprintf(writer, essChoiceMsg, essHostAddress)
		if err := askRegionSelection(cfg, scanner, writer, essRegions); err != nil {
			return err
		}
		fmt.Fprintln(writer, essAPIKeyCreateMsg)
		if err := askAPIKey(cfg, writer, passFunc); err != nil {
			return err
		}
	case eceInfraChoice:
		cfg.Host = scanner.Scan(eceHostMsg)
		if err := askAuthMechanism(cfg, scanner, writer, passFunc); err != nil {
			return err
		}
	case esspInfraChoice:
		cfg.Host = scanner.Scan(esspHostMsg)
		if err := askAPIKey(cfg, writer, passFunc); err != nil {
			return err
		}
		// For the time being the only available region for ESSP is us-west-2. Once more
		// regions have been added, this should be set in a similar way to essInfraChoice
		cfg.Region = "us-west-2"
	default:
		fmt.Fprintf(errWriter, "invalid choice, defaulting to %s", essHostAddress)
	}

	return nil
}

func askRegionSelection(cfg *Config, scanner *input.Scanner, writer io.Writer, regions map[int]string) error {
	regionChoiceRaw := scanner.Scan(regionChoiceMsg)
	fmt.Fprintln(writer)
	regionChoice, err := strconv.Atoi(regionChoiceRaw)
	if err != nil {
		return err
	}

	region, ok := regions[regionChoice]
	if !ok {
		return errors.New("invalid region choice")
	}

	cfg.Region = region

	return nil
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

	fmt.Fprintln(writer)
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
		if err := askAPIKey(cfg, writer, passFunc); err != nil {
			return err
		}
	case userPassChoice:
		cfg.User = scanner.Scan(userMsg)

		pass, err := ReadSecret(writer, passFunc, passMsg)
		if err != nil {
			return err
		}
		cfg.Pass = string(pass)
		cfg.APIKey = ""
	}

	return nil
}

func askAPIKey(cfg *Config, writer io.Writer, passFunc PassFunc) error {
	apikey, err := ReadSecret(writer, passFunc, apiKeyMsg)
	if err != nil {
		return err
	}

	cfg.APIKey = string(apikey)
	cfg.User, cfg.Pass = "", ""

	return nil
}

func validateAuth(cfg Config, writer io.Writer) error {
	a, err := NewApplication(cfg)
	if err != nil {
		return err
	}

	u, err := userapi.GetCurrent(userapi.GetCurrentParams{API: a.API})
	if err != nil {
		if _, e := deploymentapi.List(deploymentapi.ListParams{
			API: a.API,
		}); e != nil {
			// nolint
			return errors.New(invalidCredentialsMsg)
		}
		fmt.Fprint(writer, validCredentialsAlternativeMsg)
		return nil
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
