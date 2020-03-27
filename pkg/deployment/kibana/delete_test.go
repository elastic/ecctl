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

package kibana

import (
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func TestDelete(t *testing.T) {
	type args struct {
		params DeploymentParams
	}

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "fails due to parameter validation",
			args: args{params: DeploymentParams{}},
			wantErr: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				errors.New(`id "" is invalid`),
			}},
		},
		{
			name: "fails deleting",
			args: args{params: DeploymentParams{
				ID: "2c221bd86b7f48959a59ee3128d5c5e8",
				API: api.NewMock(mock.Response{
					Error: errors.New("error with API"),
				}),
			}},
			wantErr: &url.Error{
				Op:  "Get",
				URL: "https://mock-host/mock-path/clusters/kibana/2c221bd86b7f48959a59ee3128d5c5e8?convert_legacy_plans=false&show_metadata=false&show_plan_defaults=false&show_plan_logs=false&show_plans=false&show_settings=false",
				Err: errors.New("error with API"),
			},
		},
		{
			name: "succeeds deleting a stopped cluster",
			args: args{params: DeploymentParams{
				ID: "2c221bd86b7f48959a59ee3128d5c5e8",
				API: api.NewMock(
					mock.Response{Response: http.Response{
						StatusCode: 200,
						Body: mock.NewStructBody(models.KibanaClusterInfo{
							Status: ec.String("stopped"),
						}),
					}},
					mock.Response{Response: http.Response{
						StatusCode: 200,
						Body:       mock.NewStringBody("{}"),
					}},
				),
			}},
		},
		{
			name: "succeeds deleting with tracking",
			args: args{params: DeploymentParams{
				ID:                "2c221bd86b7f48959a59ee3128d5c5e8",
				Track:             true,
				TrackChangeParams: util.NewMockTrackChangeParams(""),
				API: api.NewMock(
					util.AppendTrackResponses(
						mock.Response{Response: http.Response{
							StatusCode: 200,
							Body: mock.NewStructBody(models.KibanaClusterInfo{
								Status: ec.String("stopped"),
							}),
						}},
						mock.Response{Response: http.Response{
							StatusCode: 200,
							Body:       mock.NewStringBody("{}"),
						}},
					)...,
				),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Delete(tt.args.params); !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
