---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_comment_update.html
applies_to:
  deployment:
    ece: all
---

# ecctl comment update [ecctl_comment_update]

Updates an existing resource comment.

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

:::{include} _snippets/inherited-options.md
:::


## See also [_see_also_13]

* [ecctl comment](/reference/ecctl_comment.md)	 - Manages resource comments
