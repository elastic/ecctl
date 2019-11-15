## ecctl deployment create

Creates a deployment from a file definition, allowing certain flag overrides

### Synopsis

Creates a deployment from a file definition with an automatically generated idempotency token.
On creation failure, please use the displayed idempotency token to retry the cluster creation with --request-id=<token>.

Read more about the deployment definition in https://www.elastic.co/guide/en/cloud-enterprise/current/Deployment_-_CRUD.html

```
ecctl deployment create -f <file definition.json> [flags]
```

### Examples

```
$ cat deployment_example.json
{
    "name": "my example cluster",
    "resources": {
        "apm": [
            {
                "display_name": "my apm instance",
                "ref_id": "my-apm-instance",
                "elasticsearch_cluster_ref_id": "my-es-cluster",
                "plan": {
                    "apm": {
                        "version": "6.8.4"
                    },
                    "cluster_topology": [{
                        "instance_configuration_id": "apm",
                        "size": {
                            "resource": "memory",
                            "value": 512
                        },
                        "zone_count": 1
                    }]
                }
            }
        ],
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
                            "zone_count": 1
                        }
                    ]
                }
            }
        ],
        "kibana": [
            {
                "display_name": "my kibana instance",
                "ref_id": "my-kibana-instance",
                "elasticsearch_cluster_ref_id": "my-es-cluster",
                "plan": {
                    "kibana": {
                        "version": "6.8.4"
                    },
                    "cluster_topology": [
                        {
                            "instance_configuration_id": "kibana",
                            "memory_per_node": 1024,
                            "node_count_per_zone": 1,
                            "zone_count": 1
                        }
                    ]
                }
            }
        ]
    }
}
$ ecctl deployment create -f deployment_example.json --version=7.4.1
[...]

## If th previous deployment creation failed
$ ecctl deployment create -f deployment_example.json --name adeploy --version=7.4.1
The deployment creation returned with an error, please use the displayed idempotency token
to recreate the deployment resources
Idempotency token: GMZPMRrcMYqHdmxjIQkHbdjnhPIeBElcwrHwzVlhGUSMXrEIzVXoBykSVRsKncNb
unknown error (status 500)
$ ecctl deployment create -f deployment_example.json --name adeploy --version=7.4.1 --request-id=GMZPMRrcMYqHdmxjIQkHbdjnhPIeBElcwrHwzVlhGUSMXrEIzVXoBykSVRsKncNb
[...]
```

### Options

```
  -f, --file string         JSON file that contains JSON-style domain-specific deployment definition
  -h, --help                help for create
      --name string         Overrides the deployment name
      --request-id string   Optional idempotency token - Can be found in the Stderr device when a previous deployment creation failed, for more information see the examples in the help command page
      --version string      Overrides all thee deployment's resources to the specified version
```

### Options inherited from parent commands

```
      --apikey string      API key to use to authenticate (If empty will look for EC_APIKEY environment variable)
      --config string      Config name, used to have multiple configs in $HOME/.ecctl/<env> (default "config")
      --force              Do not ask for confirmation
      --format string      Formats the output using a Go template
      --host string        Base URL to use (default "https://api.elastic-cloud.com")
      --insecure           Skips all TLS validation
      --message string     A message to set on cluster operation
      --output string      Output format [text|json] (default "text")
      --pass string        Password to use to authenticate (If empty will look for EC_PASS environment variable)
      --pprof              Enables pprofing and saves the profile to pprof-20060102150405
  -q, --quiet              Suppresses the configuration file used for the run, if any
      --region string      Elastic Cloud region
      --timeout duration   Timeout to use on all HTTP calls (default 30s)
      --trace              Enables tracing saves the trace to trace-20060102150405
      --user string        Username to use to authenticate (If empty will look for EC_USER environment variable)
      --verbose            Enable verbose mode
```

### SEE ALSO

* [ecctl deployment](ecctl_deployment.md)	 - Manages deployments

