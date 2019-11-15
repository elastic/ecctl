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
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/deployment/deputil"
	"github.com/elastic/ecctl/pkg/deployment/kibana"
	"github.com/elastic/ecctl/pkg/ecctl"
)

var kibanaListCmd = &cobra.Command{
	Use:     "list",
	Short:   "Returns the list of clusters for a region",
	PreRunE: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		size, _ := cmd.Flags().GetInt64("size")
		metadata, _ := cmd.Flags().GetBool("metadata")

		p, err := kibana.List(kibana.ListParams{
			API:     ecctl.Get().API,
			Version: cmd.Flag("version").Value.String(),
			QueryParams: deputil.QueryParams{
				Size:         size,
				ShowMetadata: metadata,
			},
		})
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format(filepath.Join("kibana", "list"), p)
	},
}

func init() {
	kibanaListCmd.Flags().BoolP("metadata", "m", false, "Shows deployment metadata")
	kibanaListCmd.Flags().Int64P("size", "s", 100, "Sets the upper limit of Kibana instances to return")
	kibanaListCmd.Flags().StringP("version", "v", "", "Filters per version")
}
