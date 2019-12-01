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
	"encoding/json"
	"errors"
	"os"
	"path"
	"reflect"

	"github.com/elastic/cloud-sdk-go/pkg/input"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/spf13/cobra"
)

var (
	// ErrNodefinitionLoaded is returned by DecodeDefinition when no reader
	// has been returned from either File or Stdin.
	ErrNodefinitionLoaded = errors.New("failed obtaining a reader from file or stdin")
)

// DecodeDefinition takes a cobra command, a flagname and the desired structure
// on which to decode the contents of either the os.Stdin or the file contents.
// If both are empty, an error is returned.
func DecodeDefinition(cmd *cobra.Command, flagname string, v interface{}) error {
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return errors.New("decode file: passed structure is not a pointer")
	}

	reader, _ := input.NewFileOrReader(os.Stdin, cmd.Flag("file").Value.String())
	if reader != nil {
		if err := json.NewDecoder(reader).Decode(&v); err != nil {
			return err
		}
		return nil
	}

	return ErrNodefinitionLoaded
}

// DecodeFile takes a filename and the pointer to a structure, opening the file
// and dumping the contents into the desired structure. Make sure a pointer is
// passed rather than the copy of a structure.
func DecodeFile(filename string, v interface{}) error {
	if reflect.ValueOf(v).Kind() != reflect.Ptr {
		return errors.New("decode file: passed structure is not a pointer")
	}

	f, err := OpenFile(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(&v)
}

// ParseQueryDSLFile parses a file that contains a query dsl json and returns the corresponding binary representation
func ParseQueryDSLFile(f string) (models.SearchRequest, error) {
	var sr models.SearchRequest

	if ext := path.Ext(f); ext != ".json" {
		return sr, errors.New("not a supported file type: only json files are currently supported")
	}

	err := DecodeFile(f, &sr)

	return sr, err
}
