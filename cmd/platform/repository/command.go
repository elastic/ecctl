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

	"github.com/elastic/cloud-sdk-go/pkg/input"
	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/platform/snaprepo"
	"github.com/elastic/ecctl/pkg/util"
)

const (
	snapshotCreateShortHelp = "Creates / updates a snapshot repository"
)

var (
	snapshotShortHelp = cmdutil.AdminReqDescription("Manages snapshot repositories")

	snapshotLongHelp = `
Manages snapshot repositories that are used by Elasticsearch clusters
to perform snapshot operations.
`[1:]

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

// Command represents the top level repository command.
var Command = &cobra.Command{
	Use:     "repository",
	Short:   snapshotShortHelp,
	Long:    snapshotLongHelp,
	PreRunE: cobra.MaximumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var platformSnapshotShowCmd = &cobra.Command{
	Use:     "show <repository name>",
	Short:   "Obtains a snapshot repository config",
	PreRunE: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := snaprepo.Get(snaprepo.GetParams{
			Params: snaprepo.Params{
				API: ecctl.Get().API,
			},
			Name: args[0],
		})
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format("", repo)
	},
}

var platformSnapshotListCmd = &cobra.Command{
	Use:     "list",
	Short:   "Lists all the snapshot repositories",
	PreRunE: cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		repos, err := snaprepo.List(snaprepo.Params{
			API: ecctl.Get().API,
		})
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format(
			filepath.Join("platform", "repositorylist"), repos,
		)
	},
}

var platformSnapshotDeleteCmd = &cobra.Command{
	Use:     "delete <repository name>",
	Short:   "Deletes a snapshot repositories",
	PreRunE: cobra.MinimumNArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		return snaprepo.Delete(snaprepo.DeleteParams{
			Params: snaprepo.Params{
				API: ecctl.Get().API,
			},
			Name: args[0],
		})
	},
}

var platformSnapshotCreateCmd = &cobra.Command{
	Use:     "create <repository name> --settings <settings file>",
	Aliases: []string{"update", "set"},
	Short:   snapshotCreateShortHelp,
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

	if err := sdkcmdutil.FileOrStdin(cmd, "settings"); err != nil {
		return err
	}

	f, err := input.NewFileOrReader(os.Stdin, configFile)
	if err != nil {
		return err
	}
	defer f.Close()

	var repoType = cmd.Flag("type").Value.String()
	config, err := parseRepoSettingsByType(f, repoType)
	if err != nil {
		return err
	}

	return snaprepo.Set(snaprepo.SetParams{
		Params: snaprepo.Params{
			API: ecctl.Get().API,
		},
		Name:   args[0],
		Config: config,
		Type:   repoType,
	})
}

func parseRepoSettingsByType(r io.Reader, t string) (util.Validator, error) {
	switch t {
	case "s3":
		return snaprepo.ParseS3Config(r)
	default:
		return snaprepo.ParseGenericConfig(r)
	}
}

func init() {
	Command.AddCommand(
		platformSnapshotListCmd,
		platformSnapshotShowCmd,
		platformSnapshotDeleteCmd,
		platformSnapshotCreateCmd,
	)

	platformSnapshotCreateCmd.Flags().String("settings", "", "Configuration file for the snapshot repository")
	platformSnapshotCreateCmd.Flags().String("type", "s3", "Repository type that will be configured")
	platformSnapshotCreateCmd.MarkFlagFilename("settings", "json", "yaml", "yml")
}
