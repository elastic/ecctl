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

package deployment

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment/deputil"
	"github.com/elastic/ecctl/pkg/util"
)

func TestResourceParams_Validate(t *testing.T) {
	var error404 = models.BasicFailedReply{Errors: []*models.BasicFailedReplyElement{
		{},
	}}
	error404Byte, _ := json.MarshalIndent(error404, "", "  ")
	type fields struct {
		API          *api.API
		DeploymentID string
		Type         string
		RefID        string
	}
	tests := []struct {
		name   string
		fields fields
		err    error
	}{
		{
			name:   "fails on empty parameters",
			fields: fields{},
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				deputil.NewInvalidDeploymentIDError(""),
				errors.New("deployment resource type cannot be empty"),
				errors.New("failed auto-discovering the resource ref id: api reference is required for command"),
				errors.New(`failed auto-discovering the resource ref id: id "" is invalid`),
			}},
		},
		{
			name: "succeeds validation when refID is populated",
			fields: fields{
				API:          api.NewMock(),
				DeploymentID: util.ValidClusterID,
				Type:         "elasticsearch",
				RefID:        "main-elasticsearch",
			},
		},
		{
			name: "returns error when autodiscovery of ref-id fails",
			fields: fields{
				API:          api.NewMock(mock.New404Response(mock.NewStructBody(error404))),
				DeploymentID: util.ValidClusterID,
				Type:         "elasticsearch",
			},
			err: &multierror.Error{Errors: []error{
				fmt.Errorf("failed auto-discovering the resource ref id: %s", string(error404Byte)),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := &ResourceParams{
				API:          tt.fields.API,
				DeploymentID: tt.fields.DeploymentID,
				Type:         tt.fields.Type,
				RefID:        tt.fields.RefID,
			}

			if err := params.Validate(); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("ResourceParams.Validate() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
