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

package instances

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func NewOverrideResponse(t *testing.T, rawCap int32) io.ReadCloser {
	b, err := json.Marshal(&models.ElasticsearchClusterInstanceSettingsOverrides{InstanceCapacity: rawCap})
	if err != nil {
		t.Fatal(err)
	}
	return ioutil.NopCloser(bytes.NewReader(b))
}

func TestGetSizeFromTopology(t *testing.T) {
	type args struct {
		topology *models.ElasticsearchClusterTopologyElement
	}
	tests := []struct {
		name string
		args args
		want TopologySize
		err  error
	}{
		{
			name: "Legacy topology",
			args: args{
				topology: &models.ElasticsearchClusterTopologyElement{
					MemoryPerNode:     2048,
					NodeConfiguration: "highio.legacy",
					NodeCountPerZone:  1,
					NodeType: &models.ElasticsearchNodeType{
						Data:   ec.Bool(true),
						Master: ec.Bool(true),
					},
				},
			},
			want: TopologySize{
				ID:   "highio.classic",
				Size: 2048,
			},
		},
		{
			name: "DNT topology",
			args: args{
				topology: &models.ElasticsearchClusterTopologyElement{
					InstanceConfigurationID: "gcp.master.classic",
					ZoneCount:               3,
					NodeType: &models.ElasticsearchNodeType{
						Data:   ec.Bool(false),
						Master: ec.Bool(true),
					},
					Size: &models.TopologySize{
						Resource: ec.String("memory"),
						Value:    ec.Int32(1024),
					},
				},
			},
			want: TopologySize{
				ID:   "gcp.master.classic",
				Size: 1024,
			},
		},
		{
			name: "DNT master classic",
			args: args{
				topology: &models.ElasticsearchClusterTopologyElement{
					InstanceConfigurationID: "aws.master.classic",
					ZoneCount:               3,
					NodeType: &models.ElasticsearchNodeType{
						Data:   ec.Bool(false),
						Master: ec.Bool(true),
					},
					Size: &models.TopologySize{
						Resource: ec.String("memory"),
						Value:    ec.Int32(0),
					},
				},
			},
			err: errors.New("couldn't obtain size from aws.master.classic"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSizeFromTopology(tt.args.topology)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("GetSizeFromTopology() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSizeFromTopology() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewOverrideRequest(t *testing.T) {
	type args struct {
		instances   []*models.ClusterInstanceInfo
		size        TopologySize
		boostFactor int32
		filter      []string
	}
	tests := []struct {
		name string
		args args
		want OverrideRequest
	}{
		{
			name: "Legacy request obtains only the matching instances",
			args: args{
				instances: []*models.ClusterInstanceInfo{
					{
						InstanceName: ec.String("instance-0000000008"),
						InstanceConfiguration: &models.ClusterInstanceConfigurationInfo{
							ID:       "aws.highio.legacy",
							Name:     "highio.legacy",
							Resource: "memory",
						},
					},
					{
						InstanceName: ec.String("instance-0000000009"),
						InstanceConfiguration: &models.ClusterInstanceConfigurationInfo{
							ID:       "aws.highio.legacy",
							Name:     "highio.legacy",
							Resource: "memory",
						},
					},
					{
						InstanceName: ec.String("instance-0000000010"),
						InstanceConfiguration: &models.ClusterInstanceConfigurationInfo{
							ID:       "aws.master.legacy",
							Name:     "master.legacy",
							Resource: "memory",
						},
					},
				},
				size: TopologySize{
					ID:   "highio.legacy",
					Size: 2048,
				},
				boostFactor: 2,
				filter: []string{
					"instance-0000000008",
					"instance-0000000009",
				},
			},
			want: OverrideRequest{
				Capacity: 2,
				Instances: []string{
					"instance-0000000008",
					"instance-0000000009",
				},
			},
		},
		{
			name: "DNT request obtains only the matching instances",
			args: args{
				instances: []*models.ClusterInstanceInfo{
					{
						InstanceName: ec.String("instance-0000000008"),
						InstanceConfiguration: &models.ClusterInstanceConfigurationInfo{
							ID:       "aws.highio.classic",
							Name:     "highio.classic",
							Resource: "memory",
						},
					},
					{
						InstanceName: ec.String("instance-0000000009"),
						InstanceConfiguration: &models.ClusterInstanceConfigurationInfo{
							ID:       "aws.highio.classic",
							Name:     "highio.classic",
							Resource: "memory",
						},
					},
					{
						InstanceName: ec.String("instance-0000000010"),
						InstanceConfiguration: &models.ClusterInstanceConfigurationInfo{
							ID:       "aws.master.classic",
							Name:     "master.classic",
							Resource: "memory",
						},
					},
				},
				size: TopologySize{
					ID:   "highio.classic",
					Size: 2048,
				},
				boostFactor: 2,
				filter: []string{
					"instance-0000000008",
					"instance-0000000009",
				},
			},
			want: OverrideRequest{
				Capacity: 2,
				Instances: []string{
					"instance-0000000008",
					"instance-0000000009",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOverrideRequest(tt.args.instances, tt.args.size, tt.args.boostFactor, tt.args.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOverrideRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOverrideCapacity(t *testing.T) {
	type args struct {
		params OverrideCapacityParams
	}
	tests := []struct {
		name string
		args args
		want []OverrideResponse
		err  error
	}{
		{
			name: "Fails due to parameter validation",
			args: args{
				params: OverrideCapacityParams{
					Instances:   nil,
					BoostFactor: 0,
					Value:       1023,
				},
			},
			err: &multierror.Error{
				Errors: []error{
					errors.New("you must specify an absolute value greater than 1024 or a boost factor greater than 1 or reset to the original value"),
					errors.New("cluster id should have a length of 32 characters"),
				},
			},
		},
		{
			name: "Fails when both absolute value AND boost factor are specified",
			args: args{
				params: OverrideCapacityParams{
					Instances:   nil,
					BoostFactor: 2,
					Value:       4096,
				},
			},
			err: &multierror.Error{
				Errors: []error{
					errors.New("you cannot specify at the same time an absolute value and a boost factor"),
					errors.New("cluster id should have a length of 32 characters"),
				},
			},
		},
		{
			name: "Fails when both absolute value or boost factor are specified and a reset is specified",
			args: args{
				params: OverrideCapacityParams{
					Instances:   nil,
					BoostFactor: 2,
					Value:       4096,
					Reset:       true,
				},
			},
			err: &multierror.Error{
				Errors: []error{
					errors.New("you cannot specify at the same time an absolute value or a boost factor and set the reset flag to true"),
					errors.New("cluster id should have a length of 32 characters"),
				},
			},
		},
		{
			name: "Success in overriding a legacy cluster",
			args: args{
				params: OverrideCapacityParams{
					ClusterParams: util.ClusterParams{
						ClusterID: "b786acd298292c2d521c0e8741761b4d",
						API: api.NewMock(mock.Response{
							Response: http.Response{
								StatusCode: 200,
								Body:       mock.NewStringBody(legacyCluster),
							},
						}, mock.Response{
							Response: http.Response{
								StatusCode: 200,
								Body:       NewOverrideResponse(t, 4096),
							},
						}),
					},
					BoostFactor: 2,
					Instances:   []string{"instance-0000000024", "instance-0000000026"},
				},
			},
			want: []OverrideResponse{{
				Instances: []string{"instance-0000000024", "instance-0000000026"},
				Capacity:  4096,
			}},
		},
		{
			name: "Success in overriding a legacy cluster with filtered instances",
			args: args{
				params: OverrideCapacityParams{
					ClusterParams: util.ClusterParams{
						ClusterID: "b786acd298292c2d521c0e8741761b4d",
						API: api.NewMock(mock.Response{
							Response: http.Response{
								StatusCode: 200,
								Body:       mock.NewStringBody(legacyCluster),
							},
						}, mock.Response{
							Response: http.Response{
								StatusCode: 200,
								Body:       NewOverrideResponse(t, 4096),
							},
						}),
					},
					BoostFactor: 2,
					Instances:   []string{"instance-0000000024"},
				},
			},
			want: []OverrideResponse{{
				Instances: []string{"instance-0000000024"},
				Capacity:  4096,
			}},
		},
		{
			name: "Success in overriding a legacy cluster and single instance override",
			args: args{
				params: OverrideCapacityParams{
					ClusterParams: util.ClusterParams{
						ClusterID: "b786acd298292c2d521c0e8741761b4d",
						API: api.NewMock(mock.Response{
							Response: http.Response{
								StatusCode: 200,
								Body:       mock.NewStringBody(legacyCluster),
							},
						}, mock.Response{
							Response: http.Response{
								StatusCode: 200,
								Body:       NewOverrideResponse(t, 4096),
							},
						}),
					},
					Instances:   []string{"instance-0000000024"},
					BoostFactor: 2,
				},
			},
			want: []OverrideResponse{{
				Instances: []string{"instance-0000000024"},
				Capacity:  4096,
			}},
		},
		{
			name: "Success in overriding a legacy ECE cluster",
			args: args{
				params: OverrideCapacityParams{
					ClusterParams: util.ClusterParams{
						ClusterID: "b786acd298292c2d521c0e8741761b4d",
						API: api.NewMock(mock.Response{
							Response: http.Response{
								StatusCode: 200,
								Body:       mock.NewStringBody(legacyECECluster),
							},
						}, mock.Response{
							Response: http.Response{
								StatusCode: 200,
								Body:       NewOverrideResponse(t, 2048),
							},
						}),
					},
					BoostFactor: 2,
					Instances:   []string{"instance-0000000000"},
				},
			},
			want: []OverrideResponse{{
				Instances: []string{"instance-0000000000"},
				Capacity:  2048,
			}},
		},
		{
			name: "Success in overriding a DNT cluster overrides different capacities",
			args: args{
				params: OverrideCapacityParams{
					ClusterParams: util.ClusterParams{
						ClusterID: "ef97cb1bee75971e19be2522eca6a021",
						API: api.NewMock(mock.Response{
							Response: http.Response{
								StatusCode: 200,
								Body:       mock.NewStringBody(dntCluster),
							},
						}, mock.Response{
							Response: http.Response{
								StatusCode: 200,
								Body:       NewOverrideResponse(t, 32768),
							},
						}, mock.Response{
							Response: http.Response{
								StatusCode: 200,
								Body:       NewOverrideResponse(t, 65536),
							},
						}, mock.Response{
							Response: http.Response{
								StatusCode: 200,
								Body:       NewOverrideResponse(t, 4096),
							},
						}),
					},
					BoostFactor: 2,
					Instances: []string{
						"instance-0000000009", "instance-0000000008",
						"instance-0000000013", "instance-0000000014",
						"instance-0000000012", "instance-0000000010", "instance-0000000011",
					},
				},
			},
			want: []OverrideResponse{
				{
					Instances: []string{"instance-0000000013", "instance-0000000014"},
					Capacity:  32768,
				},
				{
					Instances: []string{"instance-0000000009", "instance-0000000008"},
					Capacity:  65536,
				},
				{
					Instances: []string{"instance-0000000012", "instance-0000000010", "instance-0000000011"},
					Capacity:  4096,
				},
			},
		},
		{
			name: "Success in overriding a DNT cluster overrides, setting only a storage multiplier",
			args: args{
				params: OverrideCapacityParams{
					ClusterParams: util.ClusterParams{
						ClusterID: "ef97cb1bee75971e19be2522eca6a021",
						API: api.NewMock(mock.Response{
							Response: http.Response{
								StatusCode: 200,
								Body:       mock.NewStringBody(dntCluster),
							},
						}, mock.Response{
							Response: http.Response{
								StatusCode: 200,
								Body:       NewOverrideResponse(t, 32768/2),
							},
						}, mock.Response{
							Response: http.Response{
								StatusCode: 200,
								Body:       NewOverrideResponse(t, 65536/2),
							},
						}, mock.Response{
							Response: http.Response{
								StatusCode: 200,
								Body:       NewOverrideResponse(t, 4096/2),
							},
						}),
					},
					StorageMultiplier: 40,
					Instances: []string{
						"instance-0000000009", "instance-0000000008",
						"instance-0000000013", "instance-0000000014",
						"instance-0000000012", "instance-0000000010", "instance-0000000011",
					},
				},
			},
			want: []OverrideResponse{
				{
					Instances:         []string{"instance-0000000013", "instance-0000000014"},
					Capacity:          16384,
					StorageMultiplier: 40,
				},
				{
					Instances:         []string{"instance-0000000009", "instance-0000000008"},
					Capacity:          32768,
					StorageMultiplier: 40,
				},
				{
					Instances:         []string{"instance-0000000012", "instance-0000000010", "instance-0000000011"},
					Capacity:          2048,
					StorageMultiplier: 40,
				},
			},
		},
		{
			name: "Fails overriding any further than MaxInstanceCapacity",
			args: args{
				params: OverrideCapacityParams{
					ClusterParams: util.ClusterParams{
						ClusterID: "ef97cb1bee75971e19be2522eca6a021",
						API: api.NewMock(mock.Response{
							Response: http.Response{
								StatusCode: 200,
								Body:       mock.NewStringBody(dntCluster),
							},
						}, mock.Response{
							Response: http.Response{
								StatusCode: 200,
								Body:       NewOverrideResponse(t, 32768),
							},
						}, mock.Response{
							Response: http.Response{
								StatusCode: 200,
								Body:       NewOverrideResponse(t, 65536),
							},
						}, mock.Response{
							Response: http.Response{
								StatusCode: 200,
								Body:       NewOverrideResponse(t, 4096),
							},
						}),
					},
					BoostFactor: 3,
					Instances: []string{
						"instance-0000000009", "instance-0000000008",
						"instance-0000000013", "instance-0000000014",
						"instance-0000000012", "instance-0000000010", "instance-0000000011",
					},
				},
			},
			err: errors.New("instance-0000000009 instance-0000000008: capacity must not exceed 65536"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := OverrideCapacity(tt.args.params)
			if (err != nil || tt.err != nil) && err.Error() != tt.err.Error() {
				t.Errorf("OverrideCapacity() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OverrideCapacity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newSize(t *testing.T) {
	type args struct {
		params OverrideCapacityParams
		size   int32
	}
	tests := []struct {
		name        string
		args        args
		wantNewSize int32
	}{
		{
			name: "new size should be computed based on boost factor",
			args: args{
				params: OverrideCapacityParams{
					BoostFactor: 2,
				},
				size: 1024,
			},
			wantNewSize: 2048,
		},
		{
			name: "new size should be computed based on absolute value",
			args: args{
				params: OverrideCapacityParams{
					Value: 2048,
				},
				size: 1024,
			},
			wantNewSize: 2048,
		},
		{
			name: "new size should be computed based on the original value",
			args: args{
				params: OverrideCapacityParams{
					Value: 2048,
					Reset: true,
				},
				size: 1024,
			},
			wantNewSize: 1024,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newSize(tt.args.params, tt.args.size); !reflect.DeepEqual(got, tt.wantNewSize) {
				t.Errorf("newSize() = %v, want %v", got, tt.wantNewSize)
			}
		})
	}
}
