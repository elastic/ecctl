## ecctl deployment update

Updates a deployment from a file definition, allowing certain flag overrides

### Synopsis

updates a deployment from a file definition, defaulting prune_orphans=false, making the default
update action safe for partial updates, to override this behavior toggle --prune-orphans.
To track the changes toggle the --track flag.

Read more about the deployment definition in https://www.elastic.co/guide/en/cloud-enterprise/current/Deployment_-_CRUD.html

```
ecctl deployment update -f <file definition.json> [flags]
```

### Examples

```
#### Same base deployment as the create example, changing cluster_topology[0].zone_count to 3.
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
{
  "id": "f44c06c3af6f85dac05023cf243f4ab1",
  "name": "my example cluster",
  "resources": [
    {
      "id": "205745432f6345a4999dd2d77ceb1812",
      "kind": "elasticsearch",
      "ref_id": "my-es-cluster",
      "region": "ece-region"
    },
    {
      "elasticsearch_cluster_ref_id": "my-es-cluster",
      "id": "3617594f01074b76a5ca4f903f9d33ec",
      "kind": "kibana",
      "ref_id": "my-kibana-instance",
      "region": "ece-region"
    },
    {
      "elasticsearch_cluster_ref_id": "my-es-cluster",
      "id": "90c3c9566f454861b0dc935c5c7420d8",
      "kind": "apm",
      "ref_id": "my-apm-instance",
      "region": "ece-region"
    }
  ]
}
#### Setting --prune-orphans, will cause any non-specified resources to be shut down.
$ ecctl deployment update f44c06c3af6f85dac05023cf243f4ab1 -f deployment_example_update.json --prune-orphans
setting --prune-orphans to "true" will cause any resources not specified in the update request to be removed from the deployment, do you want to continue? [y/n]: y
{
  "id": "f44c06c3af6f85dac05023cf243f4ab1",
  "name": "my example cluster",
  "resources": [
    {
      "id": "205745432f6345a4999dd2d77ceb1812",
      "kind": "elasticsearch",
      "ref_id": "my-es-cluster",
      "region": "ece-region"
    }
  ],
  "shutdown_resources": {
    "apm": [
      "90c3c9566f454861b0dc935c5c7420d8"
    ],
    "appsearch": [],
    "elasticsearch": [],
    "kibana": [
      "3617594f01074b76a5ca4f903f9d33ec"
    ]
  }
}
```

### Options

```
  -f, --file string           Partial (default) or full JSON file deployment update payload
  -h, --help                  help for update
      --hide-pruned-orphans   Hides orphaned resources that were shut down (only relevant if --prune-orphans=true)
      --prune-orphans         When set to true, it will remove any resources not specified in the update request, treating the json file contents as the authoritative deployment definition
      --skip-snapshot         Skips taking an Elasticsearch snapshot prior to shutting down the deployment
  -t, --track                 Tracks the progress of the performed task
```

### Options inherited from parent commands

```
      --apikey string      API key to use to authenticate (If empty will look for EC_APIKEY environment variable)
      --config string      Config name, used to have multiple configs in $HOME/.ecctl/<env> (default "config")
      --force              Do not ask for confirmation
      --format string      Formats the output using a Go template
      --host string        Base URL to use
      --insecure           Skips all TLS validation
      --message string     A message to set on cluster operation
      --output string      Output format [text|json] (default "text")
      --pass string        Password to use to authenticate (If empty will look for EC_PASS environment variable)
      --pprof              Enables pprofing and saves the profile to pprof-20060102150405
  -q, --quiet              Suppresses the configuration file used for the run, if any
      --region string      Elasticsearch Service region
      --timeout duration   Timeout to use on all HTTP calls (default 30s)
      --trace              Enables tracing saves the trace to trace-20060102150405
      --user string        Username to use to authenticate (If empty will look for EC_USER environment variable)
      --verbose            Enable verbose mode
```

### SEE ALSO

* [ecctl deployment](ecctl_deployment.md)	 - Manages deployments

