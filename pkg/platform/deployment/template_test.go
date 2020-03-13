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

package deployment

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/client/platform_configuration_templates"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/go-openapi/strfmt"

	"github.com/elastic/ecctl/pkg/util"
)

var (
	validTemplateID   = "84e0bd6d69bb44e294809d89cea88a7e"
	invalidTemplateID = ""
)

func TestList(t *testing.T) {
	var templateListSuccess = `
	[{
	  "description": "Test default Elasticsearch trial template",
	  "id": "84e0bd6d69bb44e294809d89cea88a7e",
	  "name": "(Trial) Default Elasticsearch",
	  "system_owned": false
	},
	{
	  "description": "Test default Elasticsearch template",
	  "id": "0efbab9c368849a59fc5622ec750ba47",
	  "name": "Default Elasticsearch",
	  "system_owned": true
	}
  ]`
	tests := []struct {
		name    string
		args    ListTemplateParams
		want    *platform_configuration_templates.GetDeploymentTemplatesOK
		wantErr bool
		error   string
	}{
		{
			name: "Platform deployment templates list succeeds",
			args: ListTemplateParams{
				API: api.NewMock(mock.Response{Response: http.Response{
					Body:       mock.NewStringBody(templateListSuccess),
					StatusCode: 200,
				}}),
			},
			want: &platform_configuration_templates.GetDeploymentTemplatesOK{
				Payload: []*models.DeploymentTemplateInfo{
					{
						ID:          "84e0bd6d69bb44e294809d89cea88a7e",
						Description: "Test default Elasticsearch trial template",
						Name:        ec.String("(Trial) Default Elasticsearch"),
						SystemOwned: ec.Bool(false),
					},
					{
						ID:          "0efbab9c368849a59fc5622ec750ba47",
						Description: "Test default Elasticsearch template",
						Name:        ec.String("Default Elasticsearch"),
						SystemOwned: ec.Bool(true),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Platform deployment templates list fails",
			args: ListTemplateParams{
				API: api.NewMock(mock.Response{Error: errors.New("error")}),
			},
			wantErr: true,
			error:   `Get https://mock-host/mock-path/platform/configuration/templates/deployments?format=cluster&show_hidden=false&show_instance_configurations=false: error`,
		},
		{
			name:    "Platform deployment templates fails with an empty API",
			args:    ListTemplateParams{},
			want:    nil,
			wantErr: true,
			error:   "api reference is required for command",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListTemplates(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr != (tt.error != "") {
				t.Errorf("List() expected error message = '%v', actual error message = '%v'", tt.error, err)
			}

			if tt.wantErr && err.Error() != tt.error {
				t.Errorf("List() actual error = '%v', want error '%v'", err, tt.error)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGet(t *testing.T) {
	var sourceDate = util.ParseDate(t, "2018-04-19T18:16:57.297Z")

	var templateGetSuccess = `
	{
  "name": "(Trial) Default Elasticsearch",
  "source": {
	"user_id": "1",
	"facilitator": "adminconsole",
	"date": "2018-04-19T18:16:57.297Z",
	"admin_id": "admin",
	"action": "deployments.create-template",
	"remote_addresses": ["52.205.1.231"]
  },
  "description": "Test default Elasticsearch trial template",
  "id": "` + validTemplateID + `",
  "metadata": [{
	"key": "trial",
	"value": "true"
	}],
	"cluster_template": {
		"plan": {
			"cluster_topology": [{
				"node_type": {
					"master": true,
					"data": true
				},
				"instance_configuration_id": "default-elasticsearch",
				"size": {
					"value": 1024,
					"resource": "memory"
				}
		}],
		"elasticsearch": {
			"version": "6.2.3"
			}
		}
	},
	"system_owned": false
}`
	tests := []struct {
		name    string
		args    GetTemplateParams
		want    *models.DeploymentTemplateInfo
		wantErr bool
		error   string
	}{
		{
			name: "Platform deployment template show succeeds",
			args: GetTemplateParams{
				TemplateParams: TemplateParams{
					ID: validTemplateID,
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(templateGetSuccess),
						StatusCode: 200,
					}}),
				},
			},
			want: &models.DeploymentTemplateInfo{
				Name:        ec.String("(Trial) Default Elasticsearch"),
				ID:          validTemplateID,
				Description: "Test default Elasticsearch trial template",
				SystemOwned: ec.Bool(false),
				Metadata: []*models.MetadataItem{{

					Value: ec.String("true"),
					Key:   ec.String("trial"),
				}},
				Source: &models.ChangeSourceInfo{
					UserID:          "1",
					Facilitator:     ec.String("adminconsole"),
					Date:            &sourceDate,
					AdminID:         "admin",
					Action:          ec.String("deployments.create-template"),
					RemoteAddresses: []string{"52.205.1.231"},
				},
				ClusterTemplate: &models.DeploymentTemplateDefinitionRequest{
					Plan: &models.ElasticsearchClusterPlan{
						Elasticsearch: &models.ElasticsearchConfiguration{
							Version: "6.2.3",
						},
						ClusterTopology: []*models.ElasticsearchClusterTopologyElement{{
							InstanceConfigurationID: "default-elasticsearch",
							Size: &models.TopologySize{
								Value:    ec.Int32(1024),
								Resource: ec.String("memory"),
							},
							NodeType: &models.ElasticsearchNodeType{
								Master: ec.Bool(true),
								Data:   ec.Bool(true),
							},
						},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Platform deployment template show fails due to API error",
			args: GetTemplateParams{
				TemplateParams: TemplateParams{
					ID:  validTemplateID,
					API: api.NewMock(mock.Response{Error: errors.New("error")}),
				},
			},
			wantErr: true,
			error:   `Get https://mock-host/mock-path/platform/configuration/templates/deployments/84e0bd6d69bb44e294809d89cea88a7e?format=cluster&show_instance_configurations=false: error`,
		},
		{
			name: "Platform deployment template show fails with an empty API reference",
			args: GetTemplateParams{
				TemplateParams: TemplateParams{
					ID: validTemplateID,
				},
			},
			wantErr: true,
			error:   "api reference is required for command",
		},
		{
			name: "Platform deployment template show fails with an invalid template id",
			args: GetTemplateParams{
				TemplateParams: TemplateParams{
					ID:  invalidTemplateID,
					API: &api.API{},
				},
			},
			want:    nil,
			wantErr: true,
			error:   "invalid template ID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTemplate(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr != (tt.error != "") {
				t.Errorf("Get() expected error message = '%v', actual error message = '%v'", tt.error, err)
			}

			if tt.wantErr && err.Error() != tt.error {
				t.Errorf("Get() actual error = '%v', want error '%v'", err, tt.error)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name    string
		args    GetTemplateParams
		wantErr bool
		error   string
	}{
		{
			name: "Platform deployment template delete succeeds",
			args: GetTemplateParams{
				TemplateParams: TemplateParams{
					ID: validTemplateID,
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(`{}`),
						StatusCode: 200,
					}}),
				},
			},
			wantErr: false,
		},
		{
			name: "Platform deployment template delete fails due to API error",
			args: GetTemplateParams{
				TemplateParams: TemplateParams{
					ID:  validTemplateID,
					API: api.NewMock(mock.Response{Error: errors.New("error")}),
				},
			},
			wantErr: true,
			error:   `Delete https://mock-host/mock-path/platform/configuration/templates/deployments/84e0bd6d69bb44e294809d89cea88a7e: error`,
		},
		{
			name: "Platform deployment template delete fails with an empty API reference",
			args: GetTemplateParams{
				TemplateParams: TemplateParams{
					ID: validTemplateID,
				},
			},
			wantErr: true,
			error:   "api reference is required for command",
		},
		{
			name: "Platform deployment template delete fails with an invalid template id",
			args: GetTemplateParams{
				TemplateParams: TemplateParams{
					ID:  invalidTemplateID,
					API: &api.API{},
				},
			},
			wantErr: true,
			error:   "invalid template ID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DeleteTemplate(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr != (tt.error != "") {
				t.Errorf("Delete() expected error message = '%v', actual error message = '%v'", tt.error, err)
			}

			if tt.wantErr && err.Error() != tt.error {
				t.Errorf("Delete() actual error = '%v', want error '%v'", err, tt.error)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name    string
		args    CreateTemplateParams
		want    string
		wantErr bool
		error   string
	}{
		{
			name: "Platform deployment template create succeeds",
			args: CreateTemplateParams{
				DeploymentTemplateInfo: deploymentTemplateModel(),
				API: api.NewMock(mock.Response{Response: http.Response{
					Body:       mock.NewStringBody(`{"id": "9362b09d838640b2beea21b3343b4686"}"`),
					StatusCode: 201,
				}}),
			},
			wantErr: false,
			want:    "9362b09d838640b2beea21b3343b4686",
		},
		{
			name: "Platform deployment template create succeeds specifying template ID",
			args: CreateTemplateParams{
				DeploymentTemplateInfo: deploymentTemplateModelWithID("template-id"),
				API: api.NewMock(mock.Response{Response: http.Response{
					Body:       mock.NewStringBody(`{"id": "template-id"}"`),
					StatusCode: 201,
				}}),
			},
			wantErr: false,
			want:    "template-id",
		},
		{
			name: "Platform deployment template create fails due to API error",
			args: CreateTemplateParams{
				DeploymentTemplateInfo: deploymentTemplateModel(),
				API:                    api.NewMock(mock.Response{Error: errors.New("error")}),
			},
			wantErr: true,
			error:   `Post https://mock-host/mock-path/platform/configuration/templates/deployments: error`,
		},
		{
			name: "Platform deployment template create fails with an empty API reference",
			args: CreateTemplateParams{
				DeploymentTemplateInfo: &models.DeploymentTemplateInfo{},
			},
			wantErr: true,
			error:   "api reference is required for command",
			want:    "",
		},
		{
			name: "Platform deployment template create fails with an empty deployment template",
			args: CreateTemplateParams{
				API: &api.API{},
			},
			wantErr: true,
			error:   "deployment template is missing",
		},
		{
			name: "Platform deployment template create fails with an incomplete deployment template",
			args: CreateTemplateParams{
				API:                    &api.API{},
				DeploymentTemplateInfo: &models.DeploymentTemplateInfo{},
			},
			wantErr: true,
			error:   "validation failure list:\nname in body is required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateTemplate(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr != (tt.error != "") {
				t.Errorf("Create() expected error message = '%v', actual error message = '%v'", tt.error, err)
			}

			if tt.wantErr && err.Error() != tt.error {
				t.Errorf("Create() actual error = '%v', want error '%v'", err, tt.error)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name    string
		args    UpdateTemplateParams
		wantErr bool
		error   string
	}{
		{
			name: "Platform deployment template update succeeds",
			args: UpdateTemplateParams{
				DeploymentTemplateInfo: deploymentTemplateModel(),
				TemplateParams: TemplateParams{
					ID: validTemplateID,
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(`{}`),
						StatusCode: 200,
					}}),
				},
			},
			wantErr: false,
		},
		{
			name: "Platform deployment template update fails due to API error",
			args: UpdateTemplateParams{
				DeploymentTemplateInfo: deploymentTemplateModel(),
				TemplateParams: TemplateParams{
					ID:  validTemplateID,
					API: api.NewMock(mock.Response{Error: errors.New("error")}),
				},
			},
			wantErr: true,
			error:   `Put https://mock-host/mock-path/platform/configuration/templates/deployments/84e0bd6d69bb44e294809d89cea88a7e?create_only=false: error`,
		},
		{
			name: "Platform deployment template update fails with an empty API reference",
			args: UpdateTemplateParams{
				DeploymentTemplateInfo: deploymentTemplateModel(),
				TemplateParams: TemplateParams{
					ID: validTemplateID,
				},
			},
			wantErr: true,
			error:   "api reference is required for command",
		},
		{
			name: "Platform deployment template update fails with an invalid template id",
			args: UpdateTemplateParams{
				DeploymentTemplateInfo: deploymentTemplateModel(),
				TemplateParams: TemplateParams{
					ID:  invalidTemplateID,
					API: &api.API{},
				},
			},
			wantErr: true,
			error:   "invalid template ID",
		},
		{
			name: "Platform deployment template update fails with an empty deployment template",
			args: UpdateTemplateParams{
				TemplateParams: TemplateParams{
					ID:  validTemplateID,
					API: &api.API{},
				},
			},
			wantErr: true,
			error:   "deployment template is missing",
		},
		{
			name: "Platform deployment template update fails with an incomplete deployment template",
			args: UpdateTemplateParams{
				DeploymentTemplateInfo: &models.DeploymentTemplateInfo{},
				TemplateParams: TemplateParams{
					ID:  validTemplateID,
					API: &api.API{},
				},
			},
			wantErr: true,
			error:   "validation failure list:\nname in body is required",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UpdateTemplate(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr != (tt.error != "") {
				t.Errorf("Delete() expected error message = '%v', actual error message = '%v'", tt.error, err)
			}

			if tt.wantErr && err.Error() != tt.error {
				t.Errorf("Delete() actual error = '%v', want error '%v'", err, tt.error)
			}
		})
	}
}

func deploymentTemplateModelWithID(id string) *models.DeploymentTemplateInfo {
	var model = deploymentTemplateModel()
	model.ID = id

	return model
}
func deploymentTemplateModel() *models.DeploymentTemplateInfo {
	sourceDate, _ := strfmt.ParseDateTime("2018-04-19T18:16:57.297Z")

	template := models.DeploymentTemplateInfo{Name: ec.String("(Trial) Default Elasticsearch"),
		SystemOwned: ec.Bool(false),
		Metadata: []*models.MetadataItem{{

			Value: ec.String("true"),
			Key:   ec.String("trial"),
		}},
		Source: &models.ChangeSourceInfo{
			UserID:          "1",
			Facilitator:     ec.String("adminconsole"),
			Date:            &sourceDate,
			AdminID:         "admin",
			Action:          ec.String("deployments.create-template"),
			RemoteAddresses: []string{"52.205.1.231"},
		},
		ClusterTemplate: &models.DeploymentTemplateDefinitionRequest{
			Plan: &models.ElasticsearchClusterPlan{
				Elasticsearch: &models.ElasticsearchConfiguration{
					Version: "6.2.3",
				},
				ClusterTopology: []*models.ElasticsearchClusterTopologyElement{{
					InstanceConfigurationID: "default-elasticsearch",
					Size: &models.TopologySize{
						Value:    ec.Int32(1024),
						Resource: ec.String("memory"),
					},
					NodeType: &models.ElasticsearchNodeType{
						Master: ec.Bool(true),
						Data:   ec.Bool(true),
					},
				},
				},
			},
		}}

	return &template
}

func TestDeploymentTemplateOperationOutput(t *testing.T) {
	tests := []struct {
		templateOperation TemplateOperation
		output            string
	}{
		{
			templateOperation: TemplateOperationNone,
			output:            "None",
		},
		{
			templateOperation: TemplateOperationUpdate,
			output:            "Update",
		},
		{
			templateOperation: TemplateOperationCreate,
			output:            "Create",
		},
		{
			templateOperation: 10,
			output:            "Invalid",
		},
	}
	for _, tt := range tests {
		if tt.templateOperation.String() != tt.output {
			t.Errorf("templateOperation got %v, expected %v", tt.templateOperation, tt.output)
		}
	}
}

func TestGetRemoteDeploymentTemplates(t *testing.T) {
	var templateRemote = `
	{
  "name": "(Trial) Default Elasticsearch",
  "description": "Test default Elasticsearch trial template",
  "id": "` + validTemplateID + `",
  "metadata": [{
	"key": "trial",
	"value": "true"
	}],
	"cluster_template": {
		"plan": {
			"cluster_topology": [{
				"node_type": {
					"master": true,
					"data": true
				},
				"instance_configuration_id": "default-elasticsearch",
				"size": {
					"value": 1024,
					"resource": "memory"
				}
		}],
		"elasticsearch": {
			"version": "6.2.3"
			}
		}
	},
	"system_owned": false
}`

	var deploymentTemplateNotFound = `{
  "errors": [
	{
	  "code": "templates.template_not_found",
	  "fields": null,
	  "message": "Template Not Found"
	}
  ]
}`

	var templateRemotelist = `[
	{
	  "cluster_template": {
		"plan": {
		  "cluster_topology": [
			{
			  "instance_configuration_id": "default-elasticsearch",
			  "node_type": {
				"data": true,
				"master": true
			  },
			  "size": {
				"resource": "memory",
				"value": 1024
			  }
			}
		  ],
		  "elasticsearch": {
			"version": "6.2.3"
		  }
		}
	  },
	  "description": "Test default Elasticsearch trial template",
	  "id": "dt-1",
	  "metadata": [
		{
		  "key": "trial",
		  "value": "true"
		}
	  ],
	  "name": "(Trial) Default Elasticsearch",
	  "system_owned": false
	},
	{
	  "cluster_template": {
		"plan": {
		  "cluster_topology": [
			{
			  "instance_configuration_id": "default-elasticsearch",
			  "node_type": {
				"data": true,
				"master": true
			  },
			  "size": {
				"resource": "memory",
				"value": 2048
			  }
			}
		  ],
		  "elasticsearch": {
			"version": "6.3.1"
		  }
		}
	  },
	  "description": "Test default Elasticsearch trial template",
	  "id": "dt-2",
	  "metadata": [
		{
		  "key": "trial",
		  "value": "true"
		}
	  ],
	  "name": "(Trial) Default Elasticsearch",
	  "system_owned": false
	}
]`

	tests := []struct {
		name string
		args GetTemplateParams
		want map[string]*models.DeploymentTemplateInfo
		err  error
	}{
		{
			name: "Retrieve single deployment template",
			args: GetTemplateParams{
				TemplateParams: TemplateParams{
					ID: validTemplateID,
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(templateRemote),
						StatusCode: 200,
					}}),
				},
			},
			want: map[string]*models.DeploymentTemplateInfo{
				validTemplateID: {
					Name:        ec.String("(Trial) Default Elasticsearch"),
					ID:          validTemplateID,
					Description: "Test default Elasticsearch trial template",
					SystemOwned: ec.Bool(false),
					Metadata: []*models.MetadataItem{{

						Value: ec.String("true"),
						Key:   ec.String("trial"),
					}},
					ClusterTemplate: &models.DeploymentTemplateDefinitionRequest{
						Plan: &models.ElasticsearchClusterPlan{
							Elasticsearch: &models.ElasticsearchConfiguration{
								Version: "6.2.3",
							},
							ClusterTopology: []*models.ElasticsearchClusterTopologyElement{{
								InstanceConfigurationID: "default-elasticsearch",
								Size: &models.TopologySize{
									Value:    ec.Int32(1024),
									Resource: ec.String("memory"),
								},
								NodeType: &models.ElasticsearchNodeType{
									Master: ec.Bool(true),
									Data:   ec.Bool(true),
								},
							},
							},
						},
					},
				},
			},
		},
		{
			name: "Try to retrieve a non existing deployment template",
			args: GetTemplateParams{
				TemplateParams: TemplateParams{
					ID: validTemplateID,
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(deploymentTemplateNotFound),
						StatusCode: 404,
					}}),
				},
			},
			want: map[string]*models.DeploymentTemplateInfo{
				validTemplateID: nil,
			},
		},
		{
			name: "Retrieve two deployment template",
			args: GetTemplateParams{
				TemplateParams: TemplateParams{
					API: api.NewMock(mock.Response{Response: http.Response{
						Body:       mock.NewStringBody(templateRemotelist),
						StatusCode: 200,
					}}),
				},
			},
			want: map[string]*models.DeploymentTemplateInfo{
				"dt-1": {
					Name:        ec.String("(Trial) Default Elasticsearch"),
					ID:          "dt-1",
					Description: "Test default Elasticsearch trial template",
					SystemOwned: ec.Bool(false),
					Metadata: []*models.MetadataItem{{

						Value: ec.String("true"),
						Key:   ec.String("trial"),
					}},
					ClusterTemplate: &models.DeploymentTemplateDefinitionRequest{
						Plan: &models.ElasticsearchClusterPlan{
							Elasticsearch: &models.ElasticsearchConfiguration{
								Version: "6.2.3",
							},
							ClusterTopology: []*models.ElasticsearchClusterTopologyElement{{
								InstanceConfigurationID: "default-elasticsearch",
								Size: &models.TopologySize{
									Value:    ec.Int32(1024),
									Resource: ec.String("memory"),
								},
								NodeType: &models.ElasticsearchNodeType{
									Master: ec.Bool(true),
									Data:   ec.Bool(true),
								},
							},
							},
						},
					},
				},
				"dt-2": {
					Name:        ec.String("(Trial) Default Elasticsearch"),
					ID:          "dt-2",
					Description: "Test default Elasticsearch trial template",
					SystemOwned: ec.Bool(false),
					Metadata: []*models.MetadataItem{{

						Value: ec.String("true"),
						Key:   ec.String("trial"),
					}},
					ClusterTemplate: &models.DeploymentTemplateDefinitionRequest{
						Plan: &models.ElasticsearchClusterPlan{
							Elasticsearch: &models.ElasticsearchConfiguration{
								Version: "6.3.1",
							},
							ClusterTopology: []*models.ElasticsearchClusterTopologyElement{{
								InstanceConfigurationID: "default-elasticsearch",
								Size: &models.TopologySize{
									Value:    ec.Int32(2048),
									Resource: ec.String("memory"),
								},
								NodeType: &models.ElasticsearchNodeType{
									Master: ec.Bool(true),
									Data:   ec.Bool(true),
								},
							},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRemoteDeploymentTemplates(tt.args)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("GetRemoteDeploymentTemplates() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRemoteDeploymentTemplates() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
