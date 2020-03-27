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
	"fmt"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
)

func TestCheckErrType(t *testing.T) {
	type args struct {
		got  error
		want error
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "different error types return an error",
			args: args{
				got:  errors.New("ERROR"),
				want: errors.New("WANT THIS ERROR"),
			},
			wantErr: fmt.Errorf(formatErrType, errors.New("ERROR"), errors.New("WANT THIS ERROR")),
		},
		{
			name: "equal error types return nil",
			args: args{
				got:  errors.New("ERROR"),
				want: errors.New("ERROR"),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckErrType(tt.args.got, tt.args.want); (err != nil) != (tt.wantErr != nil) {
				t.Errorf("CheckErrType() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAppendTrackResponses(t *testing.T) {
	var success = NewSuccessfulCurrentPlan()
	type args struct {
		res []mock.Response
	}
	tests := []struct {
		name string
		args args
		want []mock.Response
	}{
		{
			name: "appends a response",
			args: args{res: []mock.Response{
				{Error: errors.New("someerror")},
			}},
			want: []mock.Response{
				{Error: errors.New("someerror")},
				success, success, success,
			},
		},
		{
			name: "appends multiple responses",
			args: args{res: []mock.Response{
				{Error: errors.New("someerror")},
				{Error: errors.New("someothererror")},
			}},
			want: []mock.Response{
				{Error: errors.New("someerror")},
				{Error: errors.New("someothererror")},
				success, success, success,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AppendTrackResponses(tt.args.res...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppendTrackResponses() = %v, want %v", got, tt.want)
			}
		})
	}
}
