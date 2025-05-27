---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_platform_instance-configuration_show.html
applies_to:
  deployment:
    ece: all
---

# ecctl platform instance-configuration show [ecctl_platform_instance-configuration_show]

Shows an instance configuration.

```
ecctl platform instance-configuration show <config id> [flags]
```


## Options [_options_90]

```
  -v, --config-version string   Instance configuration version
  -h, --help                    help for show
      --show-deleted            If set to true, allows to show deleted instance configurations
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_89]

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


## See also [_see_also_90]

* [ecctl platform instance-configuration](/reference/ecctl_platform_instance-configuration.md)	 - Manages instance configurations

