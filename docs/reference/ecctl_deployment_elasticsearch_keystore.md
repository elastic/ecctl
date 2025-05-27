---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_deployment_elasticsearch_keystore.html
applies_to:
  deployment:
    ess: all
    ece: all
---

# ecctl deployment elasticsearch keystore [ecctl_deployment_elasticsearch_keystore]

Manages Elasticsearch resource keystores.

```
ecctl deployment elasticsearch keystore [flags]
```


## Options [_options_18]

```
  -h, --help   help for keystore
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_17]

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


## See also [_see_also_18]

* [ecctl deployment elasticsearch](/reference/ecctl_deployment_elasticsearch.md)	 - Manages Elasticsearch resources
* [ecctl deployment elasticsearch keystore show](/reference/ecctl_deployment_elasticsearch_keystore_show.md)	 - Shows the settings from the Elasticsearch resource keystore
* [ecctl deployment elasticsearch keystore update](/reference/ecctl_deployment_elasticsearch_keystore_update.md)	 - Updates the contents of an Elasticsearch keystore

