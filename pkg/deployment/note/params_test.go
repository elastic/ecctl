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

package note

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"

	"github.com/elastic/ecctl/pkg/deployment"
)

func TestGetElasticsearchID(t *testing.T) {
	const getResponse = `{
  "healthy": true,
  "id": "e3dac8bf3dc64c528c295a94d0f19a77",
  "resources": {
    "elasticsearch": [{
      "id": "418017cd1c7f402cbb7a981b2004ceeb",
      "ref_id": "main-elasticsearch",
      "region": "ece-region"
    }]
  }
}`

	const getErrorResponse = `{
  "errors": null
}`

	type args struct {
		params deployment.GetParams
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
		err     error
	}{
		{
			name: "Get fails due to API failure",
			args: args{
				params: deployment.GetParams{
					DeploymentID: "e3dac8bf3dc64c528c295a94d0f19a77",
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(""),
						StatusCode: 500,
					}}),
				},
			},
			wantErr: true,
			err:     errors.New(getErrorResponse),
		},
		{
			name: "Get succeeds",
			args: args{
				params: deployment.GetParams{
					DeploymentID: "e3dac8bf3dc64c528c295a94d0f19a77",
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(getResponse),
						StatusCode: 200,
					}}),
				},
			},
			want: "418017cd1c7f402cbb7a981b2004ceeb",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getElasticsearchID(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("getElasticsearchID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.err != nil && err.Error() != tt.err.Error() {
				t.Errorf("getElasticsearchID() actual error = '%v', want error '%v'", err, tt.err)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getElasticsearchID() = %v, want %v", got, tt.want)
			}
		})
	}
}
