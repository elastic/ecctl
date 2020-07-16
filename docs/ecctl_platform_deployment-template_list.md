## ecctl platform deployment-template list

DEPRECATED (Will be removed in the next major version): Lists the platform deployment templates

### Synopsis

DEPRECATED (Will be removed in the next major version): Lists the platform deployment templates

```
ecctl platform deployment-template list [flags]
```

### Options

```
      --filter string                  Optional key:value pair that acts as a filter and excludes any template that does not have a matching metadata item associated.
  -h, --help                           help for list
      --show-instance-configurations   Shows instance configurations - only visible when using the JSON output
      --stack-version string           If present, it will cause the returned deployment templates to be adapted to return only the elements allowed in that version.
      --template-format string         If 'deployment' is specified, the deployment_template is populated in the response. If 'cluster' is specified, the cluster_template is populated in the response. (default "cluster")
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

* [ecctl platform deployment-template](ecctl_platform_deployment-template.md)	 - DEPRECATED (Will be removed in the next major version): Manages deployment templates (Available for ECE only)

