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

import (
	"fmt"
	"os"

	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/util"
)

const createLong = `Creates a deployment from a file definition with an automatically generated idempotency token.
On creation failure, please use the displayed idempotency token to retry the cluster creation with --request-id=<token>.

Read more about the deployment definition in https://www.elastic.co/guide/en/cloud-enterprise/current/Deployment_-_CRUD.html`

var createExample = `
$ cat deployment_example.json
{
    "name": "my example cluster",
    "resources": {
        "apm": [
            {
                "display_name": "my apm instance",
                "ref_id": "my-apm-instance",
                "elasticsearch_cluster_ref_id": "my-es-cluster",
                "plan": {
                    "apm": {
                        "version": "6.8.4"
                    },
                    "cluster_topology": [{
                        "instance_configuration_id": "apm",
                        "size": {
                            "resource": "memory",
                            "value": 512
                        },
                        "zone_count": 1
                    }]
                }
            }
        ],
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
                            "zone_count": 1
                        }
                    ]
                }
            }
        ],
        "kibana": [
            {
                "display_name": "my kibana instance",
                "ref_id": "my-kibana-instance",
                "elasticsearch_cluster_ref_id": "my-es-cluster",
                "plan": {
                    "kibana": {
                        "version": "6.8.4"
                    },
                    "cluster_topology": [
                        {
                            "instance_configuration_id": "kibana",
                            "memory_per_node": 1024,
                            "node_count_per_zone": 1,
                            "zone_count": 1
                        }
                    ]
                }
            }
        ]
    }
}
$ ecctl deployment create -f deployment_example.json --version=7.4.1
[...]

## If th previous deployment creation failed
$ ecctl deployment create -f deployment_example.json --name adeploy --version=7.4.1
The deployment creation returned with an error, please use the displayed idempotency token
to recreate the deployment resources
Idempotency token: GMZPMRrcMYqHdmxjIQkHbdjnhPIeBElcwrHwzVlhGUSMXrEIzVXoBykSVRsKncNb
unknown error (status 500)
$ ecctl deployment create -f deployment_example.json --name adeploy --version=7.4.1 --request-id=GMZPMRrcMYqHdmxjIQkHbdjnhPIeBElcwrHwzVlhGUSMXrEIzVXoBykSVRsKncNb
[...]`[1:]

var createCmd = &cobra.Command{
	Use:     `create -f <file definition.json>`,
	Short:   "Creates a deployment from a file definition, allowing certain flag overrides",
	Long:    createLong,
	Example: createExample,
	PreRunE: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		filename, _ := cmd.Flags().GetString("file")
		var r models.DeploymentCreateRequest
		if err := cmdutil.DecodeFile(filename, &r); err != nil {
			return err
		}

		var region string
		if ecctl.Get().Config.Region == "" {
			region = cmdutil.DefaultECERegion
		}

		name, _ := cmd.Flags().GetString("name")
		version, _ := cmd.Flags().GetString("version")

		reqID, _ := cmd.Flags().GetString("request-id")
		if reqID == "" {
			reqID = util.RandomString(64)
		}

		res, err := deployment.Create(deployment.CreateParams{
			API:       ecctl.Get().API,
			RequestID: reqID,
			Request:   &r,
			Overrides: &deployment.CreateOverrides{
				Name:    name,
				Region:  region,
				Version: version,
			},
		})

		if err != nil {
			fmt.Fprintln(os.Stderr,
				"The deployment creation returned with an error, please use the displayed idempotency token to recreate the deployment resources",
			)
			fmt.Fprintln(os.Stderr, "Idempotency token:", reqID)
			return err
		}

		return ecctl.Get().Formatter.Format("", res)
	},
}

func init() {
	Command.AddCommand(createCmd)
	createCmd.Flags().String("name", "", "Overrides the deployment name")
	createCmd.Flags().String("version", "", "Overrides all thee deployment's resources to the specified version")
	createCmd.Flags().String("request-id", "", "Optional idempotency token - Can be found in the Stderr device when a previous deployment creation failed, for more information see the examples in the help command page")
	createCmd.Flags().StringP("file", "f", "", "JSON file that contains JSON-style domain-specific deployment definition")
	createCmd.MarkFlagRequired("file")
	createCmd.MarkFlagFilename("file", "*.json")
}
