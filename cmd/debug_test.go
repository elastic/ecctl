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

package cmd

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

func deleteFile(name string, t *testing.T) {
	var err error
	var info os.FileInfo

	if info, err = os.Stat(name); err != nil {
		if _, ok := err.(*os.PathError); ok {
			return
		}
		t.Fatal(err)
	}
	println(info)
	if err := os.Remove(name); err != nil {
		t.Fatal(err)
	}
}

func Test_setupDebug(t *testing.T) {
	type args struct {
		enableTrace bool
		enablePprof bool
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "setting up both returns an error",
			args: args{true, true},
			err:  errTracePprofCannotBeEnabled,
		},
		{
			name: "setting up tracing succeeds",
			args: args{true, false},
			err:  nil,
		},
		{
			name: "setting up pprof succeeds",
			args: args{false, true},
			err:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			traceFile := fmt.Sprint("trace", "-", time.Now().Format("20060102150405"), ".out")
			pprofFile := fmt.Sprint("pprof", "-", time.Now().Format("20060102150405"), ".out")
			var firstErr error
			if firstErr = setupDebug(tt.args.enableTrace, tt.args.enablePprof); !reflect.DeepEqual(firstErr, tt.err) {
				t.Errorf("setupDebug() error = %v, wantErr %v", firstErr, tt.err)
				return
			}

			if firstErr != nil {
				return
			}
			if tt.args.enableTrace {
				if _, err := os.Stat(traceFile); err != nil {
					if _, ok := err.(*os.PathError); ok {
						t.Error("setupDebug() = didn't create a trace file")
						return
					}
				}
				deleteFile(traceFile, t)
			}
			if tt.args.enablePprof {
				if _, err := os.Stat(pprofFile); err != nil {
					if _, ok := err.(*os.PathError); ok {
						t.Error("setupDebug() = didn't create a trace file")
						return
					}
				}
				deleteFile(pprofFile, t)
			}
		})
	}
}
