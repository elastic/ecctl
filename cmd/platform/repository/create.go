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

package cmdrepository

import (
	"io"
	"os"
	"path/filepath"

	"github.com/elastic/cloud-sdk-go/pkg/api/platformapi/snaprepoapi"
	"github.com/elastic/cloud-sdk-go/pkg/input"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/util"
)

const (
	snapshotCreateShortHelp = "Creates / updates a snapshot repository"
)

var (
	snapshotCreateLongHelp = `
Creates / updates a snapshot repository using a set of settings that can be
specified as a (yaml|json) file with the --settings flag.

The available settings to set depend on the the --type flag (default s3). A
list with the supported settings for each snapshot can be found in the docs:
https://www.elastic.co/guide/en/elasticsearch/reference/current/modules-snapshots.html#_repository_plugins

The --type flag can be set to any arbitrary value if it's differs from "s3".
Only the S3 available settings are validated.
`[1:]

	snapshotCreateExamples = `
ecctl platform repository create my-snapshot-repo --settings settings.yml

ecctl platform repository update my-snapshot-repo --settings settings.yml

ecctl platform repository create custom --type fs --settings settings.yml
`[1:]
)

var platformSnapshotCreateCmd = &cobra.Command{
	Use:     "create <repository name> --settings <settings file>",
	Aliases: []string{"update", "set"},
	Short:   cmdutil.AdminReqDescription(snapshotCreateShortHelp),
	Long:    snapshotCreateLongHelp,
	Example: snapshotCreateExamples,
	PreRunE: cobra.MinimumNArgs(1),
	RunE:    setSnapshot,
}

func setSnapshot(cmd *cobra.Command, args []string) error {
	configFile := cmd.Flag("settings").Value.String()
	if !filepath.IsAbs(configFile) {
		var err error
		configFile, err = filepath.Abs(configFile)
		if err != nil {
			return err
		}
	}

	f, err := input.NewFileOrReader(os.Stdin, configFile)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	var repoType = cmd.Flag("type").Value.String()
	config, err := parseRepoSettingsByType(f, repoType)
	if err != nil {
		return err
	}

	return snaprepoapi.Set(snaprepoapi.SetParams{
		API:    ecctl.Get().API,
		Region: ecctl.Get().Config.Region,
		Name:   args[0],
		Config: config,
		Type:   repoType,
	})
}

func parseRepoSettingsByType(r io.Reader, t string) (util.Validator, error) {
	switch t {
	case "s3":
		return snaprepoapi.ParseS3Config(r)
	default:
		return snaprepoapi.ParseGenericConfig(r)
	}
}

func init() {
	Command.AddCommand(platformSnapshotCreateCmd)

	platformSnapshotCreateCmd.Flags().String("settings", "", "Configuration file for the snapshot repository")
	platformSnapshotCreateCmd.Flags().String("type", "s3", "Repository type that will be configured")
	platformSnapshotCreateCmd.MarkFlagFilename("settings", "json", "yaml", "yml")
}
