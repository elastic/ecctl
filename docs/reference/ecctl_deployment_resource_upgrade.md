---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_deployment_resource_upgrade.html
---

# ecctl deployment resource upgrade [ecctl_deployment_resource_upgrade]

Upgrades a deployment resource


## Synopsis [_synopsis_4]

Upgrades a stateless deployment resource so it matches the Elasticsearch deployment version. Only stateless resources are supported in the --kind flag

```
ecctl deployment resource upgrade <deployment id> --kind <kind> --ref-id <ref-id> [flags]
```


## Options [_options_38]

```
  -h, --help            help for upgrade
      --kind string     Required deployment resource kind (apm, appsearch, kibana, elasticsearch)
      --ref-id string   Optional deployment RefId, if not set, the RefId will be auto-discovered
  -t, --track           Tracks the progress of the performed task
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_37]

:::{include} _snippets/inherited-options.md
:::


## SEE ALSO [_see_also_38]

* [ecctl deployment resource](/reference/ecctl_deployment_resource.md)	 - Manages deployment resources

