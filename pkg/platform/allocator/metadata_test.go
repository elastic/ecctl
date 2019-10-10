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

package allocator

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"

	"github.com/elastic/ecctl/pkg/util"
)

func TestGetAllocatorMetadata(t *testing.T) {
	var getAllocatorMetadataSuccess = `[
	{
	  "key": "version",
	  "value": "2017-09-30"
	},
	{
	  "key": "instanceId",
	  "value": "i-09a0e797fb3af6864"
	},
	{
	  "key": "architecture",
	  "value": "x86_64"
	},
	{
	  "key": "instanceType",
	  "value": "i3.8xlarge"
	},
	{
	  "key": "availabilityZone",
	  "value": "us-east-1a"
	},
	{
	  "key": "pendingTime",
	  "value": "2018-05-18T13:24:21Z"
	},
	{
	  "key": "imageId",
	  "value": "ami-ba0a51c0"
	},
	{
	  "key": "privateIp",
	  "value": "172.25.61.100"
	},
	{
	  "key": "region",
	  "value": "us-east-1"
	}
  ]`

	type args struct {
		params MetadataGetParams
	}
	tests := []struct {
		name string
		args args
		want []*models.MetadataItem
		err  error
	}{
		{
			name: "Get metadata fails due to parameter validation failure (missing API)",
			args: args{
				params: MetadataGetParams{
					ID: "i-09a0e797fb3af6864",
				},
			},
			err: util.ErrAPIReq,
		},
		{
			name: "Get metadata fails due to parameter validation failure (missing ID)",
			args: args{
				params: MetadataGetParams{
					API: new(api.API),
				},
			},
			err: errors.New("allocator metadata: id cannot be empty"),
		},
		{
			name: "Get metadata fails due to API failure",
			args: args{
				params: MetadataGetParams{
					ID: "an id",
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(""),
						StatusCode: 500,
					}}),
				},
			},
			err: errors.New("unknown error (status 500)"),
		},
		{
			name: "Get metadata Succeeds",
			args: args{
				params: MetadataGetParams{
					ID: "i-09a0e797fb3af6864",
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(getAllocatorMetadataSuccess),
						StatusCode: 200,
					}}),
				},
			},
			want: []*models.MetadataItem{
				{
					Key:   ec.String("version"),
					Value: ec.String("2017-09-30"),
				},
				{
					Key:   ec.String("instanceId"),
					Value: ec.String("i-09a0e797fb3af6864"),
				},
				{
					Key:   ec.String("architecture"),
					Value: ec.String("x86_64"),
				},
				{
					Key:   ec.String("instanceType"),
					Value: ec.String("i3.8xlarge"),
				},
				{
					Key:   ec.String("availabilityZone"),
					Value: ec.String("us-east-1a"),
				},
				{
					Key:   ec.String("pendingTime"),
					Value: ec.String("2018-05-18T13:24:21Z"),
				},
				{
					Key:   ec.String("imageId"),
					Value: ec.String("ami-ba0a51c0"),
				},
				{
					Key:   ec.String("privateIp"),
					Value: ec.String("172.25.61.100"),
				},
				{
					Key:   ec.String("region"),
					Value: ec.String("us-east-1"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAllocatorMetadata(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("GetAllocatorMetadata() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllocatorMetadata() = \n%+v, want \n%+v", got, tt.want)
			}
		})
	}
}

func TestSetAllocatorMetadataItem(t *testing.T) {
	type args struct {
		params MetadataSetParams
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Set allocator metadata fails due to parameter validation (Missing API)",
			args: args{
				params: MetadataSetParams{
					ID: "an ID",
				},
			},
			err: util.ErrAPIReq,
		},
		{
			name: "Set allocator metadata fails due to parameter validation (Missing ID)",
			args: args{
				params: MetadataSetParams{
					API: new(api.API),
				},
			},
			err: errors.New("allocator metadata: id cannot be empty"),
		},
		{
			name: "Set allocator metadata fails due to parameter validation (Missing key)",
			args: args{
				params: MetadataSetParams{
					API: new(api.API),
					ID:  "an ID",
				},
			},
			err: errors.New("allocator metadata: key cannot be empty"),
		},
		{
			name: "Set allocator metadata fails due to parameter validation (Missing key value)",
			args: args{
				params: MetadataSetParams{
					API: new(api.API),
					ID:  "an ID",
					Key: "foo",
				},
			},
			err: errors.New("allocator metadata: key value cannot be empty"),
		},
		{
			name: "Set allocator metadata succeeds",
			args: args{
				params: MetadataSetParams{
					ID: "an ID",
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(""),
						StatusCode: 200,
					}}),
					Key:   "foo",
					Value: "bar",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetAllocatorMetadataItem(tt.args.params); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("SetAllocatorMetadataItem() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}

func TestDeleteAllocatorMetadataItem(t *testing.T) {
	type args struct {
		params MetadataDeleteParams
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Delete allocator metadata fails due to parameter validation (Missing API)",
			args: args{
				params: MetadataDeleteParams{
					ID: "an ID",
				},
			},
			err: util.ErrAPIReq,
		},
		{
			name: "Set allocator metadata fails due to parameter validation (Missing ID)",
			args: args{
				params: MetadataDeleteParams{
					API: new(api.API),
				},
			},
			err: errors.New("allocator metadata: id cannot be empty"),
		},
		{
			name: "Set allocator metadata fails due to parameter validation (Missing key)",
			args: args{
				params: MetadataDeleteParams{
					API: new(api.API),
					ID:  "an ID",
				},
			},
			err: errors.New("allocator metadata: key cannot be empty"),
		},
		{
			name: "Set allocator metadata succeeds",
			args: args{
				params: MetadataDeleteParams{
					ID: "an ID",
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(""),
						StatusCode: 200,
					}}),
					Key: "foo",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteAllocatorMetadataItem(tt.args.params); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("SetAllocatorMetadataItem() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
