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
	"fmt"
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

func TestShutdown(t *testing.T) {
	var error404 = models.BasicFailedReply{Errors: []*models.BasicFailedReplyElement{
		{},
	}}
	error404Byte, _ := json.MarshalIndent(error404, "", "  ")
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
			args: args{},
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				deputil.NewInvalidDeploymentIDError(""),
				errors.New("deployment resource kind cannot be empty"),
				errors.New("failed auto-discovering the resource ref id: api reference is required for command"),
				errors.New("failed auto-discovering the resource ref id: id \"\" is invalid"),
			}},
		},
		{
			name: "Returns error on kind Elasticsearch and a received API error",
			args: args{params: ShutdownParams{
				ResourceParams: deployment.ResourceParams{
					API:          api.NewMock(mock.New404Response(mock.NewStructBody(error404))),
					DeploymentID: util.ValidClusterID,
					RefID:        "elasticsearch",
					Kind:         "elasticsearch",
				},
			}},
			err: errors.New(string(error404Byte)),
		},
		{
			name: "Returns error on kind APM and a received API error",
			args: args{params: ShutdownParams{
				ResourceParams: deployment.ResourceParams{
					API:          api.NewMock(mock.New404Response(mock.NewStructBody(error404))),
					DeploymentID: util.ValidClusterID,
					RefID:        "apm",
					Kind:         "apm",
				},
			}},
			err: errors.New(string(error404Byte)),
		},
		{
			name: "Returns error on kind Kibana and a received API error",
			args: args{params: ShutdownParams{
				ResourceParams: deployment.ResourceParams{
					API:          api.NewMock(mock.New404Response(mock.NewStructBody(error404))),
					DeploymentID: util.ValidClusterID,
					RefID:        "kibana",
					Kind:         "kibana",
				},
			}},
			err: errors.New(string(error404Byte)),
		},
		{
			name: "Succeeds on kind Elasticsearch with autodiscover of the kind",
			args: args{params: ShutdownParams{
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
						mock.New200Response(mock.NewStructBody(struct{}{})),
					),
					DeploymentID: util.ValidClusterID,
					Kind:         "elasticsearch",
				},
			}},
		},
		{
			name: "Fails on kind Elasticsearch when autodiscover returns an error",
			args: args{params: ShutdownParams{
				ResourceParams: deployment.ResourceParams{
					API:          api.NewMock(mock.New404Response(mock.NewStructBody(error404))),
					DeploymentID: util.ValidClusterID,
					Kind:         "elasticsearch",
				},
			}},
			err: &multierror.Error{Errors: []error{
				fmt.Errorf("failed auto-discovering the resource ref id: %s", string(error404Byte)),
			}},
		},
		{
			name: "Succeeds on kind Elasticsearch",
			args: args{params: ShutdownParams{
				ResourceParams: deployment.ResourceParams{
					API:          api.NewMock(mock.New200Response(mock.NewStructBody(struct{}{}))),
					DeploymentID: util.ValidClusterID,
					RefID:        "elasticsearch",
					Kind:         "elasticsearch",
				},
			}},
		},
		{
			name: "Succeeds on kind APM",
			args: args{params: ShutdownParams{
				ResourceParams: deployment.ResourceParams{
					API:          api.NewMock(mock.New200Response(mock.NewStructBody(struct{}{}))),
					DeploymentID: util.ValidClusterID,
					RefID:        "apm",
					Kind:         "apm",
				},
			}},
		},
		{
			name: "Succeeds on kind Kibana",
			args: args{params: ShutdownParams{
				ResourceParams: deployment.ResourceParams{
					API:          api.NewMock(mock.New200Response(mock.NewStructBody(struct{}{}))),
					DeploymentID: util.ValidClusterID,
					RefID:        "kibana",
					Kind:         "kibana",
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Shutdown(tt.args.params); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Shutdown() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
