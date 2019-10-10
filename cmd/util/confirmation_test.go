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
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestConfirmAction(t *testing.T) {
	type args struct {
		msg    string
		reader io.Reader
	}
	tests := []struct {
		name       string
		args       args
		want       bool
		wantWriter string
	}{
		{
			name:       "returns true if y is passed",
			args:       args{msg: "are you sure? ", reader: strings.NewReader("y")},
			want:       true,
			wantWriter: "are you sure? ",
		},
		{
			name:       "returns true if Y is passed",
			args:       args{msg: "are you sure? ", reader: strings.NewReader("Y")},
			want:       true,
			wantWriter: "are you sure? ",
		},
		{
			name:       "returns true if Yes is passed",
			args:       args{msg: "are you sure? ", reader: strings.NewReader("Yes")},
			want:       true,
			wantWriter: "are you sure? ",
		},
		{
			name:       "returns false if n is passed",
			args:       args{msg: "are you sure? ", reader: strings.NewReader("n")},
			want:       false,
			wantWriter: "are you sure? ",
		},
		{
			name:       "returns false if N is passed",
			args:       args{msg: "are you sure? ", reader: strings.NewReader("N")},
			want:       false,
			wantWriter: "are you sure? ",
		},
		{
			name:       "returns false if Nope is passed",
			args:       args{msg: "are you sure? ", reader: strings.NewReader("Nope")},
			want:       false,
			wantWriter: "are you sure? ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			if got := ConfirmAction(tt.args.msg, tt.args.reader, writer); got != tt.want {
				t.Errorf("ConfirmAction() = %v, want %v", got, tt.want)
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("ConfirmAction() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}
