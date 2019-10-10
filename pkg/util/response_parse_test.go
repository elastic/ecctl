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
	"bytes"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_elasticsearch"
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_kibana"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/output"
	multierror "github.com/hashicorp/go-multierror"
)

func TestParseCUResponseParams_Validate(t *testing.T) {
	type fields struct {
		API            *api.API
		CreateResponse interface{}
		UpdateResponse interface{}
		Err            error
		TrackParams    TrackParams
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
			err: &multierror.Error{Errors: []error{
				errors.New("parse response: API cannot be empty"),
				errors.New("parse response: One of Create or Update response must be populated"),
			}},
		},
		{
			name: "Setting track to true requres an output device",
			fields: fields{
				API:            new(api.API),
				CreateResponse: "something",
				TrackParams:    TrackParams{Track: true},
			},
			err: &multierror.Error{Errors: []error{
				errors.New("track params: output device cannot be empty"),
			}},
		},
		{
			name: "Correct parameter validation with CreateResponse returns no error",
			fields: fields{
				API:            new(api.API),
				CreateResponse: "something",
				TrackParams: TrackParams{
					Track:  true,
					Output: new(output.Device),
				},
			},
			err: nil,
		},
		{
			name: "Correct parameter validation with UpdateResponse returns no error",
			fields: fields{
				API:            new(api.API),
				UpdateResponse: "something",
				TrackParams: TrackParams{
					Track:  true,
					Output: new(output.Device),
				},
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := ParseCUResponseParams{
				API:            tt.fields.API,
				CreateResponse: tt.fields.CreateResponse,
				UpdateResponse: tt.fields.UpdateResponse,
				Err:            tt.fields.Err,
				TrackParams:    tt.fields.TrackParams,
			}
			if err := params.Validate(); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("ParseCUResponseParams.Validate() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}

func TestParseCUResponse(t *testing.T) {
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
				API: &api.API{},
				CreateResponse: &clusters_elasticsearch.CreateEsClusterOK{
					Payload: &models.ClusterCrudResponse{ElasticsearchClusterID: "someID"},
				},
				TrackParams: TrackParams{Track: false},
			}},
			want: &models.ClusterCrudResponse{ElasticsearchClusterID: "someID"},
		},
		{
			name: "Passing an Elasticsearch update response succeeds",
			args: args{params: ParseCUResponseParams{
				API: &api.API{},
				UpdateResponse: &clusters_elasticsearch.UpdateEsClusterPlanAccepted{
					Payload: &models.ClusterCrudResponse{ElasticsearchClusterID: "someID"},
				},
				TrackParams: TrackParams{Track: false},
			}},
			want: &models.ClusterCrudResponse{ElasticsearchClusterID: "someID"},
		},
		{
			name: "Passing an Kibana create response succeeds",
			args: args{params: ParseCUResponseParams{
				API: &api.API{},
				CreateResponse: &clusters_kibana.CreateKibanaClusterOK{
					Payload: &models.ClusterCrudResponse{KibanaClusterID: "someID"},
				},
				TrackParams: TrackParams{Track: false},
			}},
			want: &models.ClusterCrudResponse{KibanaClusterID: "someID"},
		},
		{
			name: "Passing an Kibana update response succeeds",
			args: args{params: ParseCUResponseParams{
				API: &api.API{},
				UpdateResponse: &clusters_kibana.UpdateKibanaClusterPlanAccepted{
					Payload: &models.ClusterCrudResponse{KibanaClusterID: "someID"},
				},
				TrackParams: TrackParams{Track: false},
			}},
			want: &models.ClusterCrudResponse{KibanaClusterID: "someID"},
		},
		// Tracker tests
		{
			name: "Passing a Kibana update response succeeds and tracks the cluster",
			args: args{params: ParseCUResponseParams{
				TrackParams: TrackParams{
					PollFrequency: time.Nanosecond,
					MaxRetries:    1,
					Track:         true,
					Output:        output.NewDevice(new(bytes.Buffer)),
				},
				API: api.NewMock(
					PlanNotFound,
					PlanNotFound,
					NewFailedPlanUnknown(),
				),
				UpdateResponse: &clusters_kibana.UpdateKibanaClusterPlanAccepted{
					Payload: &models.ClusterCrudResponse{KibanaClusterID: "3ee11eb40eda22cac0cce259625c6734"},
				},
			}},
			want: &models.ClusterCrudResponse{KibanaClusterID: "3ee11eb40eda22cac0cce259625c6734"},
			err:  errors.New("cluster [3ee11eb40eda22cac0cce259625c6734][kibana] plan failed due to unknown error"),
		},
		{
			name: "Passing an Elasticsearch update response succeeds and tracks the cluster",
			args: args{params: ParseCUResponseParams{
				TrackParams: TrackParams{
					PollFrequency: time.Nanosecond,
					MaxRetries:    1,
					Track:         true,
					Output:        output.NewDevice(new(bytes.Buffer)),
				},
				API: api.NewMock(
					PlanNotFound,
					PlanNotFound,
					NewFailedPlanUnknown(),
				),
				CreateResponse: &clusters_elasticsearch.CreateEsClusterOK{
					Payload: &models.ClusterCrudResponse{ElasticsearchClusterID: "3ee11eb40eda22cac0cce259625c6734"},
				},
			}},
			want: &models.ClusterCrudResponse{ElasticsearchClusterID: "3ee11eb40eda22cac0cce259625c6734"},
			err:  errors.New("cluster [3ee11eb40eda22cac0cce259625c6734][elasticsearch] plan failed due to unknown error"),
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
