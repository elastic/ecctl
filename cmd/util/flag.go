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
	compFuncTpl = `__ecctl_valid_%s_kinds()
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
	// StatelessKinds declares the stateless deployment resource kinds
	StatelessKinds = []string{"apm", "appsearch", "kibana"}

	// StatefulKinds declares the stateful deployment resource kinds (Elasticsearch).
	StatefulKinds = []string{"elasticsearch"}

	// AllKinds is StatelessKinds appending StatefulKinds
	AllKinds = append(StatelessKinds, StatefulKinds...)

	// StatelessKindsCompFunc is the bash autocompletion function for stateless kinds.
	StatelessKindsCompFunc = fmt.Sprintf(compFuncTpl, "stateless", strings.Join(StatelessKinds, " "))

	// AllKindsCompFunc is the bash autocompletion function for all kinds.
	AllKindsCompFunc = fmt.Sprintf(compFuncTpl, "all", strings.Join(AllKinds, " "))
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

// AddKindFlag adds a kind string  flag to the specified command, with the
// resource kinds autocompletion function. It is intended to be used for any
// commands which call the deployment/resource APIs.
func AddKindFlag(cmd *cobra.Command, prefix string, all bool) *string {
	validKinds, comp := StatelessKinds, "__ecctl_valid_stateless_kinds"
	if all {
		validKinds, comp = AllKinds, "__ecctl_valid_all_kinds"
	}

	s := cmd.Flags().String("kind", "", fmt.Sprintf(
		"%s deployment resource kind (%s)", prefix, strings.Join(validKinds, ", "),
	))

	cmd.Flag("kind").Annotations = map[string][]string{cobra.BashCompCustom: {comp}}

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

// ConflictingFlags checks if both flags have been specified, and if so
// returns an error.
func ConflictingFlags(cmd *cobra.Command, first, second string) error {
	if cmd.Flag(first).Changed && cmd.Flag(second).Changed {
		return fmt.Errorf(
			`conflicting flags: "--%s" and "--%s" should not be used together"`,
			first, second,
		)
	}
	return nil
}

// MustUseAFlag checks if one or another flags are used, and if not
// returns an error.
func MustUseAFlag(cmd *cobra.Command, first, second string) error {
	if !cmd.Flag(first).Changed && !cmd.Flag(second).Changed {
		return fmt.Errorf(
			`necessary flags: one of "--%s" or "--%s" should be used"`,
			first, second,
		)
	}
	return nil
}
