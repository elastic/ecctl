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

const updateLong = `updates a deployment from a file definition, defaulting prune_orphans=false, making the default
update action safe for partial updates, to override this behavior toggle --prune-orphans.
To track the changes toggle the --track flag.

It is possible to use "--generate-payload" as: "ecctl deployment show --generate-payload ..."
to obtain a valid create / update payload from an existing deployment which can be manually edited in
before it is sent as a "--file" flag argument. See "ecctl deployment show --help" for a valid example.

Read more about the deployment definition in https://www.elastic.co/guide/en/cloud-enterprise/current/Deployment_-_CRUD.html`

var updateExample = `
#### Same base deployment as the create example, changing cluster_topology[0].zone_count to 3.
$ cat deployment_example_update.json
{
    "resources": {
        "elasticsearch": [
            {
                "display_name": "my elasticsearch cluster",
                "ref_id": "my-es-cluster",
                "plan": {
                    "deployment_template": {
                        "id": "default"
                    },
                    "elasticsearch": {
                        "version": "6.8.4"
                    },
                    "cluster_topology": [
                        {
                            "instance_configuration_id": "data.default",
                            "memory_per_node": 1024,
                            "node_count_per_zone": 1,
                            "node_type": {
                                "data": true,
                                "ingest": true,
                                "master": true,
                                "ml": false
                            },
                            "zone_count": 3
                        }
                    ]
                }
            }
        ]
    }
}
$ ecctl deployment update f44c06c3af6f85dac05023cf243f4ab1 -f deployment_example_update.json
{
  "id": "f44c06c3af6f85dac05023cf243f4ab1",
  "name": "my example cluster",
  "resources": [
    {
      "id": "205745432f6345a4999dd2d77ceb1812",
      "kind": "elasticsearch",
      "ref_id": "my-es-cluster",
      "region": "ece-region"
    },
    {
      "elasticsearch_cluster_ref_id": "my-es-cluster",
      "id": "3617594f01074b76a5ca4f903f9d33ec",
      "kind": "kibana",
      "ref_id": "my-kibana-instance",
      "region": "ece-region"
    },
    {
      "elasticsearch_cluster_ref_id": "my-es-cluster",
      "id": "90c3c9566f454861b0dc935c5c7420d8",
      "kind": "apm",
      "ref_id": "my-apm-instance",
      "region": "ece-region"
    }
  ]
}
#### Setting --prune-orphans, will cause any non-specified resources to be shut down.
$ ecctl deployment update f44c06c3af6f85dac05023cf243f4ab1 -f deployment_example_update.json --prune-orphans
setting --prune-orphans to "true" will cause any resources not specified in the update request to be removed from the deployment, do you want to continue? [y/n]: y
{
  "id": "f44c06c3af6f85dac05023cf243f4ab1",
  "name": "my example cluster",
  "resources": [
    {
      "id": "205745432f6345a4999dd2d77ceb1812",
      "kind": "elasticsearch",
      "ref_id": "my-es-cluster",
      "region": "ece-region"
    }
  ],
  "shutdown_resources": {
    "apm": [
      "90c3c9566f454861b0dc935c5c7420d8"
    ],
    "appsearch": [],
    "elasticsearch": [],
    "kibana": [
      "3617594f01074b76a5ca4f903f9d33ec"
    ]
  }
}`[1:]
