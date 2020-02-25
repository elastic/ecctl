## ecctl deployment show

Shows the specified deployment resources

### Synopsis

Shows the specified deployment resources

```
ecctl deployment show <deployment-id> [flags]
```

### Examples

```

* Shows kibana resource information from a given deployment:
  ecctl deployment show <deployment-id> --type kibana

* Shows apm resource information from a given deployment with a specified ref-id.
  ecctl deployment show <deployment-id> --type apm --ref-id apm-server
```

### Options

```
  -h, --help            help for show
  -m, --metadata        Shows the deployment metadata
      --plan-defaults   Shows the deployment plan defaults
      --plan-logs       Shows the deployment plan logs
      --plans           Shows the deployment plans
      --ref-id string   Optional deployment type RefId, if not set, the RefId will be auto-discovered
  -s, --settings        Shows the deployment settings
      --type string     Optional deployment resource type (apm, appsearch, kibana, elasticsearch)
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

* [ecctl deployment](ecctl_deployment.md)	 - Manages deployments

