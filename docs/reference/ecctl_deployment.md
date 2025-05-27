---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_deployment.html
applies_to:
  deployment:
    ess: all
    ece: all
---

# ecctl deployment [ecctl_deployment]

Manages deployments.

```
ecctl deployment [flags]
```


## Options [_options_14]

```
  -h, --help   help for deployment
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_13]

```
      --api-key string        API key to use to authenticate (If empty will look for EC_API_KEY environment variable)
      --config string         Config name, used to have multiple configs in $HOME/.ecctl/<env> (default "config")
      --force                 Do not ask for confirmation
      --format string         Formats the output using a Go template
      --host string           Base URL to use
      --insecure              Skips all TLS validation
      --message string        A message to set on cluster operation
      --output string         Output format [text|json] (default "text")
      --pass string           Password to use to authenticate (If empty will look for EC_PASS environment variable)
      --pprof                 Enables pprofing and saves the profile to pprof-20060102150405
  -q, --quiet                 Suppresses the configuration file used for the run, if any
      --region string         Elasticsearch Service region
      --timeout duration      Timeout to use on all HTTP calls (default 30s)
      --trace                 Enables tracing saves the trace to trace-20060102150405
      --user string           Username to use to authenticate (If empty will look for EC_USER environment variable)
      --verbose               Enable verbose mode
      --verbose-credentials   When set, Authorization headers on the request/response trail will be displayed as plain text
      --verbose-file string   When set, the verbose request/response trail will be written to the defined file
```


## See also [_see_also_14]

* [ecctl](/reference/ecctl.md)	 - Elastic Cloud Control
* [ecctl deployment create](/reference/ecctl_deployment_create.md)	 - Creates a deployment
* [ecctl deployment delete](/reference/ecctl_deployment_delete.md)	 - Deletes a previously shutdown deployment ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)
* [ecctl deployment elasticsearch](/reference/ecctl_deployment_elasticsearch.md)	 - Manages Elasticsearch resources
* [ecctl deployment extension](/reference/ecctl_deployment_extension.md)	 - Manages deployment extensions, such as custom plugins or bundles
* [ecctl deployment list](/reference/ecctl_deployment_list.md)	 - Lists the platform’s deployments
* [ecctl deployment plan](/reference/ecctl_deployment_plan.md)	 - Manages deployment plans
* [ecctl deployment resource](/reference/ecctl_deployment_resource.md)	 - Manages deployment resources
* [ecctl deployment restore](/reference/ecctl_deployment_restore.md)	 - Restores a previously shut down deployment and all of its associated sub-resources
* [ecctl deployment resync](/reference/ecctl_deployment_resync.md)	 - Resynchronizes the search index and cache for the selected deployment or all
* [ecctl deployment search](/reference/ecctl_deployment_search.md)	 - Performs advanced deployment search using the Elasticsearch Query DSL
* [ecctl deployment show](/reference/ecctl_deployment_show.md)	 - Shows the specified deployment resources
* [ecctl deployment shutdown](/reference/ecctl_deployment_shutdown.md)	 - Shuts down a deployment and all of its associated sub-resources
* [ecctl deployment template](/reference/ecctl_deployment_template.md)	 - Interacts with deployment template APIs
* [ecctl deployment traffic-filter](/reference/ecctl_deployment_traffic-filter.md)	 - Manages traffic filter rulesets
* [ecctl deployment update](/reference/ecctl_deployment_update.md)	 - Updates a deployment from a file definition, allowing certain flag overrides

