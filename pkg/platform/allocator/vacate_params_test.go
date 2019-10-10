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
	"reflect"
	"testing"
	"time"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/output"
	"github.com/elastic/cloud-sdk-go/pkg/sync/pool"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	multierror "github.com/hashicorp/go-multierror"
)

func TestVacateParamsValidate(t *testing.T) {
	type fields struct {
		Allocators          []string
		PreferredAllocators []string
		ClusterFilter       []string
		KindFilter          string
		PoolTimeout         pool.Timeout
		API                 *api.API
		Output              *output.Device
		TrackFrequency      time.Duration
		AllocatorDown       *bool
		Concurrency         uint16
		MaxPollRetries      uint8
	}
	tests := []struct {
		name   string
		fields fields
		err    error
	}{
		{
			name: "Accepts a correct set of parameters",
			fields: fields{
				API:         new(api.API),
				Allocators:  []string{"an allocator"},
				Concurrency: 1,
				Output:      new(output.Device),
			},
			err: nil,
		},
		{
			name: "Accepts a correct set of parameters with an elasticsearch kind filter",
			fields: fields{
				API:         new(api.API),
				Allocators:  []string{"an allocator"},
				KindFilter:  "elasticsearch",
				Concurrency: 1,
				Output:      new(output.Device),
			},
			err: nil,
		},
		{
			name: "Accepts a correct set of parameters with a kibana kind filter",
			fields: fields{
				API:         new(api.API),
				Allocators:  []string{"an allocator"},
				KindFilter:  "kibana",
				Concurrency: 1,
				Output:      new(output.Device),
			},
			err: nil,
		},
		{
			name: "Accepts a correct set of parameters with an apm kind filter",
			fields: fields{
				API:         new(api.API),
				Allocators:  []string{"an allocator"},
				KindFilter:  "apm",
				Concurrency: 1,
				Output:      new(output.Device),
			},
			err: nil,
		},
		{
			name:   "Empty parameters are not accepted",
			fields: fields{},
			err: &multierror.Error{
				Errors: []error{
					errAPIMustNotBeNil,
					errMustSpecifyAtLeast1Allocator,
					errConcurrencyCannotBeZero,
					errOutputDeviceCannotBeNil,
				},
			},
		},
		{
			name: "Cluster filter is invalid",
			fields: fields{
				API:           new(api.API),
				Allocators:    []string{"an allocator"},
				ClusterFilter: []string{"something"},
				Concurrency:   1,
				Output:        new(output.Device),
			},
			err: &multierror.Error{
				Errors: []error{
					errors.New(`cluster filter: id "something" is invalid, must be 32 characters long`),
				},
			},
		},
		{
			name: "Invalid combination of cluster filter and kind filter",
			fields: fields{
				API:           new(api.API),
				Allocators:    []string{"an allocator"},
				ClusterFilter: []string{"63d765d37613423e97b1040257cf20c8"},
				KindFilter:    "elasticsearch",
				Concurrency:   1,
				Output:        new(output.Device),
			},
			err: &multierror.Error{
				Errors: []error{
					errors.New(`only one of "clusters" or "kind" can be specified`),
				},
			},
		},
		{
			name: "Invalid combination of allocatorDown and multiple allocators",
			fields: fields{
				API:           new(api.API),
				Allocators:    []string{"an allocator", "another allocator"},
				AllocatorDown: ec.Bool(true),
				Concurrency:   1,
				Output:        new(output.Device),
			},
			err: &multierror.Error{
				Errors: []error{
					errors.New(`cannot set the AllocatorDown when multiple allocators are specified`),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := VacateParams{
				API:                 tt.fields.API,
				Allocators:          tt.fields.Allocators,
				PreferredAllocators: tt.fields.PreferredAllocators,
				ClusterFilter:       tt.fields.ClusterFilter,
				KindFilter:          tt.fields.KindFilter,
				Concurrency:         tt.fields.Concurrency,
				Output:              tt.fields.Output,
				MaxPollRetries:      tt.fields.MaxPollRetries,
				TrackFrequency:      tt.fields.TrackFrequency,
				PoolTimeout:         tt.fields.PoolTimeout,
				AllocatorDown:       tt.fields.AllocatorDown,
			}
			if err := params.Validate(); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("VacateParams.Validate() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
