---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_comment_create.html
applies_to:
  deployment:
    ece: all
---

# ecctl comment create [ecctl_comment_create]

Creates a new resource comment.

```
ecctl comment create <message> --resource-type <resource-type> --resource-id <resource-id> [flags]
```


## Options [_options_9]

```
  -h, --help                   help for create
      --resource-id string     ID of the resource that a comment belongs to.
      --resource-type string   The kind of resource that a comment belongs to. Should be one of [elasticsearch, kibana, apm, appsearch, enterprise_search, allocator, constructor, runner, proxy].
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_8]

:::{include} _snippets/inherited-options.md
:::


## See also [_see_also_9]

* [ecctl comment](/reference/ecctl_comment.md)	 - Manages resource comments
