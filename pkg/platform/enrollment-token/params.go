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

package enrollmenttoken

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/multierror"

	"github.com/elastic/ecctl/pkg/util"
)

// CreateParams is consumed by Create
type CreateParams struct {
	*api.API
	Roles    []string
	Duration time.Duration
}

// Validate ensures that there's no errors prior to performing the Create API
// call.
func (params CreateParams) Validate() error {
	var merr = multierror.NewPrefixed("enrollment-token create")
	if params.API == nil {
		merr = merr.Append(util.ErrAPIReq)
	}

	validity := int64(params.Duration.Seconds())
	if validity > math.MaxInt32 {
		merr = merr.Append(
			fmt.Errorf("validity value %d exceeds max allowed %d value in seconds", validity, math.MaxInt32),
		)
	}

	return merr.ErrorOrNil()
}

// DeleteParams is consumed by Delete
type DeleteParams struct {
	*api.API
	Token string
}

// Validate ensures that there's no errors prior to performing the Delete API
// call.
func (params DeleteParams) Validate() error {
	var merr = multierror.NewPrefixed("enrollment-token delete")
	if params.API == nil {
		merr = merr.Append(util.ErrAPIReq)
	}

	if params.Token == "" {
		merr = merr.Append(errors.New("token cannot be empty"))
	}

	return merr.ErrorOrNil()
}

// ListParams is consumed by List
type ListParams struct {
	*api.API
}

// Validate ensures that there's no errors prior to performing the List API
// call.
func (params ListParams) Validate() error {
	if params.API == nil {
		return util.ErrAPIReq
	}

	return nil
}
