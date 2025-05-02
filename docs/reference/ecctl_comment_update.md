---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_comment_update.html
---

# ecctl comment update [ecctl_comment_update]

Updates an existing resource comment ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)

```
ecctl comment update <comment id> <message> --resource-type <resource-type> --resource-id <resource-id> [flags]
```


## Options [_options_13]

```
  -h, --help                   help for update
      --resource-id string     Id of the Resource that a Comment belongs to.
      --resource-type string   The kind of Resource that a Comment belongs to. Should be one of [elasticsearch, kibana, apm, appsearch, enterprise_search, allocator, constructor, runner, proxy].
      --version string         If specified then checks for conflicts against the version stored in the persistent store.
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_12]

:::{include} /_snippets/inherited-options.md
:::


## SEE ALSO [_see_also_13]

* [ecctl comment](/reference/ecctl_comment.md)	 - Manages resource comments ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)

