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

package depresource

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/hashicorp/go-multierror"
)

func TestLatestStackVersion(t *testing.T) {
	type args struct {
		params LatestStackVersionParams
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantOut string
		err     error
	}{
		{
			name: "fails due to parameter validation",
			err: &multierror.Error{Errors: []error{
				errors.New("api reference is required for command"),
			}},
		},
		{
			name: "returns the version string when specified in the parameters",
			args: args{params: LatestStackVersionParams{
				API:     api.NewMock(),
				Version: "7.4.2",
			}},
			want: "7.4.2",
		},
		{
			name: "obtains the latest version from the API",
			args: args{params: LatestStackVersionParams{
				API: api.NewMock(mock.New200Response(mock.NewStructBody(models.StackVersionConfigs{
					Stacks: []*models.StackVersionConfig{
						{Version: "6.4.2"},
						{Version: "7.4.2"},
						{Version: "5.4.2"},
					},
				}))),
			}},
			want: "7.4.2",
		},
		{
			name: "obtains the latest version from the API and prints to the writer",
			args: args{params: LatestStackVersionParams{
				Writer: new(bytes.Buffer),
				API: api.NewMock(mock.New200Response(mock.NewStructBody(models.StackVersionConfigs{
					Stacks: []*models.StackVersionConfig{
						{Version: "6.4.2"},
						{Version: "7.4.2"},
						{Version: "5.4.2"},
					},
				}))),
			}},
			want:    "7.4.2",
			wantOut: "Obtained latest stack version: 7.4.2\n",
		},
		{
			name: "obtains an error when the list is empty",
			args: args{params: LatestStackVersionParams{
				API: api.NewMock(mock.New200Response(mock.NewStructBody(models.StackVersionConfigs{
					Stacks: []*models.StackVersionConfig{},
				}))),
			}},
			err: errors.New("version discovery: stack list is seemingly empty, something is terribly wrong"),
		},
		{
			name: "returns an error when the API call fails",
			args: args{params: LatestStackVersionParams{
				API: api.NewMock(mock.New500Response(mock.NewStructBody(errors.New("error")))),
			}},
			err: errors.New("version discovery: failed to obtain stack list, please specify a version"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LatestStackVersion(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("LatestStackVersion() error = %v, wantErr %v", err, tt.err)
				return
			}
			if got != tt.want {
				t.Errorf("LatestStackVersion() = %v, want %v", got, tt.want)
			}

			if buf, ok := tt.args.params.Writer.(*bytes.Buffer); ok {
				if out := buf.String(); out != tt.wantOut {
					t.Errorf("LatestStackVersion() buf = %v, want %v", out, tt.wantOut)
				}
			}
		})
	}
}
