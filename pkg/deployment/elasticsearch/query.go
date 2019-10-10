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
	"fmt"
	"net/http"

	multierror "github.com/hashicorp/go-multierror"
	escliapp "github.com/marclop/elasticsearch-cli/app"

	"github.com/elastic/ecctl/pkg/util"
)

// skipMaintenanceHeaders tells the EC proxy layer to still send requests to the
// underlying cluster instances even if they are in maintenance mode
var skipMaintenanceHeaders = map[string]string{
	"X-Found-Bypass-Maintenance": "true",
}

// QueryParams describes the options passed to OpenElasticsearchConsole
type QueryParams struct {
	util.AuthenticationParams
	util.ClusterParams
	RestRequest

	Client      *http.Client
	Interactive bool
	Verbose     bool
}

// Validate ensures that the params are usable by the consuming function.
func (params QueryParams) Validate() error {
	var err = new(multierror.Error)

	err = multierror.Append(err, params.ClusterParams.Validate())

	if !params.Interactive && params.RestRequest.Method == "" {
		err = multierror.Append(err, errors.New("method needs to be specified"))
	}

	return err.ErrorOrNil()
}

// RestRequest sets the one-time request that will be handled by
// elasticsearch-cli
type RestRequest struct {
	Method, Path, Body string
}

// Query performs either an API call to an elasticsearch cluster or opens an
// interactive console connected to the public facing URL.
func Query(params QueryParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	res, err := GetCluster(GetClusterParams{
		ClusterParams: params.ClusterParams,
		Metadata:      true,
	})
	if err != nil {
		return err
	}

	esCli, err := escliapp.New(&escliapp.Config{
		Host:         fmt.Sprint("https://", res.Metadata.Endpoint),
		Port:         util.DefaultESPort,
		Client:       params.Client,
		User:         params.User,
		Pass:         params.Pass,
		Verbose:      params.Verbose,
		PollInterval: util.DefaultIndexPollerRate,
		Timeout:      util.DefaultClientTimeout,
		Headers:      skipMaintenanceHeaders,
		Insecure:     params.Insecure,
	})
	if err != nil {
		return err
	}

	if params.Interactive {
		return esCli.Interactive()
	}

	return esCli.HandleCli([]string{
		params.Method,
		params.Path,
		params.Body,
	})
}
