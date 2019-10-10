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
	"os"
	"testing"
)

func Test_GetHomePath(t *testing.T) {
	type args struct {
		goos string
	}
	type fields struct {
		home        string
		homedrive   string
		homepath    string
		userprofile string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "correct substitution on darwin",
			fields: fields{
				home: "/home/myuser",
			},
			args: args{
				goos: "darwin",
			},
			want: "/home/myuser",
		},
		{
			name: "correct substitution on linux",
			fields: fields{
				home: "/home/myuser",
			},
			args: args{
				goos: "linux",
			},
			want: "/home/myuser",
		},
		{
			name: "correct substitution on windows with homedrive & homepath",
			fields: fields{
				homedrive: "C:",
				homepath:  `\Users\Myuser`,
			},
			args: args{
				goos: "windows",
			},
			want: `C:\Users\Myuser`,
		},
		{
			name: "correct substitution on windows with userprofile",
			fields: fields{
				userprofile: `C:\Users\Myuser`,
			},
			args: args{
				goos: "windows",
			},
			want: `C:\Users\Myuser`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			if tt.fields.home != "" {
				var prev = os.Getenv("HOME")
				defer os.Setenv("HOME", prev)
				os.Setenv("HOME", tt.fields.home)
			}
			if tt.fields.homedrive != "" {
				var prev = os.Getenv("HOMEDRIVE")
				defer os.Setenv("HOMEDRIVE", prev)
				os.Setenv("HOMEDRIVE", tt.fields.homedrive)
			}
			if tt.fields.homepath != "" {
				var prev = os.Getenv("HOMEPATH")
				defer os.Setenv("HOMEPATH", prev)
				os.Setenv("HOMEPATH", tt.fields.homepath)
			}
			if tt.fields.userprofile != "" {
				var prev = os.Getenv("USERPROFILE")
				defer os.Setenv("USERPROFILE", prev)
				os.Setenv("USERPROFILE", tt.fields.userprofile)
			}

			// Testcode
			if got := GetHomePath(tt.args.goos); got != tt.want {
				t.Errorf("getHomePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
