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

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/platform_infrastructure"
	"github.com/elastic/cloud-sdk-go/pkg/models"

	"github.com/elastic/ecctl/pkg/util"
)

var (
	errIDCannotBeEmpty = errors.New("proxy id cannot be empty")
)

// Params is the generic set of parameters used for any proxy call
type Params struct {
	*api.API
}

// Validate checks the parameters
func (params Params) Validate() error {
	if params.API == nil {
		return util.ErrAPIReq
	}

	return nil
}

// GetParams is the set of parameters required for retrieving a proxy
type GetParams struct {
	Params
	ID string
}

// Validate checks the parameters
func (params GetParams) Validate() error {
	if params.ID == "" {
		return errIDCannotBeEmpty
	}

	return params.Params.Validate()
}

// List gets the list of proxies for a region
func List(params Params) (*models.ProxyOverview, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	proxies, err := params.API.V1API.PlatformInfrastructure.GetProxies(
		platform_infrastructure.NewGetProxiesParams(),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return proxies.Payload, nil
}

// Get returns information about a specific proxy
func Get(params GetParams) (*models.ProxyInfo, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	proxy, err := params.API.V1API.PlatformInfrastructure.GetProxy(
		platform_infrastructure.NewGetProxyParams().
			WithProxyID(params.ID),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return proxy.Payload, nil
}
