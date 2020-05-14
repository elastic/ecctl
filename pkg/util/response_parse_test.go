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
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_elasticsearch"
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_kibana"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/multierror"
	"github.com/elastic/cloud-sdk-go/pkg/output"
	"github.com/elastic/cloud-sdk-go/pkg/plan"
	"github.com/elastic/cloud-sdk-go/pkg/plan/planutil"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
)

func TestParseCUResponseParams_Validate(t *testing.T) {
	type fields struct {
		CreateResponse interface{}
		UpdateResponse interface{}
		Err            error
		Track          bool
		planutil.TrackChangeParams
	}
	tests := []struct {
		name   string
		fields fields
		err    error
	}{
		{
			name:   "Passing an error as a parameter returns it",
			fields: fields{Err: errors.New("an error which I just passed")},
			err:    errors.New("an error which I just passed"),
		},
		{
			name:   "Empty fields trigger error",
			fields: fields{},
			err: multierror.NewPrefixed("plan tracking",
				errors.New("parse response: One of Create or Update response must be populated"),
			),
		},
		{
			name: "Correct parameter validation with CreateResponse returns no error",
			fields: fields{
				CreateResponse: "something",
				Track:          true,
				TrackChangeParams: planutil.TrackChangeParams{
					TrackChangeParams: plan.TrackChangeParams{
						API: new(api.API),
					},
					Writer: new(output.Device),
				},
			},
			err: nil,
		},
		{
			name: "Correct parameter validation with UpdateResponse returns no error",
			fields: fields{
				UpdateResponse: "something",
				Track:          true,
				TrackChangeParams: planutil.TrackChangeParams{
					TrackChangeParams: plan.TrackChangeParams{
						API: new(api.API),
					},
					Writer: new(output.Device),
				},
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := ParseCUResponseParams{
				CreateResponse:    tt.fields.CreateResponse,
				UpdateResponse:    tt.fields.UpdateResponse,
				Err:               tt.fields.Err,
				Track:             tt.fields.Track,
				TrackChangeParams: tt.fields.TrackChangeParams,
			}
			if err := params.Validate(); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("ParseCUResponseParams.Validate() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}

func TestParseCUResponse(t *testing.T) {
	var foundDeploymentIDResponse = models.DeploymentsSearchResponse{
		Deployments: []*models.DeploymentSearchResponse{
			{ID: ec.String("cbb4bc6c09684c86aa5de54c05ea1d38")},
		},
	}
	type args struct {
		params ParseCUResponseParams
	}
	tests := []struct {
		name string
		args args
		want *models.ClusterCrudResponse
		err  error
	}{
		{
			name: "Passing an error as a parameter returns it",
			args: args{params: ParseCUResponseParams{
				Err: errors.New("an error that gets returned"),
			}},
			err: errors.New("an error that gets returned"),
		},
		{
			name: "Passing an Elasticsearch create response succeeds",
			args: args{params: ParseCUResponseParams{
				CreateResponse: &clusters_elasticsearch.CreateEsClusterOK{
					Payload: &models.ClusterCrudResponse{ElasticsearchClusterID: "someID"},
				},
				Track: false,
			}},
			want: &models.ClusterCrudResponse{ElasticsearchClusterID: "someID"},
		},
		{
			name: "Passing an Elasticsearch update response succeeds",
			args: args{params: ParseCUResponseParams{
				UpdateResponse: &clusters_elasticsearch.UpdateEsClusterPlanAccepted{
					Payload: &models.ClusterCrudResponse{ElasticsearchClusterID: "someID"},
				},
				Track: false,
			}},
			want: &models.ClusterCrudResponse{ElasticsearchClusterID: "someID"},
		},
		{
			name: "Passing an Kibana create response succeeds",
			args: args{params: ParseCUResponseParams{
				CreateResponse: &clusters_kibana.CreateKibanaClusterOK{
					Payload: &models.ClusterCrudResponse{KibanaClusterID: "someID"},
				},
				Track: false,
			}},
			want: &models.ClusterCrudResponse{KibanaClusterID: "someID"},
		},
		{
			name: "Passing an Kibana update response succeeds",
			args: args{params: ParseCUResponseParams{
				UpdateResponse: &clusters_kibana.UpdateKibanaClusterPlanAccepted{
					Payload: &models.ClusterCrudResponse{KibanaClusterID: "someID"},
				},
				Track: false,
			}},
			want: &models.ClusterCrudResponse{KibanaClusterID: "someID"},
		},
		// Tracker tests
		{
			name: "Passing a Kibana update response succeeds and tracks the cluster",
			args: args{params: ParseCUResponseParams{
				Track: true,
				TrackChangeParams: planutil.TrackChangeParams{
					TrackChangeParams: plan.TrackChangeParams{
						API: api.NewMock(AppendTrackResponses(
							mock.New200StructResponse(foundDeploymentIDResponse),
						)...),
						Config: plan.TrackFrequencyConfig{
							PollFrequency: time.Nanosecond,
							MaxRetries:    1,
						},
					},
					Writer: new(output.Device),
				},
				UpdateResponse: &clusters_kibana.UpdateKibanaClusterPlanAccepted{
					Payload: &models.ClusterCrudResponse{KibanaClusterID: "3ee11eb40eda22cac0cce259625c6734"},
				},
			}},
			want: &models.ClusterCrudResponse{KibanaClusterID: "3ee11eb40eda22cac0cce259625c6734"},
		},
		{
			name: "Passing an Elasticsearch update response succeeds and tracks the cluster",
			args: args{params: ParseCUResponseParams{
				TrackChangeParams: planutil.TrackChangeParams{
					TrackChangeParams: plan.TrackChangeParams{
						API: api.NewMock(AppendTrackResponses(
							mock.New200StructResponse(foundDeploymentIDResponse),
						)...),
						Config: plan.TrackFrequencyConfig{
							PollFrequency: time.Nanosecond,
							MaxRetries:    1,
						},
					},
					Writer: new(output.Device),
				},
				CreateResponse: &clusters_elasticsearch.CreateEsClusterOK{
					Payload: &models.ClusterCrudResponse{ElasticsearchClusterID: "3ee11eb40eda22cac0cce259625c6734"},
				},
			}},
			want: &models.ClusterCrudResponse{ElasticsearchClusterID: "3ee11eb40eda22cac0cce259625c6734"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCUResponse(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("ParseCUResponse() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCUResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
