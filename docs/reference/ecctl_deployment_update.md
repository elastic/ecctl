---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_deployment_update.html
applies_to:
  deployment:
    ess: all
    ece: all
---

# ecctl deployment update [ecctl_deployment_update]

Updates a deployment from a file definition, allowing certain flag overrides.


## Synopsis [_synopsis_7]

updates a deployment from a file definition: Defaulting prune_orphans=false, making the default update action safe for partial updates. To override this behavior, toggle --prune-orphans. To track the changes toggle the --track flag.

It is possible to use "--generate-update-payload" as: "ecctl deployment show --generate-update-payload …​" to obtain a valid update payload from an existing deployment which can be manually edited in before it is sent as a "--file" flag argument. See "ecctl deployment show --help" for a valid example.

Read more about the deployment definition in [ECE Deployment APIs](https://www.elastic.co/docs/api/doc/cloud-enterprise/group/endpoint-deployments).

```
ecctl deployment update -f <file definition.json> [flags]
```


## Examples [_examples_8]

```
## Same base deployment as the create example, changing cluster_topology[0].zone_count to 3.
$ cat deployment_example_update.json
{
    "resources": {
        "elasticsearch": [
            {
                "display_name": "my elasticsearch cluster",
                "ref_id": "my-es-cluster",
                "plan": {
                    "deployment_template": {
                        "id": "default"
                    },
                    "elasticsearch": {
                        "version": "6.8.4"
                    },
                    "cluster_topology": [
                        {
                            "instance_configuration_id": "data.default",
                            "memory_per_node": 1024,
                            "node_count_per_zone": 1,
                            "node_type": {
                                "data": true,
                                "ingest": true,
                                "master": true,
                                "ml": false
                            },
                            "zone_count": 3
                        }
                    ]
                }
            }
        ]
    }
}
$ ecctl deployment update f44c06c3af6f85dac05023cf243f4ab1 -f deployment_example_update.json
...
## Setting --prune-orphans, will cause any non-specified resources to be shut down.
$ ecctl deployment update f44c06c3af6f85dac05023cf243f4ab1 -f deployment_example_update.json --prune-orphans
setting --prune-orphans to "true" will cause any resources not specified in the update request to be removed from the deployment, do you want to continue? [y/n]: y
...
```


## Options [_options_59]

```
  -f, --file string           Partial (default) or full JSON file deployment update payload
  -h, --help                  help for update
      --hide-pruned-orphans   Hides orphaned resources that were shut down (only relevant if --prune-orphans=true)
      --prune-orphans         When set to true, it will remove any resources not specified in the update request, treating the json file contents as the authoritative deployment definition
      --skip-snapshot         Skips taking an Elasticsearch snapshot prior to shutting down the deployment
  -t, --track                 Tracks the progress of the performed task
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_58]

:::{include} _snippets/inherited-options.md
:::


## See also [_see_also_59]

* [ecctl deployment](/reference/ecctl_deployment.md)	 - Manages deployments

