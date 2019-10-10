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

package cmdutil

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/elastic/ecctl/pkg/ecctl"
)

// GetInsecurePassword retrieves an insecure password from a CLI command
func GetInsecurePassword(insecure string) ([]byte, error) {
	if insecure == "" {
		return nil, errors.New("a password must be provided when using the --insecure-password flag")
	}

	fmt.Fprintln(os.Stderr, "WARNING: --insecure-password is not recommended")
	return []byte(insecure), nil
}

// PasswordVerify retrieves a password from terminal input and verifies a match
func PasswordVerify(password []byte) error {
	var writer = ecctl.Get().Config.OutputDevice
	p, err := ecctl.ReadSecret(writer, ecctl.DefaultPassFunc, "re-enter password: ")
	if err != nil {
		return err
	}

	fmt.Println()
	if !reflect.DeepEqual(password, p) {
		return errors.New("passwords do not match, action has been aborted")
	}

	return nil
}

// InsecureOrSecurePassword checks if an insecure password has been set and
// gets a password in a secure or insecure way.
func InsecureOrSecurePassword(insecure, message string, verify bool) ([]byte, error) {
	if insecure != "" {
		return GetInsecurePassword(insecure)
	}

	var writer = ecctl.Get().Config.OutputDevice
	password, err := ecctl.ReadSecret(writer, ecctl.DefaultPassFunc, message)
	if !verify {
		return password, err
	}

	return password, PasswordVerify(password)
}
