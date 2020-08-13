## ecctl platform runner

Manages platform runners (Available for ECE only)

### Synopsis

Manages platform runners (Available for ECE only)

```
ecctl platform runner [flags]
```

### Options

```
  -h, --help   help for runner
```

### Options inherited from parent commands

```
      --api-key string     API key to use to authenticate (If empty will look for EC_API_KEY environment variable)
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

* [ecctl platform](ecctl_platform.md)	 - Manages the platform (Available for ECE only)
* [ecctl platform runner list](ecctl_platform_runner_list.md)	 - Lists the existing platform runners (Available for ECE only)
* [ecctl platform runner resync](ecctl_platform_runner_resync.md)	 - Resynchronizes the search index and cache for the selected runner or all (Available for ECE only)
* [ecctl platform runner search](ecctl_platform_runner_search.md)	 - Performs advanced runner searching (Available for ECE only)
* [ecctl platform runner show](ecctl_platform_runner_show.md)	 - Shows information about the specified runner (Available for ECE only)

