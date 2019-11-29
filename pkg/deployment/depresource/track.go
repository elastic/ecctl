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
	"errors"
	"fmt"
	"sync"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/output"
	"github.com/elastic/cloud-sdk-go/pkg/plan"
	"github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

// TrackResourcesParams is consumed by TrackResources.
type TrackResourcesParams struct {
	*api.API
	Resources    []*models.DeploymentResource
	OutputDevice *output.Device
}

// Validate ensures the parameters are usable by the consuming function.
func (params TrackResourcesParams) Validate() error {
	var merr = new(multierror.Error)
	if params.API == nil {
		merr = multierror.Append(merr, util.ErrAPIReq)
	}

	if params.OutputDevice == nil {
		merr = multierror.Append(merr, errors.New("resource track: output device cannot be nil"))
	}

	return merr.ErrorOrNil()
}

// TrackResources will track changes made to a deployment's resources. Creating
// a gorouine for each cluster which is going to be tracked.
// The returned error is a multierror, containing all the encountered errors.
// This function is rather a stop-gap mechanism to be able to easily track
// responses returned by the Deployments API but it suffice going forward.
// WARNING: Does not support tracking "appsearch" resource changes.
func TrackResources(params TrackResourcesParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	if len(params.Resources) == 0 {
		return nil
	}

	var merr = new(multierror.Error)
	var errChan = make(chan error)
	var wg sync.WaitGroup

	for _, resource := range params.Resources {
		wg.Add(1)

		go func(r *models.DeploymentResource, w *sync.WaitGroup) {
			defer w.Done()
			if *r.Kind == "appsearch" {
				errChan <- fmt.Errorf("cannot track appsearch resource id %s", *r.ID)
				return
			}

			errChan <- util.TrackCluster(util.TrackClusterParams{
				TrackParams: plan.TrackParams{
					API:  params.API,
					ID:   *r.ID,
					Kind: *r.Kind,
				},
				Output: params.OutputDevice,
			})
		}(resource, &wg)
	}

	// Close the channel once all the errors have been received.
	go func(w *sync.WaitGroup) {
		w.Wait()
		close(errChan)
	}(&wg)

	for err := range errChan {
		merr = multierror.Append(merr, err)
	}

	return merr.ErrorOrNil()
}
