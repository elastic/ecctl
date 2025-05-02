---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_platform_proxy_filtered-group_update.html
---

# ecctl platform proxy filtered-group update [ecctl_platform_proxy_filtered-group_update]

Updates proxies filtered group ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)

```
ecctl platform proxy filtered-group update <filtered group id> --filters <key1=value1,key2=value2> --expected-proxies-count <int> --version <int> [flags]
```


## Options [_options_98]

```
      --expected-proxies-count int32   Expected proxies count in filtered group (default 1)
      --filters stringToString         Filters for proxies group (default [])
  -h, --help                           help for update
      --version string                 Version update for filtered group
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_97]

:::{include} /_snippets/inherited-options.md
:::


## SEE ALSO [_see_also_98]

* [ecctl platform proxy filtered-group](/reference/ecctl_platform_proxy_filtered-group.md)	 - Manages proxies filtered group ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)

