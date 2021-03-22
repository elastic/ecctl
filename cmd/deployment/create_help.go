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
	createLong = `Creates a deployment from a file definition or using the defaults from the deployment template.
Default values are provided by the deployment template, simplifying the getting started experience.
When version is not specified, the latest available stack version will automatically be used.
These are the available options:

  * Using the deployment template defaults with --deployment-template=<deployment template id>
  * Using the deployment template defaults with --deployment-template=<deployment template id> and 
    --generate-payload to use the deployment template, modify certain fields of the definition and
    use the --file flag to create the deployment.
  * File definition: --file=<file path> (shorthand: -f). You can create a definition by using the sample JSON seen here:
    https://elastic.co/guide/en/cloud/current/ec-api-deployment-crud.html#ec_create_a_deployment

As an option "--generate-payload" can be used in order to obtain the generated payload that would be sent as a request. 
Save it, update or extend the topology and create a deployment using the saved payload with the "--file" flag.`

	// nolint
	createExample = `## Create a deployment with the default values
$ ecctl deployment create --name my-deployment --deployment-template=aws-io-optimized-v2 --region=us-east-1

## To make the command wait until the resources have been created use the "--track" flag, which will output 
the current stage on which the deployment resources are in.
$ ecctl deployment create --name my-deployment --deployment-template=aws-io-optimized-v2 --region=us-east-1 --track
[...]
Deployment [b6ecbea3d5c84124b7dca457f2892086] - [Elasticsearch][b6ecbea3d5c84124b7dca457f2892086]: finished running all the plan steps (Total plan duration: 5m11.s)
Deployment [91c4d60acb804ba0a27651fac02780ec] - [Kibana][8a9d9916cd6e46a7bb0912211d76e2af]: finished running all the plan steps (Total plan duration: 4m29.58s)

## In order to use the "--deployment-template" flag, you'll need to know which deployment templates ara available to you.
You'll need to run the following command to view your deployment templates:
$ ecctl platform deployment-template list

## Use the "--generate-payload" flag to save the definition to a file for later use.
$ ecctl deployment create --name my-deployment --version=7.11.2 --generate-payload > create_example.json

## Create a deployment through the file definition.
$ ecctl deployment create --file create_example.json --track

## To retry a deployment when the previous deployment creation failed, use the request ID provided in the error response of the previous command:
$ ecctl deployment create --request-id=GMZPMRrcMYqHdmxjIQkHbdjnhPIeBElcwrHwzVlhGUSMXrEIzVXoBykSVRsKncNb`
)
