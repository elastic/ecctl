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

package depresource

import (
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/output"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/hashicorp/go-multierror"
)

// Note: This test is not fully testing the track behaviour since that covered
// in another package already. This focuses on covering the error paths.
func TestTrackResources(t *testing.T) {
	type args struct {
		params TrackResourcesParams
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "fails on parameter validation",
			err: &multierror.Error{Errors: []error{
				errors.New("api reference is required for command"),
				errors.New("resource track: output device cannot be nil"),
			}},
		},
		{
			name: "returns without an error if the resources are nil",
			args: args{params: TrackResourcesParams{
				API:          api.NewMock(),
				OutputDevice: output.NewDevice(ioutil.Discard),
			}},
		},
		{
			name: "returns with fail to track when the resource is of type appsearch",
			args: args{params: TrackResourcesParams{
				API:          api.NewMock(),
				OutputDevice: output.NewDevice(ioutil.Discard),
				Resources: []*models.DeploymentResource{
					{ID: ec.String("an id"), Kind: ec.String("appsearch")},
				},
			}},
			err: &multierror.Error{Errors: []error{
				fmt.Errorf("cannot track appsearch resource id %s", "an id"),
			}},
		},
		{
			name: "returns an error due to parameter validation",
			args: args{params: TrackResourcesParams{
				API:          api.NewMock(mock.New500Response(mock.NewStructBody(errors.New("error")))),
				OutputDevice: output.NewDevice(ioutil.Discard),
				Resources: []*models.DeploymentResource{
					{ID: ec.String("invalid id"), Kind: ec.String("elasticsearch")},
				},
			}},
			err: &multierror.Error{Errors: []error{
				fmt.Errorf("plan Track: invalid ID"),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := TrackResources(tt.args.params); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("TrackResources() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
