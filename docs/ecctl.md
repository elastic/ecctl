## ecctl

Elastic Cloud Control

### Synopsis

Elastic Cloud Control

### Options

```
      --api-key string     API key to use to authenticate (If empty will look for EC_API_KEY environment variable)
      --config string      Config name, used to have multiple configs in $HOME/.ecctl/<env> (default "config")
      --force              Do not ask for confirmation
      --format string      Formats the output using a Go template
  -h, --help               help for ecctl
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

* [ecctl auth](ecctl_auth.md)	 - Manages authentication settings
* [ecctl deployment](ecctl_deployment.md)	 - Manages deployments
* [ecctl generate](ecctl_generate.md)	 - Generates completions and docs
* [ecctl init](ecctl_init.md)	 - Creates an initial configuration file.
* [ecctl platform](ecctl_platform.md)	 - Manages the platform (Available for ECE only)
* [ecctl stack](ecctl_stack.md)	 - Manages Elastic StackPacks
* [ecctl user](ecctl_user.md)	 - Manages the platform users (Available for ECE only)
* [ecctl version](ecctl_version.md)	 - Shows ecctl version

