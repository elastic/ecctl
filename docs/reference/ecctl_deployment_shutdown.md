---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_deployment_shutdown.html
---

# ecctl deployment shutdown [ecctl_deployment_shutdown]

Shuts down a deployment and all of its associated sub-resources

```
ecctl deployment shutdown <deployment-id> [flags]
```


## Options [_options_43]

```
  -h, --help            help for shutdown
      --skip-snapshot   Skips taking an Elasticsearch snapshot prior to shutting down the deployment
  -t, --track           Tracks the progress of the performed task
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_42]

:::{include} /_snippets/inherited-options.md
:::


## SEE ALSO [_see_also_43]

* [ecctl deployment](/reference/ecctl_deployment.md)	 - Manages deployments

