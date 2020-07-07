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

package testutils

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/output"
	"github.com/stretchr/testify/assert"

	"github.com/elastic/ecctl/pkg/ecctl"
)

func Test_newConfig(t *testing.T) {
	type args struct {
		cfg MockCfg
	}
	tests := []struct {
		name string
		args args
		want ecctl.Config
	}{
		{
			name: "empty params succeed",
			args: args{cfg: MockCfg{}},
			want: ecctl.Config{
				Client:       mock.NewClient(),
				OutputDevice: output.NewDevice(new(bytes.Buffer)),
				ErrorDevice:  new(bytes.Buffer),
				Output:       "json",
				Region:       "ece-region",
				Host:         fmt.Sprintf("https://%s", api.DefaultMockHost),
				APIKey:       defaultAPIKey,
				Force:        false,
			},
		},
		{
			name: "overriding params succeed",
			args: args{cfg: MockCfg{
				Responses: []mock.Response{
					{Error: errors.New("some error")},
				},
				Out:          os.Stdout,
				Err:          os.Stderr,
				OutputFormat: "text",
				Region:       "us-east-1",
				Force:        true,
				Verbose:      true,
			}},
			want: ecctl.Config{
				Client:       mock.NewClient(mock.Response{Error: errors.New("some error")}),
				OutputDevice: output.NewDevice(os.Stdout),
				ErrorDevice:  os.Stderr,
				Output:       "text",
				Region:       "us-east-1",
				Host:         fmt.Sprintf("https://%s", api.DefaultMockHost),
				APIKey:       defaultAPIKey,
				Force:        true,
				Verbose:      true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newConfig(tt.args.cfg)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMockApp(t *testing.T) {
	cleanup := mockApp(t, MockCfg{})
	assert.NotNil(t, ecctl.Get())
	cleanup()
	assert.Nil(t, ecctl.Get())
}
