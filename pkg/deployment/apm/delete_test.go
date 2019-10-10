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

package apm

import (
	"bytes"
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/output"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func TestDelete(t *testing.T) {
	type args struct {
		params DeleteParams
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "fails due to parameter validation",
			args: args{params: DeleteParams{}},
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				errors.New(`id "" is invalid`),
			}},
		},
		{
			name: "fails due to status",
			args: args{params: DeleteParams{
				ID: "d324608c97154bdba2dff97511d40368",
				API: api.NewMock(mock.Response{Response: http.Response{
					Body: mock.NewStructBody(models.ApmInfo{
						Status: "started",
					}),
					StatusCode: 200,
				}}),
			}},
			err: errors.New("apm delete: deployment must be stopped"),
		},
		{
			name: "deletes a stopped cluster",
			args: args{params: DeleteParams{
				ID: "d324608c97154bdba2dff97511d40368",
				API: api.NewMock(
					mock.Response{Response: http.Response{
						Body: mock.NewStructBody(models.ApmInfo{
							Status: "stopped",
						}),
						StatusCode: 200,
					}},
					mock.Response{Response: http.Response{
						Body:       mock.NewStringBody("{}"),
						StatusCode: 200,
					}},
				),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Delete(tt.args.params); !reflect.DeepEqual(tt.err, err) {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}

func TestShutdown(t *testing.T) {
	type args struct {
		params ShutdownParams
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "fails due to parameter validation",
			args: args{params: ShutdownParams{}},
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				errors.New(`id "" is invalid`),
			}},
		},
		{
			name: "Succeeds without tracking",
			args: args{params: ShutdownParams{
				ID: "d324608c97154bdba2dff97511d40368",
				API: api.NewMock(mock.Response{Response: http.Response{
					Body:       mock.NewStructBody(models.ClusterCommandResponse{}),
					StatusCode: 202,
				}}),
			}},
		},
		{
			name: "Succeeds with tracking",
			args: args{params: ShutdownParams{
				TrackParams: util.TrackParams{
					Track:         true,
					Output:        output.NewDevice(new(bytes.Buffer)),
					PollFrequency: time.Millisecond,
					MaxRetries:    1,
				},
				ID: "d324608c97154bdba2dff97511d40368",
				API: api.NewMock(util.AppendTrackResponses(mock.Response{Response: http.Response{
					Body:       mock.NewStructBody(models.ClusterCommandResponse{}),
					StatusCode: 202,
				}})...),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Shutdown(tt.args.params); !reflect.DeepEqual(tt.err, err) {
				t.Errorf("Shutdown() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
