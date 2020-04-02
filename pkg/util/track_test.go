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

package util

import (
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/plan"
	"github.com/elastic/cloud-sdk-go/pkg/plan/planutil"
)

func TestSetClusterTracking(t *testing.T) {
	type args struct {
		params planutil.TrackChangeParams
		id     string
		kind   string
	}
	tests := []struct {
		name string
		args args
		want planutil.TrackChangeParams
	}{
		{
			args: args{id: "123", kind: Elasticsearch},
			want: planutil.TrackChangeParams{TrackChangeParams: plan.TrackChangeParams{
				ResourceID:       "123",
				Kind:             Elasticsearch,
				IgnoreDownstream: true,
			}},
		},
		{
			args: args{id: "123", kind: Kibana},
			want: planutil.TrackChangeParams{TrackChangeParams: plan.TrackChangeParams{
				ResourceID:       "123",
				Kind:             Kibana,
				IgnoreDownstream: true,
			}},
		},
		{
			args: args{id: "123", kind: Apm},
			want: planutil.TrackChangeParams{TrackChangeParams: plan.TrackChangeParams{
				ResourceID:       "123",
				Kind:             Apm,
				IgnoreDownstream: true,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetClusterTracking(tt.args.params, tt.args.id, tt.args.kind); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetClusterTracking() = %v, want %v", got, tt.want)
			}
		})
	}
}
