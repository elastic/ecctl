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
	"errors"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func TestReallocate(t *testing.T) {
	type args struct {
		params ReallocateParams
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Validate fails due to empty params",
			args: args{params: ReallocateParams{
				ClusterParams: util.ClusterParams{},
			}},
			err: &multierror.Error{Errors: []error{
				errors.New("cluster id should have a length of 32 characters"),
			}},
		},
		{
			name: "Fails obtaining the Elasticsearch instances",
			args: args{params: ReallocateParams{
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API:       api.NewMock(mock.New500Response(mock.NewStringBody(`{"error": "failed obtaining instances"}`))),
				},
			}},
			err: errors.New(`{"error": "failed obtaining instances"}`),
		},
		{
			name: "Fails performing the move operation",
			args: args{params: ReallocateParams{
				TrackChangeParams: util.NewMockTrackChangeParams(""),
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API: api.NewMock(
						mock.New200StructResponse(models.ElasticsearchClusterInfo{
							Topology: &models.ClusterTopologyInfo{
								Instances: []*models.ClusterInstanceInfo{
									{InstanceName: ec.String("instance-00000001")},
								},
							},
						}),
						mock.New500Response(mock.NewStringBody(`{"error": "failed moving instances"}`)),
					),
				},
			}},
			err: errors.New(`{"error": "failed moving instances"}`),
		},
		{
			name: "Succeeds when no instances are specified",
			args: args{params: ReallocateParams{
				TrackChangeParams: util.NewMockTrackChangeParams(""),
				ClusterParams: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API: api.NewMock(
						mock.New200StructResponse(models.ElasticsearchClusterInfo{
							Topology: &models.ClusterTopologyInfo{
								Instances: []*models.ClusterInstanceInfo{
									{InstanceName: ec.String("instance-00000001")},
								},
							},
						}),
						mock.New202Response(mock.NewStringBody("")),
					),
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Reallocate(tt.args.params); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Reallocate() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
