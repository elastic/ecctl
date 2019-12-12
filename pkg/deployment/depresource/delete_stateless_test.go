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

func TestDeleteStateless(t *testing.T) {
	var internalError = models.BasicFailedReply{
		Errors: []*models.BasicFailedReplyElement{
			{},
		},
	}
	internalErrorBytes, _ := json.MarshalIndent(internalError, "", "  ")
	type args struct {
		params DeleteStatelessParams
	}
	tests := []struct {
		name string
		args args
		want *models.DeploymentResourceUpgradeResponse
		err  error
	}{
		{
			name: "fails due to parameter validation",
			args: args{},
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				errors.New("id \"\" is invalid"),
				errors.New("deployment resource ref id cannot be empty"),
				errors.New("deployment resource type cannot be empty"),
			}},
		},
		{
			name: "fails due to parameter validation on invalid type",
			args: args{params: DeleteStatelessParams{
				Type: "elasticsearch",
			}},
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				errors.New("id \"\" is invalid"),
				errors.New("deployment resource ref id cannot be empty"),
				errors.New("deployment resource type \"elasticsearch\" is not supported"),
			}},
		},
		{
			name: "fails due to API error",
			args: args{params: DeleteStatelessParams{
				API:          api.NewMock(mock.New404Response(mock.NewStructBody(internalError))),
				DeploymentID: util.ValidClusterID,
				RefID:        "kibana",
				Type:         "kibana",
			}},
			err: errors.New(string(internalErrorBytes)),
		},
		{
			name: "succeeds on APM resource",
			args: args{params: DeleteStatelessParams{
				API:          api.NewMock(mock.New200Response(mock.NewStringBody(""))),
				DeploymentID: util.ValidClusterID,
				RefID:        "kibana",
				Type:         "kibana",
			}},
		},
		{
			name: "fails due to API error on APM resource",
			args: args{params: DeleteStatelessParams{
				API:          api.NewMock(mock.New404Response(mock.NewStructBody(internalError))),
				DeploymentID: util.ValidClusterID,
				RefID:        "apm",
				Type:         "apm",
			}},
			err: errors.New(string(internalErrorBytes)),
		},
		{
			name: "succeeds",
			args: args{params: DeleteStatelessParams{
				API:          api.NewMock(mock.New200Response(mock.NewStringBody(""))),
				DeploymentID: util.ValidClusterID,
				RefID:        "apm",
				Type:         "apm",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DeleteStateless(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("DeleteStateless() error = %v, wantErr %v", err, tt.err)
				return
			}
		})
	}
}
