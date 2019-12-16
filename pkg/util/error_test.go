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
	"errors"
	"reflect"
	"testing"

	"github.com/hashicorp/go-multierror"
)

func TestWrapError(t *testing.T) {
	type args struct {
		text string
		err  error
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "returns nil",
		},
		{
			name: "returns a wrapped error",
			args: args{
				text: "some",
				err:  errors.New("some error"),
			},
			err: errors.New("some: some error"),
		},
		{
			name: "returns the error when the text is empty",
			args: args{
				text: "",
				err:  errors.New("some error"),
			},
			err: errors.New("some error"),
		},
		{
			name: "returns nil when error is empty",
			args: args{
				text: "",
			},
		},
		{
			name: "returns a wrapped multierror",
			args: args{
				text: "some wrapping text",
				err: &multierror.Error{Errors: []error{
					errors.New("error 1"),
					errors.New("error 2"),
				}},
			},
			err: &multierror.Error{Errors: []error{
				errors.New("some wrapping text: error 1"),
				errors.New("some wrapping text: error 2"),
			}},
		},
		{
			name: "returns the multierror when text is empty",
			args: args{
				err: &multierror.Error{Errors: []error{
					errors.New("error 1"),
					errors.New("error 2"),
				}},
			},
			err: &multierror.Error{Errors: []error{
				errors.New("error 1"),
				errors.New("error 2"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WrapError(tt.args.text, tt.args.err); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("WrapError() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
