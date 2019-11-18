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
	"runtime"
	"testing"
)

const vTemplate = `Version:               1.0.0
Client API Version:    2.4.2
Go version:            %s
Git commit:            12345678
Built:                 Fri 15 Nov 09:24:14 2019
OS/Arch:               %s / %s
`

func TestVersionInfo_String(t *testing.T) {
	tests := []struct {
		name   string
		fields VersionInfo
		want   string
	}{
		{
			name: "Prints the version",
			fields: VersionInfo{
				Version:    "1.0.0",
				APIVersion: "2.4.2",
				Commit:     "12345678910",
				Built:      "Fri_15_Nov_09:24:14_2019",
			},
			want: fmt.Sprintf(vTemplate, runtime.Version(), runtime.GOOS, runtime.GOARCH),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.fields
			if got := v.String(); got != tt.want {
				t.Errorf("VersionInfo.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
