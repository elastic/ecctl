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

package deployment

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func TestGetParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  GetParams
		wantErr bool
		err     error
	}{
		{
			name:   "validate should return all possible errors",
			params: GetParams{},
			err: &multierror.Error{
				Errors: []error{
					errors.New("api reference is required for command"),
					errors.New(`id "" is invalid`),
				},
			},
			wantErr: true,
		},
		{
			name: "validate should return error on missing api",
			params: GetParams{
				DeploymentID: "f1d329b0fb34470ba8b18361cabdd2bc",
			},
			err: &multierror.Error{
				Errors: []error{
					errors.New("api reference is required for command"),
				},
			},
			wantErr: true,
		},
		{
			name: "validate should return error on invalid ID",
			params: GetParams{
				API: &api.API{},
			},
			err: &multierror.Error{
				Errors: []error{
					errors.New(`id "" is invalid`),
				},
			},
			wantErr: true,
		},
		{
			name: "validate should pass if all params are properly set",
			params: GetParams{
				API:          &api.API{},
				DeploymentID: "f1d329b0fb34470ba8b18361cabdd2bc",
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

func TestGet(t *testing.T) {
	const getResponse = `{
  "healthy": true,
  "id": "f1d329b0fb34470ba8b18361cabdd2bc"
}`

	const getErrorResponse = `{
  "errors": null
}`

	type args struct {
		params GetParams
	}
	tests := []struct {
		name    string
		args    args
		want    *models.DeploymentGetResponse
		wantErr bool
		err     error
	}{
		{
			name:    "Get fails due to parameter validation failure",
			args:    args{},
			wantErr: true,
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				errors.New(`id "" is invalid`),
			}},
		},
		{
			name: "Get fails due to API failure",
			args: args{
				params: GetParams{
					DeploymentID: "f1d329b0fb34470ba8b18361cabdd2bc",
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
				params: GetParams{
					DeploymentID: "f1d329b0fb34470ba8b18361cabdd2bc",
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(getResponse),
						StatusCode: 200,
					}}),
				},
			},
			want: &models.DeploymentGetResponse{
				Healthy: ec.Bool(true),
				ID:      ec.String("f1d329b0fb34470ba8b18361cabdd2bc"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.err != nil && err.Error() != tt.err.Error() {
				t.Errorf("Get() actual error = '%v', want error '%v'", err, tt.err)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAppSearch(t *testing.T) {
	const getAppSearchResponse = `{
  "elasticsearch_cluster_ref_id": "main-elasticsearch",
  "id": "3531aaf988594efa87c1aabb7caed337",
  "ref_id": "main-appsearch"
}`

	const getErrorResponse = `{
  "errors": null
}`

	type args struct {
		params GetParams
	}
	tests := []struct {
		name    string
		args    args
		want    *models.AppSearchResourceInfo
		wantErr bool
		err     error
	}{
		{
			name:    "Get fails due to parameter validation failure",
			args:    args{},
			wantErr: true,
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				errors.New(`id "" is invalid`),
			}},
		},
		{
			name: "Get fails due to API failure",
			args: args{
				params: GetParams{
					DeploymentID: "f1d329b0fb34470ba8b18361cabdd2bc",
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
				params: GetParams{
					DeploymentID: "f1d329b0fb34470ba8b18361cabdd2bc",
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(getAppSearchResponse),
						StatusCode: 200,
					}}),
				},
			},
			want: &models.AppSearchResourceInfo{
				ElasticsearchClusterRefID: ec.String("main-elasticsearch"),
				ID:                        ec.String("3531aaf988594efa87c1aabb7caed337"),
				RefID:                     ec.String("main-appsearch"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAppSearch(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.err != nil && err.Error() != tt.err.Error() {
				t.Errorf("Get() actual error = '%v', want error '%v'", err, tt.err)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetApm(t *testing.T) {
	const getApmResponse = `{
  "elasticsearch_cluster_ref_id": "main-elasticsearch",
  "id": "3531aaf988594efa87c1aabb7caed337",
  "ref_id": "main-apm"
}`

	const getErrorResponse = `{
  "errors": null
}`

	type args struct {
		params GetParams
	}
	tests := []struct {
		name    string
		args    args
		want    *models.ApmResourceInfo
		wantErr bool
		err     error
	}{
		{
			name:    "Get fails due to parameter validation failure",
			args:    args{},
			wantErr: true,
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				errors.New(`id "" is invalid`),
			}},
		},
		{
			name: "Get fails due to API failure",
			args: args{
				params: GetParams{
					DeploymentID: "f1d329b0fb34470ba8b18361cabdd2bc",
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
				params: GetParams{
					DeploymentID: "f1d329b0fb34470ba8b18361cabdd2bc",
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(getApmResponse),
						StatusCode: 200,
					}}),
				},
			},
			want: &models.ApmResourceInfo{
				ElasticsearchClusterRefID: ec.String("main-elasticsearch"),
				ID:                        ec.String("3531aaf988594efa87c1aabb7caed337"),
				RefID:                     ec.String("main-apm"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetApm(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.err != nil && err.Error() != tt.err.Error() {
				t.Errorf("Get() actual error = '%v', want error '%v'", err, tt.err)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetElasticsearch(t *testing.T) {
	const getElasticsearchResponse = `{
  "id": "f1d329b0fb34470ba8b18361cabdd2bc",
  "ref_id": "main-elasticsearch"
}`

	const getErrorResponse = `{
  "errors": null
}`

	type args struct {
		params GetParams
	}
	tests := []struct {
		name    string
		args    args
		want    *models.ElasticsearchResourceInfo
		wantErr bool
		err     error
	}{
		{
			name:    "Get fails due to parameter validation failure",
			args:    args{},
			wantErr: true,
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				errors.New(`id "" is invalid`),
			}},
		},
		{
			name: "Get fails due to API failure",
			args: args{
				params: GetParams{
					DeploymentID: "f1d329b0fb34470ba8b18361cabdd2bc",
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
				params: GetParams{
					DeploymentID: "f1d329b0fb34470ba8b18361cabdd2bc",
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(getElasticsearchResponse),
						StatusCode: 200,
					}}),
				},
			},
			want: &models.ElasticsearchResourceInfo{
				ID:    ec.String("f1d329b0fb34470ba8b18361cabdd2bc"),
				RefID: ec.String("main-elasticsearch"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetElasticsearch(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.err != nil && err.Error() != tt.err.Error() {
				t.Errorf("Get() actual error = '%v', want error '%v'", err, tt.err)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetKibana(t *testing.T) {
	const getKibanaResponse = `{
  "elasticsearch_cluster_ref_id": "main-elasticsearch",
  "id": "3531aaf988594efa87c1aabb7caed337",
  "ref_id": "main-kibana"
}`

	const getErrorResponse = `{
  "errors": null
}`

	type args struct {
		params GetParams
	}
	tests := []struct {
		name    string
		args    args
		want    *models.KibanaResourceInfo
		wantErr bool
		err     error
	}{
		{
			name:    "Get fails due to parameter validation failure",
			args:    args{},
			wantErr: true,
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				errors.New(`id "" is invalid`),
			}},
		},
		{
			name: "Get fails due to API failure",
			args: args{
				params: GetParams{
					DeploymentID: "f1d329b0fb34470ba8b18361cabdd2bc",
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
				params: GetParams{
					DeploymentID: "f1d329b0fb34470ba8b18361cabdd2bc",
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(getKibanaResponse),
						StatusCode: 200,
					}}),
				},
			},
			want: &models.KibanaResourceInfo{
				ElasticsearchClusterRefID: ec.String("main-elasticsearch"),
				ID:                        ec.String("3531aaf988594efa87c1aabb7caed337"),
				RefID:                     ec.String("main-kibana"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetKibana(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.err != nil && err.Error() != tt.err.Error() {
				t.Errorf("Get() actual error = '%v', want error '%v'", err, tt.err)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		params GetParams
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
				params: GetParams{
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
				params: GetParams{
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
			got, err := GetElasticsearchID(tt.args.params)
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
