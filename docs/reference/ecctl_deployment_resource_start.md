---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_deployment_resource_start.html
---

# ecctl deployment resource start [ecctl_deployment_resource_start]

Starts a previously stopped deployment resource

```
ecctl deployment resource start <deployment id> --kind <kind> [--all|--i <instance-id>,<instance-id>] [flags]
```


## Options [_options_35]

```
      --all                   Starts all instances of a defined resource kind
  -h, --help                  help for start
      --ignore-missing        If set and the specified instance does not exist, then quietly proceed to the next instance
  -i, --instance-id strings   Deployment instance IDs to start (e.g. instance-0000000001)
      --kind string           Required deployment resource kind (apm, appsearch, kibana, elasticsearch)
      --ref-id string         Optional deployment RefId, if not set, the RefId will be auto-discovered
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_34]

:::{include} /_snippets/inherited-options.md
:::


## SEE ALSO [_see_also_35]

* [ecctl deployment resource](/reference/ecctl_deployment_resource.md)	 - Manages deployment resources

