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

package constructor

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
)

func TestList(t *testing.T) {
	var constructorListSuccess = `
	{
		"constructors": [
		  {
			"constructor_id": "192.168.44.10",
			"status": {
			  "connected": true,
			  "maintenance_mode": false
			}
		  }
		]
	}`[1:]
	type args struct {
		params Params
	}
	tests := []struct {
		name    string
		args    args
		want    *models.ConstructorOverview
		wantErr bool
	}{
		{
			name: "Constructor list succeeds",
			args: args{params: Params{
				API: api.NewMock(mock.Response{Response: http.Response{
					Body:       mock.NewStringBody(constructorListSuccess),
					StatusCode: 200,
				}}),
			}},
			want: &models.ConstructorOverview{
				Constructors: []*models.ConstructorInfo{
					{
						ConstructorID: ec.String("192.168.44.10"),
						Status: &models.ConstructorHealthStatus{
							Connected:       ec.Bool(true),
							MaintenanceMode: ec.Bool(false),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Constructor list fails",
			args: args{params: Params{
				API: api.NewMock(mock.Response{Error: errors.New("error")}),
			}},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Constructor list with an empty API",
			args: args{params: Params{
				API: nil,
			}},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := List(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGet(t *testing.T) {
	var constructoGet = `
	{
		"constructor_id": "192.168.44.10",
		"status": {
		  "connected": true,
		  "maintenance_mode": false
		}
	  }`[1:]

	type args struct {
		params GetParams
	}
	tests := []struct {
		name    string
		args    args
		want    *models.ConstructorInfo
		wantErr bool
	}{
		{
			name: "Get constructor succeeds",
			args: args{
				params: GetParams{
					ID: "192.168.44.10",
					Params: Params{
						API: api.NewMock(mock.Response{Response: http.Response{
							Body:       mock.NewStringBody(constructoGet),
							StatusCode: 200,
						}}),
					},
				},
			},
			want: &models.ConstructorInfo{
				ConstructorID: ec.String("192.168.44.10"),
				Status: &models.ConstructorHealthStatus{
					Connected:       ec.Bool(true),
					MaintenanceMode: ec.Bool(false),
				},
			},
		},
		{
			name: "Get constructor fails",
			args: args{
				params: GetParams{
					ID: "192.168.44.10",
					Params: Params{
						API: api.NewMock(mock.Response{Error: errors.New("error")}),
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Get constructor fails due to empty params.ID",
			args: args{
				params: GetParams{
					ID:     "",
					Params: Params{new(api.API)},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnableMaintenace(t *testing.T) {
	type args struct {
		params MaintenanceParams
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Set constructor to maintenance mode succeeds",
			args: args{
				params: MaintenanceParams{
					ID: "192.168.44.10",
					Params: Params{
						API: api.NewMock(mock.Response{Response: http.Response{
							Body:       mock.NewStringBody(`{}`),
							StatusCode: 202,
						}}),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Set constructor to maintenance mode fails",
			args: args{
				params: MaintenanceParams{
					ID: "192.168.44.10",
					Params: Params{
						API: api.NewMock(mock.Response{Error: errors.New("error")}),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Set constructor to maintenance mode fails due to param validation",
			args: args{
				params: MaintenanceParams{
					Params: Params{
						API: api.NewMock(mock.Response{Error: errors.New("error")}),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := EnableMaintenace(tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("EnableMaintenace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDisableMaintenance(t *testing.T) {
	type args struct {
		params MaintenanceParams
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Set constructor maintenance mode to false succeeds",
			args: args{
				params: MaintenanceParams{
					ID: "192.168.44.10",
					Params: Params{
						API: api.NewMock(mock.Response{Response: http.Response{
							Body:       mock.NewStringBody(`{}`),
							StatusCode: 202,
						}}),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Set constructor maintenance mode to false fails",
			args: args{
				params: MaintenanceParams{
					ID: "192.168.44.10",
					Params: Params{
						API: api.NewMock(mock.Response{Error: errors.New("error")}),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Set constructor maintenance mode to false fails due to param validation",
			args: args{
				params: MaintenanceParams{
					Params: Params{
						API: api.NewMock(mock.Response{Response: http.Response{
							Body:       mock.NewStringBody(`{}`),
							StatusCode: 202,
						}}),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DisableMaintenance(tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("DisableMaintenance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
