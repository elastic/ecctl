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
	"reflect"
	"testing"
)

func TestGetInsecurePassword(t *testing.T) {
	tests := []struct {
		name    string
		arg     string
		want    []byte
		wantErr bool
		err     error
	}{
		{
			name:    "GetInsecurePassword fails due to validation",
			arg:     "",
			wantErr: true,
			err:     errors.New("a password must be provided when using the --insecure-password flag"),
		},
		{
			name: "GetInsecurePassword succeeds",
			arg:  "supersecretpass",
			want: []byte("supersecretpass"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetInsecurePassword(tt.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInsecurePassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.err != nil && !reflect.DeepEqual(err, tt.err) {
				t.Errorf("GetInsecurePassword() actual error = '%v', want error '%v'", err, tt.err)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetInsecurePassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
