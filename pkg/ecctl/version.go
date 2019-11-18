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
	"bytes"
	"fmt"
	"runtime"
	"strings"
	"text/tabwriter"
)

// VersionInfo contains detailed information about the binary.
type VersionInfo struct {
	Version    string
	APIVersion string
	Commit     string
	Built      string

	// Non Displayed fields
	Repository   string
	Organization string
}

func (v VersionInfo) String() string {
	buf := new(bytes.Buffer)

	w := tabwriter.NewWriter(buf, 2, 2, 3, ' ', 0)

	var commit = v.Commit
	if len(commit) >= 8 {
		commit = commit[:8]
	}

	fmt.Fprintln(w, "Version:\t", v.Version)
	fmt.Fprintln(w, "Client API Version:\t", v.APIVersion)
	fmt.Fprintln(w, "Go version:\t", runtime.Version())
	fmt.Fprintln(w, "Git commit:\t", commit)
	fmt.Fprintln(w, "Built:\t", strings.ReplaceAll(v.Built, "_", " "))
	fmt.Fprintln(w, "OS/Arch:\t", runtime.GOOS, "/", runtime.GOARCH)

	w.Flush()

	return buf.String()
}
