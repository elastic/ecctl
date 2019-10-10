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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	cmddeployment "github.com/elastic/ecctl/cmd/deployment"
	cmdelasticsearch "github.com/elastic/ecctl/cmd/deployment/elasticsearch"
	"github.com/elastic/ecctl/pkg/ecctl"
)

func TestPopulateValidArgs(t *testing.T) {
	type args struct {
		cmd   *cobra.Command
		cmds  []*cobra.Command
		cmdss []*cobra.Command
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "2 levels",
			args: args{
				cmd: &cobra.Command{
					Use: "ecctl",
				},
				cmds: []*cobra.Command{
					{
						Use: "command",
					},
				},
			},
		},
		{
			name: "3 levels",
			args: args{
				cmd: &cobra.Command{
					Use: "ecctl",
				},
				cmds: []*cobra.Command{
					{
						Use: "command",
						// Need to add that to make it an available subcommand
						Run: func(cmd *cobra.Command, args []string) {},
					},
				},
				cmdss: []*cobra.Command{
					{
						Use: "target",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, m := range tt.args.cmds {
				if tt.args.cmdss != nil {
					m.AddCommand(tt.args.cmdss...)
				}
			}

			tt.args.cmd.AddCommand(tt.args.cmds...)
			populateValidArgs(tt.args.cmd)

			if tt.args.cmd.ValidArgs == nil && tt.args.cmds != nil {
				t.Error("failed populating 1st level of valid args")
			}

			for _, m := range tt.args.cmds {
				fmt.Println(m.Name(), m.ValidArgs)
				if m.ValidArgs == nil && tt.args.cmdss != nil {
					t.Error("failed populating 2nd level of valid args")
				}
			}
		})
	}
}

func TestInitApp(t *testing.T) {
	type args struct {
		cmd    *cobra.Command
		client *http.Client
		config *ecctl.Config
		v      *viper.Viper
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "version command skips ecctl.Get() bootstrapping",
			args: args{
				cmd: &cobra.Command{
					Use: "version",
				},
			},
			err: nil,
		},
		{
			name: "generate command skips ecctl.Get() bootstrapping",
			args: args{
				cmd: &cobra.Command{
					Use: "help",
				},
			},
			err: nil,
		},
		{
			name: "generate command skips ecctl.Get() bootstrapping",
			args: args{
				cmd: &cobra.Command{
					Use: "generate",
				},
			},
			err: nil,
		},
		{
			name: "generate command skips ecctl.Get() bootstrapping",
			args: args{
				cmd: &cobra.Command{
					Use: "generate",
				},
			},
			err: nil,
		},
		{
			name: "fails due to empty http client",
			args: args{
				cmd: RootCmd,
			},
			err: errors.New("cmd: root http client cannot be nil"),
		},
		{
			name: "initialises rootCmd app with APIKey",
			args: args{
				cmd:    RootCmd,
				client: new(http.Client),
				config: &ecctl.Config{
					Output: "json",
					Host:   "http://localhost",
					APIKey: "some",
				},
				v: viper.New(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.config != nil {
				c := new(bytes.Buffer)
				if err := json.NewEncoder(c).Encode(tt.args.config); err != nil {
					t.Error(err)
				}

				tt.args.v.SetConfigType("json")
				if err := tt.args.v.ReadConfig(c); err != nil {
					t.Error(err)
				}
			}

			if err := initApp(tt.args.cmd, tt.args.client, tt.args.v); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("initApp() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}

func TestCheckPreRunE(t *testing.T) {
	// Test setup
	var simpleCommand = &cobra.Command{
		Use: "ecctl",
	}
	simpleCommand.AddCommand(&cobra.Command{
		Use: "dummy",
	})
	var complexCommand = &cobra.Command{
		Use: "ecctl",
	}
	var complexSubCommand = &cobra.Command{
		Use: "dummy",
	}
	complexSubCommand.AddCommand(&cobra.Command{
		Use: "complex",
	})
	complexCommand.AddCommand(complexSubCommand)

	// TC declaration
	type args struct {
		command *cobra.Command
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "dummy command with no PreRunE in subcommands fails",
			args:    args{command: simpleCommand},
			wantErr: true,
		},
		{
			name:    "dummy command with childs and PreRunE in subcommands fails",
			args:    args{command: complexCommand},
			wantErr: true,
		},
		{
			name:    "RootCmd with all PreRunE set succeeds",
			args:    args{command: RootCmd},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkPreRunE(tt.args.command); (err != nil) != tt.wantErr {
				t.Errorf("checkPreRunE() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetCommand(t *testing.T) {
	type args struct {
		command *cobra.Command
		path    []string
	}
	tests := []struct {
		name string
		args args
		want *cobra.Command
	}{
		{
			name: "Get a top level command",
			args: args{
				command: RootCmd,
				path:    []string{"deployment"},
			},
			want: cmddeployment.Command,
		},
		{
			name: "Get a 2nd level command",
			args: args{
				command: RootCmd,
				path:    []string{"deployment", "elasticsearch"},
			},
			want: cmdelasticsearch.Command,
		},
		{
			name: "Get a 3rd level command",
			args: args{
				command: RootCmd,
				path:    []string{"deployment", "elasticsearch", "plan"},
			},
			want: GetCommand(cmdelasticsearch.Command, "plan"),
		},
		{
			name: "Get a 4th level command",
			args: args{
				command: RootCmd,
				path:    []string{"deployment", "elasticsearch", "plan", "history"},
			},
			want: GetCommand(cmdelasticsearch.Command, "plan", "history"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCommand(tt.args.command, tt.args.path...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
