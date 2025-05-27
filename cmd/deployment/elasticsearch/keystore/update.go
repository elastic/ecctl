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

const updateLong = `Changes the contents of the Elasticsearch resource keystore from the specified
deployment by using the PATCH method. The payload is a partial payload where any
omitted current keystore items are not removed, unless the secrets are set to "null":
{"secrets": {"my-secret": null}}.

The contents of the specified file should be formatted to match the Elastic Cloud
API "secrets" object.
`

const updateExample = `# Set credentials for a GCS snapshot repository
$ cat gcs-creds.json
{
    "secrets": {
        "gcs.client.default.credentials_file": {
            "as_file": true,
            "value": {
                "type": "service_account",
                "project_id": "project-id",
                "private_key_id": "key-id",
                "private_key": "-----BEGIN PRIVATE KEY-----\nprivate-key\n-----END PRIVATE KEY-----\n",
                "client_email": "service-account-email",
                "client_id": "client-id",
                "auth_uri": "https://accounts.google.com/o/oauth2/auth",
                "token_uri": "https://accounts.google.com/o/oauth2/token",
                "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
                "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/service-account-email"
            }
        }
    }
}
$ ecctl deployment elasticsearch keystore set --file=gcs-creds.json <Deployment ID>
...
# Set multiple secrets in one playload
$ cat multiple.json
{
    "secrets": {
        "my-secret": {
            "value": "my-value"
        },
        "my-other-secret": {
            "value": "my-other-value"
        }
    }
}
$ ecctl deployment elasticsearch keystore set --file=multiple.json <Deployment ID>
...
`

var (
	errReadingDefPrefix  = "failed reading keystore secret definition"
	errReadingDefMessage = "provide a valid keystore secret definition using the --file flag"
)

var updateCmd = &cobra.Command{
	Use:     "update <deployment id> [--ref-id <ref-id>] {--file=<filename>.json}",
	Long:    updateLong,
	Example: updateExample,
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
	updateCmd.Flags().StringP("file", "f", "", "Required json formatted file path with the keystore secret contents.")
	updateCmd.MarkFlagFilename("file", "json")
}
