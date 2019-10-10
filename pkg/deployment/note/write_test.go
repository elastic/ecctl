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
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
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
				API: api.NewMock(mock.Response{
					Response: http.Response{
						StatusCode: http.StatusCreated,
						Status:     http.StatusText(http.StatusCreated),
						Body:       mock.NewStringBody(`{}`),
					},
				}),
				ID:      util.ValidClusterID,
				UserID:  "someid",
				Message: "note message",
				Type:    "elasticsearch",
			}},
		},
		{
			name: "Succeeds posting an Elasticsearch note (Commentator wrapped)",
			args: args{params: AddParams{
				API: api.NewMock(mock.Response{
					Response: http.Response{
						StatusCode: http.StatusCreated,
						Status:     http.StatusText(http.StatusCreated),
						Body:       mock.NewStringBody(`{}`),
					},
				}),
				ID:          util.ValidClusterID,
				UserID:      "someid",
				Message:     "note message",
				Type:        "elasticsearch",
				Commentator: &c{m: "somemessage"},
			}},
		},
		{
			name: "Succeeds posting a Kibana note",
			args: args{params: AddParams{
				API: api.NewMock(
					mock.Response{
						Response: http.Response{
							StatusCode: http.StatusOK,
							Status:     http.StatusText(http.StatusOK),
							Body: mock.NewStructBody(models.KibanaClusterInfo{
								ElasticsearchCluster: &models.TargetElasticsearchCluster{
									ElasticsearchID: ec.String("420b8b540dfc967a7a649c18e2fce4ed"),
								},
							}),
						},
					},
					mock.Response{
						Response: http.Response{
							StatusCode: http.StatusCreated,
							Status:     http.StatusText(http.StatusCreated),
							Body:       mock.NewStringBody(`{}`),
						},
					},
				),
				ID:      util.ValidClusterID,
				UserID:  "someid",
				Message: "note message",
				Type:    "kibana",
			}},
		},
		{
			name: "Fails posting an Kibana note (Fails to get deployment)",
			args: args{params: AddParams{
				API: api.NewMock(
					mock.Response{
						Response: http.Response{
							StatusCode: http.StatusNotFound,
							Status:     http.StatusText(http.StatusNotFound),
							Body:       mock.NewStringBody(`{}`),
						},
					},
				),
				ID:      util.ValidClusterID,
				UserID:  "someid",
				Message: "note message",
				Type:    "kibana",
			}},
			err: errors.New(errNull),
		},
		{
			name: "Succeeds posting an Apm note",
			args: args{params: AddParams{
				API: api.NewMock(
					mock.Response{
						Response: http.Response{
							StatusCode: http.StatusOK,
							Status:     http.StatusText(http.StatusOK),
							Body: mock.NewStructBody(models.ApmInfo{
								ElasticsearchCluster: &models.TargetElasticsearchCluster{
									ElasticsearchID: ec.String("420b8b540dfc967a7a649c18e2fce4ed"),
								},
							}),
						},
					},
					mock.Response{
						Response: http.Response{
							StatusCode: http.StatusCreated,
							Status:     http.StatusText(http.StatusCreated),
							Body:       mock.NewStringBody(`{}`),
						},
					},
				),
				ID:      util.ValidClusterID,
				UserID:  "someid",
				Message: "note message",
				Type:    "apm",
			}},
		},
		{
			name: "Fails posting an Apm note (Fails to get deployment)",
			args: args{params: AddParams{
				API: api.NewMock(
					mock.Response{
						Response: http.Response{
							StatusCode: http.StatusNotFound,
							Status:     http.StatusText(http.StatusNotFound),
							Body:       mock.NewStringBody(`{}`),
						},
					},
				),
				ID:      util.ValidClusterID,
				UserID:  "someid",
				Message: "note message",
				Type:    "apm",
			}},
			err: errors.New(errNull),
		},
		{
			name: "Fails due to parameter validation (empty params)",
			args: args{params: AddParams{}},
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				errors.New("invalid id"),
				errors.New("user id cannot be empty"),
				errors.New("invalid type : valid types are [elasticsearch kibana apm]"),
				errors.New("message cannot be empty"),
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

func TestUpdate(t *testing.T) {
	var messageDateTime = util.ParseDate(t, "2018-04-13T07:11:54.999Z")
	var simpleNoteModified = `{
		"id": "1",
		"message": "a modified message",
		"user_id": "root",
		"timestamp": "2018-04-13T07:11:54.999Z",
		"version": 2
	}`
	type args struct {
		params UpdateParams
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Note
		wantErr bool
	}{
		{
			name: "Update note succeeds",
			args: args{params: UpdateParams{
				UserID:  "marc",
				Message: "a modified message",
				Params: Params{
					NoteID: "1",
					Params: deployment.Params{
						ID: "a2c4f423c1014941b75a48292264dd25",
						API: api.NewMock(mock.Response{Response: http.Response{
							Body:       mock.NewStringBody(simpleNoteModified),
							StatusCode: 200,
						}}),
					},
				},
			}},
			want: &models.Note{
				ID:        "1",
				Message:   ec.String("a modified message"),
				UserID:    "root",
				Timestamp: messageDateTime,
			},
		},
		{
			name: "Update note fails due to client error",
			args: args{params: UpdateParams{
				Message: "a modified message",
				Params: Params{
					NoteID: "1",
					Params: deployment.Params{
						ID:  "a2c4f423c1014941b75a48292264dd25",
						API: api.NewMock(mock.Response{Error: errors.New("an error")}),
					},
				},
			}},
			wantErr: true,
		},
		{
			name: "Update note fails due to empty note ID",
			args: args{
				params: UpdateParams{
					Params: Params{
						Params: deployment.Params{
							ID:  "a2c4f423c1014941b75a48292264dd25",
							API: new(api.API),
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Update note fails due to empty message ",
			args: args{
				params: UpdateParams{
					UserID: "root",
					Params: Params{
						NoteID: "1",
						Params: deployment.Params{
							ID:  "a2c4f423c1014941b75a48292264dd25",
							API: new(api.API),
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Update(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() = %v, want %v", got, tt.want)
			}
		})
	}
}
