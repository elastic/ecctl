## ecctl deployment

Manages deployments

### Synopsis

Manages deployments

```
ecctl deployment [flags]
```

### Options

```
  -h, --help   help for deployment
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

* [ecctl](ecctl.md)	 - Elastic Cloud Control
* [ecctl deployment create](ecctl_deployment_create.md)	 - Creates a deployment
* [ecctl deployment delete](ecctl_deployment_delete.md)	 - Deletes a previously shutdown deployment (Available for ECE only)
* [ecctl deployment elasticsearch](ecctl_deployment_elasticsearch.md)	 - Manages a deployment's elasticsearch resources
* [ecctl deployment list](ecctl_deployment_list.md)	 - Lists the platform's deployments
* [ecctl deployment note](ecctl_deployment_note.md)	 - Manages a deployment's notes
* [ecctl deployment plan](ecctl_deployment_plan.md)	 - Manages deployment plans
* [ecctl deployment resource](ecctl_deployment_resource.md)	 - Manages deployment resources
* [ecctl deployment restore](ecctl_deployment_restore.md)	 - Restores a previously shut down deployment and all of its associated sub-resources
* [ecctl deployment resync](ecctl_deployment_resync.md)	 - Resynchronizes the search index and cache for the selected deployment or all
* [ecctl deployment search](ecctl_deployment_search.md)	 - Performs advanced deployment search using the Elasticsearch Query DSL
* [ecctl deployment show](ecctl_deployment_show.md)	 - Shows the specified deployment resources
* [ecctl deployment shutdown](ecctl_deployment_shutdown.md)	 - Shuts down a deployment and all of its associated sub-resources
* [ecctl deployment update](ecctl_deployment_update.md)	 - Updates a deployment from a file definition, allowing certain flag overrides

