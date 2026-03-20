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

package cmdproject

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/project"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "Lists serverless projects",
	PreRunE: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectType, _ := cmd.Flags().GetString("type")

		res, err := project.List(project.ListParams{
			API:    ecctl.Get().API,
			Host:   ecctl.Get().Config.Host,
			Type:   projectType,
			Client: ecctl.Get().Config.Client,
		})
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format(filepath.Join("project", "list"), res)
	},
}

func init() {
	Command.AddCommand(listCmd)

	listCmd.Flags().String("type", "", "Filters by project type (elasticsearch/search, observability, security)")
}
