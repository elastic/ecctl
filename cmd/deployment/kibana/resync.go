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

package cmdkibana

import (
	"fmt"

	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment/kibana"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var resyncKibanaCmd = &cobra.Command{
	Use:     "resync <cluster id>",
	Short:   "Resynchronizes the search index and cache for the selected Kibana cluster",
	PreRunE: cmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Resynchronizing Kibana cluster: %s\n", args[0])
		return kibana.Resync(kibana.DeploymentParams{
			API: ecctl.Get().API,
			ID:  args[0],
		})
	},
}
