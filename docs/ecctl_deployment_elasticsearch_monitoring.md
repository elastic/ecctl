## ecctl deployment elasticsearch monitoring

Manages monitoring for an Elasticsearch cluster

### Synopsis

Manages monitoring for an Elasticsearch cluster

```
ecctl deployment elasticsearch monitoring [flags]
```

### Options

```
  -h, --help   help for monitoring
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

* [ecctl deployment elasticsearch](ecctl_deployment_elasticsearch.md)	 - Manages Elasticsearch clusters
* [ecctl deployment elasticsearch monitoring disable](ecctl_deployment_elasticsearch_monitoring_disable.md)	 - Disables monitoring in the specified cluster
* [ecctl deployment elasticsearch monitoring enable](ecctl_deployment_elasticsearch_monitoring_enable.md)	 - Enables monitoring for the cluster by sending data to a monitoring cluster

