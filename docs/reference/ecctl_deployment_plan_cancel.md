---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_deployment_plan_cancel.html
applies_to:
  deployment:
    ess: all
    ece: all
---

# ecctl deployment plan cancel [ecctl_deployment_plan_cancel]

Cancels a resource's pending plan.

```
ecctl deployment plan cancel <deployment id> --kind <kind> [--ref-id <ref-id>] [flags]
```


## Options [_options_29]

```
  -h, --help            help for cancel
      --kind string     Required deployment resource kind (apm, appsearch, kibana, elasticsearch)
      --ref-id string   Optional deployment RefId, if not set, the RefId will be auto-discovered
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_28]

:::{include} _snippets/inherited-options.md
:::


## See also [_see_also_29]

* [ecctl deployment plan](/reference/ecctl_deployment_plan.md)	 - Manages deployment plans

