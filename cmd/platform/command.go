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

package cmdplatform

import (
	"github.com/spf13/cobra"

	cmdallocator "github.com/elastic/ecctl/cmd/platform/allocator"
	cmdconstructor "github.com/elastic/ecctl/cmd/platform/constructor"
	cmdenrollmenttoken "github.com/elastic/ecctl/cmd/platform/enrollment-token"
	cmdinstanceconfig "github.com/elastic/ecctl/cmd/platform/instance-configuration"
	cmdproxy "github.com/elastic/ecctl/cmd/platform/proxy"
	cmdrepository "github.com/elastic/ecctl/cmd/platform/repository"
	cmdrole "github.com/elastic/ecctl/cmd/platform/role"
	cmdrunner "github.com/elastic/ecctl/cmd/platform/runner"
)

// Command is the platform subcommand
var Command = &cobra.Command{
	Use:     "platform",
	Short:   "Manages the platform",
	PreRunE: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	Command.AddCommand(
		cmdallocator.Command,
		cmdconstructor.Command,
		cmdenrollmenttoken.Command,
		cmdinstanceconfig.Command,
		cmdproxy.Command,
		cmdrepository.Command,
		cmdrole.Command,
		cmdrunner.Command,
	)
}
