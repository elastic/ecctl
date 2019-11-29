## ecctl deployment elasticsearch create

Creates a deployment with (only) an Elasticsearch resource in it

### Synopsis

Creates an Elasticsearch deployment, limitting the creation scope to Elasticsearch resources.
There's a few ways to create an Elasticsearch deployment, sane default values are provided, making
the command work out of the box even when no parameters are set. When version is not specified,
the latest available stack version will automatically be used. These are the available options:

  * Simplified flags: --zones <zone count> --size <node memory in MB>
  * Advanced flags: --topology-element <json obj> (shorthand: -e).
    Note that the flag can be specified multiple times for complex topologies.
    The JSON object has the following format:
    {
      "name": "["data", "master", "ml"]" # type string
      "size": 1024 # type int32
      "zone_count": 1 # type int32
    }
  * File definition: --file=<file path> (shorthand: -f). The definition can be found in:
    https://www.elastic.co/guide/en/cloud-enterprise/current/definitions.html#ElasticsearchPayload

As an option "--generate-payload" can be used in order to obtain the generated ElasticsearchPayload
that would be sent as a request, save it, update or extend the topology and create an Elasticsearch
deployment using the saved payload with the "--file" flag.

```
ecctl deployment elasticsearch create {--file|--size <int> --version <string> --zones <string>|--topology-element <obj>} [flags]
```

### Examples

```
## Create a single node cluster. The command will exit after the API response has been returned, 
## without waiting until the deployment resources have been created. 
$ ecctl deployment elasticsearch create --name example-cluster --size 2048
Obtained latest stack version: 7.4.2
{
  "created": true,
  "id": "439fdd1d1b6e4713b6a86847f5a6a199",
  "name": "example-cluster",
  "resources": [
      "id": "a95f5c474fba482b989f790a1f8475b3",
      "kind": "elasticsearch",
      "ref_id": "elasticsearch",
      "region": "ece-region"
    }
  ]
}

## To make the command wait until the resources have been created use the "--track" flag, which will
## output the current stage on which the deployment resources are in.
$ deployment elasticsearch create --name example-cluster --size 2048 --track
Obtained latest stack version: 7.4.2
[...]
Cluster [6be1a417f8bc408cafead7d9db724886][Elasticsearch]: running step "wait-until-running" (Plan duration 1.348361695s)...
Cluster [6be1a417f8bc408cafead7d9db724886][Elasticsearch]: running step "verify-non-split" (Plan duration 51.296428148s)...
Cluster [6be1a417f8bc408cafead7d9db724886][Elasticsearch]: running step "set-quorum-size" (Plan duration 57.381950576s)...
Cluster [6be1a417f8bc408cafead7d9db724886][Elasticsearch]: running step "set-maintenance" (Plan duration 58.296756321s)...
Cluster [6be1a417f8bc408cafead7d9db724886][Elasticsearch]: running step "apply-hot-warm-default-allocation" (Plan duration 1m3.285855089s)...
Cluster [6be1a417f8bc408cafead7d9db724886][Elasticsearch]: finished running all the plan steps (Total plan duration: 1m4.756486638s)

## Additionally, a more advanced topology can be created through "--topology-element" or shorthand "-e".
$ ecctl deployment elasticsearch create --name my-cluster --topology-element '{"size": 1024, "zone_count": 2, "name": "data"}' --topology-element '{"size": 1024, "zone_count": 1, "name": "ml"}' --generate-payload
Obtained latest stack version: 7.4.2
{
  "plan": {
    "cluster_topology": [
      {
        "instance_configuration_id": "data.default",
        "node_type": {
          "data": true,
          "ingest": true,
          "master": true
        },
        "size": {
          "resource": "memory",
          "value": 1024
        },
        "zone_count": 2
      },
      {
        "instance_configuration_id": "ml",
        "node_type": {
          "data": false,
          "ingest": false,
          "master": false,
          "ml": true
        },
        "size": {
          "resource": "memory",
          "value": 1024
        },
        "zone_count": 1
      }
    ],
    "deployment_template": {
      "id": "default"
    },
    "elasticsearch": {
      "version": "7.4.2"
    }
  },
  "ref_id": "elasticsearch",
  "region": "ece-region"
}

## Save the definition to a file for later use.
$ ecctl deployment elasticsearch create --name my-cluster --size 1024 --track --generate-payload --zones 2 > elasticsearch_create_example.json
Obtained latest stack version: 7.4.2

## Create an Elasticsearch deployment through the file definition and track the progress
$ ecctl deployment elasticsearch create --file elasticsearch_create_example.json --track
[...]
Cluster [9c96d8a0df1c47f8a45cd154fc0e3c83][Elasticsearch]: running step "wait-until-running" (Plan duration 3.165747696s)...
Cluster [9c96d8a0df1c47f8a45cd154fc0e3c83][Elasticsearch]: running step "verify-non-split" (Plan duration 1m2.476847682s)...
Cluster [9c96d8a0df1c47f8a45cd154fc0e3c83][Elasticsearch]: running step "set-quorum-size" (Plan duration 1m7.575588825s)...
Cluster [9c96d8a0df1c47f8a45cd154fc0e3c83][Elasticsearch]: running step "set-maintenance" (Plan duration 1m8.464692293s)...
Cluster [9c96d8a0df1c47f8a45cd154fc0e3c83][Elasticsearch]: running step "apply-hot-warm-default-allocation" (Plan duration 1m13.631385049s)...
Cluster [9c96d8a0df1c47f8a45cd154fc0e3c83][Elasticsearch]: running step "apply-plan-settings" (Plan duration 1m14.335030452s)...
Cluster [9c96d8a0df1c47f8a45cd154fc0e3c83][Elasticsearch]: running step "post-plan-cleanup" (Plan duration 1m15.463009785s)...
Cluster [9c96d8a0df1c47f8a45cd154fc0e3c83][Elasticsearch]: finished running all the plan steps (Total plan duration: 1m16.89854434s)

## Create the deployment piping through the file contents tracking the creation progress
$ cat elasticsearch_create_example.json | dev-cli deployment elasticsearch create --track
[...]
Cluster [7dcaeb621dba4248b6a4efc8080a055c][Elasticsearch]: running step "wait-until-running" (Plan duration 3.955507371s)...
Cluster [7dcaeb621dba4248b6a4efc8080a055c][Elasticsearch]: running step "verify-non-split" (Plan duration 1m2.434546366s)...
Cluster [7dcaeb621dba4248b6a4efc8080a055c][Elasticsearch]: running step "set-quorum-size" (Plan duration 1m7.269306003s)...
Cluster [7dcaeb621dba4248b6a4efc8080a055c][Elasticsearch]: running step "set-maintenance" (Plan duration 1m10.321987044s)...
Cluster [7dcaeb621dba4248b6a4efc8080a055c][Elasticsearch]: running step "apply-hot-warm-default-allocation" (Plan duration 1m15.337019401s)...
Cluster [7dcaeb621dba4248b6a4efc8080a055c][Elasticsearch]: running step "apply-plan-settings" (Plan duration 1m16.346500871s)...
Cluster [7dcaeb621dba4248b6a4efc8080a055c][Elasticsearch]: running step "post-plan-cleanup" (Plan duration 1m18.334419179s)...
Cluster [7dcaeb621dba4248b6a4efc8080a055c][Elasticsearch]: finished running all the plan steps (Total plan duration: 1m20.043525071s)
```

### Options

```
      --deployment-template string     Deployment template ID on which to base the deployment from (default "default")
      --file string                    ElasticsearchPayload file definition. See help for more information
      --generate-payload               Returns the ElasticsearchPayload without actually creating the deployment resources
  -h, --help                           help for create
      --name string                    Optional name for the Elasticsearch deployment
      --plugin strings                 Additional plugins to add to the Elasticsearch deployment
      --ref-id string                  RefId for the Elasticsearch deployment (default "elasticsearch")
      --size int32                     Memory (RAM) in MB that each of the deployment nodes will have (default 4096)
  -e, --topology-element stringArray   Topology element definition. See help for more information
  -t, --track                          Tracks the progress of the performed task
      --version string                 Version to use, if not specified, the latest available stack version will be used
      --zones int32                    Number of zones the deployment will span (default 1)
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
      --timeout duration   Timeout to use on all HTTP calls (default 30s)
      --trace              Enables tracing saves the trace to trace-20060102150405
      --user string        Username to use to authenticate (If empty will look for EC_USER environment variable)
      --verbose            Enable verbose mode
```

### SEE ALSO

* [ecctl deployment elasticsearch](ecctl_deployment_elasticsearch.md)	 - Manages Elasticsearch clusters

