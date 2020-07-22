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

package cmdenrollmenttoken

import (
	"fmt"

	"github.com/elastic/cloud-sdk-go/pkg/api/platformapi/enrollmenttokenapi"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var deleteTokenCmd = &cobra.Command{
	Use:     "delete <enrollment-token>",
	Short:   cmdutil.AdminReqDescription("Deletes an enrollment token"),
	PreRunE: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := enrollmenttokenapi.Delete(enrollmenttokenapi.DeleteParams{
			API:    ecctl.Get().API,
			Region: ecctl.Get().Config.Region,
			Token:  args[0],
		})
		if err != nil {
			return err
		}
		fmt.Printf("Token %s deleted.\n", args[0])
		return nil
	},
}

func init() {
	Command.AddCommand(deleteTokenCmd)
}
