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

package formatter

import (
	"encoding/json"
	"fmt"
	"io"
)

// NewJSON acts as the factory for formatter.JSON
func NewJSON(output io.Writer) *JSON {
	return &JSON{output}
}

// JSON formats into text
type JSON struct {
	// Where formatter.JSON will output the result
	o io.Writer
}

// format formats the data according to the template
// using json.MarshalIndent
func (f *JSON) format(data interface{}) error {
	r, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(f.o, string(r))
	return nil
}

// Name obtains the name of the formatter
func (f *JSON) Name() string { return "json" }

// Format is used from the cmd for conveniency
// it receives a path and the data to be formatted.
//
// It currently doesn't have any effect on the JSON
// formatting due to `MarshalIndent` being used
func (f *JSON) Format(path string, data interface{}) error {
	return f.format(data)
}
