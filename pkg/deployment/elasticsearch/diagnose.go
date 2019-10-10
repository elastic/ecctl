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
	"bytes"
	"errors"
	"io"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_elasticsearch"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

// DiagnoseParams is consumed by Diagnose.
type DiagnoseParams struct {
	util.ClusterParams
	Writer io.Writer
}

// Validate ensures the parameters are usable by the consumer.
func (params DiagnoseParams) Validate() error {
	err := new(multierror.Error)
	if params.Writer == nil {
		err = multierror.Append(err, errors.New("writer cannot be nil"))
	}

	err = multierror.Append(err, params.ClusterParams.Validate())

	return err.ErrorOrNil()
}

// Diagnose generates a diagnostic bundle and writes it to the specified
// Writer.
func Diagnose(params DiagnoseParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	r, err := params.V1API.ClustersElasticsearch.GenerateEsClusterDiagnostics(
		clusters_elasticsearch.NewGenerateEsClusterDiagnosticsParams().
			WithClusterID(params.ClusterID),
		params.AuthWriter,
	)
	if err != nil {
		return api.UnwrapError(err)
	}

	return util.ReturnErrOnly(
		io.Copy(params.Writer, bytes.NewReader(r.Payload)),
	)
}
