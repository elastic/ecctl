## ecctl platform allocator

Manages allocators (Available for ECE only)

### Synopsis

Manages allocators (Available for ECE only)

```
ecctl platform allocator [flags]
```

### Options

```
  -h, --help   help for allocator
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

* [ecctl platform](ecctl_platform.md)	 - Manages the platform (Available for ECE only)
* [ecctl platform allocator list](ecctl_platform_allocator_list.md)	 - Returns all allocators that have instances or are connected to the platform. Use --all flag or --output json to show all. Use --query to match any of the allocators properties. (Available for ECE only)
* [ecctl platform allocator maintenance](ecctl_platform_allocator_maintenance.md)	 - Sets the allocator in Maintenance mode (Available for ECE only)
* [ecctl platform allocator metadata](ecctl_platform_allocator_metadata.md)	 - Manages an allocator's metadata
* [ecctl platform allocator search](ecctl_platform_allocator_search.md)	 - Performs advanced allocator searching (Available for ECE only)
* [ecctl platform allocator show](ecctl_platform_allocator_show.md)	 - Returns information about the allocator (Available for ECE only)
* [ecctl platform allocator vacate](ecctl_platform_allocator_vacate.md)	 - Moves all the resources from the specified allocator (Available for ECE only)

