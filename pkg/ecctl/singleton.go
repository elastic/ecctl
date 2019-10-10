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
	"errors"
)

var application *App

// Get obtains the initialized instance of the application
func Get() *App { return application }

// Instance instantiates a new Application from a config, if it's already
// been instantiated, an error is returned.
func Instance(c Config) (*App, error) {
	if application != nil {
		return nil, errors.New("app: application already instantiated")
	}

	var err error
	application, err = NewApplication(c)

	return application, err
}
