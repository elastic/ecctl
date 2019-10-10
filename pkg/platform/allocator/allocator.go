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

package allocator

import (
	"errors"
	"strings"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/platform_infrastructure"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"

	"github.com/elastic/ecctl/pkg/util"
)

// GetParams is used to get an allocator
type GetParams struct {
	*api.API
	ID string
}

// Validate ensures that the parameters are correct
func (params GetParams) Validate() error {
	if params.API == nil {
		return util.ErrAPIReq
	}
	if params.ID == "" {
		return errors.New("get: id cannot be empty")
	}
	return nil
}

// Get obtains an allocator from an ID
func Get(params GetParams) (*models.AllocatorInfo, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.API.V1API.PlatformInfrastructure.GetAllocator(
		platform_infrastructure.NewGetAllocatorParams().
			WithAllocatorID(params.ID),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}

// ListParams is used to list allocators
type ListParams struct {
	*api.API
	Query string
	// Expected format is key:value slice. i.e. [key:val, key:value]
	FilterTags string
	ShowAll    bool
}

// Validate ensures that the parameters are correct
func (params ListParams) Validate() error {
	if params.API == nil {
		return util.ErrAPIReq
	}

	return nil
}

// List obtains the full list of allocators
func List(params ListParams) (*models.AllocatorOverview, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.API.V1API.PlatformInfrastructure.GetAllocators(
		platform_infrastructure.NewGetAllocatorsParams().
			WithQ(ec.String(params.Query)),
		params.AuthWriter,
	)

	if err != nil {
		return nil, api.UnwrapError(err)
	}
	if !params.ShowAll {
		for _, z := range res.Payload.Zones {
			z.Allocators = FilterConnectedOrWithInstances(z.Allocators)
		}
	}

	for _, z := range res.Payload.Zones {
		z.Allocators = FilterByTag(tagsToMap(params.FilterTags), z.Allocators)
	}

	return res.Payload, nil
}

// MaintenanceParams is used to set / unset maintenance mode
type MaintenanceParams struct {
	*api.API
	ID string
}

// Validate ensures that the parameters are correct
func (params MaintenanceParams) Validate() error {
	if params.API == nil {
		return util.ErrAPIReq
	}
	if params.ID == "" {
		return errors.New("maintenance: id cannot be empty")
	}
	return nil
}

// StartMaintenance sets an allocator to maintenance mode
func StartMaintenance(params MaintenanceParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	return util.ReturnErrOnly(
		params.API.V1API.PlatformInfrastructure.StartAllocatorMaintenanceMode(
			platform_infrastructure.NewStartAllocatorMaintenanceModeParams().
				WithAllocatorID(params.ID),
			params.AuthWriter,
		),
	)
}

// StopMaintenance unsets an allocator to maintenance mode
func StopMaintenance(params MaintenanceParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	return util.ReturnErrOnly(
		params.API.V1API.PlatformInfrastructure.StopAllocatorMaintenanceMode(
			platform_infrastructure.NewStopAllocatorMaintenanceModeParams().
				WithAllocatorID(params.ID),
			params.AuthWriter,
		),
	)
}

// FilterByTag filters a list of allocators based on metadata tags and returns a new matched list
func FilterByTag(tags map[string]string, allocators []*models.AllocatorInfo) []*models.AllocatorInfo {
	var matchedAllocators []*models.AllocatorInfo

	if len(tags) > 0 {
		for _, a := range allocators {
			matches := 0
			for k, v := range tags {
				for _, m := range a.Metadata {
					if *m.Key == k && *m.Value == v {
						matches++
						break
					}
				}
			}
			if matches == len(tags) {
				matchedAllocators = append(matchedAllocators, a)
			}
		}
		return matchedAllocators
	}

	return allocators
}

// FilterConnectedOrWithInstances filters a list of allocators and returns only the connected ones or those who have more than one instance
func FilterConnectedOrWithInstances(allocators []*models.AllocatorInfo) []*models.AllocatorInfo {
	var matchedAllocators []*models.AllocatorInfo

	for _, a := range allocators {
		if *a.Status.Connected || len(a.Instances) > 0 {
			matchedAllocators = append(matchedAllocators, a)
		}
	}
	return matchedAllocators
}

func tagsToMap(filterArgs string) map[string]string {
	filterArgs = strings.ReplaceAll(filterArgs, "[", "")
	filterArgs = strings.ReplaceAll(filterArgs, "]", "")
	tags := strings.Split(filterArgs, ",")
	var tagsMap = make(map[string]string)

	for _, t := range tags {
		tag := strings.Split(t, ":")
		if len(tag) == 2 {
			tagKey := tag[0]
			tagValue := tag[1]

			tagsMap[tagKey] = tagValue
		}
	}

	return tagsMap
}
