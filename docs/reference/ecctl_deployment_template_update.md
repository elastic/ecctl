---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_deployment_template_update.html
applies_to:
  deployment:
    ece: all
---

# ecctl deployment template update [ecctl_deployment_template_update]

Updates an existing deployment template.

```
ecctl deployment template update --template-id <template id> --file <definition.json> [flags]
```


## Options [_options_49]

```
  -f, --file string                    Deployment template definition.
  -h, --help                           help for update
      --hide-instance-configurations   Hides instance configurations - only visible when using the JSON output.
      --template-id string             Required template ID to update.
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_48]

:::{include} _snippets/inherited-options.md
:::


## See also [_see_also_49]

* [ecctl deployment template](/reference/ecctl_deployment_template.md)	 - Interacts with deployment template APIs

