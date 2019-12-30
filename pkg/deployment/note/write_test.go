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
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment"
	"github.com/elastic/ecctl/pkg/util"
)

var errNull = `{
  "errors": null
}`

type c struct{ m string }

func (cm *c) Message(s string) string {
	return fmt.Sprint(cm.m, " ", s)
}

func (cm *c) Set(s string) {}

func TestAdd(t *testing.T) {
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

	type args struct {
		params AddParams
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Succeeds posting an Elasticsearch note",
			args: args{params: AddParams{
				Params: deployment.Params{
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(getResponse),
						StatusCode: 200,
					}},
						mock.Response{Response: http.Response{
							StatusCode: http.StatusCreated,
							Status:     http.StatusText(http.StatusCreated),
							Body:       mock.NewStringBody(`{}`),
						},
						}),
					ID: "e3dac8bf3dc64c528c295a94d0f19a77",
				},
				UserID:  "someid",
				Message: "note message",
			}},
		},
		{
			name: "Succeeds posting an Elasticsearch note (Commentator wrapped)",
			args: args{params: AddParams{
				Params: deployment.Params{
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(getResponse),
						StatusCode: 200,
					}},
						mock.Response{
							Response: http.Response{
								StatusCode: http.StatusCreated,
								Status:     http.StatusText(http.StatusCreated),
								Body:       mock.NewStringBody(`{}`),
							},
						}),
					ID: "e3dac8bf3dc64c528c295a94d0f19a77",
				},
				UserID:      "someid",
				Message:     "note message",
				Commentator: &c{m: "somemessage"},
			}},
		},
		{
			name: "Fails posting note (Fails to get deployment)",
			args: args{params: AddParams{
				Params: deployment.Params{
					API: api.NewMock(
						mock.Response{
							Response: http.Response{
								StatusCode: http.StatusNotFound,
								Status:     http.StatusText(http.StatusNotFound),
								Body:       mock.NewStringBody(`{}`),
							},
						},
					),
					ID: util.ValidClusterID,
				},
				UserID:  "someid",
				Message: "note message",
			}},
			err: errors.New(errNull),
		},
		{
			name: "Fails due to parameter validation (empty params)",
			args: args{params: AddParams{}},
			err: &multierror.Error{Errors: []error{
				errors.New("user id cannot be empty"),
				errors.New("note message cannot be empty"),
				util.ErrAPIReq,
				errors.New(`id "" is invalid`),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Add(tt.args.params); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
