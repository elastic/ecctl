---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_project_create.html
applies_to:
  deployment:
    ess: all
---
# ecctl project create [ecctl_project_create]

Creates a serverless project

```
ecctl project create [flags]
```

## Options [_options_68a]

```
  -h, --help            help for create
      --name string     Project name (required)
      --region string   Region ID (e.g. aws-us-east-1) (required)
      --tier string     Product tier (observability: complete/logs_essentials, security: complete/essentials)
      --type string     Project type (elasticsearch/search, observability, security) (required)
```

## Options inherited from parent commands [_options_inherited_from_parent_commands_68a]

:::{include} _snippets/inherited-options.md
:::

## See also [_see_also_68a]

* [ecctl project](/reference/ecctl_project.md)	 - Manages serverless projects

