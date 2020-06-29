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

package cmdauthkey

import (
	"os"
	"sync"
	"time"

	userauthapi "github.com/elastic/cloud-sdk-go/pkg/api/userapi/authapi"
	"github.com/elastic/cloud-sdk-go/pkg/multierror"
	"github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/ecctl"
)

var deleteCmd = &cobra.Command{
	Use:     "delete <key id> <key id> ...",
	Short:   "Deletes one or more existing API keys for the specified user",
	PreRunE: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			var msg = "Multiple keys will be deleted, do you want to continue? [y/n]: "
			if !cmdutil.ConfirmAction(msg, os.Stdin, ecctl.Get().Config.OutputDevice) {
				return nil
			}
		}

		var merr = multierror.NewPrefixed("delete key")
		var wg sync.WaitGroup
		for i := range args {
			wg.Add(1)
			go func(index int) {
				merr = merr.Append(
					userauthapi.DeleteKey(userauthapi.DeleteKeyParams{
						API: ecctl.Get().API,
						ID:  args[index],
					}),
				)
				wg.Done()
			}(i)

			// Only delete a key per second as anything below causes a 500
			<-time.After(time.Second)
		}

		wg.Wait()
		return merr.ErrorOrNil()
	},
}
