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

package cmddeployment

import (
	"github.com/spf13/cobra"

	cmdapm "github.com/elastic/ecctl/cmd/deployment/apm"
	cmdelasticsearch "github.com/elastic/ecctl/cmd/deployment/elasticsearch"
	cmdkibana "github.com/elastic/ecctl/cmd/deployment/kibana"
	cmddeploymentnote "github.com/elastic/ecctl/cmd/deployment/note"
)

// Command is the deployment subcommand
var Command = &cobra.Command{
	Use:     "deployment",
	Short:   "Manages deployments",
	PreRunE: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	Command.AddCommand(
		cmddeploymentnote.Command,
		cmdelasticsearch.Command,
		cmdkibana.Command,
		cmdapm.Command,
	)
}
