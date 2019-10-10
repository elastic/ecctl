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
	"reflect"
	"testing"

	multierror "github.com/hashicorp/go-multierror"
)

func TestValidateRoles(t *testing.T) {
	tests := []struct {
		name    string
		arg     []string
		wantErr bool
		err     error
	}{
		{
			name: "validate should return an error when ece_platform_admin is used along other roles",
			arg:  []string{platformAdminRole, platformViewerRole},
			err: &multierror.Error{
				Errors: []error{
					errors.New("user: ece_platform_admin cannot be used in conjunction with other roles"),
				},
			},
			wantErr: true,
		},
		{
			name: "validate should return an error when ece_platform_admin is used along other roles",
			arg:  []string{deploymentsManagerRole, deploymentsViewerRole},
			err: &multierror.Error{
				Errors: []error{
					errors.New("user: only one of ece_deployment_manager or ece_deployment_viewer can be chosen"),
				},
			},
			wantErr: true,
		},
		{
			name:    "validate should pass if all params are properly set",
			arg:     []string{platformViewerRole, deploymentsManagerRole},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRoles(tt.arg)

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.err == nil {
				t.Errorf("Validate() expected errors = '%v' but no errors returned", tt.err)
			}

			if tt.wantErr && err.Error() != tt.err.Error() {
				t.Errorf("Validate() expected errors = '%v' but got %v", tt.err, err)
			}
		})
	}
}

func TestHasBothDeploymentRoles(t *testing.T) {
	tests := []struct {
		name string
		arg  []string
		want bool
	}{
		{
			name: "should return true when both deployment roles are present",
			arg:  []string{deploymentsManagerRole, deploymentsViewerRole},
			want: true,
		},
		{
			name: "should return false if both deployment roles are not present",
			arg:  []string{deploymentsManagerRole, platformViewerRole, platformAdminRole},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasBothDeploymentRoles(tt.arg)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
