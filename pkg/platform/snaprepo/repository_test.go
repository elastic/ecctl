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

package snaprepo

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func TestGet(t *testing.T) {
	var getSnapshotSuccess = `
	{
		"repository_name": "my_snapshot_repo",
		"config": {
			"region":"us-west-1",
			"bucket":"mybucket",
			"access_key":"anaccesskey",
			"secret_key":"asecretkey"
		}
	}
`[1:]

	type args struct {
		params GetParams
	}
	tests := []struct {
		name    string
		args    args
		want    *models.RepositoryConfig
		wantErr error
	}{
		{
			name: "Getting a snapshot repository succeeds",
			args: args{
				params: GetParams{
					Params: Params{
						API: api.NewMock(mock.Response{Response: http.Response{
							StatusCode: 200,
							Body:       mock.NewStringBody(getSnapshotSuccess),
						}}),
					},
					Name: "my_snapshot_repo",
				},
			},
			want: &models.RepositoryConfig{
				RepositoryName: ec.String("my_snapshot_repo"),
				Config: map[string]interface{}{
					"region":     "us-west-1",
					"bucket":     "mybucket",
					"access_key": "anaccesskey",
					"secret_key": "asecretkey",
				},
			},
		},
		{
			name: "Getting a snapshot repository fails when api returns an error",
			args: args{
				params: GetParams{
					Params: Params{
						API: api.NewMock(mock.Response{Error: errors.New("ERROR")}),
					},
					Name: "my_snapshot_repo",
				},
			},
			wantErr: &url.Error{
				Op:  "Get",
				URL: "https://mock-host/mock-path/platform/configuration/snapshots/repositories/my_snapshot_repo",
				Err: errors.New("ERROR"),
			},
		},
		{
			name: "Getting a snapshot repository fails when parameters are invalid",
			args: args{
				params: GetParams{
					Params: Params{API: nil},
					Name:   "my_snapshot_repo",
				},
			},
			wantErr: errAPICannotBeNil,
		},
		{
			name: "Getting a snapshot repository fails when parameters are invalid",
			args: args{
				params: GetParams{
					Params: Params{API: new(api.API)},
				},
			},
			wantErr: errNameCannotBeEmpty,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.args.params)

			if err := util.CheckErrType(err, tt.wantErr); err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type args struct {
		params DeleteParams
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Delete succeeds",
			args: args{
				params: DeleteParams{
					Name: "my_repo",
					Params: Params{
						API: api.NewMock(mock.Response{Response: http.Response{
							StatusCode: 200,
							Body:       mock.NewStringBody(`{}`),
						}}),
					},
				},
			},
		},
		{
			name: "Delete fails on 404",
			args: args{
				params: DeleteParams{
					Name: "my_repo",
					Params: Params{
						API: api.NewMock(mock.New404Response(mock.NewStringBody(`{"error": "some error"}`))),
					},
				},
			},
			wantErr: errors.New(`{"error": "some error"}`),
		},
		{
			name: "Delete fails on invalid params",
			args: args{
				params: DeleteParams{
					Params: Params{API: new(api.API)},
				},
			},
			wantErr: errNameCannotBeEmpty,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Delete(tt.args.params)
			if err := util.CheckErrType(err, tt.wantErr); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestList(t *testing.T) {
	var listSnapshotsSuccess = `
	{
		"configs": [
			{
				"config": {
					"access_key": "myaccesskey",
					"base_path": "apath",
					"bucket": "mybucket",
					"canned_acl": "private",
					"compress": true,
					"protocol": "http",
					"region": "us-east-1",
					"secret_key": "mysupersecretkey",
					"server_side_encryption": true,
					"storage_class": "standard"
				},
				"repository_name": "my_repo_1"
			},
			{
				"config": {
					"access_key": "myaccesskey",
					"bucket": "mybucket",
					"region": "us-east-1",
					"secret_key": "mysupersecretkey",
					"server_side_encryption": true,
					"storage_class": "standard"
				},
				"repository_name": "my_repo_2"
			}
		]
	}
`[1:]

	type args struct {
		params Params
	}
	tests := []struct {
		name    string
		args    args
		want    *models.RepositoryConfigs
		wantErr error
	}{
		{
			name: "List Succeeds",
			args: args{
				params: Params{
					API: api.NewMock(mock.Response{Response: http.Response{
						StatusCode: 200,
						Body:       mock.NewStringBody(listSnapshotsSuccess),
					}}),
				},
			},
			want: &models.RepositoryConfigs{
				Configs: []*models.RepositoryConfig{
					{
						RepositoryName: ec.String("my_repo_1"),
						Config: map[string]interface{}{
							"region":                 "us-east-1",
							"bucket":                 "mybucket",
							"access_key":             "myaccesskey",
							"secret_key":             "mysupersecretkey",
							"base_path":              "apath",
							"compress":               true,
							"server_side_encryption": true,
							"canned_acl":             "private",
							"storage_class":          "standard",
							"protocol":               "http",
						},
					},
					{
						RepositoryName: ec.String("my_repo_2"),
						Config: map[string]interface{}{
							"region":                 "us-east-1",
							"bucket":                 "mybucket",
							"access_key":             "myaccesskey",
							"secret_key":             "mysupersecretkey",
							"server_side_encryption": true,
							"storage_class":          "standard",
						},
					},
				},
			},
		},
		{
			name:    "List fails when parameters are not valid",
			args:    args{},
			wantErr: errAPICannotBeNil,
		},
		{
			name: "List fails when the API returns an error",
			args: args{
				params: Params{
					API: api.NewMock(mock.Response{Error: errors.New("ERROR")}),
				},
			},
			wantErr: &url.Error{
				Op:  "Get",
				URL: "https://mock-host/mock-path/platform/configuration/snapshots/repositories",
				Err: errors.New("ERROR"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := List(tt.args.params)
			d := json.NewEncoder(os.Stdout)
			d.SetIndent("", "    ")
			d.Encode(tt.want)

			if err := util.CheckErrType(err, tt.wantErr); err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSet(t *testing.T) {
	type args struct {
		params SetParams
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "Set succeeds",
			args: args{
				params: SetParams{
					Name: "snapshot_repo_name",
					Type: "s3",
					Config: S3Config{
						Region:    "us-east-1",
						Bucket:    "mybucket",
						AccessKey: "myaccesskey",
						SecretKey: "mysecretkey",
					},
					Params: Params{
						API: api.NewMock(mock.Response{Response: http.Response{
							StatusCode: 200,
							Body:       mock.NewStringBody(`{}`),
						}}),
					},
				},
			},
		},
		{
			name:    "Set fails due to unset name",
			args:    args{},
			wantErr: errNameCannotBeEmpty,
		},
		{
			name: "Set fails due to unset config",
			args: args{
				params: SetParams{
					Name: "name",
					Type: "s3",
				},
			},
			wantErr: errConfigMustBeSet,
		},
		{
			name: "Set fails due to invalid config",
			args: args{
				params: SetParams{
					Name:   "name",
					Type:   "s3",
					Config: new(S3Config),
				},
			},
			wantErr: multierror.Append(nil, []error{
				errRegionCannotBeEmpty,
				errBucketCannotBeEmpty,
				errAccessKeyCannotBeEmpty,
				errSecretKeyCannotBeEmpty,
			}...),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Set(tt.args.params)
			if err := util.CheckErrType(err, tt.wantErr); err != nil {
				t.Error(err)
			}
		})
	}
}
