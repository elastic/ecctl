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
	"os"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/auth"
	"github.com/elastic/cloud-sdk-go/pkg/output"
	"github.com/stretchr/testify/assert"
)

func Test_newAPIConfig(t *testing.T) {
	apiKey := auth.APIKey("somekey")
	type args struct {
		cfg Config
	}
	tests := []struct {
		name string
		args args
		want api.Config
		err  string
	}{
		{
			name: "fails on invalid config",
			args: args{},
			err:  "invalid configuration options specified: 4 errors occurred:\n\t* api_key or user and pass must be specified\n\t* error device must not be nil\n\t* output device must not be nil\n\t* output must be one either json or text\n\n",
		},
		{
			name: "succeeds without verbose",
			args: args{cfg: Config{
				APIKey:       "somekey",
				Output:       "text",
				OutputDevice: output.NewDevice(new(bytes.Buffer)),
				ErrorDevice:  new(bytes.Buffer),
			}},
			want: api.Config{
				AuthWriter: &apiKey,
				VerboseSettings: api.VerboseSettings{
					Device:     output.NewDevice(new(bytes.Buffer)),
					RedactAuth: true,
				},
				ErrorDevice: new(bytes.Buffer),
			},
		},
		{
			name: "succeeds with verbose",
			args: args{cfg: Config{
				APIKey:             "somekey",
				Output:             "text",
				OutputDevice:       output.NewDevice(new(bytes.Buffer)),
				ErrorDevice:        new(bytes.Buffer),
				Verbose:            true,
				VerboseCredentials: true,
				VerboseFile:        "request.log",
			}},
			want: api.Config{
				AuthWriter: &apiKey,
				VerboseSettings: api.VerboseSettings{
					Verbose: true,
				},
				ErrorDevice: new(bytes.Buffer),
			},
		},
		{
			name: "fails creating a file on an inexistent path",
			args: args{cfg: Config{
				APIKey:             "somekey",
				Output:             "text",
				OutputDevice:       output.NewDevice(new(bytes.Buffer)),
				ErrorDevice:        new(bytes.Buffer),
				Verbose:            true,
				VerboseCredentials: true,
				VerboseFile:        "/some/path/no/exist/request.log",
			}},
			err: "failed creating verbose file \"/some/path/no/exist/request.log\": open /some/path/no/exist/request.log: no such file or directory",
		},
		{
			name: "fails on invalid credentials",
			args: args{cfg: Config{
				APIKey:       "some",
				User:         "some",
				Pass:         "some",
				Output:       "text",
				OutputDevice: output.NewDevice(new(bytes.Buffer)),
				ErrorDevice:  new(bytes.Buffer),
			}},
			err: "invalid configuration options specified: 1 error occurred:\n\t* cannot specify both api_key and user / pass\n\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newAPIConfig(tt.args.cfg)
			if err != nil {
				assert.EqualError(t, err, tt.err)
			}

			if tt.args.cfg.VerboseFile != "" && got.Device != nil {
				if f, ok := got.Device.(*os.File); ok {
					assert.Equal(t, tt.args.cfg.VerboseFile, f.Name())
					defer func() {
						f.Close()
						os.Remove(f.Name())
					}()

					// not possible to assert the *os.File due to different
					// File Descriptors.
					tt.want.Device = f
				}
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
