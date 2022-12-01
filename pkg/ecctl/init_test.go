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
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/multierror"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/spf13/viper"
)

var emptyPassFunc = func(int) ([]byte, error) { return nil, nil }

func copyFixtures(t *testing.T, fixturePath string) func() {
	var cleanupfiles []string
	err := filepath.Walk(fixturePath, func(path string, info os.FileInfo, _ error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(info.Name()) != ".orig" {
			return nil
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var newFile = strings.Replace(path, ".orig", "", 1)
		err = os.WriteFile(newFile, b, 0660)
		if err != nil {
			return err
		}

		cleanupfiles = append(cleanupfiles, newFile)

		return nil
	})
	if err != nil {
		t.Error(err)
	}

	return func() {
		for _, file := range cleanupfiles {
			func(f string) {
				if err := os.Remove(f); err != nil {
					t.Error(err)
				}
			}(file)
		}

		// Cleanup JSON files as well (Written by viper)
		jsonFiles, err := filepath.Glob(filepath.Join(fixturePath, "*.json"))
		if err != nil {
			t.Error(err)
		}
		for _, f := range jsonFiles {
			if err := os.Remove(f); err != nil {
				t.Error(err)
			}
		}
	}
}

func TestInitConfig(t *testing.T) {
	const testFiles = "./test_files"
	defer copyFixtures(t, testFiles)()

	var emptyViper = viper.New()
	emptyViper.AddConfigPath(testFiles)

	var hostConfigViper = viper.New()
	hostConfigViper.AddConfigPath(testFiles)
	hostConfigViper.SetConfigName("hostconfig")
	var hostConfigContents = `{
  "host": "https://localhost"
}`
	var apikeyConfig = viper.New()
	apikeyConfig.AddConfigPath(testFiles)
	apikeyConfig.SetConfigName("apikey")
	var apikeyConfigContents = `{
  "host": "https://localhost",
  "api_key": "[REDACTED]",
  "verbose": true
}`
	var userPassConfig = viper.New()
	userPassConfig.AddConfigPath(testFiles)
	userPassConfig.SetConfigName("userpass")
	var userPassConfigContents = `{
  "user": "marc",
  "pass": "[REDACTED]",
  "host": "https://localhost"
}`

	var emptyViperToCreateConfig = viper.New()
	emptyViperToCreateConfig.AddConfigPath(testFiles)

	var emptyViperToCreateConfigUserPass = viper.New()
	emptyViperToCreateConfigUserPass.AddConfigPath(testFiles)

	var userPassConfigToModify = viper.New()
	userPassConfigToModify.AddConfigPath(testFiles)
	userPassConfigToModify.SetConfigName("userpassmodif")
	var userPassConfigToModifyContents = `{
  "user": "marc",
  "pass": "[REDACTED]",
  "host": "https://localhost"
}`

	var apiKeyConfigToModify = viper.New()
	apiKeyConfigToModify.AddConfigPath(testFiles)
	apiKeyConfigToModify.SetConfigName("apikeymodif")
	var apiKeyConfigToModifyContents = `{
  "host": "https://localhost",
  "api_key": "[REDACTED]",
  "output": "json"
}`

	type args struct {
		params InitConfigParams
	}
	tests := []struct {
		name         string
		args         args
		err          error
		wantOutput   string
		wantSettings map[string]interface{}
	}{
		{
			name: "fail on parameter validation",
			err: multierror.NewPrefixed("invalid init configuration",
				errors.New("viper instance cannot be nil"),
				errors.New("input reader cannot be nil"),
				errors.New("output writer cannot be nil"),
				errors.New("error writer cannot be nil"),
				errors.New("password read function cannot be nil"),
				errors.New("http client cannot be nil"),
			),
		},
		{
			name: "doesn't find a config file and user skips the creation of one",
			args: args{params: InitConfigParams{
				Viper:            emptyViper,
				Reader:           strings.NewReader("n\n"),
				Writer:           new(bytes.Buffer),
				ErrWriter:        new(bytes.Buffer),
				PasswordReadFunc: emptyPassFunc,
				Client:           new(http.Client),
			}},
			wantOutput: disclaimer + missingConfigMsg,
		},
		{
			name: "finds a config file, prints the contents and user skips the creation of one",
			args: args{params: InitConfigParams{
				Viper:            hostConfigViper,
				Reader:           strings.NewReader("n\n"),
				Writer:           new(bytes.Buffer),
				ErrWriter:        new(bytes.Buffer),
				PasswordReadFunc: emptyPassFunc,
				Client:           new(http.Client),
			}},
			wantSettings: map[string]interface{}{
				"host": "https://localhost",
			},
			wantOutput: disclaimer +
				fmt.Sprintf(settingsPathMsg, "test_files/hostconfig.yaml") +
				hostConfigContents + "\n" + existingConfigMsg,
		},
		{
			name: "finds a config file, prints the contents with obscured api_key and user skips the creation of one",
			args: args{params: InitConfigParams{
				Viper:            apikeyConfig,
				Reader:           strings.NewReader("n\n"),
				Writer:           new(bytes.Buffer),
				ErrWriter:        new(bytes.Buffer),
				PasswordReadFunc: emptyPassFunc,
				Client:           new(http.Client),
			}},
			wantSettings: map[string]interface{}{
				"api_key": "someapikey",
				"host":    "https://localhost",
				"verbose": true,
			},
			wantOutput: disclaimer +
				fmt.Sprintf(settingsPathMsg, "test_files/apikey.yaml") +
				apikeyConfigContents + "\n" + existingConfigMsg,
		},
		{
			name: "finds a config file prints the contents with obscured pass and user skips the creation of one",
			args: args{params: InitConfigParams{
				Viper:            userPassConfig,
				Reader:           strings.NewReader("n\n"),
				Writer:           new(bytes.Buffer),
				ErrWriter:        new(bytes.Buffer),
				PasswordReadFunc: emptyPassFunc,
				Client:           new(http.Client),
			}},
			wantSettings: map[string]interface{}{
				"host": "https://localhost",
				"pass": "marcmarc",
				"user": "marc",
			},
			wantOutput: disclaimer +
				fmt.Sprintf(settingsPathMsg, "test_files/userpass.yaml") +
				userPassConfigContents + "\n" + existingConfigMsg,
		},
		{
			name: "doesn't find a config file and user creates a new one",
			args: args{params: InitConfigParams{
				Viper:    emptyViperToCreateConfig,
				FilePath: filepath.Join(testFiles, "newConfig"),
				Reader: io.MultiReader(
					strings.NewReader("y\n"),
					strings.NewReader("1\n"),
					strings.NewReader("1\n"),
					strings.NewReader("1\n"),
					strings.NewReader("anapikey\n"),
					strings.NewReader("1\n"),
				),
				Writer:    new(bytes.Buffer),
				ErrWriter: new(bytes.Buffer),
				PasswordReadFunc: func(int) ([]byte, error) {
					return []byte("somekey"), nil
				},
				Client: mock.NewClient(mock.New200Response(mock.NewStructBody(models.User{
					UserName: ec.String("anacleto"),
				}))),
			}},
			wantSettings: map[string]interface{}{
				"api_key":  "somekey",
				"host":     "https://api.elastic-cloud.com",
				"insecure": true,
				"output":   "text",
				"region":   "gcp-us-central1",
			},
			wantOutput: disclaimer + missingConfigMsg + hostChoiceMsg + "\n" +
				fmt.Sprintf(essChoiceMsg, essHostAddress) + regionChoiceMsg + "\n" +
				essAPIKeyCreateMsg + "\n" +
				apiKeyMsg + "\n" + formatChoiceMsg + "\n" + "\n" +
				fmt.Sprintf(validCredentialsMsg, "anacleto") + finalMsg + "\n",
		},
		{
			name: "doesn't find a config file and user creates a new one with user/pass",
			args: args{params: InitConfigParams{
				Viper:    emptyViperToCreateConfigUserPass,
				FilePath: filepath.Join(testFiles, "newConfigUserPass"),
				Reader: io.MultiReader(
					strings.NewReader("y\n"),
					strings.NewReader("2\n"),
					strings.NewReader("https://ahost\n"),
					strings.NewReader("2\n"),
					strings.NewReader("auser\n"),
					strings.NewReader("1\n"),
				),
				Writer:    new(bytes.Buffer),
				ErrWriter: new(bytes.Buffer),
				PasswordReadFunc: func(int) ([]byte, error) {
					return []byte("apassword"), nil
				},
				Client: mock.NewClient(
					mock.New200Response(mock.NewStructBody(models.TokenResponse{
						Token: ec.String("atoken"),
					})),
					mock.New200Response(mock.NewStructBody(models.User{
						UserName: ec.String("auser"),
					})),
				),
			}},
			wantSettings: map[string]interface{}{
				"host":     "https://ahost",
				"insecure": true,
				"output":   "text",
				"pass":     "apassword",
				"user":     "auser",
			},
			wantOutput: disclaimer + missingConfigMsg + hostChoiceMsg + "\n" + eceHostMsg +
				authChoiceMsg + "\n" + userMsg + passMsg + "\n" + formatChoiceMsg +
				"\n" + "\n" + fmt.Sprintf(validCredentialsMsg, "auser") + finalMsg + "\n",
		},
		{
			name: "doesn't find a config file and user creates a new one with user/pass, GET user fails, but deployment list succeeds",
			args: args{params: InitConfigParams{
				Viper:    emptyViperToCreateConfigUserPass,
				FilePath: filepath.Join(testFiles, "newConfigUserPass"),
				Reader: io.MultiReader(
					strings.NewReader("y\n"),
					strings.NewReader("2\n"),
					strings.NewReader("https://ahost\n"),
					strings.NewReader("2\n"),
					strings.NewReader("auser\n"),
					strings.NewReader("1\n"),
				),
				Writer:    new(bytes.Buffer),
				ErrWriter: new(bytes.Buffer),
				PasswordReadFunc: func(int) ([]byte, error) {
					return []byte("apassword"), nil
				},
				Client: mock.NewClient(
					mock.New200Response(mock.NewStructBody(models.TokenResponse{
						Token: ec.String("atoken"),
					})),
					mock.New404Response(mock.NewStructBody(models.User{
						UserName: ec.String("auser"),
					})),
					mock.New200Response(mock.NewStructBody(models.DeploymentsListResponse{})),
				),
			}},
			wantSettings: map[string]interface{}{
				"host":     "https://ahost",
				"insecure": true,
				"output":   "text",
				"pass":     "apassword",
				"user":     "auser",
			},
			wantOutput: disclaimer + missingConfigMsg + hostChoiceMsg + "\n" + eceHostMsg +
				authChoiceMsg + "\n" + userMsg + passMsg + "\n" + formatChoiceMsg +
				"\n" + "\n" + validCredentialsAlternativeMsg + finalMsg + "\n",
		},
		{
			name: "doesn't find a config file and user creates a new one with user/pass, and returns error on API test",
			args: args{params: InitConfigParams{
				Viper:    emptyViperToCreateConfigUserPass,
				FilePath: filepath.Join(testFiles, "newConfigUserPass"),
				Reader: io.MultiReader(
					strings.NewReader("y\n"),
					strings.NewReader("2\n"),
					strings.NewReader("https://ahost\n"),
					strings.NewReader("2\n"),
					strings.NewReader("auser\n"),
					strings.NewReader("1\n"),
				),
				Writer:    new(bytes.Buffer),
				ErrWriter: new(bytes.Buffer),
				PasswordReadFunc: func(int) ([]byte, error) {
					return []byte("apassword"), nil
				},
				Client: mock.NewClient(
					mock.New200Response(mock.NewStructBody(models.TokenResponse{
						Token: ec.String("atoken"),
					})),
					mock.New404Response(mock.NewStructBody(models.User{
						UserName: ec.String("auser"),
					})),
					mock.New404Response(mock.NewStructBody(models.DeploymentsListResponse{})),
				),
			}},
			wantSettings: map[string]interface{}{
				"host":     "https://ahost",
				"insecure": true,
				"output":   "text",
				"pass":     "apassword",
				"user":     "auser",
			},
			err: errors.New(invalidCredentialsMsg),
			wantOutput: disclaimer + missingConfigMsg + hostChoiceMsg + "\n" + eceHostMsg +
				authChoiceMsg + "\n" + userMsg + passMsg + "\n" + formatChoiceMsg +
				"\n" + "\n",
		},
		{
			name: "finds a config file and user changes the values",
			args: args{params: InitConfigParams{
				Viper:    userPassConfigToModify,
				FilePath: filepath.Join(testFiles, "doesnt_matter"),
				Reader: io.MultiReader(
					strings.NewReader("y\n"),
					strings.NewReader("3\n"),
					strings.NewReader("https://ahost\n"),
					strings.NewReader("1\n"),
					strings.NewReader("anapikey\n"),
					strings.NewReader("1\n"),
				),
				Writer:    new(bytes.Buffer),
				ErrWriter: new(bytes.Buffer),
				PasswordReadFunc: func(int) ([]byte, error) {
					return []byte("somekey"), nil
				},
				Client: mock.NewClient(mock.New200Response(mock.NewStructBody(models.User{
					UserName: ec.String("anacleto"),
				}))),
			}},
			wantSettings: map[string]interface{}{
				"api_key":  "somekey",
				"host":     "https://ahost",
				"insecure": true,
				"output":   "text",
				"region":   "us-west-2",
			},
			wantOutput: disclaimer +
				fmt.Sprintf(settingsPathMsg, "test_files/userpassmodif.yaml") +
				userPassConfigToModifyContents + "\n" + existingConfigMsg + hostChoiceMsg +
				"\n" + esspHostMsg + apiKeyMsg + "\n" + formatChoiceMsg +
				"\n" + "\n" + fmt.Sprintf(validCredentialsMsg, "anacleto") + finalMsg + "\n",
		},
		{
			name: "finds a config file and user changes the values, from api_key to user/pass",
			args: args{params: InitConfigParams{
				Viper:    apiKeyConfigToModify,
				FilePath: filepath.Join(testFiles, "doesnt_matter"),
				Reader: io.MultiReader(
					strings.NewReader("y\n"),
					strings.NewReader("2\n"),
					strings.NewReader("https://ahost\n"),
					strings.NewReader("2\n"),
					strings.NewReader("auser\n"),
					strings.NewReader("1\n"),
				),
				Writer:    new(bytes.Buffer),
				ErrWriter: new(bytes.Buffer),
				PasswordReadFunc: func(int) ([]byte, error) {
					return []byte("apassword"), nil
				},
				Client: mock.NewClient(
					mock.New200Response(mock.NewStructBody(models.TokenResponse{
						Token: ec.String("atoken"),
					})),
					mock.New200Response(mock.NewStructBody(models.User{
						UserName: ec.String("auser"),
					})),
				),
			}},
			wantSettings: map[string]interface{}{
				"host":     "https://ahost",
				"insecure": true,
				"output":   "text",
				"pass":     "apassword",
				"user":     "auser",
			},
			wantOutput: disclaimer +
				fmt.Sprintf(settingsPathMsg, "test_files/apikeymodif.yaml") +
				apiKeyConfigToModifyContents + "\n" + existingConfigMsg + hostChoiceMsg +
				"\n" + eceHostMsg + authChoiceMsg + "\n" + userMsg +
				passMsg + "\n" + formatChoiceMsg + "\n" + "\n" + fmt.Sprintf(validCredentialsMsg, "auser") + finalMsg + "\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantSettings == nil {
				tt.wantSettings = make(map[string]interface{})
			}
			if err := InitConfig(tt.args.params); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("InitConfig() error = %v, wantErr %v", err, tt.err)
			}

			if buf, ok := tt.args.params.Writer.(*bytes.Buffer); ok {
				p, err := filepath.Abs(".")
				if err != nil {
					t.Error(err)
				}

				got := strings.Replace(buf.String(), p+string(os.PathSeparator), "", 1)
				if got != tt.wantOutput {
					t.Errorf("InitConfig() output = %v, wantOutput %v", got, tt.wantOutput)
				}

				settings := tt.args.params.Viper.AllSettings()
				if !reflect.DeepEqual(tt.wantSettings, settings) {
					t.Errorf("InitConfig() settings = %v, wantSettings %v", settings, tt.wantSettings)
				}
			}
		})
	}
}
