---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_deployment_resource_shutdown.html
---

# ecctl deployment resource shutdown [ecctl_deployment_resource_shutdown]

Shuts down a deployment resource by its kind and ref-id


## Synopsis [_synopsis_3]

Shuts down a deployment resource kind (APM, Appsearch, Elasticsearch, Kibana). Shutting down a kind doesn’t necessarily shut down the deployment itself but rather a specific  resource.

```
ecctl deployment resource shutdown <deployment id> --kind <kind> --ref-id <ref-id> [flags]
```


## Options [_options_33]

```
  -h, --help            help for shutdown
      --hide            Optionally hides the deployment resource from being listed by default
      --kind string     Required deployment resource kind (apm, appsearch, kibana, elasticsearch)
      --ref-id string   Optional deployment RefId, auto-discovered if not specified
      --skip-snapshot   Optional flag to toggle skipping the resource snapshot before shutting it down
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_32]

:::{include} _snippets/inherited-options.md
:::


## SEE ALSO [_see_also_33]

* [ecctl deployment resource](/reference/ecctl_deployment_resource.md)	 - Manages deployment resources

