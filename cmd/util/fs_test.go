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
	"io/ioutil"
	"os"
	"reflect"
	"syscall"
	"testing"
)

//nolint
func createTempFile(t *testing.T, name string) (*os.File, func()) {
	f, err := ioutil.TempFile("", name)
	if err != nil {
		t.Fatal(err)
	}
	return f, func() { os.RemoveAll(f.Name()) }
}

func TestOpenFile(t *testing.T) {
	type args struct {
		name     string
		filename string
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Openfile succeeds",
			args: args{
				name: "something",
			},
		},
		{
			name: "Openfile fails on unexisting file (stat)",
			args: args{
				name:     "something",
				filename: "something_unexisting",
			},
			err: &os.PathError{
				Op:   "stat",
				Path: "something_unexisting",
				Err:  syscall.ENOENT,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, cleanup := createTempFile(t, tt.args.name)
			defer cleanup()
			var name = tt.args.filename
			if name == "" {
				name = f.Name()
			}
			if _, err := OpenFile(name); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("OpenFile() error = %v, wantErr %v", err, tt.err)
				return
			}
		})
	}
}
