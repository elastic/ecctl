## ecctl platform proxy filtered-group

Manages proxies filtered group

### Synopsis

Manages proxies filtered group

```
ecctl platform proxy filtered-group [flags]
```

### Options

```
  -h, --help   help for filtered-group
```

### Options inherited from parent commands

```
      --apikey string      API key to use to authenticate (If empty will look for EC_APIKEY environment variable)
      --config string      Config name, used to have multiple configs in $HOME/.ecctl/<env> (default "config")
      --force              Do not ask for confirmation
      --format string      Formats the output using a Go template
      --host string        Base URL to use
      --insecure           Skips all TLS validation
      --message string     A message to set on cluster operation
      --output string      Output format [text|json] (default "text")
      --pass string        Password to use to authenticate (If empty will look for EC_PASS environment variable)
      --pprof              Enables pprofing and saves the profile to pprof-20060102150405
  -q, --quiet              Suppresses the configuration file used for the run, if any
      --region string      Elasticsearch Service region
      --timeout duration   Timeout to use on all HTTP calls (default 30s)
      --trace              Enables tracing saves the trace to trace-20060102150405
      --user string        Username to use to authenticate (If empty will look for EC_USER environment variable)
      --verbose            Enable verbose mode
```

### SEE ALSO

* [ecctl platform proxy](ecctl_platform_proxy.md)	 - Manages proxies (Available for ECE only)
* [ecctl platform proxy filtered-group create](ecctl_platform_proxy_filtered-group_create.md)	 - Creates proxies filtered group
* [ecctl platform proxy filtered-group delete](ecctl_platform_proxy_filtered-group_delete.md)	 - Deletes proxies filtered group
* [ecctl platform proxy filtered-group list](ecctl_platform_proxy_filtered-group_list.md)	 - Returns all proxies filtered groups in the platform
* [ecctl platform proxy filtered-group show](ecctl_platform_proxy_filtered-group_show.md)	 - Shows details for proxies filtered group
* [ecctl platform proxy filtered-group update](ecctl_platform_proxy_filtered-group_update.md)	 - Updates proxies filtered group

