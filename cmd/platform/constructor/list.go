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
	"path/filepath"

	"github.com/elastic/cloud-sdk-go/pkg/api/platformapi/constructorapi"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var listConstructorsCmd = &cobra.Command{
	Use:     "list",
	Short:   cmdutil.AdminReqDescription("Returns all of the constructors in the platform"),
	PreRunE: cobra.NoArgs,
	RunE:    listConstructors,
}

func listConstructors(cmd *cobra.Command, args []string) error {
	a, err := constructorapi.List(constructorapi.ListParams{
		API:    ecctl.Get().API,
		Region: ecctl.Get().Config.Region,
	})
	if err != nil {
		return err
	}

	return ecctl.Get().Formatter.Format(filepath.Join("constructor", "list"), a)
}

func init() {
	Command.AddCommand(listConstructorsCmd)
}
