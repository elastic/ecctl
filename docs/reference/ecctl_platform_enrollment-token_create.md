---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_platform_enrollment-token_create.html
applies_to:
  deployment:
    ess: all
    ece: all
---

# ecctl platform enrollment-token create [ecctl_platform_enrollment-token_create]

Creates an enrollment token for roles.

```
ecctl platform enrollment-token create --role <ROLE> [flags]
```


## Examples [_examples_10]

```
  ecctl [globalFlags] enrollment-token create --role coordinator
  ecctl [globalFlags] enrollment-token create --role coordinator --role proxy
  ecctl [globalFlags] enrollment-token create --role allocator --validity 120s
  ecctl [globalFlags] enrollment-token create --role allocator --validity 2h {ece-icon} (Available for ECE only)
```


## Options [_options_81]

```
  -h, --help                help for create
  -r, --role stringArray    Role(s) to associate the tokens with.
  -v, --validity duration   Time in seconds for which this token is valid. Currently this will make the token ephemeral (persistent: false)
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_80]

:::{include} _snippets/inherited-options.md
:::


## See also [_see_also_81]

* [ecctl platform enrollment-token](/reference/ecctl_platform_enrollment-token.md)	 - Manages tokens
