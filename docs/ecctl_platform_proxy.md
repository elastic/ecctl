## ecctl platform proxy

Manages proxies (Available for ECE only)

```
ecctl platform proxy [flags]
```

### Options

```
  -h, --help   help for proxy
```

### Options inherited from parent commands

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

### SEE ALSO

* [ecctl platform](ecctl_platform.md)	 - Manages the platform (Available for ECE only)
* [ecctl platform proxy filtered-group](ecctl_platform_proxy_filtered-group.md)	 - Manages proxies filtered group (Available for ECE only)
* [ecctl platform proxy list](ecctl_platform_proxy_list.md)	 - Returns all of the proxies in the platform (Available for ECE only)
* [ecctl platform proxy settings](ecctl_platform_proxy_settings.md)	 - Manages proxies settings (Available for ECE only)
* [ecctl platform proxy show](ecctl_platform_proxy_show.md)	 - Returns information about the proxy with given id (Available for ECE only)

