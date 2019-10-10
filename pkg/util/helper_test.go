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
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
)

func TestClusterParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		fields  *ClusterParams
		wantErr bool
	}{
		{
			name: "It should not error if all params exist and cluster id has correct length",
			fields: &ClusterParams{
				ClusterID: "320b7b540dfc967a7a649c18e2fce4ed",
				API:       new(api.API),
			},
			wantErr: false,
		},
		{
			name: "It should error if cluster id has invalid length",
			fields: &ClusterParams{
				ClusterID: "myclusterid",
			},
			wantErr: true,
		},
		{
			name: "It should error if cluster id is empty",
			fields: &ClusterParams{
				ClusterID: "",
			},
			wantErr: true,
		},
		{
			name: "It should error if API is nil",
			fields: &ClusterParams{
				ClusterID: "320b7b540dfc967a7a649c18e2fce4ed",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fields.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("ClusterParams.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
