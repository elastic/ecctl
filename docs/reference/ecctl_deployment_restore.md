---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_deployment_restore.html
applies_to:
  deployment:
    ess: all
    ece: all
---

# ecctl deployment restore [ecctl_deployment_restore]

Restores a previously shut down deployment and all of its associated sub-resources.


## Synopsis [_synopsis_5]

Use --restore-snapshot to automatically restore the latest available Elasticsearch snapshot (if applicable)

```
ecctl deployment restore <deployment-id> [flags]
```


## Examples [_examples_4]

```
$ ecctl deployment restore 5c17ad7c8df73206baa54b6e2829d9bc
{
  "id": "5c17ad7c8df73206baa54b6e2829d9bc"
}
```

## Options [_options_39]

```
  -h, --help               help for restore
      --restore-snapshot   Restores snapshots for those resources that allow it (Elasticsearch)
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_38]

:::{include} _snippets/inherited-options.md
:::


## See also [_see_also_39]

* [ecctl deployment](/reference/ecctl_deployment.md)	 - Manages deployments

