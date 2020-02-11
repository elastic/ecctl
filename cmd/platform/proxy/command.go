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

package cmdproxy

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	cmdfilteredgroup "github.com/elastic/ecctl/cmd/platform/proxy/filteredgroup"
	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/platform/proxy"
)

const (
	proxyListMessage = `Returns all of the proxies in the platform`
	proxyShowMessage = `Returns information about the proxy with given id`
)

// Command represents the proxy command
var Command = &cobra.Command{
	Use:     "proxy",
	Short:   fmt.Sprintf("Manages proxies %v", cmdutil.PlatformAdminRequired),
	PreRunE: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func listProxies(cmd *cobra.Command, args []string) error {
	a, err := proxy.List(proxy.Params{
		API: ecctl.Get().API,
	})
	if err != nil {
		return err
	}

	return ecctl.Get().Formatter.Format(filepath.Join("proxy", "list"), a)
}

var showProxyCmd = &cobra.Command{
	Use:     "show <proxy id>",
	Short:   proxyShowMessage,
	PreRunE: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := proxy.Get(proxy.GetParams{
			Params: proxy.Params{
				API: ecctl.Get().API,
			},
			ID: args[0],
		})
		if err != nil {
			return err
		}
		return ecctl.Get().Formatter.Format(filepath.Join("proxy", "show"), a)
	},
}

var listProxiesCmd = &cobra.Command{
	Use:     "list",
	Short:   proxyListMessage,
	PreRunE: cobra.MaximumNArgs(0),
	RunE:    listProxies,
}

func init() {
	Command.AddCommand(
		listProxiesCmd,
		showProxyCmd,
		cmdfilteredgroup.Command,
	)
}
