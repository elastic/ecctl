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
	"bytes"
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func TestDiagnose(t *testing.T) {
	type args struct {
		params DiagnoseParams
		buf    *bytes.Buffer
	}
	tests := []struct {
		name     string
		args     args
		wantRead string
		err      error
	}{
		{
			name: "succeeds",
			args: args{
				params: DiagnoseParams{
					ClusterParams: util.ClusterParams{
						ClusterID: util.ValidClusterID,
						API: api.NewMock(mock.Response{
							Response: http.Response{
								Header: http.Header{
									"Content-Type": []string{"application/zip"},
								},
								StatusCode: 200,
								Body:       mock.NewStringBody("something diagnosed"),
							},
						}),
					},
				},
				buf: new(bytes.Buffer),
			},
			wantRead: "something diagnosed",
		},
		{
			name: "fails due to API error",
			args: args{
				params: DiagnoseParams{
					ClusterParams: util.ClusterParams{
						ClusterID: util.ValidClusterID,
						API: api.NewMock(mock.Response{
							Response: http.Response{
								Header: http.Header{
									"Content-Type": []string{"application/zip"},
								},
								StatusCode: 500,
								Body:       mock.NewStringBody("something diagnosed"),
							},
						}),
					},
				},
				buf: new(bytes.Buffer),
			},
			err: errors.New("unknown error (status 500)"),
		},
		{
			name: "empty params return error",
			args: args{params: DiagnoseParams{
				ClusterParams: util.ClusterParams{},
			}},
			err: &multierror.Error{Errors: []error{
				errors.New("writer cannot be nil"),
				errors.New("cluster id should have a length of 32 characters"),
			}},
		},
		{
			name: "missing API returns error",
			args: args{params: DiagnoseParams{
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
				},
			}},
			err: &multierror.Error{Errors: []error{
				errors.New("writer cannot be nil"),
				errors.New("api reference is required for command"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.buf != nil {
				tt.args.params.Writer = tt.args.buf
			}
			if err := Diagnose(tt.args.params); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Diagnose() error = %v, wantErr %v", err, tt.err)
			}

			if tt.wantRead != "" {
				if tt.args.buf == nil {
					t.Errorf("Diagnose() error = empty buffer")
					return
				}
				if tt.args.buf.String() != tt.wantRead {
					t.Errorf("Diagnose() read = %v, wantRead %v", tt.args.buf.String(), tt.wantRead)
				}
			}
		})
	}
}
