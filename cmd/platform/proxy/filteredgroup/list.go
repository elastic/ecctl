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

package cmdfilteredgroup

import (
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/platform/proxy/filteredgroup"
)

const (
	filteredGroupsUse = `Returns all proxies filtered groups in the platform`
)

func listProxyFilteredGroups(cmd *cobra.Command, args []string) error {
	response, err := filteredgroup.List(filteredgroup.CommonParams{
		API: ecctl.Get().API,
		ID:  "all",
	})
	if err != nil {
		return err
	}

	return ecctl.Get().Formatter.Format(filepath.Join("filtered-group", "list"), response)
}

var platformProxyFilteredGroupsListCmd = &cobra.Command{
	Use:     "list",
	Short:   filteredGroupsUse,
	PreRunE: cobra.MaximumNArgs(0),
	RunE:    listProxyFilteredGroups,
}

func init() {
	Command.AddCommand(
		platformProxyFilteredGroupsListCmd,
	)
}
