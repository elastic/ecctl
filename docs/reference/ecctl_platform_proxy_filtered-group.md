---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_platform_proxy_filtered-group.html
applies_to:
  deployment:
    ece: all
---

# ecctl platform proxy filtered-group [ecctl_platform_proxy_filtered-group]

Manages proxies filtered group.

```
ecctl platform proxy filtered-group [flags]
```


## Options [_options_93]

```
  -h, --help   help for filtered-group
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_92]

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


## SEE ALSO [_see_also_93]

* [ecctl platform proxy](/reference/ecctl_platform_proxy.md) - Manages proxies
* [ecctl platform proxy filtered-group create](/reference/ecctl_platform_proxy_filtered-group_create.md) - Creates proxies filtered group
* [ecctl platform proxy filtered-group delete](/reference/ecctl_platform_proxy_filtered-group_delete.md) - Deletes proxies filtered group
* [ecctl platform proxy filtered-group list](/reference/ecctl_platform_proxy_filtered-group_list.md) - Returns all proxies filtered groups in the platform
* [ecctl platform proxy filtered-group show](/reference/ecctl_platform_proxy_filtered-group_show.md) - Shows details for proxies filtered group
* [ecctl platform proxy filtered-group update](/reference/ecctl_platform_proxy_filtered-group_update.md) - Updates proxies filtered group