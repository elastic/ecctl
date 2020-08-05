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

package cmddeploymentnote

import (
	"path/filepath"

	"github.com/elastic/cloud-sdk-go/pkg/api"
	"github.com/elastic/cloud-sdk-go/pkg/api/deploymentapi"
	"github.com/elastic/cloud-sdk-go/pkg/api/deploymentapi/noteapi"
	"github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"
	"github.com/spf13/cobra"

	"github.com/elastic/ecctl/pkg/ecctl"
	"github.com/elastic/ecctl/pkg/util"
)

// Command represents the deployment note subcommand.
var Command = &cobra.Command{
	Use:     "note",
	Short:   "Manages a deployment's notes",
	PreRunE: cobra.MaximumNArgs(0),
	Run:     func(cmd *cobra.Command, args []string) { cmd.Help() },
}

var deploymentNoteCreateCmd = &cobra.Command{
	Use:     "create <deployment id> --comment <comment content>",
	Aliases: []string{"add"},
	Short:   "Adds a note to a deployment",
	PreRunE: cmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		comment, _ := cmd.Flags().GetString("comment")
		esID, err := getElasticsearchID(args[0], ecctl.Get().API)
		if err != nil {
			return err
		}

		return noteapi.Add(noteapi.AddParams{
			Params: noteapi.Params{
				API:    ecctl.Get().API,
				Region: ecctl.Get().Config.Region,
				ID:     esID,
			},
			Message: comment,
			UserID:  ecctl.Get().Config.User,
		})
	},
}

var deploymentNoteListCmd = &cobra.Command{
	Use:     "list <deployment id>",
	Short:   "Lists the deployment notes",
	PreRunE: cmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		esID, err := getElasticsearchID(args[0], ecctl.Get().API)
		if err != nil {
			return err
		}

		res, err := noteapi.List(noteapi.Params{
			API:    ecctl.Get().API,
			Region: ecctl.Get().Config.Region,
			ID:     esID,
		})
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format(filepath.Join("deployment", "notelist"), res)
	},
}

var deploymentNoteUpdateCmd = &cobra.Command{
	Use:     "update <deployment id> --id <note id> --comment <comment content>",
	Short:   "Updates the deployment notes",
	PreRunE: cmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		comment, _ := cmd.Flags().GetString("comment")
		noteID, _ := cmd.Flags().GetString("id")
		esID, err := getElasticsearchID(args[0], ecctl.Get().API)
		if err != nil {
			return err
		}

		return util.ReturnErrOnly(
			noteapi.Update(noteapi.UpdateParams{
				Message: comment,
				UserID:  ecctl.Get().Config.User,
				NoteID:  noteID,
				Params: noteapi.Params{
					API:    ecctl.Get().API,
					Region: ecctl.Get().Config.Region,
					ID:     esID,
				},
			}),
		)
	},
}

var deploymentNoteShowCmd = &cobra.Command{
	Use:     "show <deployment id> --id <note id>",
	Short:   "Shows a deployment note",
	PreRunE: cmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		noteID, _ := cmd.Flags().GetString("id")
		esID, err := getElasticsearchID(args[0], ecctl.Get().API)
		if err != nil {
			return err
		}

		res, err := noteapi.Get(noteapi.GetParams{
			NoteID: noteID,
			Params: noteapi.Params{
				API:    ecctl.Get().API,
				Region: ecctl.Get().Config.Region,
				ID:     esID,
			},
		})
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format(filepath.Join("deployment", "noteshow"), res)
	},
}

func init() {
	Command.AddCommand(
		deploymentNoteCreateCmd,
		deploymentNoteListCmd,
		deploymentNoteUpdateCmd,
		deploymentNoteShowCmd,
	)

	deploymentNoteCreateCmd.Flags().String("comment", "", "Content of your deployment note")
	deploymentNoteCreateCmd.MarkFlagRequired("comment")
	deploymentNoteUpdateCmd.Flags().String("comment", "", "Content of your deployment note")
	deploymentNoteUpdateCmd.MarkFlagRequired("comment")
	deploymentNoteUpdateCmd.Flags().String("id", "", "Note ID")
	deploymentNoteUpdateCmd.MarkFlagRequired("id")
	deploymentNoteShowCmd.Flags().String("id", "", "Note ID")
	deploymentNoteShowCmd.MarkFlagRequired("id")
}

// The notes endpoint only allows elasticsearch IDs, but the endpoint
// `deployments/DEPLOYMENT_ID/notes` endpoint is not specific to elasticsearch.
// For the time being the command layer will fetch the ES ID until the endpoint
// accepts deployment IDs.
func getElasticsearchID(deploymentID string, ecAPI *api.API) (string, error) {
	esID, err := deploymentapi.GetElasticsearchID(deploymentapi.GetParams{
		API:          ecAPI,
		DeploymentID: deploymentID,
	})
	if err != nil {
		return "", err
	}

	return esID, nil
}
