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

package cmddeployment

const (
	// nolint
	createLong = `Creates a deployment which can be defined through flags or from a file definition.
Sane default values are provided, making the command work out of the box even when no parameters are set. 
When version is not specified, the latest available stack version will automatically be used. 
These are the available options:

  * Simplified flags to set size and zone count for each instance type (Elasticsearch, Kibana, APM, Enterprise Search and App Search). 
  * Advanced flag for different Elasticsearch node types: --es-node-topology <json obj> (shorthand: -e).
    Note that the flag can be specified multiple times for complex topologies.
    The JSON object has the following format:
    {
      "node_type": "["data", "master", "ml"]" # type string
      "size": "1g" # type string
      "zone_count": 1 # type int32
    }
  * File definition: --file=<file path> (shorthand: -f). You can create a definition by using the sample JSON seen here:
    https://elastic.co/guide/en/cloud/current/ec-api-deployment-crud.html#ec_create_a_deployment

As an option "--generate-payload" can be used in order to obtain the generated payload that would be sent as a request. 
Save it, update or extend the topology and create a deployment using the saved payload with the "--file" flag.`

	// nolint
	createExample = `## Create a deployment with the default values for Elasticsearch, a Kibana instance with a modified size, 
and a default APM instance. While Elasticsearch and Kibana come enabled by default, APM, Enterprise Search and App Search need to be 
enabled by using the "--apm", "--enterprise_search" and "--appsearch" flags. The command will exit after the API response has been returned, without 
waiting until the deployment resources have been created. 
$ ecctl deployment create --name my-deployment --zones 2 --kibana-size 2g --apm --apm-size 0.5g

## To make the command wait until the resources have been created use the "--track" flag, which will output 
the current stage on which the deployment resources are in.
$ deployment create --name my-deployment --track
[...]
Deployment [b6ecbea3d5c84124b7dca457f2892086] - [Elasticsearch][b6ecbea3d5c84124b7dca457f2892086]: finished running all the plan steps (Total plan duration: 5m11.s)
Deployment [91c4d60acb804ba0a27651fac02780ec] - [Kibana][8a9d9916cd6e46a7bb0912211d76e2af]: finished running all the plan steps (Total plan duration: 4m29.58s)

## Additionally, a more advanced topology for Elasticsearch can be created through "--es-node-topology" or shorthand "-e".
The following command will create a deployment with two 1GB Elasticsearch instances of the type "data" and 
one 1GB Elasticsearch instance of the type "ml".
$ ecctl deployment create --name my-deployment --es-node-topology '{"size": "1g", "zone_count": 2, "node_type": "data"}' --es-node-topology '{"size": "1g", "zone_count": 1, "node_type": "ml"}'

## In order to use the "--deployment-template" flag, you'll need to know which deployment templates ara available to you.
You'll need to run the following command to view your deployment templates:
$ ecctl platform deployment-template list

## Use the "--generate-payload" flag to save the definition to a file for later use.
$ ecctl deployment create --name my-deployment --size 1g --track --generate-payload --zones 2 > create_example.json

## Create a deployment through the file definition.
$ ecctl deployment create --file create_example.json --track

## To retry a deployment when the previous deployment creation failed, use the request ID provided in the error response of the previous command:
$ ecctl deployment create --request-id=GMZPMRrcMYqHdmxjIQkHbdjnhPIeBElcwrHwzVlhGUSMXrEIzVXoBykSVRsKncNb`
)
