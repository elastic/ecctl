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
	"bytes"
	"fmt"
	"testing"
)

var bufferRef = &bytes.Buffer{}

func TestNew(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name       string
		args       args
		wantParent string
		wantChild  string
	}{
		{
			`New receives "text" as a parameter, returns a chain with two formatters ("text" and "json")`,
			args{
				"text",
			},
			"text",
			"json",
		},
		{
			`New receives "something" as a parameter returns a single formatter in the chain`,
			args{
				"something",
			},
			"json",
			"",
		},
		{
			`New receives "json" as a parameter returns a single formatter in the chain`,
			args{
				"json",
			},
			"json",
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := bufferRef
			got := New(o, tt.args.name)
			if tt.wantParent != "" && got.Name() != tt.wantParent {
				t.Errorf("got parent chain = %s, want = %s", got.Name(), tt.wantParent)
			}
			if tt.wantChild != "" && got.next.Name() != tt.wantChild {
				t.Errorf("got child chain = %s, want = %s", got.next.Name(), tt.wantChild)
			}
		})
	}
}

type mockFormatter struct {
	err error
}

func (m *mockFormatter) Format(string, interface{}) error { return m.err }

func (m *mockFormatter) Name() string { return "mockFormatter" }

func TestChainFormatter_Format(t *testing.T) {
	type fields struct {
		next      Chainer
		formatter Formatter
	}
	type args struct {
		path string
		data interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"In a single chain: Format returns error = nil, doesn't go down the chain",
			fields{
				formatter: &mockFormatter{nil},
			},
			args{
				path: "a",
				data: "b",
			},
			false,
		},
		{
			"In a single chain: Format returns an error when the last member of the chain returns error",
			fields{
				formatter: &mockFormatter{
					fmt.Errorf("error"),
				},
			},
			args{
				path: "a",
				data: "b",
			},
			true,
		},
		{
			"In a two formatter chain: First formatter returns an error, but the next returns error = nil",
			fields{
				formatter: &mockFormatter{
					fmt.Errorf("error"),
				},
				next: NewChain(&mockFormatter{nil}),
			},
			args{
				path: "a",
				data: "b",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Chain{
				next:      tt.fields.next,
				formatter: tt.fields.formatter,
			}
			if err := f.Format(tt.args.path, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("ChainFormatter.Format() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
