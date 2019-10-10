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

package ecctl

import (
	"bytes"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/output"
	multierror "github.com/hashicorp/go-multierror"
)

func TestConfigValidate(t *testing.T) {
	type fields struct {
		User         string
		Pass         string
		Host         string
		APIKey       string
		Region       string
		Output       string
		Verbose      bool
		Message      string
		Format       string
		OutputDevice *output.Device
		ErrorDevice  io.Writer
		Client       *http.Client
	}
	tests := []struct {
		name   string
		fields fields
		err    error
	}{
		{
			name: "Validate fails when output is not valid",
			fields: fields{
				Output: "INVALID OUTPUT",
				APIKey: "dummy",
			},
			err: &multierror.Error{Errors: []error{
				errInvalidOutputFormat,
				errInvalidOutputDevice,
				errInvalidErrorDevice,
			}},
		},
		{
			name: "Validate fails when output = json and custom format",
			fields: fields{
				Output: JSONOutput,
				Format: "{{ .Field }}",
				APIKey: "dummy",
			},
			err: &multierror.Error{Errors: []error{
				errCannotSpecifyJSONOutputAndCustomFormat,
				errInvalidOutputDevice,
				errInvalidErrorDevice,
			}},
		},
		{
			name: "Validate fails due to specifying both user / pass and APIKey",
			fields: fields{
				Output:       JSONOutput,
				Region:       "ece-region",
				APIKey:       "dummy",
				User:         "dummy",
				Pass:         "dummypass",
				OutputDevice: output.NewDevice(new(bytes.Buffer)),
			},
			err: &multierror.Error{Errors: []error{
				errInvalidBothAuthenticaitonSettings,
				errInvalidErrorDevice,
			}},
		},
		{
			name: "Validate fails due to empty credentials",
			fields: fields{
				Output:       JSONOutput,
				Region:       "ece-region",
				OutputDevice: output.NewDevice(new(bytes.Buffer)),
			},
			err: &multierror.Error{Errors: []error{
				errInvalidEmptyAuthenticaitonSettings,
				errInvalidErrorDevice,
			}},
		},
		{
			name: "Validate succeeds",
			fields: fields{
				Output:       JSONOutput,
				Region:       "ece-region",
				APIKey:       "dummy",
				OutputDevice: output.NewDevice(new(bytes.Buffer)),
				ErrorDevice:  new(bytes.Buffer),
			},
		},
		{
			name: "Validate succeeds with user/password",
			fields: fields{
				Output:       JSONOutput,
				Region:       "ece-region",
				User:         "dummy",
				Pass:         "dummypass",
				OutputDevice: output.NewDevice(new(bytes.Buffer)),
				ErrorDevice:  new(bytes.Buffer),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Config{
				User:         tt.fields.User,
				Pass:         tt.fields.Pass,
				Host:         tt.fields.Host,
				Region:       tt.fields.Region,
				APIKey:       tt.fields.APIKey,
				Output:       tt.fields.Output,
				Verbose:      tt.fields.Verbose,
				Message:      tt.fields.Message,
				Format:       tt.fields.Format,
				OutputDevice: tt.fields.OutputDevice,
				ErrorDevice:  tt.fields.ErrorDevice,
				Client:       tt.fields.Client,
			}
			if err := c.Validate(); !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
