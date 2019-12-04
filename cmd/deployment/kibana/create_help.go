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

package cmdkibana

const (
	kibanaCreateLong = `Creates a Kibana deployment, limitting the creation scope to Kibana resources.
There's a few ways to create an Kibana deployment, sane default values are provided, making
the command work out of the box even when no parameters are set. When version is not specified,
the matching elasticsearch deployment version will be used. These are the available options:

  * Simplified flags: --zones <zone count> --size <node memory in MB>
  * File definition: --file=<file path> (shorthand: -f). The definition can be found in:
    https://www.elastic.co/guide/en/cloud-enterprise/current/definitions.html#KibanaPayload

As an option "--generate-payload" can be used in order to obtain the generated KibanaPayload
that would be sent as a request, save it, update or extend the topology and create an Kibana
deployment using the saved payload with the "--file" flag.`

	kibanaCreateExample = `## Create a single node cluster. The command will exit after the API response has been returned, 
## without waiting until the deployment resources have been created. To make the command wait until
the resources have been created use the "--track" flag.
$ ecctl deployment kibana create --id=fc2fe19a8c4f43d385004febe4e63900 --track
{
  "id": "fc2fe19a8c4f43d385004febe4e63900",
  "name": "fc2fe19a8c4f43d385004febe4e63900",
  "resources": [
    {
      "elasticsearch_cluster_ref_id": "elasticsearch",
      "id": "8d37cff942fb482ba7a3cac10eea943a",
      "kind": "kibana",
      "ref_id": "kibana",
      "region": "ece-region"
    },
    {
      "cloud_id": "fc2fe19a8c4f43d385004febe4e63900:MTkyLjE2OC40NC4xMC5pcC5lcy5pbzo5MjQzJDJjOGRkMWQxYmExNDQyZGU4NDU2NmY0ZTkxY2M1YWJmJDhkMzdjZmY5NDJmYjQ4MmJhN2EzY2FjMTBlZWE5NDNh",
      "id": "2c8dd1d1ba1442de84566f4e91cc5abf",
      "kind": "elasticsearch",
      "ref_id": "elasticsearch",
      "region": "ece-region"
    }
  ]
}
Cluster [8d37cff942fb482ba7a3cac10eea943a][Kibana]: running step "wait-until-running" (Plan duration 1.480086699s)...
Cluster [2c8dd1d1ba1442de84566f4e91cc5abf][Elasticsearch]: finished running all the plan steps (Total plan duration: 1.598400189s)
Cluster [8d37cff942fb482ba7a3cac10eea943a][Kibana]: running step "set-maintenance" (Plan duration 1m2.277750264s)...
Cluster [8d37cff942fb482ba7a3cac10eea943a][Kibana]: finished running all the plan steps (Total plan duration: 1m7.544473245s)

## Save the definition to a file for later use.
$ ecctl deployment kibana create --generate-payload --id fc2fe19a8c4f43d385004febe4e63900 --zones 2 --size 2048 > kibana_create_example.json

## Create the deployment piping through the file contents tracking the creation progress
$ cat kibana_create_example.json | dev-cli deployment kibana create --track --id fc2fe19a8c4f43d385004febe4e63900
[...]`
)
