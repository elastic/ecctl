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

package instanceconfig

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
)

func TestNewConfig(t *testing.T) {
	type args struct {
		source io.Reader
	}
	tests := []struct {
		name string
		args args
		want *models.InstanceConfiguration
		err  error
	}{
		{
			name: "NewConfig succeeds",
			args: args{
				source: strings.NewReader(newConfigKibanaInstanceConfig),
			},
			want: &models.InstanceConfiguration{
				ID:                "kibana",
				Description:       "Instance configuration to be used for Kibana",
				Name:              ec.String("kibana"),
				InstanceType:      ec.String("kibana"),
				StorageMultiplier: float64(4),
				NodeTypes:         []string{},
				DiscreteSizes: &models.DiscreteSizes{
					DefaultSize: ec.Int32(1024),
					Resource:    ec.String("memory"),
					Sizes: []int32{
						1024,
						2048,
						4096,
						8192,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfig(tt.args.source)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("NewConfig() error = %+v, wantErr %+v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
