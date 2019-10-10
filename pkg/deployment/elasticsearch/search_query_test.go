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

package elasticsearch

import (
	"errors"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/models"
)

func TestNewQueryByUserID(t *testing.T) {
	type args struct {
		params NewQueryParams
	}
	tests := []struct {
		name string
		args args
		want []*models.QueryContainer
		err  error
	}{
		{
			name: "fails due to param validation",
			args: args{params: NewQueryParams{}},
			err:  errors.New("query user id cannot be empty"),
		},
		{
			name: "returns a formed query",
			args: args{params: NewQueryParams{UserID: "marc"}},
			want: []*models.QueryContainer{{
				Term: map[string]models.TermQuery{
					"settings.metadata.owner_id": {
						Value: "marc",
					},
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewQueryByUserID(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("NewQueryByUserID() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewQueryByUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewStartedClusterQuery(t *testing.T) {
	type args struct {
		includeInitializing bool
	}
	tests := []struct {
		name string
		args args
		want []*models.QueryContainer
	}{
		{
			name: "with includeIntializing",
			args: args{includeInitializing: true},
			want: []*models.QueryContainer{
				{
					Term: map[string]models.TermQuery{
						"status": {
							Value: "started",
						},
					},
				},
				{
					Term: map[string]models.TermQuery{
						"status": {
							Value: "initializing",
						},
					},
				},
			},
		},
		{
			name: "without includeIntializing",
			args: args{includeInitializing: false},
			want: []*models.QueryContainer{{
				Term: map[string]models.TermQuery{
					"status": {
						Value: "started",
					},
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStartedClusterQuery(tt.args.includeInitializing); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStartedClusterQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
