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

package testutils

import (
	"net/http"
	"os"

	"github.com/elastic/cloud-sdk-go/pkg/output"

	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/util"
)

// MockInitApp initiates a mocked app
func MockInitApp() error {
	var config = ecctl.Config{
		Client:       new(http.Client),
		OutputDevice: output.NewDevice(os.Stdout),
		ErrorDevice:  os.Stderr,
		Output:       "json",
		Host:         "http://somehost",
		APIKey:       "helloiamakey",
	}

	return util.ReturnErrOnly(ecctl.Instance(config))
}
