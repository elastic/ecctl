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
	"io"
	"os"
	"runtime/pprof"
	"runtime/trace"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var errTracePprofCannotBeEnabled = errors.New("only one of trace or pprof can be enabled at a time")

// setupDebug either enables tracing or pprofing to a file
func setupDebug(enableTrace, enablePprof bool) error {
	if enableTrace && enablePprof {
		return errTracePprofCannotBeEnabled
	}
	if enableTrace {
		f, err := createDebugFile("trace")
		if err != nil {
			return err
		}
		return trace.Start(f)
	}

	if enablePprof {
		f, err := createDebugFile("pprof")
		if err != nil {
			return err
		}
		return pprof.StartCPUProfile(f)
	}

	return nil
}

func createDebugFile(name string) (io.WriteCloser, error) {
	f, err := os.Create(
		fmt.Sprint(name, "-", time.Now().Format("20060102150405"), ".out"))
	if err != nil {
		return nil, err
	}
	return f, nil
}

func stopDebug(v *viper.Viper) {
	if v.GetBool("trace") {
		trace.Stop()
	}
	if v.GetBool("pprof") {
		pprof.StopCPUProfile()
	}
}
