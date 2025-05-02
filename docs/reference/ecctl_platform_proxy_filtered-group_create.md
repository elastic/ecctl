---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_platform_proxy_filtered-group_create.html
---

# ecctl platform proxy filtered-group create [ecctl_platform_proxy_filtered-group_create]

Creates proxies filtered group ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)

```
ecctl platform proxy filtered-group create <filtered group id> --filters <key1=value1,key2=value2> --expected-proxies-count <int> [flags]
```


## Options [_options_94]

```
      --expected-proxies-count int32   Expected proxies count in filtered group (default 1)
      --filters stringToString         Filters for proxies group (default [])
  -h, --help                           help for create
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_93]

:::{include} /_snippets/inherited-options.md
:::


## SEE ALSO [_see_also_94]

* [ecctl platform proxy filtered-group](/reference/ecctl_platform_proxy_filtered-group.md)	 - Manages proxies filtered group ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)

