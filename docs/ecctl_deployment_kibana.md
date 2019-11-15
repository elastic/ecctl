## ecctl deployment kibana

Manages Kibana instances

### Synopsis

Manages Kibana instances

```
ecctl deployment kibana [flags]
```

### Options

```
  -h, --help   help for kibana
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

* [ecctl deployment](ecctl_deployment.md)	 - Manages deployments
* [ecctl deployment kibana create](ecctl_deployment_kibana_create.md)	 - Creates a Kibana instance
* [ecctl deployment kibana delete](ecctl_deployment_kibana_delete.md)	 - Deletes a Kibana instance
* [ecctl deployment kibana enable](ecctl_deployment_kibana_enable.md)	 - Enables a kibana instance in the selected deployment
* [ecctl deployment kibana list](ecctl_deployment_kibana_list.md)	 - Returns the list of clusters for a region
* [ecctl deployment kibana plan](ecctl_deployment_kibana_plan.md)	 - Manages Kibana plans
* [ecctl deployment kibana reallocate](ecctl_deployment_kibana_reallocate.md)	 - Reallocates Kibana instances
* [ecctl deployment kibana restart](ecctl_deployment_kibana_restart.md)	 - Restarts a Kibana instance
* [ecctl deployment kibana resync](ecctl_deployment_kibana_resync.md)	 - Resynchronizes the search index and cache for the selected Kibana instance
* [ecctl deployment kibana show](ecctl_deployment_kibana_show.md)	 - Returns the cluster information
* [ecctl deployment kibana start](ecctl_deployment_kibana_start.md)	 - Starts a Kibana instance
* [ecctl deployment kibana stop](ecctl_deployment_kibana_stop.md)	 - Downscales a Kibana instance
* [ecctl deployment kibana upgrade](ecctl_deployment_kibana_upgrade.md)	 - Upgrades the Kibana instance to the same version as the Elasticsearch one

