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

package cmdutil

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/multierror"
	"github.com/elastic/cloud-sdk-go/pkg/output"
	"github.com/elastic/cloud-sdk-go/pkg/plan/planutil"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"

	"github.com/elastic/ecctl/pkg/formatter"
)

func TestTrack(t *testing.T) {
	var outputBuf = new(bytes.Buffer)
	var secondOutputBuf = new(bytes.Buffer)
	var thirdOutputBuf = new(bytes.Buffer)
	var wantOutTokenValue = `{
  "token": "tokenvalue"
}
`
	type args struct {
		params TrackParams
		buf    *bytes.Buffer
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		err     error
	}{
		{
			name: "returns an error when Formatter is not specified",
			args: args{params: TrackParams{}},
			err: multierror.NewPrefixed("plan tracker",
				errors.New("formatter cannot be nil"),
			),
		},
		{
			name: "returns a formatted structure when track is set to false",
			args: args{params: TrackParams{
				Track:     false,
				Formatter: formatter.New(outputBuf, "text"),
				Response: models.TokenResponse{
					Token: ec.String("tokenvalue"),
				},
			}, buf: outputBuf},
			wantOut: wantOutTokenValue,
		},
		{
			name: "returns an error when track is set to false and the formatter errors",
			args: args{params: TrackParams{
				Track: false,
				Formatter: formatter.NewText(&formatter.TextConfig{
					Output: io.Discard,
				}),
				Response: models.TokenResponse{
					Token: ec.String("tokenvalue"),
				},
			}},
			err: errors.New(`template: no template "default" associated with template "text"`),
		},
		{
			name: "returns a formatted structure when track is set to true and returns tracking error",
			args: args{params: TrackParams{
				Track:     true,
				Formatter: formatter.New(secondOutputBuf, "text"),
				Response: models.TokenResponse{
					Token: ec.String("tokenvalue"),
				},
			}, buf: secondOutputBuf},
			wantOut: wantOutTokenValue,
			err: multierror.NewPrefixed("plan track change",
				errors.New("API cannot be nil"),
				errors.New("one of DeploymentID or ResourceID must be specified"),
				errors.New("kind cannot be empty"),
			),
		},
		{
			name: "prints the first error when track is set to true and also returns the formatter errors",
			args: args{params: TrackParams{
				TrackChangeParams: planutil.TrackChangeParams{
					Writer: output.NewDevice(thirdOutputBuf),
				},
				Track: true,
				Formatter: formatter.NewText(&formatter.TextConfig{
					Output: thirdOutputBuf,
				}),
				Response: models.TokenResponse{
					Token: ec.String("tokenvalue"),
				},
			}, buf: thirdOutputBuf},
			wantOut: `template: no template "default" associated with template "text"` + "\n",
			err: multierror.NewPrefixed("plan track change",
				errors.New("API cannot be nil"),
				errors.New("one of DeploymentID or ResourceID must be specified"),
				errors.New("kind cannot be empty"),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Track(tt.args.params); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Track() error = %v, wantErr %v", err, tt.err)
			}
			if tt.args.buf != nil {
				if got := tt.args.buf.String(); got != tt.wantOut {
					t.Errorf("Track() output = %v, wantOut %v", got, tt.wantOut)
				}
			}
		})
	}
}
