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

package snaprepo

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
)

var (
	errConfigMustNotBeEmpty = errors.New("config must not be empty")
)

// ParseGenericConfig reads the contents of an io.Reader and tries to parse its
// contents as YAML or JSON, returns an error if parsing fails in both formats.
func ParseGenericConfig(input io.Reader) (GenericConfig, error) {
	var config GenericConfig
	if input == nil {
		return config, errReaderCannotBeNil
	}

	var buf = new(bytes.Buffer)
	if _, err := buf.ReadFrom(input); err != nil {
		return config, err
	}

	if err := yaml.Unmarshal(buf.Bytes(), &config); err == nil {
		return config, nil
	}

	if err := json.Unmarshal(buf.Bytes(), &config); err != nil {
		return config, errors.Wrap(err, errFailedParsingConfig.Error())
	}
	return config, nil
}

// GenericConfig wraps a map[string]interface{} type so it implements the
// common.Validator interface.
type GenericConfig map[string]interface{}

// Validate checks that the length of the map is >= 1
func (c GenericConfig) Validate() error {
	if len(c) < 1 {
		return errConfigMustNotBeEmpty
	}
	return nil
}
