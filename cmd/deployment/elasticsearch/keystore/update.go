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

package cmdeskeystore

import (
	"errors"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/elastic/cloud-sdk-go/pkg/api/deploymentapi/eskeystoreapi"
	"github.com/elastic/cloud-sdk-go/pkg/models"
	sdkcmdutil "github.com/elastic/cloud-sdk-go/pkg/util/cmdutil"

	"github.com/elastic/ecctl/pkg/ecctl"
)

const updateLong = `Changes the contents of the Elasticsearch resource keystore from the
specified deployment by using the PATCH method. The payload is a partial payload where
any ignored current keystore items are not removed, unless the secrets are
set to "null": {"secrets": {"my-secret": null}}.`

var (
	errReadingDefPrefix  = "failed reading keystore secret definition"
	errReadingDefMessage = "provide a valid keystore secret definition using the --file flag"
)

var updateCmd = &cobra.Command{
	Use:     "update <deployment id> [--ref-id <ref-id>] {--file=<filename>.json}",
	Long:    updateLong,
	Aliases: []string{"set"},
	Short:   "Updates the contents of an Elasticsearch keystore",
	PreRunE: sdkcmdutil.MinimumNArgsAndUUID(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var contents models.KeystoreContents
		if err := sdkcmdutil.DecodeDefinition(cmd, "file", &contents); err != nil {
			if errors.Is(err, io.EOF) {
				return fmt.Errorf("%s: %s", errReadingDefPrefix, errReadingDefMessage)
			}
			return fmt.Errorf("%s: %s", errReadingDefPrefix, err.Error())
		}

		refID, _ := cmd.Flags().GetString("ref-id")
		res, err := eskeystoreapi.Update(eskeystoreapi.UpdateParams{
			API:          ecctl.Get().API,
			DeploymentID: args[0],
			RefID:        refID,
			Contents:     &contents,
		})
		if err != nil {
			return err
		}

		return ecctl.Get().Formatter.Format("deployment/eskeystore_show", res)
	},
}

func init() {
	Command.AddCommand(updateCmd)
	updateCmd.Flags().String("ref-id", "", "Optional ref_id to use for the Elasticsearch resource, auto-discovered if not specified.")
	updateCmd.Flags().StringP("file", "p", "", "Required json formatted file path with the keystore secret contents.")
	updateCmd.MarkFlagFilename("file", "json")
}
