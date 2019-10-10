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

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/client/platform_infrastructure"
	"github.com/elastic/cloud-sdk-go/pkg/models"

	"github.com/elastic/ecctl/pkg/util"
)

// MetadataSetParams is is used to set a single allocator metadata key
type MetadataSetParams struct {
	*api.API
	ID, Key, Value string
}

// Validate ensures that the parameters are correct
func (params MetadataSetParams) Validate() error {
	if params.API == nil {
		return util.ErrAPIReq
	}
	if params.ID == "" {
		return errors.New("allocator metadata: id cannot be empty")
	}
	if params.Key == "" {
		return errors.New("allocator metadata: key cannot be empty")
	}
	if params.Value == "" {
		return errors.New("allocator metadata: key value cannot be empty")
	}
	return nil
}

// SetAllocatorMetadataItem sets a single metadata item to a given allocators metadata
func SetAllocatorMetadataItem(params MetadataSetParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	_, err := params.API.V1API.PlatformInfrastructure.SetAllocatorMetadataItem(
		platform_infrastructure.NewSetAllocatorMetadataItemParams().
			WithAllocatorID(params.ID).
			WithKey(params.Key).
			WithBody(&models.MetadataItemValue{Value: &params.Value}),
		params.AuthWriter,
	)

	if err != nil {
		return api.UnwrapError(err)
	}

	return nil
}

// MetadataDeleteParams is used to delete a single metadata key
type MetadataDeleteParams struct {
	*api.API
	ID, Key string
}

// Validate ensures that the parameters are correct
func (params MetadataDeleteParams) Validate() error {
	if params.API == nil {
		return util.ErrAPIReq
	}
	if params.ID == "" {
		return errors.New("allocator metadata: id cannot be empty")
	}
	if params.Key == "" {
		return errors.New("allocator metadata: key cannot be empty")
	}
	return nil
}

// DeleteAllocatorMetadataItem delete a single metadata item to a given allocators metadata
func DeleteAllocatorMetadataItem(params MetadataDeleteParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	_, err := params.API.V1API.PlatformInfrastructure.DeleteAllocatorMetadataItem(
		platform_infrastructure.NewDeleteAllocatorMetadataItemParams().
			WithAllocatorID(params.ID).
			WithKey(params.Key),
		params.AuthWriter,
	)

	if err != nil {
		return api.UnwrapError(err)
	}

	return nil
}

// MetadataGetParams is used to retrieve allocator metadata
type MetadataGetParams struct {
	*api.API
	ID string
}

// Validate ensures that the parameters are correct
func (params MetadataGetParams) Validate() error {
	if params.API == nil {
		return util.ErrAPIReq
	}
	if params.ID == "" {
		return errors.New("allocator metadata: id cannot be empty")
	}
	return nil
}

// GetAllocatorMetadata Retrieves the metadata for a given allocator
func GetAllocatorMetadata(params MetadataGetParams) ([]*models.MetadataItem, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	res, err := params.API.V1API.PlatformInfrastructure.GetAllocatorMetadata(
		platform_infrastructure.NewGetAllocatorMetadataParams().
			WithAllocatorID(params.ID),
		params.AuthWriter,
	)
	if err != nil {
		return nil, api.UnwrapError(err)
	}

	return res.Payload, nil
}
