## ecctl platform

Manages the platform

### Synopsis

Manages the platform

```
ecctl platform [flags]
```

### Options

```
  -h, --help   help for platform
```

### Options inherited from parent commands

```
      --apikey string      API key to use to authenticate (If empty will look for EC_APIKEY environment variable)
      --config string      Config name, used to have multiple configs in $HOME/.ecctl/<env> (default "config")
      --force              Do not ask for confirmation
      --format string      Formats the output using a Go template
      --host string        Base URL to use (default "https://api.elastic-cloud.com")
      --insecure           Skips all TLS validation
      --message string     A message to set on cluster operation
      --output string      Output format [text|json] (default "text")
      --pass string        Password to use to authenticate (If empty will look for EC_PASS environment variable)
      --pprof              Enables pprofing and saves the profile to pprof-20060102150405
  -q, --quiet              Suppresses the configuration file used for the run, if any
      --region string      Elastic Cloud region
      --timeout duration   Timeout to use on all HTTP calls (default 30s)
      --trace              Enables tracing saves the trace to trace-20060102150405
      --user string        Username to use to authenticate (If empty will look for EC_USER environment variable)
      --verbose            Enable verbose mode
```

### SEE ALSO

* [ecctl](ecctl.md)	 - Elastic Cloud Control
* [ecctl platform allocator](ecctl_platform_allocator.md)	 - Manages allocators
* [ecctl platform allocators](ecctl_platform_allocators.md)	 - Returns all allocators that have instances or are connected in the region. Use --all flag or --output json to show all. Use --query to match any of the allocators properties.
* [ecctl platform constructor](ecctl_platform_constructor.md)	 - Manages constructors
* [ecctl platform constructors](ecctl_platform_constructors.md)	 - Returns all of the constructors in the region
* [ecctl platform deployment-template](ecctl_platform_deployment-template.md)	 - Manages deployment templates
* [ecctl platform enrollment-token](ecctl_platform_enrollment-token.md)	 - Manages tokens
* [ecctl platform enrollment-tokens](ecctl_platform_enrollment-tokens.md)	 - Retrieves a list of persistent enrollment tokens
* [ecctl platform info](ecctl_platform_info.md)	 - Shows information about the platform
* [ecctl platform instance-configuration](ecctl_platform_instance-configuration.md)	 - Manages instance configurations
* [ecctl platform proxies](ecctl_platform_proxies.md)	 - Returns all of the proxies in the region
* [ecctl platform proxy](ecctl_platform_proxy.md)	 - Manages proxies
* [ecctl platform repository](ecctl_platform_repository.md)	 - Manages snapshot repositories
* [ecctl platform role](ecctl_platform_role.md)	 - Manages platform roles
* [ecctl platform stack](ecctl_platform_stack.md)	 - Manages Elastic StackPacks
* [ecctl platform stacks](ecctl_platform_stacks.md)	 - Lists Elastic StackPacks

