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
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/multierror"
	planmock "github.com/elastic/cloud-sdk-go/pkg/plan/mock"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"

	"github.com/elastic/ecctl/cmd/util/testutils"
)

//go:embed "testdata/template-aws-io-optimized-v2.json"
var awsIoOptimisedTemplate []byte

func Test_createCmd(t *testing.T) {
	var deploymentID = ec.RandomResourceID()
	var esID = ec.RandomResourceID()
	var kibanaID = ec.RandomResourceID()
	var azureCreateResponse = models.DeploymentCreateResponse{
		Created: ec.Bool(true),
		ID:      ec.String(deploymentID),
		Name:    ec.String("search-dev-azure-westeurope"),
		Resources: []*models.DeploymentResource{
			{
				ID:     ec.String(esID),
				Kind:   ec.String("elasticsearch"),
				RefID:  ec.String("main-elasticsearch"),
				Region: ec.String("azure-westeurope"),
				Credentials: &models.ClusterCredentials{
					Username: ec.String("myuser"),
					Password: ec.String("mypass"),
				},
			},
			{
				ID:     ec.String(kibanaID),
				Kind:   ec.String("kibana"),
				RefID:  ec.String("main-kibana"),
				Region: ec.String("azure-westeurope"),
			},
		},
	}
	azureCreateResponseBytes, err := json.MarshalIndent(azureCreateResponse, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	var defaultCreateResponse = models.DeploymentCreateResponse{
		Created: ec.Bool(true),
		ID:      ec.String(deploymentID),
		Name:    ec.String("some-deployment"),
		Resources: []*models.DeploymentResource{
			{
				ID:     ec.String(esID),
				Kind:   ec.String("elasticsearch"),
				RefID:  ec.String("main-elasticsearch"),
				Region: ec.String("ece-region"),
				Credentials: &models.ClusterCredentials{
					Username: ec.String("myuser"),
					Password: ec.String("mypass"),
				},
			},
			{
				ID:     ec.String(kibanaID),
				Kind:   ec.String("kibana"),
				RefID:  ec.String("main-kibana"),
				Region: ec.String("ece-region"),
			},
		},
	}
	defaultCreateResponseBytes, err := json.MarshalIndent(defaultCreateResponse, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	defaultCreateResponseBytes = append(defaultCreateResponseBytes, []byte("\n")...)

	var awsDeploymentID = ec.RandomResourceID()
	var awsESID = ec.RandomResourceID()
	var awsKibanaID = ec.RandomResourceID()
	var awsAPMID = ec.RandomResourceID()
	var awsCreateResponse = models.DeploymentCreateResponse{
		Created: ec.Bool(true),
		ID:      ec.String(awsDeploymentID),
		Name:    ec.String("test-create-aws"),
		Resources: []*models.DeploymentResource{
			{
				ID:     ec.String(awsESID),
				Kind:   ec.String("elasticsearch"),
				RefID:  ec.String("main-elasticsearch"),
				Region: ec.String("us-east-1"),
				Credentials: &models.ClusterCredentials{
					Username: ec.String("myuser"),
					Password: ec.String("mypass"),
				},
			},
			{
				ID:     ec.String(awsKibanaID),
				Kind:   ec.String("kibana"),
				RefID:  ec.String("main-kibana"),
				Region: ec.String("us-east-1"),
			},
			{
				ID:          ec.String(awsAPMID),
				Kind:        ec.String("apm"),
				RefID:       ec.String("main-apm"),
				Region:      ec.String("us-east-1"),
				SecretToken: "some-secret-token",
			},
		},
	}
	awsCreateResponseBytes, err := json.MarshalIndent(awsCreateResponse, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	var awsTrackOutput = fmt.Sprintf(`Deployment [%s] - [Elasticsearch][%s]: running step "waiting-for-some-step" (Plan duration )...
Deployment [%s] - [Kibana][%s]: running step "waiting-for-some-step" (Plan duration )...
Deployment [%s] - [Apm][%s]: running step "waiting-for-some-step" (Plan duration )...`+"\n"+
		"\x1b[92;mDeployment [%s] - [Elasticsearch][%s]: finished running all the plan steps\x1b[0m (Total plan duration )\n"+
		"\x1b[92;mDeployment [%s] - [Kibana][%s]: finished running all the plan steps\x1b[0m (Total plan duration )\n"+
		"\x1b[92;mDeployment [%s] - [Apm][%s]: finished running all the plan steps\x1b[0m (Total plan duration )\n",
		awsDeploymentID, awsESID, awsDeploymentID, awsKibanaID, awsDeploymentID, awsAPMID,
		awsDeploymentID, awsESID, awsDeploymentID, awsKibanaID, awsDeploymentID, awsAPMID,
	)

	var awsCreateResponses = []mock.Response{
		{
			Response: http.Response{
				StatusCode: 201,
				Body:       mock.NewStructBody(awsCreateResponse),
			},
			Assert: &mock.RequestAssertion{
				Method: "POST",
				Header: api.DefaultWriteMockHeaders,
				Body:   mock.NewStringBody(`{"name":"test-create-aws","resources":{"apm":[{"elasticsearch_cluster_ref_id":"main-elasticsearch","plan":{"apm":{"version":"7.8.0"},"cluster_topology":[{"instance_configuration_id":"aws.apm.r4","size":{"resource":"memory","value":512},"zone_count":1}]},"ref_id":"main-apm","region":"us-east-1"}],"appsearch":null,"elasticsearch":[{"plan":{"cluster_topology":[{"instance_configuration_id":"aws.data.highio.i3","node_roles":null,"node_type":{"data":true,"ingest":true,"master":true},"size":{"resource":"memory","value":1024},"zone_count":2}],"deployment_template":{"id":"aws-io-optimized"},"elasticsearch":{"version":"7.8.0"}},"ref_id":"main-elasticsearch","region":"us-east-1","settings":{"dedicated_masters_threshold":6}}],"enterprise_search":null,"integrations_server":null,"kibana":[{"elasticsearch_cluster_ref_id":"main-elasticsearch","plan":{"cluster_topology":[{"instance_configuration_id":"aws.kibana.r4","size":{"resource":"memory","value":1024},"zone_count":1}],"kibana":{"version":"7.8.0"}},"ref_id":"main-kibana","region":"us-east-1"}]}}` + "\n"),
				Path:   "/api/v1/deployments",
				Host:   api.DefaultMockHost,
				Query: url.Values{
					"request_id": {"some_other_request_id"},
				},
			},
		},
		mock.New200StructResponse(planmock.Generate(planmock.GenerateConfig{
			ID: awsDeploymentID,
			Elasticsearch: []planmock.GeneratedResourceConfig{
				{
					ID: awsESID,
					PendingLog: planmock.NewPlanStepLog(
						planmock.NewPlanStep("waiting-for-some-step", "pending"),
					),
				},
			},
			Kibana: []planmock.GeneratedResourceConfig{
				{
					ID: awsKibanaID,
					PendingLog: planmock.NewPlanStepLog(
						planmock.NewPlanStep("waiting-for-some-step", "pending"),
					),
				},
			},
			Apm: []planmock.GeneratedResourceConfig{
				{
					ID: awsAPMID,
					PendingLog: planmock.NewPlanStepLog(
						planmock.NewPlanStep("waiting-for-some-step", "pending"),
					),
				},
			},
		})),
		mock.New500Response(mock.NewStringBody("some error")),
		mock.New200StructResponse(planmock.Generate(planmock.GenerateConfig{
			ID: awsDeploymentID,
			Elasticsearch: []planmock.GeneratedResourceConfig{
				{
					ID: awsESID,
					CurrentLog: planmock.NewPlanStepLog(
						planmock.NewPlanStep("plan-completed", "success"),
					),
				},
			},
			Kibana: []planmock.GeneratedResourceConfig{
				{
					ID: awsKibanaID,
					CurrentLog: planmock.NewPlanStepLog(
						planmock.NewPlanStep("plan-completed", "success"),
					),
				},
			},
			Apm: []planmock.GeneratedResourceConfig{
				{
					ID: awsAPMID,
					CurrentLog: planmock.NewPlanStepLog(
						planmock.NewPlanStep("plan-completed", "success"),
					),
				},
			},
		})),
		mock.New200StructResponse(planmock.Generate(planmock.GenerateConfig{
			ID: awsDeploymentID,
			Elasticsearch: []planmock.GeneratedResourceConfig{
				{
					ID: awsESID,
					CurrentLog: planmock.NewPlanStepLog(
						planmock.NewPlanStep("plan-completed", "success"),
					),
				},
			},
			Kibana: []planmock.GeneratedResourceConfig{
				{
					ID: awsKibanaID,
					CurrentLog: planmock.NewPlanStepLog(
						planmock.NewPlanStep("plan-completed", "success"),
					),
				},
			},
			Apm: []planmock.GeneratedResourceConfig{
				{
					ID: awsAPMID,
					CurrentLog: planmock.NewPlanStepLog(
						planmock.NewPlanStep("plan-completed", "success"),
					),
				},
			},
		})),
	}

	tests := []struct {
		name string
		args testutils.Args
		want testutils.Assertion
	}{
		{
			name: "unexisting file returns invalid argument",
			args: testutils.Args{
				Cmd: createCmd,
				Args: []string{
					"create", "--file=unexisting.json",
				},
			},
			want: testutils.Assertion{
				Err: multierror.NewPrefixed("failed reading the file definition",
					errors.New("invalid argument"),
					errors.New("could not read the specified file, please make sure it exists"),
				).Error(),
			},
		},
		{
			name: "succeeds creating a deployment with default values and without tracking (aws-io-optimized-v2 DT)",
			args: testutils.Args{
				Cmd: createCmd,
				Args: []string{
					"create",
					"--request-id=some_request_id",
					"--version=7.8.0",
					"--deployment-template=aws-io-optimized-v2",
					"--name=with_default",
				},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					{
						Response: http.Response{
							StatusCode: 200,
							Body:       mock.NewByteBody(awsIoOptimisedTemplate),
						},
						Assert: &mock.RequestAssertion{
							Method: "GET",
							Header: api.DefaultReadMockHeaders,
							Path:   "/api/v1/deployments/templates/aws-io-optimized-v2",
							Host:   api.DefaultMockHost,
							Query: url.Values{
								"region":                       {"ece-region"},
								"show_instance_configurations": {"true"},
								"stack_version":                {"7.8.0"},
							},
						},
					},
					{
						Response: http.Response{
							StatusCode: 201,
							Body:       mock.NewStructBody(defaultCreateResponse),
						},
						Assert: &mock.RequestAssertion{
							Method: "POST",
							Header: api.DefaultWriteMockHeaders,
							Body:   mock.NewStringBody(`{"name":"with_default","resources":{"apm":[{"elasticsearch_cluster_ref_id":"main-elasticsearch","plan":{"apm":{"version":"7.8.0"},"cluster_topology":[{"instance_configuration_id":"aws.apm.r5d","size":{"resource":"memory","value":512},"zone_count":1}]},"ref_id":"main-apm","region":"us-east-1"}],"appsearch":null,"elasticsearch":[{"plan":{"autoscaling_enabled":false,"cluster_topology":[{"id":"coordinating","instance_configuration_id":"aws.coordinating.m5d","node_roles":null,"node_type":{"data":false,"ingest":true,"master":false},"size":{"resource":"memory","value":0},"topology_element_control":{"min":{"resource":"memory","value":0}},"zone_count":2},{"autoscaling_max":{"resource":"memory","value":118784},"elasticsearch":{"node_attributes":{"data":"hot"}},"id":"hot_content","instance_configuration_id":"aws.data.highio.i3","node_roles":null,"node_type":{"data":true,"ingest":true,"master":true},"size":{"resource":"memory","value":8192},"topology_element_control":{"min":{"resource":"memory","value":1024}},"zone_count":2},{"autoscaling_max":{"resource":"memory","value":118784},"elasticsearch":{"node_attributes":{"data":"warm"}},"id":"warm","instance_configuration_id":"aws.data.highstorage.d3","node_roles":null,"node_type":{"data":true,"ingest":false,"master":false},"size":{"resource":"memory","value":0},"topology_element_control":{"min":{"resource":"memory","value":0}},"zone_count":2},{"autoscaling_max":{"resource":"memory","value":59392},"elasticsearch":{"node_attributes":{"data":"cold"}},"id":"cold","instance_configuration_id":"aws.data.highstorage.d3","node_roles":null,"node_type":{"data":true,"ingest":false,"master":false},"size":{"resource":"memory","value":0},"topology_element_control":{"min":{"resource":"memory","value":0}},"zone_count":1},{"id":"master","instance_configuration_id":"aws.master.r5d","node_roles":null,"node_type":{"data":false,"ingest":false,"master":true},"size":{"resource":"memory","value":0},"topology_element_control":{"min":{"resource":"memory","value":0}},"zone_count":3},{"autoscaling_max":{"resource":"memory","value":61440},"autoscaling_min":{"resource":"memory","value":0},"id":"ml","instance_configuration_id":"aws.ml.m5d","node_roles":null,"node_type":{"data":false,"ingest":false,"master":false,"ml":true},"size":{"resource":"memory","value":0},"topology_element_control":{"min":{"resource":"memory","value":0}},"zone_count":1}],"deployment_template":{"id":"aws-io-optimized-v2"},"elasticsearch":{"version":"7.8.0"}},"ref_id":"main-elasticsearch","region":"us-east-1","settings":{"dedicated_masters_threshold":6}}],"enterprise_search":[{"elasticsearch_cluster_ref_id":"main-elasticsearch","plan":{"cluster_topology":[{"instance_configuration_id":"aws.enterprisesearch.m5d","node_type":{"appserver":true,"connector":true,"worker":true},"size":{"resource":"memory","value":0},"zone_count":2}],"enterprise_search":{"version":"7.8.0"}},"ref_id":"main-enterprise_search","region":"us-east-1"}],"integrations_server":null,"kibana":[{"elasticsearch_cluster_ref_id":"main-elasticsearch","plan":{"cluster_topology":[{"instance_configuration_id":"aws.kibana.r5d","size":{"resource":"memory","value":1024},"zone_count":1}],"kibana":{"version":"7.8.0"}},"ref_id":"main-kibana","region":"us-east-1"}]}}` + "\n"),
							Path:   "/api/v1/deployments",
							Host:   api.DefaultMockHost,
							Query: url.Values{
								"request_id": {"some_request_id"},
							},
						},
					},
				}},
			},
			want: testutils.Assertion{
				Stdout: string(defaultCreateResponseBytes),
			},
		},
		{
			name: "succeeds creating a deployment with default values and without tracking (aws-io-optimized-v2 DT) with node_roles",
			args: testutils.Args{
				Cmd: createCmd,
				Args: []string{
					"create",
					"--request-id=some_request_id",
					"--version=7.11.2",
					"--deployment-template=aws-io-optimized-v2",
					"--name=with_default",
				},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					{
						Response: http.Response{
							StatusCode: 200,
							Body:       mock.NewByteBody(awsIoOptimisedTemplate),
						},
						Assert: &mock.RequestAssertion{
							Method: "GET",
							Header: api.DefaultReadMockHeaders,
							Path:   "/api/v1/deployments/templates/aws-io-optimized-v2",
							Host:   api.DefaultMockHost,
							Query: url.Values{
								"region":                       {"ece-region"},
								"show_instance_configurations": {"true"},
								"stack_version":                {"7.11.2"},
							},
						},
					},
					{
						Response: http.Response{
							StatusCode: 201,
							Body:       mock.NewStructBody(defaultCreateResponse),
						},
						Assert: &mock.RequestAssertion{
							Method: "POST",
							Header: api.DefaultWriteMockHeaders,
							Body:   mock.NewStringBody(`{"name":"with_default","resources":{"apm":[{"elasticsearch_cluster_ref_id":"main-elasticsearch","plan":{"apm":{"version":"7.11.2"},"cluster_topology":[{"instance_configuration_id":"aws.apm.r5d","size":{"resource":"memory","value":512},"zone_count":1}]},"ref_id":"main-apm","region":"us-east-1"}],"appsearch":null,"elasticsearch":[{"plan":{"autoscaling_enabled":false,"cluster_topology":[{"id":"coordinating","instance_configuration_id":"aws.coordinating.m5d","node_roles":["ingest","remote_cluster_client"],"size":{"resource":"memory","value":0},"topology_element_control":{"min":{"resource":"memory","value":0}},"zone_count":2},{"autoscaling_max":{"resource":"memory","value":118784},"elasticsearch":{"node_attributes":{"data":"hot"}},"id":"hot_content","instance_configuration_id":"aws.data.highio.i3","node_roles":["master","ingest","remote_cluster_client","data_hot","transform","data_content"],"size":{"resource":"memory","value":8192},"topology_element_control":{"min":{"resource":"memory","value":1024}},"zone_count":2},{"autoscaling_max":{"resource":"memory","value":118784},"elasticsearch":{"node_attributes":{"data":"warm"}},"id":"warm","instance_configuration_id":"aws.data.highstorage.d3","node_roles":["data_warm","remote_cluster_client"],"size":{"resource":"memory","value":0},"topology_element_control":{"min":{"resource":"memory","value":0}},"zone_count":2},{"autoscaling_max":{"resource":"memory","value":59392},"elasticsearch":{"node_attributes":{"data":"cold"}},"id":"cold","instance_configuration_id":"aws.data.highstorage.d3","node_roles":["data_cold","remote_cluster_client"],"size":{"resource":"memory","value":0},"topology_element_control":{"min":{"resource":"memory","value":0}},"zone_count":1},{"id":"master","instance_configuration_id":"aws.master.r5d","node_roles":["master"],"size":{"resource":"memory","value":0},"topology_element_control":{"min":{"resource":"memory","value":0}},"zone_count":3},{"autoscaling_max":{"resource":"memory","value":61440},"autoscaling_min":{"resource":"memory","value":0},"id":"ml","instance_configuration_id":"aws.ml.m5d","node_roles":["ml","remote_cluster_client"],"size":{"resource":"memory","value":0},"topology_element_control":{"min":{"resource":"memory","value":0}},"zone_count":1}],"deployment_template":{"id":"aws-io-optimized-v2"},"elasticsearch":{"version":"7.11.2"}},"ref_id":"main-elasticsearch","region":"us-east-1","settings":{"dedicated_masters_threshold":6}}],"enterprise_search":[{"elasticsearch_cluster_ref_id":"main-elasticsearch","plan":{"cluster_topology":[{"instance_configuration_id":"aws.enterprisesearch.m5d","node_type":{"appserver":true,"connector":true,"worker":true},"size":{"resource":"memory","value":0},"zone_count":2}],"enterprise_search":{"version":"7.11.2"}},"ref_id":"main-enterprise_search","region":"us-east-1"}],"integrations_server":null,"kibana":[{"elasticsearch_cluster_ref_id":"main-elasticsearch","plan":{"cluster_topology":[{"instance_configuration_id":"aws.kibana.r5d","size":{"resource":"memory","value":1024},"zone_count":1}],"kibana":{"version":"7.11.2"}},"ref_id":"main-kibana","region":"us-east-1"}]}}` + "\n"),
							Path:   "/api/v1/deployments",
							Host:   api.DefaultMockHost,
							Query: url.Values{
								"request_id": {"some_request_id"},
							},
						},
					},
				}},
			},
			want: testutils.Assertion{
				Stdout: string(defaultCreateResponseBytes),
			},
		},
		{
			name: "succeeds creating an Azure deployment with payload and without tracking",
			args: testutils.Args{
				Cmd: createCmd,
				Args: []string{
					"create", "--file=testdata/create-azure.json",
					"--request-id=some_request_id",
				},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					{
						Response: http.Response{
							StatusCode: 201,
							Body:       mock.NewStructBody(azureCreateResponse),
						},
						Assert: &mock.RequestAssertion{
							Method: "POST",
							Header: api.DefaultWriteMockHeaders,
							Body:   mock.NewStringBody(`{"name":"search-dev-azure-westeurope","resources":{"apm":null,"appsearch":null,"elasticsearch":[{"plan":{"cluster_topology":[{"elasticsearch":{},"instance_configuration_id":"azure.data.highio.l32sv2","node_roles":null,"node_type":{"data":true,"ingest":true,"master":true},"size":{"resource":"memory","value":1024},"zone_count":2}],"deployment_template":{"id":"azure-io-optimized"},"elasticsearch":{"version":"7.8.0"}},"ref_id":"main-elasticsearch","region":"azure-westeurope","settings":{"dedicated_masters_threshold":6}}],"enterprise_search":null,"integrations_server":null,"kibana":[{"elasticsearch_cluster_ref_id":"main-elasticsearch","plan":{"cluster_topology":[{"instance_configuration_id":"azure.kibana.e32sv3","size":{"resource":"memory","value":1024},"zone_count":1}],"kibana":{"version":"7.8.0"}},"ref_id":"main-kibana","region":"azure-westeurope"}]}}` + "\n"),
							Path:   "/api/v1/deployments",
							Host:   api.DefaultMockHost,
							Query: url.Values{
								"request_id": {"some_request_id"},
							},
						},
					},
				}},
			},
			want: testutils.Assertion{
				Stdout: string(azureCreateResponseBytes) + "\n",
			},
		},
		{
			name: "succeeds creating an AWS deployment with payload with tracking",
			args: testutils.Args{
				Cmd: createCmd,
				Args: []string{
					"create", "--file=testdata/create-aws.json", "--track",
					"--request-id=some_other_request_id",
				},
				Cfg: testutils.MockCfg{Responses: awsCreateResponses, OutputFormat: "text"},
			},
			want: testutils.Assertion{
				Stdout: string(awsCreateResponseBytes) + "\n" + awsTrackOutput,
			},
		},
		{
			name: "existing file tries to create deployment with payload and fails without tracking",
			args: testutils.Args{
				Cmd: createCmd,
				Args: []string{
					"create", "--file=testdata/create-azure.json",
					"--request-id=some_request_id",
				},
				Cfg: testutils.MockCfg{
					Responses: []mock.Response{
						{
							Response: http.Response{
								StatusCode: 404,
								Body:       mock.NewStringBody(`{"error": "some"}`),
							},
							Assert: &mock.RequestAssertion{
								Method: "POST",
								Header: api.DefaultWriteMockHeaders,
								Body:   mock.NewStringBody(`{"name":"search-dev-azure-westeurope","resources":{"apm":null,"appsearch":null,"elasticsearch":[{"plan":{"cluster_topology":[{"elasticsearch":{},"instance_configuration_id":"azure.data.highio.l32sv2","node_roles":null,"node_type":{"data":true,"ingest":true,"master":true},"size":{"resource":"memory","value":1024},"zone_count":2}],"deployment_template":{"id":"azure-io-optimized"},"elasticsearch":{"version":"7.8.0"}},"ref_id":"main-elasticsearch","region":"azure-westeurope","settings":{"dedicated_masters_threshold":6}}],"enterprise_search":null,"integrations_server":null,"kibana":[{"elasticsearch_cluster_ref_id":"main-elasticsearch","plan":{"cluster_topology":[{"instance_configuration_id":"azure.kibana.e32sv3","size":{"resource":"memory","value":1024},"zone_count":1}],"kibana":{"version":"7.8.0"}},"ref_id":"main-kibana","region":"azure-westeurope"}]}}` + "\n"),
								Path:   "/api/v1/deployments",
								Host:   api.DefaultMockHost,
								Query: url.Values{
									"request_id": {"some_request_id"},
								},
							},
						},
					},
				},
			},
			want: testutils.Assertion{
				Err: `{"error": "some"}`,
				Stderr: "The deployment creation returned with an error. Use the displayed request ID to recreate the deployment resources" +
					"\n" + "Request ID: some_request_id" + "\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutils.RunCmdAssertion(t, tt.args, tt.want)
			tt.args.Cmd.ResetFlags()
			defer initFlags()
		})
	}
}
