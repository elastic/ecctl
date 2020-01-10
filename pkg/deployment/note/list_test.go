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
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment"
	"github.com/elastic/ecctl/pkg/util"
)

func TestList(t *testing.T) {
	var messageDateTime = util.ParseDate(t, "2018-04-13T07:11:54.999Z")
	const listNotesPayload = `
	{
		"version": 2,
		"notes": [
			{
				"id": "1",
				"message": "a message",
				"user_id": "root",
                "timestamp": "2018-04-13T07:11:54.999Z"
			},
			{
				"id": "2",
				"message": "another message",
				"user_id": "marc"
			}
		]
}`
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
		params Params
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Notes
		wantErr error
	}{
		{
			name: "List notes succeeds",
			args: args{
				params: Params{deployment.Params{
					ID: "e3dac8bf3dc64c528c295a94d0f19a77",
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(getResponse),
						StatusCode: 200,
					}},
						mock.Response{Response: http.Response{
							Body:       mock.NewStringBody(listNotesPayload),
							StatusCode: 200,
						}}),
				}},
			},
			want: &models.Notes{
				Notes: []*models.Note{
					{
						ID:        "1",
						Message:   ec.String("a message"),
						UserID:    "root",
						Timestamp: messageDateTime,
					},
					{
						ID:      "2",
						Message: ec.String("another message"),
						UserID:  "marc",
					},
				},
			},
		},
		{
			name: "List notes fails when an error is received (fails to get deployment)",
			args: args{
				params: Params{deployment.Params{
					ID: "a2c4f423c1014941b75a48292264dd25",
					API: api.NewMock(mock.Response{
						Response: http.Response{
							StatusCode: http.StatusNotFound,
							Status:     http.StatusText(http.StatusNotFound),
							Body:       mock.NewStringBody(`{}`),
						},
					}),
				}},
			},
			wantErr: errors.New(errNull),
		},
		{
			name: "List notes fails when an api error is received",
			args: args{
				params: Params{deployment.Params{
					ID: "a2c4f423c1014941b75a48292264dd25",
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(getResponse),
						StatusCode: 200,
					}},
						mock.Response{
							Response: http.Response{
								StatusCode: http.StatusNotFound,
								Status:     http.StatusText(http.StatusNotFound),
								Body:       mock.NewStringBody(`{}`),
							},
						}),
				}},
			},
			wantErr: errors.New(errNull),
		},
		{
			name: "List fails due to validation",
			args: args{params: Params{}},
			wantErr: &multierror.Error{Errors: []error{
				errors.New("api reference is required for command"),
				errors.New(`id "" is invalid`),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := List(tt.args.params)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}
