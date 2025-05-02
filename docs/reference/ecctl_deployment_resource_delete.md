---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_deployment_resource_delete.html
---

# ecctl deployment resource delete [ecctl_deployment_resource_delete]

Deletes a previously shut down deployment resource

```
ecctl deployment resource delete <deployment id> --kind <kind> --ref-id <ref-id> [flags]
```


## Options [_options_31]

```
  -h, --help            help for delete
      --kind string     Required stateless deployment resource kind (apm, appsearch, kibana)
      --ref-id string   Optional deployment RefId, auto-discovered if not specified
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_30]

:::{include} /_snippets/inherited-options.md
:::


## SEE ALSO [_see_also_31]

* [ecctl deployment resource](/reference/ecctl_deployment_resource.md)	 - Manages deployment resources

