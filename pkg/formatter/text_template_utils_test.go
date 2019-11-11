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
	"reflect"
	"testing"

	"github.com/elastic/ecctl/pkg/platform/snaprepo"
)

func Test_rpadTrim(t *testing.T) {
	type args struct {
		s string
		n int
	}
	tests := []struct {
		name  string
		args  args
		want  string
		panic interface{}
	}{
		{
			name: "rpadTrim returns a trimmed & padded string",
			args: args{"mystring", 5},
			want: "mys  ",
		},
		{
			name: "rpadTrim returns a trimmed & padded string",
			args: args{"my", 3},
			want: "m  ",
		},
		{
			name: "rpadTrim returns a padded string",
			args: args{"my", 10},
			want: "my        ",
		},
		{
			name:  "rpadTrim panics on too short of a string",
			args:  args{"myaa", 2},
			want:  "m  ",
			panic: "padding is too small, needs to be at least 3",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				panicText := recover()
				if !reflect.DeepEqual(panicText, tt.panic) {
					t.Errorf("rpadTrim panic() = %v, want %v", panicText, tt.panic)
				}
			}()
			if got := rpadTrim(tt.args.s, tt.args.n); got != tt.want {
				t.Errorf("rpadTrim() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_toS3TypeConfig(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name string
		args args
		want snaprepo.S3TypeConfig
	}{
		{
			name: "Parses an interface that complies",
			args: args{
				map[string]interface{}{
					"type": "s3",
					"settings": map[string]interface{}{
						"region": "us-east-1",
					},
				},
			},
			want: snaprepo.S3TypeConfig{
				Type: "s3",
				Settings: snaprepo.S3Config{
					Region: "us-east-1",
				},
			},
		},
		{
			name: "Parses an interface that complies",
			args: args{
				"omg what is this thing",
			},
			want: snaprepo.S3TypeConfig{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toS3TypeConfig(tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toS3TypeConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_equal(t *testing.T) {
	type args struct {
		x, y interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "compare equal returns false",
			args: args{x: "1", y: "2"},
			want: false,
		},
		{
			name: "compare equal returns true",
			args: args{x: "1", y: "1"},
			want: true,
		},
		{
			name: "compare equal returns false on other types",
			args: args{x: "1", y: 234},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := equal(tt.args.x, tt.args.y); got != tt.want {
				t.Errorf("equal() = %v, want %v", got, tt.want)
			}
		})
	}
}
