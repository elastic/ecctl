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

	// List of current "live" resources,
	Resources []*models.DeploymentResource
	// Orphaned are the resources that will become orphaned after a deployment
	// update.
	Orphaned *models.Orphaned

	OutputDevice *output.Device
}

type resource struct {
	id, kind string
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

	var resources = make([]resource, 0, len(params.Resources))
	for _, r := range params.Resources {
		resources = append(resources, resource{
			id:   *r.ID,
			kind: *r.Kind,
		})
	}

	if params.Orphaned != nil {
		for _, id := range params.Orphaned.Apm {
			resources = append(resources, resource{
				id:   id,
				kind: "apm",
			})
		}
		for _, id := range params.Orphaned.Appsearch {
			resources = append(resources, resource{
				id:   id,
				kind: "appsearch",
			})
		}
		for _, id := range params.Orphaned.Kibana {
			resources = append(resources, resource{
				id:   id,
				kind: "kibana",
			})
		}
		for _, es := range params.Orphaned.Elasticsearch {
			resources = append(resources, resource{
				id:   *es.ID,
				kind: "elasticsearch",
			})
		}
	}

	var merr = new(multierror.Error)
	var errChan = make(chan error)
	var wg sync.WaitGroup
	for _, r := range resources {
		wg.Add(1)
		go trackResource(params.API, params.OutputDevice, r, &wg, errChan)
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

func trackResource(a *api.API, out *output.Device, r resource, w *sync.WaitGroup, errChan chan error) {
	defer w.Done()
	if r.kind == "appsearch" {
		errChan <- fmt.Errorf("cannot track appsearch resource id \"%s\"", r.id)
		return
	}

	errChan <- util.TrackCluster(util.TrackClusterParams{
		TrackParams: plan.TrackParams{
			API:  a,
			ID:   r.id,
			Kind: r.kind,
		},
		Output: out,
	})
}
