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

var defaultTemplateResponse = models.DeploymentTemplateInfoV2{
	ID: ec.String("default"),
	DeploymentTemplate: &models.DeploymentCreateRequest{
		Resources: &models.DeploymentCreateResources{
			Apm: []*models.ApmPayload{
				{
					Plan: &models.ApmPlan{
						ClusterTopology: []*models.ApmTopologyElement{
							{
								Size: &models.TopologySize{
									Resource: ec.String("memory"),
									Value:    ec.Int32(1024),
								},
								ZoneCount: 1,
							},
						},
					},
				},
			},
			EnterpriseSearch: []*models.EnterpriseSearchPayload{
				{
					Plan: &models.EnterpriseSearchPlan{
						ClusterTopology: []*models.EnterpriseSearchTopologyElement{
							{
								Size: &models.TopologySize{
									Resource: ec.String("memory"),
									Value:    ec.Int32(1024),
								},
								ZoneCount: 1,
							},
						},
					},
				},
			},
			Appsearch: []*models.AppSearchPayload{
				{
					Plan: &models.AppSearchPlan{
						ClusterTopology: []*models.AppSearchTopologyElement{
							{
								Size: &models.TopologySize{
									Resource: ec.String("memory"),
									Value:    ec.Int32(1024),
								},
								ZoneCount: 1,
							},
						},
					},
				},
			},
			Kibana: []*models.KibanaPayload{
				{
					Plan: &models.KibanaClusterPlan{
						ClusterTopology: []*models.KibanaClusterTopologyElement{
							{
								Size: &models.TopologySize{
									Resource: ec.String("memory"),
									Value:    ec.Int32(1024),
								},
								ZoneCount: 1,
							},
						},
					},
				},
			},
			Elasticsearch: []*models.ElasticsearchPayload{
				{
					Plan: &models.ElasticsearchClusterPlan{
						ClusterTopology: defaultESTopologies,
					},
				},
			},
		},
	},
}

var defaultESTopologies = []*models.ElasticsearchClusterTopologyElement{
	{
		InstanceConfigurationID: "default.data",
		Size: &models.TopologySize{
			Resource: ec.String("memory"),
			Value:    ec.Int32(1024),
		},
		NodeType: &models.ElasticsearchNodeType{
			Data: ec.Bool(true),
		},
	},
	{
		InstanceConfigurationID: "default.master",
		Size: &models.TopologySize{
			Resource: ec.String("memory"),
			Value:    ec.Int32(1024),
		},
		NodeType: &models.ElasticsearchNodeType{
			Master: ec.Bool(true),
		},
	},
	{
		InstanceConfigurationID: "default.ml",
		Size: &models.TopologySize{
			Resource: ec.String("memory"),
			Value:    ec.Int32(1024),
		},
		NodeType: &models.ElasticsearchNodeType{
			Ml: ec.Bool(true),
		},
	},
}

func Test_createCmd(t *testing.T) {
	var deploymentID = ec.RandomResourceID()
	var esID = ec.RandomResourceID()
	var kibanaID = ec.RandomResourceID()
	var apmID = ec.RandomResourceID()
	var appsearchID = ec.RandomResourceID()
	var enterprisesearchID = ec.RandomResourceID()
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

	var overrideCreateResponse = models.DeploymentCreateResponse{
		Created: ec.Bool(true),
		ID:      ec.String(deploymentID),
		Name:    ec.String("some-deployment-with-overrides"),
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
			{
				ID:     ec.String(apmID),
				Kind:   ec.String("apm"),
				RefID:  ec.String("main-apm"),
				Region: ec.String("ece-region"),
			},
			{
				ID:     ec.String(appsearchID),
				Kind:   ec.String("appsearch"),
				RefID:  ec.String("main-appsearch"),
				Region: ec.String("ece-region"),
			},
			{
				ID:     ec.String(enterprisesearchID),
				Kind:   ec.String("enterprise_search"),
				RefID:  ec.String("main-enterprise_search"),
				Region: ec.String("ece-region"),
			},
		},
	}
	overrideCreateResponseBytes, err := json.MarshalIndent(overrideCreateResponse, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	overrideCreateResponseBytes = append(overrideCreateResponseBytes, []byte("\n")...)

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
				Body:   mock.NewStringBody(`{"name":"test-create-aws","resources":{"apm":[{"elasticsearch_cluster_ref_id":"main-elasticsearch","plan":{"apm":{"version":"7.8.0"},"cluster_topology":[{"instance_configuration_id":"aws.apm.r4","size":{"resource":"memory","value":512},"zone_count":1}]},"ref_id":"main-apm","region":"us-east-1"}],"appsearch":null,"elasticsearch":[{"plan":{"cluster_topology":[{"instance_configuration_id":"aws.data.highio.i3","node_type":{"data":true,"ingest":true,"master":true},"size":{"resource":"memory","value":1024},"zone_count":2}],"deployment_template":{"id":"aws-io-optimized"},"elasticsearch":{"version":"7.8.0"}},"ref_id":"main-elasticsearch","region":"us-east-1","settings":{"dedicated_masters_threshold":6}}],"enterprise_search":null,"kibana":[{"elasticsearch_cluster_ref_id":"main-elasticsearch","plan":{"cluster_topology":[{"instance_configuration_id":"aws.kibana.r4","size":{"resource":"memory","value":1024},"zone_count":1}],"kibana":{"version":"7.8.0"}},"ref_id":"main-kibana","region":"us-east-1"}]}}` + "\n"),
				Path:   "/api/v1/deployments",
				Host:   api.DefaultMockHost,
				Query: url.Values{
					"request_id":    {"some_other_request_id"},
					"validate_only": {"false"},
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
				),
			},
		},
		{
			name: "succeeds creating a deployment with default values and without tracking",
			args: testutils.Args{
				Cmd: createCmd,
				Args: []string{
					"create", "--request-id=some_request_id",
				},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					{
						Response: http.Response{
							StatusCode: 200,
							Body: mock.NewStructBody([]models.DeploymentTemplateInfoV2{
								defaultTemplateResponse,
							}),
						},
						Assert: &mock.RequestAssertion{
							Method: "GET",
							Header: api.DefaultReadMockHeaders,
							Path:   "/api/v1/deployments/templates",
							Host:   api.DefaultMockHost,
							Query: url.Values{
								"region":                       {"ece-region"},
								"show_hidden":                  {"false"},
								"show_instance_configurations": {"true"},
							},
						},
					},
					{
						Response: http.Response{
							StatusCode: 200,
							Body: mock.NewStructBody(models.StackVersionConfigs{
								Stacks: []*models.StackVersionConfig{{Version: "7.8.0"}},
							}),
						},
						Assert: &mock.RequestAssertion{
							Host:   api.DefaultMockHost,
							Header: api.DefaultReadMockHeaders,
							Path:   "/api/v1/regions/ece-region/stack/versions",
							Method: "GET",
							Query: url.Values{
								"show_deleted":  {"false"},
								"show_unusable": {"false"},
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
							Body:   mock.NewStringBody(`{"resources":{"apm":null,"appsearch":null,"elasticsearch":[{"plan":{"cluster_topology":[{"instance_configuration_id":"default.data","node_type":{"data":true},"size":{"resource":"memory","value":4096},"zone_count":1}],"deployment_template":{"id":"default"},"elasticsearch":{"version":"7.8.0"}},"ref_id":"main-elasticsearch","region":"ece-region"}],"enterprise_search":null,"kibana":[{"elasticsearch_cluster_ref_id":"main-elasticsearch","plan":{"cluster_topology":[{"size":{"resource":"memory","value":1024},"zone_count":1}],"kibana":{"version":"7.8.0"}},"ref_id":"main-kibana","region":"ece-region"}]}}` + "\n"),
							Path:   "/api/v1/deployments",
							Host:   api.DefaultMockHost,
							Query: url.Values{
								"request_id":    {"some_request_id"},
								"validate_only": {"false"},
							},
						},
					},
				}},
			},
			want: testutils.Assertion{
				Stderr: "Obtained latest stack version: 7.8.0\n",
				Stdout: string(defaultCreateResponseBytes),
			},
		},
		{
			name: "succeeds creating a deployment with default values and without tracking (Normal GetCall)",
			args: testutils.Args{
				Cmd: createCmd,
				Args: []string{
					"create", "--request-id=some_request_id", "--version=7.8.0",
					"--dt-as-list=false",
				},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					{
						Response: http.Response{
							StatusCode: 200,
							Body:       mock.NewStructBody(defaultTemplateResponse),
						},
						Assert: &mock.RequestAssertion{
							Method: "GET",
							Header: api.DefaultReadMockHeaders,
							Path:   "/api/v1/deployments/templates/default",
							Host:   api.DefaultMockHost,
							Query: url.Values{
								"region":                       {"ece-region"},
								"show_instance_configurations": {"true"},
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
							Body:   mock.NewStringBody(`{"resources":{"apm":null,"appsearch":null,"elasticsearch":[{"plan":{"cluster_topology":[{"instance_configuration_id":"default.data","node_type":{"data":true},"size":{"resource":"memory","value":4096},"zone_count":1}],"deployment_template":{"id":"default"},"elasticsearch":{"version":"7.8.0"}},"ref_id":"main-elasticsearch","region":"ece-region"}],"enterprise_search":null,"kibana":[{"elasticsearch_cluster_ref_id":"main-elasticsearch","plan":{"cluster_topology":[{"size":{"resource":"memory","value":1024},"zone_count":1}],"kibana":{"version":"7.8.0"}},"ref_id":"main-kibana","region":"ece-region"}]}}` + "\n"),
							Path:   "/api/v1/deployments",
							Host:   api.DefaultMockHost,
							Query: url.Values{
								"request_id":    {"some_request_id"},
								"validate_only": {"false"},
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
			name: "succeeds creating a deployment with overrides and without tracking",
			args: testutils.Args{
				Cmd: createCmd,
				Args: []string{
					"create", "--apm", "--appsearch", "--enterprise_search", "--request-id=some_request_id", "--version=7.8.0",
				},
				Cfg: testutils.MockCfg{Responses: []mock.Response{
					{
						Response: http.Response{
							StatusCode: 200,
							Body: mock.NewStructBody([]models.DeploymentTemplateInfoV2{
								defaultTemplateResponse,
							}),
						},
						Assert: &mock.RequestAssertion{
							Method: "GET",
							Header: api.DefaultReadMockHeaders,
							Path:   "/api/v1/deployments/templates",
							Host:   api.DefaultMockHost,
							Query: url.Values{
								"region":                       {"ece-region"},
								"show_hidden":                  {"false"},
								"show_instance_configurations": {"true"},
							},
						},
					},
					{
						Response: http.Response{
							StatusCode: 201,
							Body:       mock.NewStructBody(overrideCreateResponse),
						},
						Assert: &mock.RequestAssertion{
							Method: "POST",
							Header: api.DefaultWriteMockHeaders,
							Body:   mock.NewStringBody(`{"resources":{"apm":[{"elasticsearch_cluster_ref_id":"main-elasticsearch","plan":{"apm":{"version":"7.8.0"},"cluster_topology":[{"size":{"resource":"memory","value":512},"zone_count":1}]},"ref_id":"main-apm","region":"ece-region"}],"appsearch":[{"elasticsearch_cluster_ref_id":"main-elasticsearch","plan":{"appsearch":{"version":"7.8.0"},"cluster_topology":[{"size":{"resource":"memory","value":2048},"zone_count":1}]},"ref_id":"main-appsearch","region":"ece-region"}],"elasticsearch":[{"plan":{"cluster_topology":[{"instance_configuration_id":"default.data","node_type":{"data":true},"size":{"resource":"memory","value":4096},"zone_count":1}],"deployment_template":{"id":"default"},"elasticsearch":{"version":"7.8.0"}},"ref_id":"main-elasticsearch","region":"ece-region"}],"enterprise_search":[{"elasticsearch_cluster_ref_id":"main-elasticsearch","plan":{"cluster_topology":[{"size":{"resource":"memory","value":4096},"zone_count":1}],"enterprise_search":{"version":"7.8.0"}},"ref_id":"main-enterprise_search","region":"ece-region"}],"kibana":[{"elasticsearch_cluster_ref_id":"main-elasticsearch","plan":{"cluster_topology":[{"size":{"resource":"memory","value":1024},"zone_count":1}],"kibana":{"version":"7.8.0"}},"ref_id":"main-kibana","region":"ece-region"}]}}` + "\n"),
							Path:   "/api/v1/deployments",
							Host:   api.DefaultMockHost,
							Query: url.Values{
								"request_id":    {"some_request_id"},
								"validate_only": {"false"},
							},
						},
					},
				}},
			},
			want: testutils.Assertion{
				Stdout: string(overrideCreateResponseBytes),
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
							Body:   mock.NewStringBody(`{"name":"search-dev-azure-westeurope","resources":{"apm":null,"appsearch":null,"elasticsearch":[{"plan":{"cluster_topology":[{"elasticsearch":{},"instance_configuration_id":"azure.data.highio.l32sv2","node_type":{"data":true,"ingest":true,"master":true},"size":{"resource":"memory","value":1024},"zone_count":2}],"deployment_template":{"id":"azure-io-optimized"},"elasticsearch":{"version":"7.8.0"}},"ref_id":"main-elasticsearch","region":"azure-westeurope","settings":{"dedicated_masters_threshold":6}}],"enterprise_search":null,"kibana":[{"elasticsearch_cluster_ref_id":"main-elasticsearch","plan":{"cluster_topology":[{"instance_configuration_id":"azure.kibana.e32sv3","size":{"resource":"memory","value":1024},"zone_count":1}],"kibana":{"version":"7.8.0"}},"ref_id":"main-kibana","region":"azure-westeurope"}]}}` + "\n"),
							Path:   "/api/v1/deployments",
							Host:   api.DefaultMockHost,
							Query: url.Values{
								"request_id":    {"some_request_id"},
								"validate_only": {"false"},
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
								Body:   mock.NewStringBody(`{"name":"search-dev-azure-westeurope","resources":{"apm":null,"appsearch":null,"elasticsearch":[{"plan":{"cluster_topology":[{"elasticsearch":{},"instance_configuration_id":"azure.data.highio.l32sv2","node_type":{"data":true,"ingest":true,"master":true},"size":{"resource":"memory","value":1024},"zone_count":2}],"deployment_template":{"id":"azure-io-optimized"},"elasticsearch":{"version":"7.8.0"}},"ref_id":"main-elasticsearch","region":"azure-westeurope","settings":{"dedicated_masters_threshold":6}}],"enterprise_search":null,"kibana":[{"elasticsearch_cluster_ref_id":"main-elasticsearch","plan":{"cluster_topology":[{"instance_configuration_id":"azure.kibana.e32sv3","size":{"resource":"memory","value":1024},"zone_count":1}],"kibana":{"version":"7.8.0"}},"ref_id":"main-kibana","region":"azure-westeurope"}]}}` + "\n"),
								Path:   "/api/v1/deployments",
								Host:   api.DefaultMockHost,
								Query: url.Values{
									"request_id":    {"some_request_id"},
									"validate_only": {"false"},
								},
							},
						},
					},
				},
			},
			want: testutils.Assertion{
				Err: errors.New(`{"error": "some"}`),
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
