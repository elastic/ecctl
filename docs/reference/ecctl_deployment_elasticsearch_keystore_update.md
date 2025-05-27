---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_deployment_elasticsearch_keystore_update.html
applies_to:
  deployment:
    ece: all
---

# ecctl deployment elasticsearch keystore update [ecctl_deployment_elasticsearch_keystore_update]

Updates the contents of an Elasticsearch keystore.


## Synopsis [_synopsis_2]

Changes the contents of the Elasticsearch resource keystore from the specified deployment by using the PATCH method. The payload is a partial payload where any omitted current keystore items are not removed, unless the secrets are set to "null": {"secrets": {"my-secret": null}}.

The contents of the specified file should be formatted to match the Elasticsearch Service API "KeystoreContents" model.

```
ecctl deployment elasticsearch keystore update <deployment id> [--ref-id <ref-id>] {--file=<filename>.json} [flags]
```


## Examples [_examples_2]

```
# Set credentials for a GCS snapshot repository
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
```


## Options [_options_20]

```
  -f, --file string     Required json formatted file path with the keystore secret contents.
  -h, --help            help for update
      --ref-id string   Optional ref_id to use for the Elasticsearch resource, auto-discovered if not specified.
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_19]

```
      --api-key string        API key to use to authenticate (If empty will look for EC_API_KEY environment variable)
      --config string         Config name, used to have multiple configs in $HOME/.ecctl/<env> (default "config")
      --force                 Do not ask for confirmation
      --format string         Formats the output using a Go template
      --host string           Base URL to use
      --insecure              Skips all TLS validation
      --message string        A message to set on cluster operation
      --output string         Output format [text|json] (default "text")
      --pass string           Password to use to authenticate (If empty will look for EC_PASS environment variable)
      --pprof                 Enables pprofing and saves the profile to pprof-20060102150405
  -q, --quiet                 Suppresses the configuration file used for the run, if any
      --region string         Elasticsearch Service region
      --timeout duration      Timeout to use on all HTTP calls (default 30s)
      --trace                 Enables tracing saves the trace to trace-20060102150405
      --user string           Username to use to authenticate (If empty will look for EC_USER environment variable)
      --verbose               Enable verbose mode
      --verbose-credentials   When set, Authorization headers on the request/response trail will be displayed as plain text
      --verbose-file string   When set, the verbose request/response trail will be written to the defined file
```


## See also [_see_also_20]

* [ecctl deployment elasticsearch keystore](/reference/ecctl_deployment_elasticsearch_keystore.md)	 - Manages Elasticsearch resource keystores
