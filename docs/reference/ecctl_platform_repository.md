---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_platform_repository.html
---

# ecctl platform repository [ecctl_platform_repository]

Manages snapshot repositories ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)


## Synopsis [_synopsis_10]

Manages snapshot repositories that are used by Elasticsearch clusters to perform snapshot operations.

```
ecctl platform repository [flags]
```


## Options [_options_104]

```
  -h, --help   help for repository
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_103]

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


## SEE ALSO [_see_also_104]

* [ecctl platform](/reference/ecctl_platform.md)	 - Manages the platform ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)
* [ecctl platform repository create](/reference/ecctl_platform_repository_create.md)	 - Creates / updates a snapshot repository ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)
* [ecctl platform repository delete](/reference/ecctl_platform_repository_delete.md)	 - Deletes a snapshot repositories ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)
* [ecctl platform repository list](/reference/ecctl_platform_repository_list.md)	 - Lists all the snapshot repositories ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)
* [ecctl platform repository show](/reference/ecctl_platform_repository_show.md)	 - Obtains a snapshot repository config ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)

