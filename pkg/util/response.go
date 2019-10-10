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

import (
	"github.com/elastic/cloud-sdk-go/pkg/client"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
)

// SetRawJSON prepares the API transport to send raw JSON
//
// While setting the global JSONMime producer to runtime.TextProducer() looks
// scary, runtime.TextProducer() falls back on swag.WriteJSON(data) if the
// object we're marshalling is a struct or slice. So it actually works in most
// cases.
//
// This is only necessary because the Cloud Swagger Spec lists some JSON blob
// data as strings. When those are changed to JSON Objects, this can go away.
func SetRawJSON(c *client.Rest) {
	// Don't configure transport when mocking
	if t, ok := c.Transport.(*httptransport.Runtime); ok {
		t.Producers[runtime.JSONMime] = runtime.TextProducer()
		c.SetTransport(t)
	}
}

// UnsetRawJSON un-does what SetRawJSON does :)
func UnsetRawJSON(c *client.Rest) {
	// Don't configure transport when mocking
	if t, ok := c.Transport.(*httptransport.Runtime); ok {
		t.Producers[runtime.JSONMime] = runtime.JSONProducer()
		c.SetTransport(t)
	}
}
