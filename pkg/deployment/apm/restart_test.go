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
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func TestRestart(t *testing.T) {
	type args struct {
		params RestartParams
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "fails due to parameter validation",
			args: args{params: RestartParams{}},
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				errors.New(`id "" is invalid`),
			}},
		},
		{
			name: "Succeeds without tracking",
			args: args{params: RestartParams{
				ID: "d324608c97154bdba2dff97511d40368",
				API: api.NewMock(mock.Response{Response: http.Response{
					Body:       mock.NewStructBody(models.ClusterCommandResponse{}),
					StatusCode: 202,
				}}),
			}},
		},
		{
			name: "Succeeds with tracking",
			args: args{params: RestartParams{
				TrackChangeParams: util.NewMockTrackChangeParams("d324608c97154bdba2dff97511d40368"),
				ID:                "d324608c97154bdba2dff97511d40368",
				API: api.NewMock(util.AppendTrackResponses(mock.Response{Response: http.Response{
					Body:       mock.NewStructBody(models.ClusterCommandResponse{}),
					StatusCode: 202,
				}})...),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Restart(tt.args.params); !reflect.DeepEqual(tt.err, err) {
				t.Errorf("Restart() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
