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

package kibana

import (
	"errors"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	multierror "github.com/hashicorp/go-multierror"
)

func TestDeploymentParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  *DeploymentParams
		wantErr bool
		err     error
	}{
		{
			name:   "validate should return all possible errors",
			params: &DeploymentParams{},
			err: &multierror.Error{
				Errors: []error{
					errors.New("api reference is required for command"),
					errors.New(`id "" is invalid`),
				},
			},
			wantErr: true,
		},
		{
			name: "validate should pass if all params are properly set",
			params: &DeploymentParams{
				ID:  "5c641576747442eba0ebd67944ccbe10",
				API: &api.API{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.err == nil {
				t.Errorf("Validate() expected errors = '%v' but no errors returned", tt.err)
			}

			if tt.wantErr && !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Validate() expected errors = '%v' but got %v", tt.err, err)
			}
		})
	}
}
