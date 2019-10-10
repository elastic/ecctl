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

package user

import (
	"github.com/elastic/cloud-sdk-go/pkg/util/slice"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

// Available user roles
const platformAdminRole = "ece_platform_admin"
const platformViewerRole = "ece_platform_viewer"
const deploymentsManagerRole = "ece_deployment_manager"
const deploymentsViewerRole = "ece_deployment_viewer"

// ValidateRoles ensures the parameters are usable by the consuming function.
func ValidateRoles(roles []string) error {
	var merr = new(multierror.Error)

	if len(roles) > 1 && slice.HasString(roles, platformAdminRole) {
		merr = multierror.Append(merr, errors.Errorf("user: %v cannot be used in conjunction with other roles", platformAdminRole))
	}

	if hasBothDeploymentRoles(roles) {
		merr = multierror.Append(merr, errors.Errorf("user: only one of %v or %v can be chosen",
			deploymentsManagerRole, deploymentsViewerRole))
	}

	return merr.ErrorOrNil()
}

func hasBothDeploymentRoles(roles []string) bool {
	return slice.HasString(roles, deploymentsManagerRole) &&
		slice.HasString(roles, deploymentsViewerRole)
}
