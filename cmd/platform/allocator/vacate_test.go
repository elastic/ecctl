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

package cmdallocator

import (
	"errors"
	"testing"
)

func TestValidateSkipDataMigration(t *testing.T) {
	type args struct {
		resources []string
		moveOnly  bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "Succeeds if a resource has been set and moveOnly is set to true",
			args: args{
				resources: []string{"0f32fa44e"},
				moveOnly:  true,
			},
			wantErr: false,
		},
		{
			name: "Returns an error if a resource has been set and moveOnly is set to false",
			args: args{
				resources: []string{"0f32fa44e"},
				moveOnly:  false,
			},
			wantErr: true,
			err:     errors.New("skip data migration is not available if there are no resource IDs specified or move-only is set to false"),
		},
		{
			name: "Returns an error if no resource has been set and moveOnly is set to true",
			args: args{
				resources: []string{},
				moveOnly:  true,
			},
			wantErr: true,
			err:     errors.New("skip data migration is not available if there are no resource IDs specified or move-only is set to false"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSkipDataMigration(tt.args.resources, tt.args.moveOnly)

			if (err != nil) != tt.wantErr {
				t.Errorf("validateSkipDataMigration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.err == nil {
				t.Errorf("validateSkipDataMigration() expected errors = '%v' but no errors returned", tt.err)
			}

			if tt.wantErr && err.Error() != tt.err.Error() {
				t.Errorf("validateSkipDataMigration() expected errors = '%v' but got %v", tt.err, err)
			}
		})
	}
}
