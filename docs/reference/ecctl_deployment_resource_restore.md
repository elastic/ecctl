---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_deployment_resource_restore.html
---

# ecctl deployment resource restore [ecctl_deployment_resource_restore]

Restores a previously shut down deployment resource

```
ecctl deployment resource restore <deployment id> --kind <kind> --ref-id <ref-id> [flags]
```


## Options [_options_32]

```
  -h, --help               help for restore
      --kind string        Required deployment resource kind (apm, appsearch, kibana, elasticsearch)
      --ref-id string      Optional deployment RefId, auto-discovered if not specified
      --restore-snapshot   Optional flag to toggle restoring a snapshot for an Elasticsearch resource. It has no effect for other resources
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_31]

:::{include} /_snippets/inherited-options.md
:::


## SEE ALSO [_see_also_32]

* [ecctl deployment resource](/reference/ecctl_deployment_resource.md)	 - Manages deployment resources

