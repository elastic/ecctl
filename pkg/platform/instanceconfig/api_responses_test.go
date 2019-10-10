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

package instanceconfig

const listInstanceConfigsSuccess = `[{
  "description": "Instance configuration to be used for a higher disk/memory ratio",
  "discrete_sizes": {
    "default_size": 1024,
    "resource": "memory",
    "sizes": [
      1024,
      2048,
      4096,
      8192,
      16384,
      32768,
      65536,
      131072,
      262144
    ]
  },
  "id": "data.highstorage",
  "instance_type": "elasticsearch",
  "name": "data.highstorage",
  "node_types": [
    "data",
    "ingest",
    "master"
  ],
  "storage_multiplier": 32
}, {
  "description": "Instance configuration to be used for Kibana",
  "discrete_sizes": {
    "default_size": 1024,
    "resource": "memory",
    "sizes": [
      1024,
      2048,
      4096,
      8192
    ]
  },
  "id": "kibana",
  "instance_type": "kibana",
  "name": "kibana",
  "node_types": null,
  "storage_multiplier": 4
}]`

const getInstanceConfigsSuccess = `{
  "description": "Instance configuration to be used for a higher disk/memory ratio",
  "discrete_sizes": {
    "default_size": 1024,
    "resource": "memory",
    "sizes": [
      1024,
      2048,
      4096,
      8192,
      16384,
      32768,
      65536,
      131072,
      262144
    ]
  },
  "id": "data.highstorage",
  "instance_type": "elasticsearch",
  "name": "data.highstorage",
  "node_types": [
    "data",
    "ingest",
    "master"
  ],
  "storage_multiplier": 32
}`

const getInstanceConfigsSuccessKibana = `{
  "description": "Instance configuration to be used for Kibana",
  "discrete_sizes": {
    "default_size": 1024,
    "resource": "memory",
    "sizes": [
      1024,
      2048,
      4096,
      8192
    ]
  },
  "id": "kibana",
  "instance_type": "kibana",
  "name": "kibana",
  "node_types": null,
  "storage_multiplier": 4
}`

const newConfigKibanaInstanceConfig = `{
  "description": "Instance configuration to be used for Kibana",
  "discrete_sizes": {
    "default_size": 1024,
    "resource": "memory",
    "sizes": [
      1024,
      2048,
      4096,
      8192
    ]
  },
  "id": "kibana",
  "instance_type": "kibana",
  "name": "kibana",
  "node_types": [],
  "storage_multiplier": 4
}`
