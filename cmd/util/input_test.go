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
	"errors"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/spf13/cobra"
)

func TestFileOrStdin(t *testing.T) {
	cmdWithFileFlag := &cobra.Command{
		Use: "something",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	cmdWithFileFlag.Flags().String("file", "", "something")
	cmdWithFileFlagPopulated := &cobra.Command{
		Use: "something",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	cmdWithFileFlagPopulated.Flags().String("file", "somevalue", "something")
	type args struct {
		cmd      *cobra.Command
		name     string
		stdin    *os.File
		create   bool
		populate bool
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Returns an error when stdin and filename is empty",
			args: args{
				cmd:  cmdWithFileFlag,
				name: "file",
			},
			err: errors.New("empty stdin and file definition, need one of the two to be populated"),
		},
		{
			name: "Returns an error when stdin and file flag are populated",
			args: args{
				cmd:      cmdWithFileFlagPopulated,
				name:     "file",
				create:   true,
				populate: true,
			},
			err: errors.New("non empty stdin and file definition, need one of the two to be populated"),
		},
		{
			name: "Returns no error when stdin is populated",
			args: args{
				cmd:      cmdWithFileFlag,
				name:     "file",
				create:   true,
				populate: true,
			},
		},
		{
			name: "Returns no error when filename is populated",
			args: args{
				cmd:  cmdWithFileFlagPopulated,
				name: "file",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.create {
				f, deleteFile := createTempFile(t, tt.args.name)
				if tt.args.populate {
					var oldstdin = os.Stdin
					os.Stdin = f
					defer func() { os.Stdin = oldstdin }()
				}
				defer deleteFile()
				tt.args.stdin = f
			}
			if err := FileOrStdin(tt.args.cmd, tt.args.name); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("FileOrStdin() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}

func TestParseBoolP(t *testing.T) {
	cmdWithEmptyBoolFlag := &cobra.Command{
		Use: "something",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	cmdWithEmptyBoolFlag.Flags().String("boolean", "", "something")

	cmdWithTrueBoolFlag := &cobra.Command{
		Use: "something",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	cmdWithTrueBoolFlag.Flags().String("boolean", "true", "something")

	cmdWithFalseBoolFlag := &cobra.Command{
		Use: "something",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	cmdWithFalseBoolFlag.Flags().String("boolean", "false", "something")

	cmdWithFauxBoolFlagValue := &cobra.Command{
		Use: "something",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	cmdWithFauxBoolFlagValue.Flags().String("boolean", "faux", "something")

	type args struct {
		cmd  *cobra.Command
		name string
	}
	tests := []struct {
		name string
		args args
		want *bool
		err  error
	}{
		{
			name: "empty string means nil bool",
			args: args{cmd: cmdWithEmptyBoolFlag, name: "boolean"},
		},
		{
			name: "fails to obtain flag returns an error",
			args: args{cmd: cmdWithEmptyBoolFlag, name: "unexisting"},
			err:  errors.New("flag accessed but not defined: unexisting"),
		},
		{
			name: "fails to parse bool flag returns an error",
			args: args{cmd: cmdWithFauxBoolFlagValue, name: "boolean"},
			err: &strconv.NumError{
				Func: "ParseBool",
				Num:  "faux",
				Err:  strconv.ErrSyntax,
			},
		},
		{
			name: "true string means true bool",
			args: args{cmd: cmdWithTrueBoolFlag, name: "boolean"},
			want: ec.Bool(true),
		},
		{
			name: "false string means false bool",
			args: args{cmd: cmdWithFalseBoolFlag, name: "boolean"},
			want: ec.Bool(false),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseBoolP(tt.args.cmd, tt.args.name)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("ParseBoolP() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseBoolP() = %v, want %v", got, tt.want)
			}
		})
	}
}
