---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_comment_delete.html
---

# ecctl comment delete [ecctl_comment_delete]

Deletes a resource comment ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)

```
ecctl comment delete <comment id> --resource-type <resource-type> --resource-id <resource-id> [flags]
```


## Options [_options_10]

```
  -h, --help                   help for delete
      --resource-id string     ID of the resource that the comment belongs to.
      --resource-type string   The kind of resource that a comment belongs to. Should be one of [elasticsearch, kibana, apm, appsearch, enterprise_search, allocator, constructor, runner, proxy].
      --version string         If specified then checks for conflicts against the version stored in the persistent store.
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_9]

:::{include} /_snippets/inherited-options.md
:::


## SEE ALSO [_see_also_10]

* [ecctl comment](/reference/ecctl_comment.md)	 - Manages resource comments ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)

