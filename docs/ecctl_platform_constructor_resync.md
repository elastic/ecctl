## ecctl platform constructor resync

Resynchronizes the search index and cache for the selected constructor or all (Available for ECE only)

### Synopsis

Resynchronizes the search index and cache for the selected constructor or all (Available for ECE only)

```
ecctl platform constructor resync {<constructor id> | --all} [flags]
```

### Options

```
      --all    Resynchronizes the search index for all constructors
  -h, --help   help for resync
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

* [ecctl platform constructor](ecctl_platform_constructor.md)	 - Manages constructors (Available for ECE only)

