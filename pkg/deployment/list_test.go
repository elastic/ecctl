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
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"

	"github.com/elastic/ecctl/pkg/util"
)

var deploymentList = `{
  "deployments": [
    {
      "id": "709551ecd21143adbb7a37bbf36c4321",
      "name": "admin-console-elasticsearch",
      "resources": [
        {
          "cloud_id": "admin-console-elasticsearch:MTkyLjE2OC40NC4xMC5pcC5lcy5pbzo5MjQzJDcwOTU1MWVjZDIxMTQzYWRiYjdhMzdiYmYzNmM0MzIxJA==",
          "id": "709551ecd21143adbb7a37bbf36c4321",
          "kind": "elasticsearch",
          "ref_id": "elasticsearch",
          "region": "ece-region"
        }
      ]
    },
    {
      "id": "83c2027d92c44aeba34378f44cf08137",
      "name": "logging-and-metrics",
      "resources": [
        {
          "cloud_id": "logging-and-metrics:MTkyLjE2OC40NC4xMC5pcC5lcy5pbzo5MjQzJDgzYzIwMjdkOTJjNDRhZWJhMzQzNzhmNDRjZjA4MTM3JGQxNDBiNmUzNWIwNDQ1YTlhZmQ0NDIwZjZjMWIzYjIx",
          "id": "83c2027d92c44aeba34378f44cf08137",
          "kind": "elasticsearch",
          "ref_id": "elasticsearch",
          "region": "ece-region"
        },
        {
          "elasticsearch_cluster_ref_id": "elasticsearch",
          "id": "d140b6e35b0445a9afd4420f6c1b3b21",
          "kind": "kibana",
          "ref_id": "kibana",
          "region": "ece-region"
        }
      ]
    },
    {
      "id": "89a54cee4a2847c291f520ac00760886",
      "name": "security-cluster",
      "resources": [
        {
          "cloud_id": "security-cluster:MTkyLjE2OC40NC4xMC5pcC5lcy5pbzo5MjQzJDg5YTU0Y2VlNGEyODQ3YzI5MWY1MjBhYzAwNzYwODg2JA==",
          "id": "89a54cee4a2847c291f520ac00760886",
          "kind": "elasticsearch",
          "ref_id": "elasticsearch",
          "region": "ece-region"
        }
      ]
    }
  ]
}
`

func TestList(t *testing.T) {
	type args struct {
		params ListParams
	}
	tests := []struct {
		name string
		args args
		want *models.DeploymentsListResponse
		err  error
	}{
		{
			name: "fails on parameter validation",
			err:  util.ErrAPIReq,
		},
		{
			name: "fails on API error",
			args: args{params: ListParams{
				API: api.NewMock(mock.New500Response(mock.NewStringBody(`{"error": "some error"}`))),
			}},
			err: errors.New(`{"error": "some error"}`),
		},
		{
			name: "Succeeds",
			args: args{params: ListParams{
				API: api.NewMock(mock.New200Response(mock.NewStringBody(deploymentList))),
			}},
			want: &models.DeploymentsListResponse{Deployments: []*models.DeploymentsListingData{
				{
					ID:   ec.String("709551ecd21143adbb7a37bbf36c4321"),
					Name: ec.String("admin-console-elasticsearch"),
					Resources: []*models.DeploymentResource{
						{
							CloudID: "admin-console-elasticsearch:MTkyLjE2OC40NC4xMC5pcC5lcy5pbzo5MjQzJDcwOTU1MWVjZDIxMTQzYWRiYjdhMzdiYmYzNmM0MzIxJA==",
							ID:      ec.String("709551ecd21143adbb7a37bbf36c4321"),
							Kind:    ec.String("elasticsearch"),
							RefID:   ec.String("elasticsearch"),
							Region:  ec.String("ece-region"),
						},
					},
				},
				{
					ID:   ec.String("83c2027d92c44aeba34378f44cf08137"),
					Name: ec.String("logging-and-metrics"),
					Resources: []*models.DeploymentResource{
						{
							CloudID: "logging-and-metrics:MTkyLjE2OC40NC4xMC5pcC5lcy5pbzo5MjQzJDgzYzIwMjdkOTJjNDRhZWJhMzQzNzhmNDRjZjA4MTM3JGQxNDBiNmUzNWIwNDQ1YTlhZmQ0NDIwZjZjMWIzYjIx",
							ID:      ec.String("83c2027d92c44aeba34378f44cf08137"),
							Kind:    ec.String("elasticsearch"),
							RefID:   ec.String("elasticsearch"),
							Region:  ec.String("ece-region"),
						},
						{
							ElasticsearchClusterRefID: "elasticsearch",
							ID:                        ec.String("d140b6e35b0445a9afd4420f6c1b3b21"),
							Kind:                      ec.String("kibana"),
							RefID:                     ec.String("kibana"),
							Region:                    ec.String("ece-region"),
						},
					},
				},
				{
					ID:   ec.String("89a54cee4a2847c291f520ac00760886"),
					Name: ec.String("security-cluster"),
					Resources: []*models.DeploymentResource{
						{
							CloudID: "security-cluster:MTkyLjE2OC40NC4xMC5pcC5lcy5pbzo5MjQzJDg5YTU0Y2VlNGEyODQ3YzI5MWY1MjBhYzAwNzYwODg2JA==",
							ID:      ec.String("89a54cee4a2847c291f520ac00760886"),
							Kind:    ec.String("elasticsearch"),
							RefID:   ec.String("elasticsearch"),
							Region:  ec.String("ece-region"),
						},
					},
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := List(tt.args.params)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("List() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
