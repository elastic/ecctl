## ecctl deployment apm

Manages APM deployments

### Synopsis

Manages APM deployments

```
ecctl deployment apm [flags]
```

### Options

```
  -h, --help   help for apm
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

* [ecctl deployment](ecctl_deployment.md)	 - Manages deployments
* [ecctl deployment apm create](ecctl_deployment_apm_create.md)	 - Creates an Apm instance
* [ecctl deployment apm delete](ecctl_deployment_apm_delete.md)	 - Deletes an APM deployment
* [ecctl deployment apm list](ecctl_deployment_apm_list.md)	 - Lists the APM deployments
* [ecctl deployment apm plan](ecctl_deployment_apm_plan.md)	 - Manages APM plans
* [ecctl deployment apm restart](ecctl_deployment_apm_restart.md)	 - Restarts an APM deployment
* [ecctl deployment apm resync](ecctl_deployment_apm_resync.md)	 - Resynchronizes the search index and cache for the selected APM deployment
* [ecctl deployment apm show](ecctl_deployment_apm_show.md)	 - Shows the specified APM deployment
* [ecctl deployment apm stop](ecctl_deployment_apm_stop.md)	 - Stops an APM deployment
* [ecctl deployment apm upgrade](ecctl_deployment_apm_upgrade.md)	 - Upgrades an APM deployment to the Elasticsearch deployment version

