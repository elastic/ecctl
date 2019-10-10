## ecctl deployment elasticsearch plan update

Applies or validates the provided plan and tracks the resulting change attempt

### Synopsis

Applies the provided plan and tracks the resulting change attempt.
If --validate is set, the plan will simply be validated by the API failing if the plan is invalid. 
The plan can either be provided via the --file parameter or piped via stdin. 
The plan must be a valid ElasticsearchClusterPlan and will be validated against the specified cluster by the API.

```
ecctl deployment elasticsearch plan update <cluster id> [flags]
```

### Options

```
  -f, --file string   Provide the location of a file containing the plan JSON
  -h, --help          help for update
      --track         Tracks the progress of the performed task (default true)
  -v, --validate      Only validate the plan
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

* [ecctl deployment elasticsearch plan](ecctl_deployment_elasticsearch_plan.md)	 - Manages Elasticsearch plans

