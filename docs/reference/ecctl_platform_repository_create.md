---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_platform_repository_create.html
applies_to:
  deployment:
    ece: all
---

# ecctl platform repository create [ecctl_platform_repository_create]

Creates or updates a snapshot repository.


## Synopsis [_synopsis_11]

Creates / updates a snapshot repository using a set of settings that can be specified as a (yaml|json) file with the --settings flag.

The available settings to set depend on the the --type flag (default s3). A list with the supported settings for each snapshot can be found in the docs: [Snapshot and restore](docs-content://deploy-manage/tools/snapshot-and-restore.md)

The --type flag can be set to any arbitrary value if it’s differs from "s3". Only the S3 available settings are validated.

```
ecctl platform repository create <repository name> --settings <settings file> [flags]
```


## Examples [_examples_12]

```
ecctl platform repository create my-snapshot-repo --settings settings.yml

ecctl platform repository update my-snapshot-repo --settings settings.yml

ecctl platform repository create custom --type fs --settings settings.yml
```


## Options [_options_105]

```
  -h, --help              help for create
      --settings string   Configuration file for the snapshot repository
      --type string       Repository type that will be configured (default "s3")
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_104]

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


## See also [_see_also_105]

* [ecctl platform repository](/reference/ecctl_platform_repository.md) - Manages snapshot repositories

