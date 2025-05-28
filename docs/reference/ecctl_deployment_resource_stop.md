---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_deployment_resource_stop.html
applies_to:
  deployment:
    ess: all
    ece: all
---

# ecctl deployment resource stop [ecctl_deployment_resource_stop]

Stops a deployment resource.

```
ecctl deployment resource stop <deployment id> --kind <kind> [--all|--i <instance-id>,<instance-id>] [flags]
```


## Options [_options_37]

```
      --all                   Stops all instances of a defined resource kind
  -h, --help                  help for stop
      --ignore-missing        If set and the specified instance does not exist, then quietly proceed to the next instance
  -i, --instance-id strings   Deployment instance IDs to stop (e.g. instance-0000000001)
      --kind string           Required deployment resource kind (apm, appsearch, kibana, elasticsearch)
      --ref-id string         Optional deployment RefId, if not set, the RefId will be auto-discovered
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_36]

:::{include} _snippets/inherited-options.md
:::


## See also [_see_also_37]

* [ecctl deployment resource](/reference/ecctl_deployment_resource.md)	 - Manages deployment resources

