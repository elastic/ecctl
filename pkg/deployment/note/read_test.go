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

	"github.com/elastic/ecctl/pkg/deployment"
	"github.com/elastic/ecctl/pkg/util"
)

func TestList(t *testing.T) {
	var messageDateTime = util.ParseDate(t, "2018-04-13T07:11:54.999Z")
	var listNotesPayload = `
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
}`[1:]
	type args struct {
		params ListParams
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Notes
		wantErr bool
	}{
		{
			name: "List notes succeeds",
			args: args{
				params: ListParams{Params: deployment.Params{
					ID: "a2c4f423c1014941b75a48292264dd25",
					API: api.NewMock(mock.Response{Response: http.Response{
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
			name: "List notes fails when an error is received",
			args: args{
				params: ListParams{Params: deployment.Params{
					ID:  "a2c4f423c1014941b75a48292264dd25",
					API: api.NewMock(mock.Response{Error: errors.New("an error")}),
				}},
			},
			wantErr: true,
		},
		{
			name: "List fails when deployment ID is empty",
			args: args{
				params: ListParams{
					Params: deployment.Params{
						API: new(api.API),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := List(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGet(t *testing.T) {
	var messageDateTime = util.ParseDate(t, "2018-04-13T07:11:54.999Z")
	var simpleNote = `{
		"id": "1",
		"message": "a message",
		"user_id": "root",
		"timestamp": "2018-04-13T07:11:54.999Z"
	}`

	type args struct {
		params GetParams
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Note
		wantErr bool
	}{
		{
			name: "Get note succeeds",
			args: args{params: GetParams{Params: Params{
				NoteID: "1",
				Params: deployment.Params{
					ID: "a2c4f423c1014941b75a48292264dd25",
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(simpleNote),
						StatusCode: 200,
					}}),
				},
			}}},
			want: &models.Note{
				ID:        "1",
				Message:   ec.String("a message"),
				UserID:    "root",
				Timestamp: messageDateTime,
			},
		},
		{
			name: "Get note fails due to client error",
			args: args{
				params: GetParams{Params: Params{
					NoteID: "1",
					Params: deployment.Params{
						ID:  "a2c4f423c1014941b75a48292264dd25",
						API: api.NewMock(mock.Response{Error: errors.New("an error")}),
					},
				}},
			},
			wantErr: true,
		},
		{
			name: "Get note fails due to empty note ID",
			args: args{
				params: GetParams{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
