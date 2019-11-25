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

package elasticsearch

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

const secretGet = `{
  "secrets": {
    "s3.client.foobar.access_key": {
      "as_file": false
    },
    "s3.client.foobar.secret_key": {
      "as_file": false
    }
  }
}`

func TestSetKeystore(t *testing.T) {
	const secretContents = `{
		"secrets": {
			"s3.client.foobar.access_key": {
				"value": "AKIXAIQFKXPHIFXSILUWPA",
				"as_file": false
			},
			"s3.client.foobar.secret_key": {
				"value": "18qXOpY2zGlApay1237dLXh+LG1X5LUNWjTHq5X1SWjf++m+p0"
			}
		}
	}`
	var req models.KeystoreContents
	json.Unmarshal([]byte(secretContents), &req)

	type args struct {
		params SetKeystoreParams
	}
	tests := []struct {
		name string
		args args
		want *models.KeystoreContents
		err  error
	}{
		{
			name: "fails due to param validation",
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				util.ErrClusterLength,
				errors.New("elasticsarch keystore: set requires the keystore contents to be specified"),
			}},
		},
		{
			name: "fails due to API error",
			args: args{params: SetKeystoreParams{
				API:       api.NewMock(mock.New500Response(mock.NewStringBody("error"))),
				ClusterID: util.ValidClusterID,
				Request:   &req,
			}},
			err: errors.New("error"),
		},
		{
			name: "succeeds",
			args: args{params: SetKeystoreParams{
				API:       api.NewMock(mock.New202Response(mock.NewStringBody(secretGet))),
				ClusterID: util.ValidClusterID,
				Request:   &req,
			}},
			want: &models.KeystoreContents{
				Secrets: map[string]models.KeystoreSecret{
					"s3.client.foobar.access_key": {
						AsFile: ec.Bool(false),
					},
					"s3.client.foobar.secret_key": {
						AsFile: ec.Bool(false),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SetKeystore(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("SetKeystore() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetKeystore() = %v, want %v", got, tt.want)
			}
		})
	}
}
