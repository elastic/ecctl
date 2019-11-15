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

package filteredgroup

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
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

func TestShow(t *testing.T) {
	var proxiesFilteredGroup = `
	{
      "expected_proxies_count": 5,
      "filters": [
        {
          "key": "proxyType",
          "value": "main-nextgen"
        }
      ],
      "id": "test1"
	}`
	type args struct {
		params CommonParams
	}
	tests := []struct {
		name string
		args args
		want *models.ProxiesFilteredGroup
		err  error
	}{
		{
			name: "Proxies filtered group show succeeds",
			args: args{params: CommonParams{
				API: api.NewMock(mock.Response{Response: http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
					Body:       ioutil.NopCloser(strings.NewReader(proxiesFilteredGroup)),
				}}),
				ID: "test1",
			}},
			want: &models.ProxiesFilteredGroup{
				ExpectedProxiesCount: ec.Int32(5),
				Filters: []*models.ProxiesFilter{
					{
						Key:   ec.String("proxyType"),
						Value: ec.String("main-nextgen"),
					},
				},
				ID: *ec.String("test1"),
			},
		},
		{
			name: "Proxies filtered group show fails with 403 Forbidden",
			args: args{params: CommonParams{
				API: api.NewMock(mock.New500Response(mock.NewStringBody(`{"error": "some error"}`))),
				ID:  "test1",
			}},
			err: errors.New(`{"error": "some error"}`),
		},
		{
			name: "Proxies filtered group show fails due to empty API",
			args: args{params: CommonParams{
				API: nil,
				ID:  "test1",
			}},
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Get(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	var proxiesFilteredGroup = `
	{
      "expected_proxies_count": 15,
      "filters": [
        {
          "key": "proxyType",
          "value": "main-nextgen"
        }
      ],
      "id": "test2"
	}`
	type args struct {
		params CreateParams
	}
	tests := []struct {
		name string
		args args
		want *models.ProxiesFilteredGroup
		err  error
	}{
		{
			name: "Proxies filtered group create succeeds",
			args: args{params: CreateParams{
				CommonParams: CommonParams{
					API: api.NewMock(mock.Response{Response: http.Response{
						StatusCode: http.StatusOK,
						Status:     http.StatusText(http.StatusOK),
						Body:       ioutil.NopCloser(strings.NewReader(proxiesFilteredGroup)),
					}}),
					ID: "test2",
				},
				Filters: map[string]string{
					"proxyType": "main-nextgen",
				},
				ExpectedProxiesCount: 15,
			}},
			want: &models.ProxiesFilteredGroup{
				ExpectedProxiesCount: ec.Int32(15),
				Filters: []*models.ProxiesFilter{
					{
						Key:   ec.String("proxyType"),
						Value: ec.String("main-nextgen"),
					},
				},
				ID: *ec.String("test2"),
			},
		},
		{
			name: "Proxies filtered group create fails with 403 Forbidden",
			args: args{params: CreateParams{
				CommonParams: CommonParams{
					API: api.NewMock(mock.Response{Response: http.Response{
						StatusCode: http.StatusForbidden,
						Status:     http.StatusText(http.StatusForbidden),
						Body:       mock.NewStringBody(`{"error": "some error"}`),
					}}),
					ID: "test2",
				},
				Filters: map[string]string{
					"proxyType": "main-nextgen",
				},
				ExpectedProxiesCount: 15,
			}},
			err: errors.New(`{"error": "some error"}`),
		},
		{
			name: "Proxies filtered group create expected proxies count greater than 0 fails",
			args: args{params: CreateParams{
				CommonParams: CommonParams{
					API: api.NewMock(mock.Response{Response: http.Response{
						StatusCode: http.StatusOK,
						Status:     http.StatusText(http.StatusOK),
						Body:       mock.NewStringBody(`{}`),
					}}),
					ID: "test2",
				},
				Filters: map[string]string{
					"proxyType": "main-nextgen",
				},
				ExpectedProxiesCount: 0,
			}},
			err: &multierror.Error{Errors: []error{
				errExpectedProxiesCountCannotBeLesserThanZero,
			}},
		},
		{
			name: "Proxies filtered group create empty filters fails",
			args: args{params: CreateParams{
				CommonParams: CommonParams{
					API: api.NewMock(mock.Response{Response: http.Response{
						StatusCode: http.StatusOK,
						Status:     http.StatusText(http.StatusOK),
						Body:       mock.NewStringBody(`{}`),
					}}),
					ID: "test2",
				},
				Filters:              map[string]string{},
				ExpectedProxiesCount: 15,
			}},
			err: &multierror.Error{Errors: []error{
				errFiltersCannotBeEmpty,
			}},
		},
		{
			name: "Proxies filtered group create fails due to empty API",
			args: args{params: CreateParams{
				CommonParams: CommonParams{
					API: nil,
					ID:  "test2",
				},
				Filters: map[string]string{
					"proxyType": "main-nextgen",
				},
				ExpectedProxiesCount: 15,
			}},
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Create(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	var EmptyBody = ``
	type args struct {
		params CommonParams
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "Proxies filtered group delete succeeds",
			args: args{params: CommonParams{
				API: api.NewMock(mock.Response{Response: http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
					Body:       ioutil.NopCloser(strings.NewReader(EmptyBody)),
				}}),
				ID: "test2",
			}},
			err: nil,
		},
		{
			name: "Proxies filtered group delete fails with 403 Forbidden",
			args: args{params: CommonParams{
				API: api.NewMock(mock.New500Response(mock.NewStringBody(`{"error": "some error"}`))),
				ID:  "test1",
			}},
			err: errors.New(`{"error": "some error"}`),
		},
		{
			name: "Proxies filtered group delete fails due to empty API",
			args: args{params: CommonParams{
				API: nil,
				ID:  "test1",
			}},
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Delete(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	var proxiesFilteredGroup = `
	{
      "expected_proxies_count": 15,
      "filters": [
        {
          "key": "proxyType",
          "value": "main-nextgen"
        }
      ],
      "id": "test2"
	}`
	type args struct {
		params UpdateParams
	}
	tests := []struct {
		name string
		args args
		want *models.ProxiesFilteredGroup
		err  error
	}{
		{
			name: "Proxies filtered group update succeeds",
			args: args{params: UpdateParams{
				CreateParams: CreateParams{
					CommonParams: CommonParams{
						API: api.NewMock(mock.Response{Response: http.Response{
							StatusCode: http.StatusOK,
							Status:     http.StatusText(http.StatusOK),
							Body:       ioutil.NopCloser(strings.NewReader(proxiesFilteredGroup)),
						}}),
						ID: "test2",
					},
					Filters: map[string]string{
						"proxyType": "main-nextgen",
					},
					ExpectedProxiesCount: 15,
				},
				Version: 1,
			}},
			want: &models.ProxiesFilteredGroup{
				ExpectedProxiesCount: ec.Int32(15),
				Filters: []*models.ProxiesFilter{
					{
						Key:   ec.String("proxyType"),
						Value: ec.String("main-nextgen"),
					},
				},
				ID: *ec.String("test2"),
			},
		},
		{
			name: "Proxies filtered group update fails with 403 Forbidden",
			args: args{params: UpdateParams{
				CreateParams: CreateParams{
					CommonParams: CommonParams{
						API: api.NewMock(mock.Response{Response: http.Response{
							StatusCode: http.StatusForbidden,
							Status:     http.StatusText(http.StatusForbidden),
							Body:       mock.NewStringBody(`{"error": "some error"}`),
						}}),
						ID: "test2",
					},
					Filters: map[string]string{
						"proxyType": "main-nextgen",
					},

					ExpectedProxiesCount: 15,
				},
				Version: 1,
			},
			},
			err: errors.New(`{"error": "some error"}`),
		},
		{
			name: "Proxies filtered group update expected proxies count greater than 0 fails",
			args: args{params: UpdateParams{
				CreateParams: CreateParams{
					CommonParams: CommonParams{
						API: api.NewMock(mock.Response{Response: http.Response{
							StatusCode: http.StatusOK,
							Status:     http.StatusText(http.StatusOK),
							Body:       mock.NewStringBody(`{}`),
						}}),
						ID: "test2",
					},
					Filters: map[string]string{
						"proxyType": "main-nextgen",
					},

					ExpectedProxiesCount: 0,
				},
				Version: 1,
			},
			},
			err: &multierror.Error{Errors: []error{
				errExpectedProxiesCountCannotBeLesserThanZero,
			}},
		},
		{
			name: "Proxies filtered group update empty filters fails",
			args: args{params: UpdateParams{
				CreateParams: CreateParams{
					CommonParams: CommonParams{
						API: api.NewMock(mock.Response{Response: http.Response{
							StatusCode: http.StatusOK,
							Status:     http.StatusText(http.StatusOK),
							Body:       mock.NewStringBody(`{}`),
						}}),
						ID: "test2",
					},
					Filters: map[string]string{},

					ExpectedProxiesCount: 15,
				},
				Version: 1,
			},
			},
			err: &multierror.Error{Errors: []error{
				errFiltersCannotBeEmpty,
			}},
		},
		{
			name: "Proxies filtered group update version less than 0 fails",
			args: args{params: UpdateParams{
				CreateParams: CreateParams{
					CommonParams: CommonParams{
						API: api.NewMock(mock.Response{Response: http.Response{
							StatusCode: http.StatusOK,
							Status:     http.StatusText(http.StatusOK),
							Body:       mock.NewStringBody(`{}`),
						}}),
						ID: "test2",
					},
					Filters: map[string]string{
						"proxyType": "main-nextgen",
					},
					ExpectedProxiesCount: 15,
				},
				Version: -1,
			},
			},
			err: &multierror.Error{Errors: []error{
				errVersionCannotBeLesserTahZero,
			}},
		},
		{
			name: "Proxies filtered group update fails due to empty API",
			args: args{params: UpdateParams{
				CreateParams: CreateParams{
					CommonParams: CommonParams{
						API: nil,
						ID:  "test2",
					},
					Filters: map[string]string{
						"proxyType": "main-nextgen",
					},

					ExpectedProxiesCount: 15,
				},
				Version: 1,
			},
			},
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Update(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList(t *testing.T) {
	var proxiesHealth = `
	{
	  "filtered_groups": [
		{
		  "group": {
			"id": "main-proxies",
			"filters": [
			  {
				"key": "proxyType",
				"value": "main-nextgen"
			  }
			],
			"expected_proxies_count": 1
		  },
		  "observed_proxies_count": 7,
		  "status": "Green"
		},
		{
		  "group": {
			"id": "apm-proxies",
			"filters": [
			  {
				"key": "proxyType",
				"value": "apm"
			  }
			],
			"expected_proxies_count": 1
		  },
		  "observed_proxies_count": 0,
		  "status": "Red"
		}
	  ],
	  "expected_proxies_count": 2,
	  "status": "Yellow",
	  "allocations": [
		{
		  "allocations_type": "apm",
		  "max_allocations": 82,
		  "proxies_at_max_allocations": 10
		},
		{
		  "allocations_type": "elasticsearch",
		  "max_allocations": 293,
		  "proxies_at_max_allocations": 3
		},
		{
		  "allocations_type": "sitesearch",
		  "max_allocations": 0,
		  "proxies_at_max_allocations": 10
		},
		{
		  "allocations_type": "appsearch",
		  "max_allocations": 0,
		  "proxies_at_max_allocations": 10
		},
		{
		  "allocations_type": "enterprisesearch",
		  "max_allocations": 0,
		  "proxies_at_max_allocations": 10
		},
		{
		  "allocations_type": "kibana",
		  "max_allocations": 247,
		  "proxies_at_max_allocations": 3
		}
	  ],
	  "observed_proxies_count": 10
	}`
	type args struct {
		params CommonParams
	}
	tests := []struct {
		name string
		args args
		want []*models.ProxiesFilteredGroupHealth
		err  error
	}{
		{
			name: "Proxies filtered group list succeeds",
			args: args{params: CommonParams{
				API: api.NewMock(mock.Response{Response: http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
					Body:       ioutil.NopCloser(strings.NewReader(proxiesHealth)),
				}}),
				ID: "all",
			}},
			want: []*models.ProxiesFilteredGroupHealth{
				{
					Group: &models.ProxiesFilteredGroup{
						ExpectedProxiesCount: ec.Int32(1),
						Filters: []*models.ProxiesFilter{
							{
								Key:   ec.String("proxyType"),
								Value: ec.String("main-nextgen"),
							},
						},
						ID: *ec.String("main-proxies"),
					},
					ObservedProxiesCount: ec.Int32(7),
					Status:               ec.String("Green"),
				},
				{
					Group: &models.ProxiesFilteredGroup{
						ExpectedProxiesCount: ec.Int32(1),
						Filters: []*models.ProxiesFilter{
							{
								Key:   ec.String("proxyType"),
								Value: ec.String("apm"),
							},
						},
						ID: *ec.String("apm-proxies"),
					},
					ObservedProxiesCount: ec.Int32(0),
					Status:               ec.String("Red"),
				},
			},
		},
		{
			name: "Proxies filtered group list fails with 403 Forbidden",
			args: args{params: CommonParams{
				API: api.NewMock(mock.New500Response(mock.NewStringBody(`{"error": "some error"}`))),
				ID:  "all",
			}},
			err: errors.New(`{"error": "some error"}`),
		},
		{
			name: "Proxies filtered group list fails due to empty API",
			args: args{params: CommonParams{
				API: nil,
				ID:  "all",
			}},
			err: &multierror.Error{Errors: []error{
				util.ErrAPIReq,
			}},
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
