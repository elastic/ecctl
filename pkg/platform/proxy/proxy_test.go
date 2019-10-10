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

package proxy

import (
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"

	"github.com/elastic/ecctl/pkg/util"
)

func TestList(t *testing.T) {
	var proxyList = `
	{
		"proxies": [
		  {
            "healthy": true,
            "host_ip": "",
            "metadata": {},
            "proxy_id": "87b2c433c761",
            "public_hostname": ""
		  }
		]
	}`[1:]
	type args struct {
		params Params
	}
	tests := []struct {
		name string
		args args
		want *models.ProxyOverview
		err  error
	}{
		{
			name: "Proxy list succeeds",
			args: args{params: Params{
				API: api.NewMock(mock.Response{Response: http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
					Body:       ioutil.NopCloser(strings.NewReader(proxyList)),
				}}),
			}},
			want: &models.ProxyOverview{
				Proxies: []*models.ProxyInfo{
					{
						Healthy:        ec.Bool(true),
						HostIP:         ec.String(""),
						Metadata:       make(map[string]interface{}),
						ProxyID:        ec.String("87b2c433c761"),
						PublicHostname: ec.String(""),
					},
				},
			},
		},
		{
			name: "Proxy list fails",
			args: args{params: Params{
				API: api.NewMock(mock.Response{Response: http.Response{
					StatusCode: http.StatusForbidden,
					Status:     http.StatusText(http.StatusForbidden),
					Body:       mock.NewStringBody(`{}`),
				}}),
			}},
			err: errors.New("unknown error (status 403)"),
		},
		{
			name: "Proxy list fails due to empty API",
			args: args{params: Params{
				API: nil,
			}},
			err: util.ErrAPIReq,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := List(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("List() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGet(t *testing.T) {
	var proxyGet = `
	{
            "healthy": true,
            "host_ip": "",
            "metadata": {},
            "proxy_id": "87b2c433c761",
            "public_hostname": ""
    }`[1:]

	type args struct {
		params GetParams
	}
	tests := []struct {
		name string
		args args
		want *models.ProxyInfo
		err  error
	}{
		{
			name: "Get proxy succeeds",
			args: args{params: GetParams{
				ID: "87b2c433c761",
				Params: Params{
					API: api.NewMock(mock.Response{Response: http.Response{
						StatusCode: http.StatusOK,
						Status:     http.StatusText(http.StatusOK),
						Body:       ioutil.NopCloser(strings.NewReader(proxyGet)),
					}}),
				},
			},
			},
			want: &models.ProxyInfo{
				Healthy:        ec.Bool(true),
				HostIP:         ec.String(""),
				Metadata:       make(map[string]interface{}),
				ProxyID:        ec.String("87b2c433c761"),
				PublicHostname: ec.String(""),
			},
		},
		{
			name: "Proxy get fails",
			args: args{params: GetParams{
				ID: "87b2c433c761",
				Params: Params{
					API: api.NewMock(mock.Response{Response: http.Response{
						StatusCode: http.StatusForbidden,
						Status:     http.StatusText(http.StatusForbidden),
						Body:       mock.NewStringBody(`{}`),
					}}),
				},
			}},
			err: errors.New("unknown error (status 403)"),
		},
		{
			name: "Get proxy fails due to empty params.ID",
			args: args{params: GetParams{
				ID:     "",
				Params: Params{new(api.API)},
			},
			},
			err: errors.New("proxy id cannot be empty"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("List() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}
