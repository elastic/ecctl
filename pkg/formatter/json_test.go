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

package formatter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
)

type testMarshal struct {
	A string `json:"a"`
	B string `json:"b"`
	C string `json:"c"`
}

func TestJSON_format(t *testing.T) {
	type fields struct {
		o io.Writer
	}
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    string
	}{
		{
			"test Marshal",
			fields{
				&bytes.Buffer{},
			},
			args{
				&testMarshal{
					"DataA",
					"DataB",
					"DataC",
				},
			},
			false,
			`{
  "a": "DataA",
  "b": "DataB",
  "c": "DataC"
}
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &JSON{
				o: tt.fields.o,
			}
			if err := f.format(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("JSON.format() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				got := tt.fields.o.(*bytes.Buffer).String()
				if got != tt.want {
					t.Errorf("JSON.format() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestJSON_Format(t *testing.T) {
	type fields struct {
		o io.Writer
	}
	type args struct {
		path string
		data interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Format action with no format",
			fields{
				&bytes.Buffer{},
			},
			args{
				path: "allocator/list",
				data: &models.AllocatorOverview{
					Zones: []*models.AllocatorZoneInfo{
						{
							Allocators: []*models.AllocatorInfo{
								{
									AllocatorID: ec.String("MyAllocatorID"),
									Capacity: &models.AllocatorCapacity{
										Memory: &models.AllocatorCapacityMemory{
											Total: ec.Int32(4096),
											Used:  ec.Int32(2048),
										},
									},
									Status: &models.AllocatorHealthStatus{
										Connected:       ec.Bool(true),
										MaintenanceMode: ec.Bool(false),
									},
									ZoneID: ec.String("zone1"),
									HostIP: ec.String("1.1.1.1"),
								},
							},
							ZoneID: ec.String("zone1"),
						},
						{
							Allocators: []*models.AllocatorInfo{
								{
									AllocatorID: ec.String("MyAllocatorID"),
									Capacity: &models.AllocatorCapacity{
										Memory: &models.AllocatorCapacityMemory{
											Total: ec.Int32(4096),
											Used:  ec.Int32(2048),
										},
									},
									Status: &models.AllocatorHealthStatus{
										Connected:       ec.Bool(true),
										MaintenanceMode: ec.Bool(false),
									},
									ZoneID: ec.String("zone2"),
									HostIP: ec.String("1.1.1.2"),
									Instances: []*models.AllocatedInstanceStatus{
										{},
									},
								},
							},
							ZoneID: ec.String("zone2"),
						},
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &JSON{
				o: tt.fields.o,
			}
			if err := f.Format(tt.args.path, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("JSON.Format() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				got := tt.fields.o.(*bytes.Buffer).String()
				res, _ := json.MarshalIndent(tt.args.data, "", "  ")
				if got != fmt.Sprintln(string(res)) {
					t.Errorf("JSON.Format() got = %v, want %v", got, string(res))
				}
			}
		})
	}
}
