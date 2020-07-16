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
	"fmt"
)

// PlatformAdminRequired is an additional helper text for commands
const PlatformAdminRequired = "(Available for ECE only)"

// DeprecatedText is an additional helper text for commands
const DeprecatedText = "DEPRECATED (Will be removed in the next major version):"

// AdminReqDescription adds a text about required admin permissions to a string
func AdminReqDescription(desc string) string {
	return fmt.Sprintf("%s %v", desc, PlatformAdminRequired)
}

// DeprecatedDescription adds a text about deprecation to a string
func DeprecatedDescription(desc string) string {
	return fmt.Sprintf("%v %s", DeprecatedText, desc)
}
