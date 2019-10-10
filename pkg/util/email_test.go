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
	"testing"
)

func TestValidateEmail(t *testing.T) {
	type args struct {
		prefix string
		email  string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name:    "ValidateEmail should return an error if string is not in a valid format",
			args:    args{"user", "hi"},
			err:     errors.New("user: hi is not a valid email address format"),
			wantErr: true,
		},
		{
			name:    "ValidateEmail should return an error when entered password is too long",
			args:    args{"user", "46434164864164964163546874654944645654554665467846146948196846163549464341648641649641635468746549446456545546654678461469481968461635494643416486416496416354687465494464565455466546784614694819684616354965465465476354685746354687484765454ss65@example.com"},
			err:     errors.New("an email address must not exceed 254 characters"),
			wantErr: true,
		},
		{
			name:    "ValidateEmail should pass if string is a valid email address",
			args:    args{"user", "hi@example.com"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEmail(tt.args.prefix, tt.args.email)

			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.err == nil {
				t.Errorf("ValidateEmail() expected errors = '%v' but no errors returned", tt.err)
			}

			if tt.wantErr && err.Error() != tt.err.Error() {
				t.Errorf("ValidateEmail() expected errors = '%v' but got %v", tt.err, err)
			}
		})
	}
}
