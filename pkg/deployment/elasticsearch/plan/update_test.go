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

package plan

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func TestUpdate(t *testing.T) {
	type args struct {
		params UpdateParams
	}
	tests := []struct {
		name string
		args args
		want *models.ClusterCrudResponse
		err  error
	}{
		{
			name: "fails due to parameter validation",
			args: args{params: UpdateParams{}},
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
				errors.New(`id "" is invalid`),
			}},
		},
		{
			name: "Successfully validates a plan without tracking",
			args: args{params: UpdateParams{
				ID:           "d324608c97154bdba2dff97511d40368",
				ValidateOnly: true,
				API: api.NewMock(mock.Response{Response: http.Response{
					Body: mock.NewStructBody(models.ClusterCrudResponse{
						ElasticsearchClusterID: "d324608c97154bdba2dff97511d40368",
					}),
					StatusCode: 200,
				}}),
			}},
			want: &models.ClusterCrudResponse{
				ElasticsearchClusterID: "d324608c97154bdba2dff97511d40368",
			},
		},
		{
			name: "Successfully validates a plan, ignores requested tracking",
			args: args{params: UpdateParams{
				ID:                "d324608c97154bdba2dff97511d40368",
				ValidateOnly:      true,
				Track:             true,
				TrackChangeParams: util.NewMockTrackChangeParams(""),
				API: api.NewMock(mock.Response{Response: http.Response{
					Body: mock.NewStructBody(models.ClusterCrudResponse{
						ElasticsearchClusterID: "d324608c97154bdba2dff97511d40368",
					}),
					StatusCode: 200,
				}}),
			}},
			want: &models.ClusterCrudResponse{
				ElasticsearchClusterID: "d324608c97154bdba2dff97511d40368",
			},
		},
		{
			name: "Successfully submits the new plan without tracking",
			args: args{params: UpdateParams{
				ID: "d324608c97154bdba2dff97511d40368",
				API: api.NewMock(mock.Response{Response: http.Response{
					Body: mock.NewStructBody(models.ClusterCrudResponse{
						ElasticsearchClusterID: "d324608c97154bdba2dff97511d40368",
					}),
					StatusCode: 202,
				}}),
			}},
			want: &models.ClusterCrudResponse{
				ElasticsearchClusterID: "d324608c97154bdba2dff97511d40368",
			},
		},
		{
			name: "Successfully submits the new plan with tracking",
			args: args{params: UpdateParams{
				ID:                "d324608c97154bdba2dff97511d40368",
				Track:             true,
				TrackChangeParams: util.NewMockTrackChangeParams(""),
				API: api.NewMock(util.AppendTrackResponses(mock.Response{Response: http.Response{
					Body: mock.NewStructBody(models.ClusterCrudResponse{
						ElasticsearchClusterID: "d324608c97154bdba2dff97511d40368",
					}),
					StatusCode: 202,
				}})...),
			}},
			want: &models.ClusterCrudResponse{
				ElasticsearchClusterID: "d324608c97154bdba2dff97511d40368",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Update(tt.args.params)
			if !reflect.DeepEqual(tt.err, err) {
				t.Errorf("UpdatePlan() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdatePlan() = %v, want %v", got, tt.want)
			}
		})
	}
}
