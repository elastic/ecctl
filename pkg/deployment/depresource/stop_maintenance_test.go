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
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func TestStopMaintenanceMode(t *testing.T) {
	var internalError = models.BasicFailedReply{
		Errors: []*models.BasicFailedReplyElement{
			{},
		},
	}
	internalErrorBytes, _ := json.MarshalIndent(internalError, "", "  ")

	type args struct {
		params StopParams
	}

	tests := []struct {
		name string
		args args
		want models.DeploymentResourceCommandResponse
		err  error
	}{
		{
			name: "fails due to parameter validation",
			args: args{},
			err: &multierror.Error{Errors: []error{
				errors.New("api reference is required for command"),
				errors.New(`id "" is invalid`),
				errors.New(`"" is not a valid resource type. Accepted resource types are: [elasticsearch kibana apm appsearch]`),
				errors.New("deployment stop: a ref_id must be provided"),
			}},
		},
		{
			name: "fails due to API error",
			args: args{params: StopParams{
				API:          api.NewMock(mock.New404Response(mock.NewStructBody(internalError))),
				DeploymentID: util.ValidClusterID,
				Type:         "elasticsearch",
				RefID:        "main-elasticsearch",
			}},
			err: errors.New(string(internalErrorBytes)),
		},
		{
			name: "succeeds",
			args: args{params: StopParams{
				API:          api.NewMock(mock.New200Response(mock.NewStringBody(""))),
				DeploymentID: util.ValidClusterID,
				Type:         "elasticsearch",
				RefID:        "main-elasticsearch",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StopMaintenanceMode(tt.args.params)
			if tt.err != nil && err.Error() != tt.err.Error() {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStopInstancesMaintenanceMode(t *testing.T) {
	var internalError = models.BasicFailedReply{
		Errors: []*models.BasicFailedReplyElement{
			{},
		},
	}
	internalErrorBytes, _ := json.MarshalIndent(internalError, "", "  ")

	type args struct {
		params StopInstancesParams
	}

	tests := []struct {
		name string
		args args
		want models.DeploymentResourceCommandResponse
		err  error
	}{
		{
			name: "fails due to parameter validation",
			args: args{},
			err: &multierror.Error{Errors: []error{
				errors.New("deployment stop: at least 1 instance ID must be provided"),
				errors.New("api reference is required for command"),
				errors.New(`id "" is invalid`),
				errors.New(`"" is not a valid resource type. Accepted resource types are: [elasticsearch kibana apm appsearch]`),
				errors.New("deployment stop: a ref_id must be provided"),
			}},
		},
		{
			name: "fails due to API error",
			args: args{params: StopInstancesParams{
				StopParams: StopParams{
					API:          api.NewMock(mock.New404Response(mock.NewStructBody(internalError))),
					DeploymentID: util.ValidClusterID,
					Type:         "elasticsearch",
					RefID:        "main-elasticsearch",
				},
				InstanceIDs: []string{"instance-0000000001", "instance-0000000002"},
			}},
			err: errors.New(string(internalErrorBytes)),
		},
		{
			name: "succeeds",
			args: args{params: StopInstancesParams{
				StopParams: StopParams{
					API:          api.NewMock(mock.New200Response(mock.NewStringBody(""))),
					DeploymentID: util.ValidClusterID,
					Type:         "elasticsearch",
					RefID:        "main-elasticsearch",
				},
				InstanceIDs: []string{"instance-0000000001", "instance-0000000002"},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StopInstancesMaintenanceMode(tt.args.params)
			if tt.err != nil && err.Error() != tt.err.Error() {
				t.Errorf("Stop() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Stop() = %v, want %v", got, tt.want)
			}
		})
	}
}
