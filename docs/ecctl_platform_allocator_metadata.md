## ecctl platform allocator metadata

Manages an allocator's metadata

### Synopsis

Manages an allocator's metadata

```
ecctl platform allocator metadata [flags]
```

### Options

```
  -h, --help   help for metadata
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

* [ecctl platform allocator](ecctl_platform_allocator.md)	 - Manages allocators (Available for ECE only)
* [ecctl platform allocator metadata delete](ecctl_platform_allocator_metadata_delete.md)	 - Deletes a single metadata item from a given allocators metadata
* [ecctl platform allocator metadata set](ecctl_platform_allocator_metadata_set.md)	 - Sets or updates a single metadata item to a given allocators metadata
* [ecctl platform allocator metadata show](ecctl_platform_allocator_metadata_show.md)	 - Shows allocator metadata

