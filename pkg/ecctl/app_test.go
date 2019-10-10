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

package ecctl

import (
	"errors"
	"reflect"
	"testing"

	multierror "github.com/hashicorp/go-multierror"
)

func Test_newAuthWriter(t *testing.T) {
	type args struct {
		c Config
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "fails on empty config",
			args: args{},
			err: &multierror.Error{Errors: []error{
				errors.New("auth: Username must not be empty"),
				errors.New("auth: Password must not be empty"),
			}},
		},
		{
			name: "fails on invalid basicAuth credentials",
			args: args{
				c: Config{},
			},
			err: &multierror.Error{Errors: []error{
				errors.New("auth: Username must not be empty"),
				errors.New("auth: Password must not be empty"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := newAuthWriter(tt.args.c)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("newAuthWriter() error = %v, wantErr %v", err, tt.err)
				return
			}
		})
	}
}
