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

package cmdutil

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/deployment/elasticsearch/instances"
	"github.com/elastic/ecctl/pkg/util"
)

const (
	compFuncTpl = `__ecctl_valid_%s_types()
{
   COMPREPLY=( %s )
}
`
)

const (
	maxPollRetriesFlag = "max-poll-retries"
	pollFrequencyFlag  = "poll-frequency"
)

var (
	// StatelessTypes declares the stateless deployment resource types
	StatelessTypes = []string{"apm", "appsearch", "kibana"}

	// StatefulTypes declares the stateful deployment resource types (Elasticsearch).
	StatefulTypes = []string{"elasticsearch"}

	// AllTypes is StatelessTypes appending StatefulTypes
	AllTypes = append(StatelessTypes, StatefulTypes...)

	// StatelessTypesCompFunc is the bash autocompletion function for stateless types.
	StatelessTypesCompFunc = fmt.Sprintf(compFuncTpl, "stateless", strings.Join(StatelessTypes, " "))

	// AllTypesCompFunc is the bash autocompletion function for all types.
	AllTypesCompFunc = fmt.Sprintf(compFuncTpl, "all", strings.Join(AllTypes, " "))
)

// GetInstances tries to obtain a slice with the elasticsearch cluster
// instance names either from the current cluster topology or from the
// specified --instance flag in the cobra.Command.
func GetInstances(cmd *cobra.Command, params util.ClusterParams, flagName string) ([]string, error) {
	if all, _ := cmd.Flags().GetBool("all"); all {
		return instances.List(params)
	}
	return cmd.Flags().GetStringSlice(flagName)
}

// AddTypeFlag adds a type string  flag to the specified command, with the
// resource types autocompletion function. It is intended to be used for any
// commands which call the deployment/resource APIs.
func AddTypeFlag(cmd *cobra.Command, prefix string, all bool) *string {
	validTypes, comp := StatelessTypes, "__ecctl_valid_stateless_types"
	if all {
		validTypes, comp = AllTypes, "__ecctl_valid_all_types"
	}

	s := cmd.Flags().String("type", "", fmt.Sprintf(
		"%s deployment resource type (%s)", prefix, strings.Join(validTypes, ", "),
	))

	cmd.Flag("type").Annotations = map[string][]string{cobra.BashCompCustom: {comp}}

	return s
}

// AddTrackFlags adds flags which control the tracking frequency to the passed
// command reference.
func AddTrackFlags(cmd *cobra.Command) {
	cmd.Flags().Int(maxPollRetriesFlag, util.DefaultRetries, "Optional maximum plan tracking retries")
	cmd.Flags().Duration(pollFrequencyFlag, util.DefaultPollFrequency, "Optional polling frequency to check for plan change updates")
}

// GetTrackSettings obtains the currently set tracking settings, the first
// return value being the MaxPollRetries and the second one the poll frequency.
func GetTrackSettings(cmd *cobra.Command) (int, time.Duration) {
	maxPollRetries, _ := cmd.Flags().GetInt(maxPollRetriesFlag)
	pollFrequency, _ := cmd.Flags().GetDuration(pollFrequencyFlag)
	return maxPollRetries, pollFrequency
}
