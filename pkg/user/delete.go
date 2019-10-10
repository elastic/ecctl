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
	"errors"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/users"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

// DeleteParams is consumed by Delete
type DeleteParams struct {
	*api.API

	UserName string
}

// Validate ensures the parameters are usable by the consuming function.
func (params DeleteParams) Validate() error {
	var merr = new(multierror.Error)
	if params.API == nil {
		merr = multierror.Append(merr, util.ErrAPIReq)
	}

	if params.UserName == "" {
		merr = multierror.Append(merr, errors.New("user: delete requires a username"))
	}

	return merr.ErrorOrNil()
}

// Delete deletes a user given a username
func Delete(params DeleteParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	return util.ReturnErrOnly(params.V1API.Users.DeleteUser(
		users.NewDeleteUserParams().
			WithUserName(params.UserName),
		params.AuthWriter,
	))
}
