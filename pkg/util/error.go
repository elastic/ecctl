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

package util

import (
	"fmt"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/hashicorp/go-multierror"
)

// ReturnErrOnly is used to strip the first return argument of a function call
func ReturnErrOnly(_ interface{}, e error) error {
	return api.UnwrapError(e)
}

// WrapError takes the error that is passed and either wraps the error with the
// text or each of the multierror errors with it.
func WrapError(text string, err error) error {
	if merr, ok := err.(*multierror.Error); ok {
		var newMerrs = make([]error, 0, len(merr.Errors))
		for _, e := range merr.Errors {
			newMerrs = append(newMerrs, wrap(e, text))
		}
		merr.Errors = newMerrs
		return merr.ErrorOrNil()
	}

	return wrap(err, text)
}

func wrap(err error, text string) error {
	if text != "" {
		return fmt.Errorf("%s: %s", text, err)
	}
	return err
}
