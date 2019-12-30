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

	"github.com/spf13/cobra"

	cmdutil "github.com/elastic/ecctl/cmd/util"
	"github.com/elastic/ecctl/pkg/deployment"
	"github.com/elastic/ecctl/pkg/deployment/note"
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
	Use:     "create <deployment id> --message <message content> --type [elasticsearch|kibana|apm]",
	Aliases: []string{"add"},
	Short:   "Adds a note to a deployment",
	PreRunE: cmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return note.Add(note.AddParams{
			Params: deployment.Params{
				API: ecctl.Get().API,
				ID:  args[0],
			},
			Message: cmd.Flag("message").Value.String(),
			UserID:  ecctl.Get().Config.User,
		})
	},
}

var deploymentNoteListCmd = &cobra.Command{
	Use:     "list <deployment id>",
	Short:   "Lists the deployment notes",
	PreRunE: cmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := note.List(note.ListParams{
			Params: deployment.Params{
				API: ecctl.Get().API,
				ID:  args[0],
			},
		})
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format(filepath.Join("deployment", "notelist"), res)
	},
}

var deploymentNoteUpdateCmd = &cobra.Command{
	Use:     "update <deployment id> --id <note id> --message <message content>",
	Short:   "Updates the deployment notes",
	PreRunE: cmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return util.ReturnErrOnly(
			note.Update(note.UpdateParams{
				Message: cmd.Flag("message").Value.String(),
				UserID:  ecctl.Get().Config.User,
				Params: note.Params{
					NoteID: cmd.Flag("id").Value.String(),
					Params: deployment.Params{
						API: ecctl.Get().API,
						ID:  args[0],
					},
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
		res, err := note.Get(note.GetParams{
			Params: note.Params{
				NoteID: cmd.Flag("id").Value.String(),
				Params: deployment.Params{
					API: ecctl.Get().API,
					ID:  args[0],
				},
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

	deploymentNoteUpdateCmd.Flags().String("id", "", "Note ID")
	deploymentNoteShowCmd.Flags().String("id", "", "Note ID")
}
