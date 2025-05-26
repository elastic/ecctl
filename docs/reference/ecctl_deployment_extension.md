---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_deployment_extension.html
applies_to:
  deployment:
    ess: all
    ece: all
---

# ecctl deployment extension [ecctl_deployment_extension]

Manages deployment extensions, such as custom plugins or bundles.

```
ecctl deployment extension [flags]
```


## Options [_options_21]

```
  -h, --help   help for extension
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_20]

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


## SEE ALSO [_see_also_21]

* [ecctl deployment](/reference/ecctl_deployment.md)	 - Manages deployments
* [ecctl deployment extension create](/reference/ecctl_deployment_extension_create.md)	 - Creates an extension
* [ecctl deployment extension delete](/reference/ecctl_deployment_extension_delete.md)	 - Deletes a deployment extension
* [ecctl deployment extension list](/reference/ecctl_deployment_extension_list.md)	 - Lists all deployment extensions
* [ecctl deployment extension show](/reference/ecctl_deployment_extension_show.md)	 - Shows information about a deployment extension
* [ecctl deployment extension update](/reference/ecctl_deployment_extension_update.md)	 - Updates an extension

