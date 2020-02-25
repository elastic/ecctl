## ecctl deployment appsearch

Manages AppSearch deployments

### Synopsis

Manages AppSearch deployments

```
ecctl deployment appsearch [flags]
```

### Options

```
  -h, --help   help for appsearch
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
* [ecctl deployment appsearch create](ecctl_deployment_appsearch_create.md)	 - Creates an AppSearch instance
* [ecctl deployment appsearch delete](ecctl_deployment_appsearch_delete.md)	 - Deletes a previously shut down AppSearch deployment resource
* [ecctl deployment appsearch show](ecctl_deployment_appsearch_show.md)	 - Shows the specified AppSearch deployment
* [ecctl deployment appsearch shutdown](ecctl_deployment_appsearch_shutdown.md)	 - Shuts down an AppSearch deployment
* [ecctl deployment appsearch upgrade](ecctl_deployment_appsearch_upgrade.md)	 - Upgrades an AppSearch deployment to the Elasticsearch deployment version

