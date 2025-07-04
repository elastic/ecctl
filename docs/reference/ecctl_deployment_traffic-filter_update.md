---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_deployment_traffic-filter_update.html
applies_to:
  deployment:
    ess: all
    ece: all
---

# ecctl deployment traffic-filter update [ecctl_deployment_traffic-filter_update]

Updates a network security policy or traffic filter ruleset.

```
ecctl deployment traffic-filter update <traffic-filter id> {--file <file-path> | --generate-payload} [flags]
```


## Examples [_examples_7]

```
* Return the current traffic filter state as a valid update payload.
  ecctl deployment traffic-filter update <traffic-filter id> --generate-payload > update.json

* After editing the file with your new values pass it as an argument to the --file flag.
  ecctl deployment traffic-filter update <traffic-filter id> --file update.json
```


## Options [_options_58]

```
      --file string        Path to the file containing the update JSON definition.
      --generate-payload   Outputs JSON which can be used as an argument for the --file flag.
  -h, --help               help for update
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_57]

:::{include} _snippets/inherited-options.md
:::


## See also [_see_also_58]

* [ecctl deployment traffic-filter](/reference/ecctl_deployment_traffic-filter.md)	 - Manages network security policies or traffic filter rulesets

