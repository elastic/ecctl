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
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	// ErrEmptyStdinAndFile thrown when either the Sdin or the file are empty.
	ErrEmptyStdinAndFile = errors.New(
		"empty stdin and file definition, need one of the two to be populated",
	)
)

// FileOrStdin returns an error when the followong scenarios happen:
// * No Stdin and Flag value are empty.
// * Unredable Stdin
// * Both Stdin and Flag value present.
func FileOrStdin(cmd *cobra.Command, name string) error {
	var stat, err = os.Stdin.Stat()
	var filenameProvided = cmd.Flag(name).Value.String() != ""
	var unreadableStdin = err != nil
	var stdinProvided = (stat.Mode() & os.ModeCharDevice) == 0

	if !filenameProvided && unreadableStdin {
		return errors.New("unreadable stdin and empty file definition")
	}

	if !filenameProvided && !stdinProvided {
		return ErrEmptyStdinAndFile
	}

	if filenameProvided && stdinProvided {
		return errors.New("non empty stdin and file definition, need one of the two to be populated")
	}

	return nil
}

// ParseBoolP parses a string flag that is meant to be a boolean with 3 values:
// * true
// * false
// * nil
func ParseBoolP(cmd *cobra.Command, name string) (*bool, error) {
	flagRaw, err := cmd.Flags().GetString(name)
	if err != nil {
		return nil, err
	}

	var flagVal *bool
	ss, err := strconv.ParseBool(flagRaw)
	if err != nil && flagRaw != "" {
		return nil, err
	}

	if flagRaw != "" && err == nil {
		return &ss, nil
	}

	return flagVal, nil
}
