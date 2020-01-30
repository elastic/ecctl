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

package cmdutil

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/mock"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/elastic/cloud-sdk-go/pkg/util/ec"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/elastic/ecctl/pkg/util"
)

func TestGetInstances(t *testing.T) {
	cmdWithSliceFlag := &cobra.Command{
		Use: "something",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	cmdWithSliceFlag.Flags().StringSlice("instance", []string{"1", "2", "3"}, "instance")
	cmdWithSliceFlag.Flags().Bool("all", false, "all")
	cmdWithAllFlag := &cobra.Command{
		Use: "something",
		Run: func(cmd *cobra.Command, args []string) {},
	}
	cmdWithAllFlag.Flags().Bool("all", true, "all")
	type args struct {
		cmd      *cobra.Command
		params   util.ClusterParams
		flagName string
	}
	tests := []struct {
		name string
		args args
		want []string
		err  error
	}{
		{
			name: "obtains the instances explicitly set",
			args: args{
				cmd:      cmdWithSliceFlag,
				flagName: "instance",
			},
			want: []string{"1", "2", "3"},
		},
		{
			name: "obtains the instances from the cluster topology",
			args: args{
				cmd: cmdWithAllFlag,
				params: util.ClusterParams{
					ClusterID: util.ValidClusterID,
					API: api.NewMock(mock.Response{Response: http.Response{
						StatusCode: 200,
						Body: mock.NewStructBody(models.ElasticsearchClusterInfo{
							Topology: &models.ClusterTopologyInfo{
								Instances: []*models.ClusterInstanceInfo{
									{InstanceName: ec.String("instance-000000")},
								},
							},
						}),
					}}),
				},
			},
			want: []string{"instance-000000"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetInstances(tt.args.cmd, tt.args.params, tt.args.flagName)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("GetInstances() error = %v, wantErr %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetInstances() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddTypeFlag(t *testing.T) {
	var wantSomethingAssert = &flag.Flag{
		Name:  "type",
		Usage: "Optional deployment resource type (apm, appsearch, kibana)",
		Annotations: map[string][]string{
			cobra.BashCompCustom: {"__ecctl_valid_stateless_types"},
		},
	}

	var wantSomethingRequiredAssert = &flag.Flag{
		Name:  "type",
		Usage: "Required deployment resource type (apm, appsearch, kibana, elasticsearch)",
		Annotations: map[string][]string{
			cobra.BashCompCustom: {"__ecctl_valid_all_types"},
		},
	}

	type args struct {
		cmd    *cobra.Command
		prefix string
		all    bool
	}
	tests := []struct {
		name string
		args args
		want *flag.Flag
	}{
		{
			name: "Annotates the type flag with stateless types",
			args: args{
				cmd: &cobra.Command{
					Use: "something",
					Run: func(cmd *cobra.Command, args []string) {},
				},
				prefix: "Optional",
				all:    false,
			},
			want: wantSomethingAssert,
		},
		{
			name: "Annotates the type flag with all types",
			args: args{
				cmd: &cobra.Command{
					Use: "somethingrequired",
					Run: func(cmd *cobra.Command, args []string) {},
				},
				prefix: "Required",
				all:    true,
			},
			want: wantSomethingRequiredAssert,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddTypeFlag(tt.args.cmd, tt.args.prefix, tt.args.all)
			got := tt.args.cmd.Flag("type")
			got.Value = nil
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddTypeFlag() got = \n%+v, want \n%+v", got, tt.want)
			}
		})
	}
}
