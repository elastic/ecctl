## ecctl deployment apm list

Lists the APM deployments

### Synopsis

Lists the APM deployments

```
ecctl deployment apm list [flags]
```

### Options

```
  -h, --help            help for list
      --hidden          Shows hidden deployments
  -m, --metadata        Shows deployment metadata
      --plan-defaults   Shows deployment plan defaults
      --plan-logs       Shows deployment plan logs
      --plans           Shows deployment plans
      --query string    Custom Elasticsearch query to filter deployment
  -s, --size int        Limits the number of deployments to the size value (default 100)
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

* [ecctl deployment apm](ecctl_deployment_apm.md)	 - Manages APM deployments

