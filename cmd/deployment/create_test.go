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
	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/stretchr/testify/assert"

	"github.com/elastic/ecctl/cmd/util/testutils"
)

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
				Stderr: `http transport warning: failed converting *mock.RoundTripper to *http.Transport` + "\n",
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
				Stderr: `http transport warning: failed converting *mock.RoundTripper to *http.Transport` + "\n",
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
				Stderr: `http transport warning: failed converting *mock.RoundTripper to *http.Transport` + "\n",
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
				Stderr: `http transport warning: failed converting *mock.RoundTripper to *http.Transport` + "\n" +
					"The deployment creation returned with an error. Use the displayed request ID to recreate the deployment resources" +
					"\n" + "Request ID: some_request_id" + "\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testutils.RunCmdAssertion(t, tt.args, tt.want)
		})
	}
}

func Test_returnErrOnHidden(t *testing.T) {
	type args struct {
		err    error
		hidden bool
	}
	tests := []struct {
		name string
		args args
		err  error
	}{
		{
			name: "an error is returned",
			args: args{
				err: errors.New("some"),
			},
			err: errors.New("some"),
		},
		{
			name: "an error is returned when hidden is false and == sdkcmdutil.ErrNodefinitionLoaded",
			args: args{
				err: sdkcmdutil.ErrNodefinitionLoaded,
			},
		},
		{
			name: "an error is returned when hidden is true and == sdkcmdutil.ErrNodefinitionLoaded",
			args: args{
				err:    sdkcmdutil.ErrNodefinitionLoaded,
				hidden: true,
			},
			err: sdkcmdutil.ErrNodefinitionLoaded,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := returnErrOnHidden(tt.args.err, tt.args.hidden)
			if !assert.Equal(t, tt.err, err) {
				t.Error(err)
			}
		})
	}
}
