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
	"fmt"
	"io"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	// DefaultPassFunc can be used to consume a password from a file descriptor.
	// References terminal.ReadPassword which obscures the user input.
	DefaultPassFunc = terminal.ReadPassword
)

// ReadSecret obtains a secret by reading the PassFunc passing the stdin file
// descriptor. If passFunc is empty, it defaults to terminal.
func ReadSecret(w io.Writer, passFunc PassFunc, msg string) ([]byte, error) {
	fmt.Fprint(w, msg)
	b, err := passFunc(syscall.Stdin)
	fmt.Fprintln(w)
	if err != nil {
		return nil, err
	}

	return b, nil
}
