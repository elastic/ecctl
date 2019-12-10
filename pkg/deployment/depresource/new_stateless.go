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

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/util"
)

// NewStateless is consumed by NewKibana.
type NewStateless struct {
	*api.API

	// Required: DeploymentID.
	DeploymentID string

	// Required region name.
	Region string

	// Optional ElasticsearchRefID.
	ElasticsearchRefID string

	// Optional name. If not specified it defaults to the autogeneratd ID.
	Name string

	// Optional RefID for the for the stateless deployment resource.
	RefID string

	// Optional: Version is the stateless resource Version. If not set it'll
	// automatically be set to the Elasticsearch RefID version.
	Version string

	// Stateless deployment resource size.
	Size      int32
	ZoneCount int32

	// Optional Deployment Template ID.
	TemplateID string
}

// Validate ensures the parameters are usable by the consuming function.
func (params *NewStateless) Validate() error {
	var merr = new(multierror.Error)

	if params.API == nil {
		merr = multierror.Append(merr, util.ErrAPIReq)
	}

	if len(params.DeploymentID) != 32 {
		merr = multierror.Append(merr, util.ErrDeploymentID)
	}

	if params.Region == "" {
		merr = multierror.Append(merr, errors.New("deployment topology: region cannot be empty"))
	}

	return merr.ErrorOrNil()
}

func (params *NewStateless) fillDefaults(refID string) {
	if params.RefID == "" {
		params.RefID = refID
	}
}

// When either of TemplateID or ElasticsearchRefID are not set, obtain the
// values from the running deployment. Overriding either or both when not set.
// Receives a pointer to the structure, meaning it modifies its values.
func getTemplateAndRefID(params *NewStateless) error {
	if params.TemplateID == "" || params.ElasticsearchRefID == "" {
		res, err := GetDeploymentInfo(GetDeploymentInfoParams{
			API:          params.API,
			DeploymentID: params.DeploymentID,
		})
		if err != nil {
			return err
		}
		if params.TemplateID == "" {
			params.TemplateID = res.DeploymentTemplate
		}
		if params.ElasticsearchRefID == "" {
			params.ElasticsearchRefID = res.RefID
		}
		return nil
	}
	return nil
}
