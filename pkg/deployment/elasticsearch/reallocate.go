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
	"strings"
	"time"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/clusters_elasticsearch"
	"github.com/elastic/cloud-sdk-go/pkg/output"
	"github.com/elastic/cloud-sdk-go/pkg/plan"
	multierror "github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/deployment/note"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/util"
)

const reallocateMessage = "reallocated %s %s by ecctl"

// ReallocateParams is used by ReallocateES as a config struct
type ReallocateParams struct {
	util.ClusterParams
	Instances []string

	// Commenting user
	User          string
	InstancesDown *bool
	OutputDevice  *output.Device
	PollFrequency time.Duration
	MaxRetries    uint8
}

// Validate ensures that the parameters are usable by the consuming function.
func (params ReallocateParams) Validate() error {
	var err = multierror.Append(new(multierror.Error), params.ClusterParams.Validate())

	if params.OutputDevice == nil {
		err = multierror.Append(err, errors.New("output device cannot be nil"))
	}

	return err.ErrorOrNil()
}

// Reallocate will reallocate the Elasticsearch cluster instances, if no
// Instances are specified, all of the instances will be moved. The operation
// will attempt to create a cluster comment to leave context in the cluster
func Reallocate(params ReallocateParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	if len(params.Instances) == 0 || params.Instances[0] == "" {
		params.Instances = make([]string, 0)
		res, err := GetCluster(GetClusterParams{
			ClusterParams: params.ClusterParams,
			Plans:         true,
		})

		if err != nil {
			return err
		}

		for _, i := range res.Topology.Instances {
			params.Instances = append(params.Instances, *i.InstanceName)
		}
	}

	// We don't really care if the operation succeeds or not, since adding a
	// note to the deployment is not critical and is only a nice to have.
	//nolint
	note.Add(note.AddParams{
		API:         params.API,
		ID:          params.ClusterID,
		Type:        "elasticsearch",
		Message:     fmt.Sprintf(reallocateMessage, "elasticsearch", strings.Join(params.Instances, " ")),
		UserID:      params.User,
		Commentator: ecctl.GetOperationInstance(),
	})

	if _, err := params.V1API.ClustersElasticsearch.MoveEsClusterInstances(
		clusters_elasticsearch.NewMoveEsClusterInstancesParams().
			WithClusterID(params.ClusterID).
			WithInstanceIds(params.Instances).
			WithInstancesDown(params.InstancesDown),
		params.AuthWriter,
	); err != nil {
		return api.UnwrapError(err)
	}

	return util.TrackCluster(util.TrackClusterParams{
		Output: params.OutputDevice,
		TrackParams: plan.TrackParams{
			API:           params.API,
			ID:            params.ClusterID,
			Kind:          "elasticsearch",
			PollFrequency: params.PollFrequency,
			MaxRetries:    params.MaxRetries,
		},
	})
}
