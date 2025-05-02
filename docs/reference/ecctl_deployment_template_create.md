---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_deployment_template_create.html
---

# ecctl deployment template create [ecctl_deployment_template_create]

Creates a new deployment template ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)

```
ecctl deployment template create --file <definition.json> [flags]
```


## Options [_options_45]

```
  -f, --file string                    Deployment template definition.
  -h, --help                           help for create
      --hide-instance-configurations   Hides instance configurations - only visible when using the JSON output.
      --template-id string             Optional deployment template ID. Otherwise the deployment template will be created with an auto-generated ID.
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_44]

:::{include} /_snippets/inherited-options.md
:::


## SEE ALSO [_see_also_45]

* [ecctl deployment template](/reference/ecctl_deployment_template.md)	 - Interacts with deployment template APIs

