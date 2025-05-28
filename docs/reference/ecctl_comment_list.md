---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_comment_list.html
applies_to:
  deployment:
    ece: all
---

# ecctl comment list [ecctl_comment_list]

Lists all resource comments.

```
ecctl comment list --resource-type <resource-type> --resource-id <resource-id> [flags]
```


## Options [_options_11]

```
  -h, --help                   help for list
      --resource-id string     Id of the resource that a comment belongs to.
      --resource-type string   The kind of resource that a comment belongs to. Should be one of [elasticsearch, kibana, apm, appsearch, enterprise_search, allocator, constructor, runner, proxy].
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_10]

:::{include} _snippets/inherited-options.md
:::


## See also [_see_also_11]

* [ecctl comment](/reference/ecctl_comment.md)	 - Manages resource comments
