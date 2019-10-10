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

package util

import (
	"testing"
)

func TestUnderscoreToDashes(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "String uses underscore",
			args: args{"example_string"},
			want: "example-string",
		},
		{
			name: "String is camelcase",
			args: args{"ExampleString"},
			want: "example-string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UnderscoreToDashes(tt.args.str); got != tt.want {
				t.Errorf("UnderscoreToDashes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDashesToUnderscore(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "String uses dash",
			args: args{"example-string"},
			want: "example_string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DashesToUnderscore(tt.args.str); got != tt.want {
				t.Errorf("DashesToUnderscore() = %v, want %v", got, tt.want)
			}
		})
	}
}
