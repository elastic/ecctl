## ecctl deployment resource delete

Deletes a previously shut down deployment resource

### Synopsis

Deletes a previously shut down deployment resource

```
ecctl deployment resource delete <deployment id> --kind <kind> --ref-id <ref-id> [flags]
```

### Options

```
  -h, --help            help for delete
      --kind string     Required stateless deployment resource kind (apm, appsearch, kibana)
      --ref-id string   Optional deployment RefId, auto-discovered if not specified
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

* [ecctl deployment resource](ecctl_deployment_resource.md)	 - Manages deployment resources

