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
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func TestUpgradeStatelessParams_FillDefaults(t *testing.T) {
	tests := []struct {
		name    string
		params  UpgradeStatelessParams
		wantErr bool
		err     error
	}{
		{
			name: "fill defaults should return error on missing api",
			params: UpgradeStatelessParams{
				DeploymentID: "f1d329b0fb34470ba8b18361cabdd2bc",
				Type:         "elasticsearch",
			},
			err: &multierror.Error{
				Errors: []error{
					errors.New("api reference is required for command"),
				},
			},
			wantErr: true,
		},
		{
			name: "fillDefaults should succeed when RefID is set",
			params: UpgradeStatelessParams{
				RefID: "main-elasticsearch",
			},
			wantErr: false,
		},
		{
			name: "validate should pass if all params are properly set and RefID is empty",
			params: UpgradeStatelessParams{
				API: api.NewMock(
					mock.New200Response(mock.NewStructBody(models.DeploymentGetResponse{
						Healthy: ec.Bool(true),
						ID:      ec.String("3531aaf988594efa87c1aabb7caed337"),
						Resources: &models.DeploymentResources{
							Elasticsearch: []*models.ElasticsearchResourceInfo{{
								ID:    ec.String("3531aaf988594efa87c1aabb7caed337"),
								RefID: ec.String("elasticsearch"),
							}},
						},
					})),
				),
				DeploymentID: "f1d329b0fb34470ba8b18361cabdd2bc",
				Type:         "elasticsearch",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.fillDefaults()

			if (err != nil) != tt.wantErr {
				t.Errorf("fillDefaults() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.err == nil {
				t.Errorf("fillDefaults() expected errors = '%v' but no errors returned", tt.err)
			}

			if tt.wantErr && err.Error() != tt.err.Error() {
				t.Errorf("fillDefaults() expected errors = '%v' but got %v", tt.err, err)
			}
		})
	}
}

func TestUpgradeStateless(t *testing.T) {
	var internalError = models.BasicFailedReply{
		Errors: []*models.BasicFailedReplyElement{
			{},
		},
	}
	internalErrorBytes, _ := json.MarshalIndent(internalError, "", "  ")
	type args struct {
		params UpgradeStatelessParams
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
				errors.New("deployment resource type cannot be empty"),
			}},
		},
		{
			name: "fails due to API error",
			args: args{params: UpgradeStatelessParams{
				API:          api.NewMock(mock.New404Response(mock.NewStructBody(internalError))),
				DeploymentID: util.ValidClusterID,
				Type:         "kibana",
			}},
			err: errors.New(string(internalErrorBytes)),
		},
		{
			name: "succeeds",
			args: args{params: UpgradeStatelessParams{
				API:          api.NewMock(mock.New202Response(mock.NewStringBody(""))),
				DeploymentID: util.ValidClusterID,
				Type:         "kibana",
				RefID:        "main-kibana",
			}},
			want: new(models.DeploymentResourceUpgradeResponse),
		},
		{
			name: "succeeds when RefID is not set",
			args: args{params: UpgradeStatelessParams{
				API: api.NewMock(
					mock.New200Response(mock.NewStructBody(models.DeploymentGetResponse{
						Healthy: ec.Bool(true),
						ID:      ec.String("3531aaf988594efa87c1aabb7caed337"),
						Resources: &models.DeploymentResources{
							Elasticsearch: []*models.ElasticsearchResourceInfo{{
								ID:    ec.String("3531aaf988594efa87c1aabb7caed337"),
								RefID: ec.String("elasticsearch"),
							}},
						},
					})),
					mock.New200Response(mock.NewStringBody("")),
				),
				DeploymentID: util.ValidClusterID,
				Type:         "elasticsearch",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UpgradeStateless(tt.args.params)
			if tt.err != nil && err.Error() != tt.err.Error() {
				t.Errorf("UpgradeStateless() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpgradeStateless() = %v, want %v", got, tt.want)
			}
		})
	}
}
