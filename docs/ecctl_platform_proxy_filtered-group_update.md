## ecctl platform proxy filtered-group update

Updates proxies filtered group (Available for ECE only)

```
ecctl platform proxy filtered-group update <filtered group id> --filters <key1=value1,key2=value2> --expected-proxies-count <int> --version <int> [flags]
```

### Options

```
      --expected-proxies-count int32   Expected proxies count in filtered group (default 1)
      --filters stringToString         Filters for proxies group (default [])
  -h, --help                           help for update
      --version string                 Version update for filtered group
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
      --verbose-credentials   When set, it will not redact the Authorization headers on the request/response trail
      --verbose-file string   If set, it will write the verbose request/response trail to this file
```

### SEE ALSO

* [ecctl platform proxy filtered-group](ecctl_platform_proxy_filtered-group.md)	 - Manages proxies filtered group (Available for ECE only)

