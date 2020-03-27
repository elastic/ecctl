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

package cmdutil

import (
	"errors"
	"fmt"

	"github.com/elastic/cloud-sdk-go/pkg/plan/planutil"
	"github.com/hashicorp/go-multierror"

	"github.com/elastic/ecctl/pkg/formatter"
)

// TrackParams is consumed by Track.
type TrackParams struct {
	planutil.TrackChangeParams

	// Formatter used to print the structure.
	Formatter formatter.Formatter

	// When set to true, it tracks the progress of the resource change with the
	// specified TrackResourcesParams.
	Track bool

	// Template is the template name which the formatter will use.
	Template string

	// Response will be printed using the formatter and template name.
	Response interface{}
}

// Validate ensures the parameters are usable by the consuming function.
func (params TrackParams) Validate() error {
	var merr = new(multierror.Error)

	if params.Formatter == nil {
		merr = multierror.Append(merr, errors.New("track: formatter cannot be nil"))
	}

	return merr.ErrorOrNil()
}

// Track will either print and track the parameter specified Response.
// If the formatter is not specified an error will be returned.
func Track(params TrackParams) error {
	if err := params.Validate(); err != nil {
		return err
	}

	if err := params.Formatter.Format(params.Template, params.Response); err != nil {
		if !params.Track {
			return err
		}
		if params.TrackChangeParams.Writer != nil {
			fmt.Fprintln(params.Writer, err)
		}
	}

	if !params.Track {
		return nil
	}

	return planutil.TrackChange(params.TrackChangeParams)
}
