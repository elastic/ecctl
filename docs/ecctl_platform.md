## ecctl platform

Manages the platform (Available for ECE only)

### Synopsis

Manages the platform (Available for ECE only)

```
ecctl platform [flags]
```

### Options

```
  -h, --help   help for platform
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

* [ecctl](ecctl.md)	 - Elastic Cloud Control
* [ecctl platform allocator](ecctl_platform_allocator.md)	 - Manages allocators (Available for ECE only)
* [ecctl platform constructor](ecctl_platform_constructor.md)	 - Manages constructors (Available for ECE only)
* [ecctl platform enrollment-token](ecctl_platform_enrollment-token.md)	 - Manages tokens (Available for ECE only)
* [ecctl platform info](ecctl_platform_info.md)	 - Shows information about the platform (Available for ECE only)
* [ecctl platform instance-configuration](ecctl_platform_instance-configuration.md)	 - Manages instance configurations (Available for ECE only)
* [ecctl platform proxy](ecctl_platform_proxy.md)	 - Manages proxies (Available for ECE only)
* [ecctl platform repository](ecctl_platform_repository.md)	 - Manages snapshot repositories (Available for ECE only)
* [ecctl platform role](ecctl_platform_role.md)	 - Manages platform roles (Available for ECE only)
* [ecctl platform runner](ecctl_platform_runner.md)	 - Manages platform runners (Available for ECE only)

