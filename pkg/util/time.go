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

package util

import "time"

const (
	// DefaultClientTimeout client http timeout for console requests
	DefaultClientTimeout = 15
)

// GetTimeoutFromSize computes a time.Duration from the size, if the computed
// value is smaller than the minimum timeout that is returned instead
func GetTimeoutFromSize(size int64) time.Duration {
	var minTimeout = DefaultClientTimeout * time.Second
	var timeout = time.Duration(size/100*2) * time.Second

	if timeout < minTimeout {
		return minTimeout
	}

	return timeout
}
