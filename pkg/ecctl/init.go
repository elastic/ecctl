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
	_ = iota
	essInfraChoice
	eceInfraChoice
	esspInfraChoice
)

const (
	_ = iota
	gcpUsCentral1Choice
	gcpUsEast4Choice
	gcpUsWest1Choice
	gcpNorthamericaNortheast1Choice
	gcpAustraliaSoutheast1Choice
	gcpEuropeWest1Choice
	gcpEuropeWest2Choice
	gcpEuropeWest3Choice
	gcpAsiaNortheast1Choice
	gcpAsiaSouth1Choice
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

	apiKeyMsg = "Paste your API Key and press enter: "
	userMsg   = "Type in your username: "
	//nolint
	passMsg = "Type in your password: "

	validCredentialsMsg = "Your credentials seem to be valid, and show you're authenticated as \"%s\".\n\n"
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
  [2] us-east4 (N. Virginia)
  [3] us-west1 (Oregon)
  [4] northamerica-northeast1 (Montreal)
  [5] australia-southeast1 (Sydney)
  [6] europe-west1 (Belgium)
  [7] europe-west2 (London)
  [8] europe-west3 (Frankfurt)
  [9] asia-northeast1 (Tokyo)
  [10] asia-south1 (Mumbai)

  AWS
  [11] us-east-1 (N. Virginia)
  [12] us-west-1 (N. California)
  [13] us-west-2 (Oregon)
  [14] eu-central-1 (Frankfurt)
  [15] eu-west-2 (London)
  [16] eu-west-1 (Ireland)
  [17] ap-northeast-1 (Tokyo)
  [18] ap-southeast-1 (Singapore)
  [19] ap-southeast-2 (Sydney)
  [20] sa-east-1 (São Paulo)

  Azure
  [21] eastus2 (Virginia)
  [22] westus2 (Washington)
  [23] westeurope (Netherlands)
  [24] japaneast (Tokyo)
  [25] southeastasia (Singapore)

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

	// Insecure is set to true by default to allow API calls against HTTPS
	// endpoints with self-signed certificates.
	cfg.Insecure = true

	if err := askInfraSelection(&cfg, scanner, params.Writer, params.ErrWriter); err != nil {
		return err
	}

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

func askInfraSelection(cfg *Config, scanner *input.Scanner, writer, errWriter io.Writer) error {
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
		if err := askRegionSelection(cfg, scanner, writer); err != nil {
			return err
		}
	case eceInfraChoice:
		cfg.Host = scanner.Scan(eceHostMsg)
	case esspInfraChoice:
		cfg.Host = scanner.Scan(esspHostMsg)
	default:
		fmt.Fprintf(errWriter, "invalid choice, defaulting to %s", essHostAddress)
	}

	return nil
}

func askRegionSelection(cfg *Config, scanner *input.Scanner, writer io.Writer) error {
	regionChoiceRaw := scanner.Scan(regionChoiceMsg)
	fmt.Fprintln(writer)
	regionChoice, err := strconv.Atoi(regionChoiceRaw)
	if err != nil {
		return err
	}

	switch regionChoice {
	case gcpUsCentral1Choice:
		cfg.Region = "gcp-us-central1"
	case gcpUsEast4Choice:
		cfg.Region = "gcp-us-east4"
	case gcpUsWest1Choice:
		cfg.Region = "gcp-us-west1"
	case gcpNorthamericaNortheast1Choice:
		cfg.Region = "gcp-northamerica-northeast1"
	case gcpAustraliaSoutheast1Choice:
		cfg.Region = "gcp-australia-southeast1"
	case gcpEuropeWest1Choice:
		cfg.Region = "gcp-europe-west1"
	case gcpEuropeWest2Choice:
		cfg.Region = "gcp-europe-west2"
	case gcpEuropeWest3Choice:
		cfg.Region = "gcp-europe-west3"
	case gcpAsiaNortheast1Choice:
		cfg.Region = "gcp-asia-northeast1"
	case gcpAsiaSouth1Choice:
		cfg.Region = "gcp-asia-south1"
	case awsUsEast1Choice:
		cfg.Region = "us-east-1"
	case awsUsWest1Choice:
		cfg.Region = "us-west-1"
	case awsUsWest2Choice:
		cfg.Region = "us-west-2"
	case awsEuCentral1Choice:
		cfg.Region = "aws-eu-central-1"
	case awsEuWest2Choice:
		cfg.Region = "aws-eu-west-2"
	case awsEuWest1Choice:
		cfg.Region = "eu-west-1"
	case awsApNortheast1Choice:
		cfg.Region = "ap-northeast-1"
	case awsApSoutheast1Choice:
		cfg.Region = "ap-southeast-1"
	case awsApSoutheast2Choice:
		cfg.Region = "ap-southeast-2"
	case awsSaEast1Choice:
		cfg.Region = "sa-east-1"
	case azureEastUs2Choice:
		cfg.Region = "azure-eastus2"
	case azureWestUs2Choice:
		cfg.Region = "azure-westus2"
	case azureWestEuropeChoice:
		cfg.Region = "azure-westeurope"
	case azureJapanEastChoice:
		cfg.Region = "azure-japaneast"
	case azureSouthEastAsiaChoice:
		cfg.Region = "azure-southeastasia"
	default:
		return errors.New("invalid region choice")
	}

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
