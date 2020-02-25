## ecctl deployment elasticsearch plan

Manages Elasticsearch plans

### Synopsis

Manages Elasticsearch plans

```
ecctl deployment elasticsearch plan [flags]
```

### Options

```
  -h, --help   help for plan
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

* [ecctl deployment elasticsearch](ecctl_deployment_elasticsearch.md)	 - Manages Elasticsearch clusters
* [ecctl deployment elasticsearch plan cancel](ecctl_deployment_elasticsearch_plan_cancel.md)	 - Cancels the pending plan
* [ecctl deployment elasticsearch plan history](ecctl_deployment_elasticsearch_plan_history.md)	 - Lists the plan history
* [ecctl deployment elasticsearch plan monitor](ecctl_deployment_elasticsearch_plan_monitor.md)	 - Monitors the pending plan
* [ecctl deployment elasticsearch plan reapply](ecctl_deployment_elasticsearch_plan_reapply.md)	 - Reapplies the latest plan attempt, resetting all transient settings
* [ecctl deployment elasticsearch plan update](ecctl_deployment_elasticsearch_plan_update.md)	 - Applies or validates the provided plan and tracks the resulting change attempt

