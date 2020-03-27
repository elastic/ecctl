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

package apm

import (
	"errors"
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

func TestUpgrade(t *testing.T) {
	type args struct {
		params UpgradeParams
	}
	tests := []struct {
		name string
		args args
		want *models.ClusterUpgradeInfo
		err  error
	}{
		{
			name: "fails due to parameter validation",
			args: args{params: UpgradeParams{}},
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				errors.New(`id "" is invalid`),
			}},
		},
		{
			name: "Succeeds without tracking",
			args: args{params: UpgradeParams{
				ID: "d324608c97154bdba2dff97511d40368",
				API: api.NewMock(mock.Response{Response: http.Response{
					Body: mock.NewStructBody(models.ClusterUpgradeInfo{
						ClusterID:      ec.String("d324608c97154bdba2dff97511d40368"),
						ClusterVersion: ec.String("6.4.3"),
					}),
					StatusCode: 202,
				}}),
			}},
			want: &models.ClusterUpgradeInfo{
				ClusterID:      ec.String("d324608c97154bdba2dff97511d40368"),
				ClusterVersion: ec.String("6.4.3"),
			},
		},
		{
			name: "Succeeds with tracking",
			args: args{params: UpgradeParams{
				TrackChangeParams: util.NewMockTrackChangeParams("d324608c97154bdba2dff97511d40368"),
				ID:                "d324608c97154bdba2dff97511d40368",
				API: api.NewMock(util.AppendTrackResponses(mock.Response{Response: http.Response{
					Body: mock.NewStructBody(models.ClusterUpgradeInfo{
						ClusterID:      ec.String("d324608c97154bdba2dff97511d40368"),
						ClusterVersion: ec.String("6.4.2"),
					}),
					StatusCode: 202,
				}})...),
			}},
			want: &models.ClusterUpgradeInfo{
				ClusterID:      ec.String("d324608c97154bdba2dff97511d40368"),
				ClusterVersion: ec.String("6.4.2"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Upgrade(tt.args.params)
			if !reflect.DeepEqual(tt.err, err) {
				t.Errorf("Upgrade() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Upgrade() = %v, want %v", got, tt.want)
			}
		})
	}
}
