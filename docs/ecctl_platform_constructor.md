## ecctl platform constructor

Manages constructors (Available for ECE only)

### Synopsis

Manages constructors (Available for ECE only)

```
ecctl platform constructor [flags]
```

### Options

```
  -h, --help   help for constructor
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
* [ecctl platform constructor list](ecctl_platform_constructor_list.md)	 - Returns all of the constructors in the platform (Available for ECE only)
* [ecctl platform constructor maintenance](ecctl_platform_constructor_maintenance.md)	 - Sets/un-sets a constructor's maintenance mode (Available for ECE only)
* [ecctl platform constructor resync](ecctl_platform_constructor_resync.md)	 - Resynchronizes the search index and cache for the selected constructor or all (Available for ECE only)
* [ecctl platform constructor show](ecctl_platform_constructor_show.md)	 - Returns information about the constructor with given ID (Available for ECE only)

