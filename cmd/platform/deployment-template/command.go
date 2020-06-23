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

package cmddeploymentdemplate

import (
	"encoding/json"
	"os"
	"path"

	"github.com/elastic/cloud-sdk-go/pkg/input"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
)

const (
	showInstanceConfigurations = "show-instance-configurations"
	stackVersion               = "stack-version"
	filter                     = "filter"
)

// Command represents the top level deployment-template command.
var Command = &cobra.Command{
	Use:     "deployment-template",
	Short:   cmdutil.AdminReqDescription("Manages deployment templates"),
	PreRunE: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func parseTemplateFile(fp string) (*models.DeploymentTemplateInfo, error) {
	if fp != "" {
		if ext := path.Ext(fp); ext != ".json" {
			return nil, errors.New("unsupported file type, only json template files are currently supported")
		}
	}

	templateFile, err := input.NewFileOrReader(os.Stdin, fp)
	if err != nil {
		return nil, err
	}
	defer templateFile.Close()

	var templateConfiguration models.DeploymentTemplateInfo
	if err := json.NewDecoder(templateFile).Decode(&templateConfiguration); err != nil {
		return nil, err
	}

	return &templateConfiguration, nil
}
