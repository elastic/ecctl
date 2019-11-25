## ecctl deployment elasticsearch

Manages Elasticsearch clusters

### Synopsis

Manages Elasticsearch clusters

```
ecctl deployment elasticsearch [flags]
```

### Options

```
  -h, --help   help for elasticsearch
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
* [ecctl deployment elasticsearch console](ecctl_deployment_elasticsearch_console.md)	 - Starts an interactive console with the cluster
* [ecctl deployment elasticsearch create](ecctl_deployment_elasticsearch_create.md)	 - Creates an Elasticsearch cluster
* [ecctl deployment elasticsearch delete](ecctl_deployment_elasticsearch_delete.md)	 - Deletes an Elasticsearch cluster
* [ecctl deployment elasticsearch diagnose](ecctl_deployment_elasticsearch_diagnose.md)	 - Generates a diagnostics bundle for the cluster
* [ecctl deployment elasticsearch instances](ecctl_deployment_elasticsearch_instances.md)	 - Manages elasticsearch at the instance level
* [ecctl deployment elasticsearch keystore](ecctl_deployment_elasticsearch_keystore.md)	 - Manages an Elasticsearch cluster's keystore
* [ecctl deployment elasticsearch list](ecctl_deployment_elasticsearch_list.md)	 - Returns the list of Elasticsearch clusters
* [ecctl deployment elasticsearch monitoring](ecctl_deployment_elasticsearch_monitoring.md)	 - Manages monitoring for an Elasticsearch cluster
* [ecctl deployment elasticsearch plan](ecctl_deployment_elasticsearch_plan.md)	 - Manages Elasticsearch plans
* [ecctl deployment elasticsearch reallocate](ecctl_deployment_elasticsearch_reallocate.md)	 - Reallocates the Elasticsearch cluster nodes
* [ecctl deployment elasticsearch restart](ecctl_deployment_elasticsearch_restart.md)	 - Restarts an Elasticsearch cluster
* [ecctl deployment elasticsearch resync](ecctl_deployment_elasticsearch_resync.md)	 - Resynchronizes the search index and cache for the selected Elasticsearch cluster
* [ecctl deployment elasticsearch search](ecctl_deployment_elasticsearch_search.md)	 - Performs advanced clusters searching
* [ecctl deployment elasticsearch show](ecctl_deployment_elasticsearch_show.md)	 - Displays information about the specified cluster
* [ecctl deployment elasticsearch shutdown](ecctl_deployment_elasticsearch_shutdown.md)	 - Shuts down an Elasticsearch cluster, so that it no longer contains any running instances
* [ecctl deployment elasticsearch start](ecctl_deployment_elasticsearch_start.md)	 - Starts a stopped Elasticsearch cluster
* [ecctl deployment elasticsearch upgrade](ecctl_deployment_elasticsearch_upgrade.md)	 - Upgrades the cluster to the specified version

