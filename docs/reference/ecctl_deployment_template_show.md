---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_deployment_template_show.html
---

# ecctl deployment template show [ecctl_deployment_template_show]

Displays a deployment template ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)

```
ecctl deployment template show --template-id <template id> [flags]
```


## Options [_options_48]

```
  -h, --help                           help for show
      --hide-instance-configurations   Hides instance configurations - only visible when using the JSON output
      --stack-version string           Optional filter to only return deployment templates which are valid for the specified stack version.
      --template-id string             Required template ID to update.
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_47]

:::{include} /_snippets/inherited-options.md
:::


## SEE ALSO [_see_also_48]

* [ecctl deployment template](/reference/ecctl_deployment_template.md)	 - Interacts with deployment template APIs

