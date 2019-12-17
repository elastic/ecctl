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

	"github.com/elastic/ecctl/pkg/deployment"
	"github.com/elastic/ecctl/pkg/deployment/deputil"
	"github.com/elastic/ecctl/pkg/util"
)

func TestRestore(t *testing.T) {
	var internalError = models.BasicFailedReply{
		Errors: []*models.BasicFailedReplyElement{
			{},
		},
	}
	internalErrorBytes, _ := json.MarshalIndent(internalError, "", "  ")
	type args struct {
		params RestoreParams
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "fails due to param validation",
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				deputil.NewInvalidDeploymentIDError(""),
				errors.New("deployment resource type cannot be empty"),
				errors.New("failed auto-discovering the resource ref id: api reference is required for command"),
				errors.New("failed auto-discovering the resource ref id: id \"\" is invalid"),
			}},
		},
		{
			name: "fails due to restore Kibana due to API error",
			args: args{params: RestoreParams{
				ResourceParams: deployment.ResourceParams{
					API:          api.NewMock(mock.New404Response(mock.NewStructBody(internalError))),
					DeploymentID: util.ValidClusterID,
					RefID:        "kibana",
					Type:         "kibana",
				},
			}},
			err: errors.New(string(internalErrorBytes)),
		},
		{
			name: "fails to restore APM due to API error",
			args: args{params: RestoreParams{
				ResourceParams: deployment.ResourceParams{
					API:          api.NewMock(mock.New404Response(mock.NewStructBody(internalError))),
					DeploymentID: util.ValidClusterID,
					RefID:        "apm",
					Type:         "apm",
				},
			}},
			err: errors.New(string(internalErrorBytes)),
		},
		{
			name: "fails due to restore Elasticsearch due to API error",
			args: args{params: RestoreParams{
				ResourceParams: deployment.ResourceParams{
					API:          api.NewMock(mock.New404Response(mock.NewStructBody(internalError))),
					DeploymentID: util.ValidClusterID,
					RefID:        "elasticsearch",
					Type:         "elasticsearch",
				},
			}},
			err: errors.New(string(internalErrorBytes)),
		},
		{
			name: "Succeeds restoring Kibana",
			args: args{params: RestoreParams{
				ResourceParams: deployment.ResourceParams{
					API:          api.NewMock(mock.New200Response(mock.NewStringBody(""))),
					DeploymentID: util.ValidClusterID,
					RefID:        "kibana",
					Type:         "kibana",
				},
			}},
		},
		{
			name: "Succeeds restoring APM",
			args: args{params: RestoreParams{
				ResourceParams: deployment.ResourceParams{
					API:          api.NewMock(mock.New200Response(mock.NewStringBody(""))),
					DeploymentID: util.ValidClusterID,
					RefID:        "apm",
					Type:         "apm",
				},
			}},
		},
		{
			name: "Succeeds restoring Elasticsearch",
			args: args{params: RestoreParams{
				ResourceParams: deployment.ResourceParams{
					API:          api.NewMock(mock.New200Response(mock.NewStringBody(""))),
					DeploymentID: util.ValidClusterID,
					RefID:        "elasticsearch",
					Type:         "elasticsearch",
				},
			}},
		},
		{
			name: "Succeeds restoring Elasticsearch with refID autodiscovery",
			args: args{params: RestoreParams{
				ResourceParams: deployment.ResourceParams{
					API: api.NewMock(
						mock.New200Response(mock.NewStructBody(models.DeploymentGetResponse{
							Healthy: ec.Bool(true),
							ID:      ec.String(util.ValidClusterID),
							Resources: &models.DeploymentResources{
								Elasticsearch: []*models.ElasticsearchResourceInfo{{
									ID:    ec.String(util.ValidClusterID),
									RefID: ec.String("elasticsearch"),
								}},
							},
						})),
						mock.New200Response(mock.NewStringBody("")),
					),
					DeploymentID: util.ValidClusterID,
					Type:         "elasticsearch",
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Restore(tt.args.params); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Restore() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
