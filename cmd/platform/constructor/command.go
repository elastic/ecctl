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

package cmdconstructor

import (
	"fmt"
	"path/filepath"

	"github.com/elastic/cloud-sdk-go/pkg/api/platformapi/constructor"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
)

const (
	constructorListMessage        = `Returns all of the constructors in the platform`
	constructorShowMessage        = `Returns information about the constructor with given ID`
	constructorMaintenanceMessage = `Sets/un-sets a constructor's maintenance mode`
)

// Command represents the constructor command
var Command = &cobra.Command{
	Use:     "constructor",
	Short:   cmdutil.AdminReqDescription("Manages constructors"),
	PreRunE: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func listConstructors(cmd *cobra.Command, args []string) error {
	a, err := constructor.List(constructor.Params{
		API: ecctl.Get().API,
	})
	if err != nil {
		return err
	}

	return ecctl.Get().Formatter.Format(filepath.Join("constructor", "list"), a)
}

var showConstructorCmd = &cobra.Command{
	Use:     "show <constructor id>",
	Short:   constructorShowMessage,
	PreRunE: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := constructor.Get(constructor.GetParams{
			Params: constructor.Params{
				API: ecctl.Get().API,
			},
			ID: args[0],
		})
		if err != nil {
			return err
		}
		return ecctl.Get().Formatter.Format(filepath.Join("constructor", "show"), a)
	},
}

var maintenanceConstructorCmd = &cobra.Command{
	Use:     "maintenance <constructor id>",
	Short:   constructorMaintenanceMessage,
	PreRunE: cobra.MinimumNArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		unset, _ := cmd.Flags().GetBool("unset")
		fmt.Printf("Setting contructor %s maintenance to %t\n", args[0], !unset)
		var params = constructor.MaintenanceParams{
			Params: constructor.Params{
				API: ecctl.Get().API,
			},
			ID: args[0],
		}

		if unset {
			return constructor.DisableMaintenance(params)
		}
		return constructor.EnableMaintenace(params)
	},
}

var listConstructorsCmd = &cobra.Command{
	Use:     "list",
	Short:   constructorListMessage,
	PreRunE: cobra.MaximumNArgs(0),
	RunE:    listConstructors,
}

func init() {
	Command.AddCommand(
		listConstructorsCmd,
		showConstructorCmd,
		maintenanceConstructorCmd,
	)

	maintenanceConstructorCmd.Flags().Bool("unset", false, "Unset constructor maintenance mode")
}
