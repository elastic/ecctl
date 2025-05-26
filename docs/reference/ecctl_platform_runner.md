---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_platform_runner.html
applies_to:
  deployment:
    ece: all
---

# ecctl platform runner [ecctl_platform_runner]

Manages platform runners.

```
ecctl platform runner [flags]
```


## Options [_options_115]

```
  -h, --help   help for runner
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_114]

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


## SEE ALSO [_see_also_115]

* [ecctl platform](/reference/ecctl_platform.md) - Manages the platform
* [ecctl platform runner list](/reference/ecctl_platform_runner_list.md) - Lists the existing platform runners
* [ecctl platform runner resync](/reference/ecctl_platform_runner_resync.md) - Resynchronizes the search index and cache for the selected runner or all
* [ecctl platform runner search](/reference/ecctl_platform_runner_search.md) - Performs advanced runner searching
* [ecctl platform runner show](/reference/ecctl_platform_runner_show.md) - Shows information about the specified runner 

