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

package testutils

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

const (
	durationExpr = `duration.*\)`
	durationRepl = "duration )"
)

// Args represent are required args to test a command.
type Args struct {
	Cmd  *cobra.Command
	Args []string
	Cfg  MockCfg
}

// Assertion to use for tests.
type Assertion struct {
	Err    error
	Stdout string
	Stderr string
}

// TestCommand tests a `*cobra.Command` and uses the testing.T struct to
// return any unmatched assertions.
func TestCommand(t *testing.T, args Args, assertion Assertion) {
	cfg := fillDefaults(args.Cfg)
	cfg.Out = new(bytes.Buffer)
	cfg.Err = new(bytes.Buffer)

	defer MockApp(t, cfg)()

	// Set the command arguments.
	args.Cmd.Root().SetArgs(args.Args)

	// Necessary to reduce the clutter on the output
	args.Cmd.SilenceUsage = true
	args.Cmd.SilenceErrors = true

	// Set the Out and Err to the mocked devices.
	args.Cmd.SetOutput(cfg.Out)
	args.Cmd.SetErr(cfg.Err)

	if err := args.Cmd.Execute(); !assert.Equal(t, assertion.Err, err) {
		t.Error(err)
	}

	if buf, ok := args.Cmd.OutOrStdout().(*bytes.Buffer); ok {
		var got = buf.String()

		// When the output contains the `--track` flag, removes the non-
		// assertable duration time.
		if strings.Contains(strings.Join(args.Args, " "), "--track") {
			got = regexp.MustCompile(durationExpr).ReplaceAllString(
				got, durationRepl,
			)
		}

		if got != assertion.Stdout {
			t.Errorf(`"Got stdout "%s" != want "%s"`, got, assertion.Stdout)
		}
	}

	if buf, ok := args.Cmd.ErrOrStderr().(*bytes.Buffer); ok {
		if got := buf.String(); got != assertion.Stderr {
			t.Errorf(`"Got stderr "%s" != want "%s"`, got, assertion.Stderr)
		}
	}
}
