---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_deployment_traffic-filter_create.html
applies_to:
  deployment:
    ess: all
    ece: all
---

# ecctl deployment traffic-filter create [ecctl_deployment_traffic-filter_create]

Creates network security policies or traffic filter rulesets.

```
ecctl deployment traffic-filter create --region <region> --name <filter name> --type <filter type> --source <filter source>,<filter source>  [flags]
```


## Options [_options_54]

```
      --description string   Optional description for the traffic filter.
  -h, --help                 help for create
      --include-by-default   If set, any future eligible deployments will have this filter applied automatically.
      --name string          Name for the traffic filter.
      --source strings       List of allowed traffic filter sources. Can be IP addresses, CIDR masks, or VPC endpoint IDs
      --type string          Type of traffic filter. Can be one of [ip, vpce])
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_53]

:::{include} _snippets/inherited-options.md
:::


## See also [_see_also_54]

* [ecctl deployment traffic-filter](/reference/ecctl_deployment_traffic-filter.md) - Manages network security policies or traffic filter rulesets

