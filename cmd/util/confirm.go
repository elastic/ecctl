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
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/elastic/cloud-sdk-go/pkg/input"

	"github.com/elastic/ecctl/pkg/ecctl"
)

// ActionConfirm asks you to confirm before performing an action
func ActionConfirm(actionRaw, msg string) (*bool, error) {
	var action *bool

	parsedAction, err := strconv.ParseBool(actionRaw)
	if err != nil && actionRaw != "" {
		return nil, err
	}

	if actionRaw != "" && err == nil {
		action = &parsedAction
	}

	if action != nil && !ecctl.Get().Config.Force {
		scanner := input.NewScanner(os.Stdin, ecctl.Get().Config.OutputDevice)
		if confirm := strings.ToLower(scanner.Scan(
			msg,
		)); !strings.HasPrefix(confirm, "y") {
			return nil, errors.New("action has been aborted")
		}
		fmt.Fprintln(ecctl.Get().Config.OutputDevice, "continuing...")
	}

	return action, nil
}
