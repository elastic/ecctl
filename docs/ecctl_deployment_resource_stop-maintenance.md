## ecctl deployment resource stop-maintenance

Stops maintenance mode on a deployment resource

### Synopsis

Stops maintenance mode on a deployment resource

```
ecctl deployment resource stop-maintenance <deployment id> --type <type> [--all|--i <instance-id>,<instance-id>] [flags]
```

### Options

```
      --all                   Stops maintenace mode on all instances of a defined resource type
  -h, --help                  help for stop-maintenance
      --ignore-missing        If set and the specified instance does not exist, then quietly proceed to the next instance
  -i, --instance-id strings   Deployment instance IDs to use (e.g. instance-0000000001)
      --ref-id string         Optional deployment RefId, if not set, the RefId will be auto-discovered
      --type string           Deployment resource type to use (elasticsearch, kibana, apm, or appsearch)
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
      --timeout duration   Timeout to use on all HTTP calls (default 30s)
      --trace              Enables tracing saves the trace to trace-20060102150405
      --user string        Username to use to authenticate (If empty will look for EC_USER environment variable)
      --verbose            Enable verbose mode
```

### SEE ALSO

* [ecctl deployment resource](ecctl_deployment_resource.md)	 - Manages deployment resources

