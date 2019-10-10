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
	"io"
	"strings"

	"github.com/elastic/cloud-sdk-go/pkg/input"
	"github.com/elastic/cloud-sdk-go/pkg/output"
)

// ConfirmAction is used to print a message to the writer and read the input
// from the reader. If the first letter is "y" or "Y" it will return true.
func ConfirmAction(msg string, reader io.Reader, writer io.Writer) bool {
	var device = output.NewDevice(writer)
	if t, ok := writer.(*output.Device); ok {
		device = t
	}
	var s = input.NewScanner(reader, device)
	return strings.HasPrefix(strings.ToLower(s.Scan(msg)), "y")
}
