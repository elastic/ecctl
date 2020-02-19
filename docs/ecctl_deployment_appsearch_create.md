## ecctl deployment appsearch create

Creates an AppSearch instance

### Synopsis

Creates an AppSearch deployment, limitting the creation scope to AppSearch resources.
There are a few ways to create an AppSearch deployment, sane default values are provided, making
the command work out of the box even when no parameters are set. When version is not specified,
the matching elasticsearch deployment version will be used. These are the available options:

  * Simplified flags: --zones <zone count> --size <node memory in MB>
  * File definition: --file=<file path> (shorthand: -f). The definition can be found in:
    https://www.elastic.co/guide/en/cloud-enterprise/current/definitions.html#AppSearchPayload

As an option, "--generate-payload" can be used in order to obtain the generated AppSearchPayload
that would be sent as a request, save it, update or extend the topology and create an AppSearch
deployment using the saved payload with the "--file" flag.

```
ecctl deployment appsearch create --id <deployment-id> [flags]
```

### Examples

```
## Create a single AppSearch server. The command will exit after the API response has been returned, 
## without waiting until the deployment resources have been created. To make the command wait until
the resources have been created use the "--track" flag.
$ ecctl deployment appsearch create --id=a57f8b7ce54c4afb90ce3755d1e94000 --track
{
  "id": "a57f8b7ce54c4afb90ce3755d1e94000",
  "name": "a57f8b7ce54c4afb90ce3755d1e94000",
  "resources": [
    {
      "elasticsearch_cluster_ref_id": "elasticsearch",
      "id": "53d104a432a648f68ec76d52ecb521d5",
      "kind": "appsearch",
      "ref_id": "main-appsearch",
      "region": "ece-region"
    },
    {
      "elasticsearch_cluster_ref_id": "elasticsearch",
      "id": "39e4a65fc2b14651b666aaff18a13b8f",
      "kind": "kibana",
      "ref_id": "main-kibana",
      "region": "ece-region"
    },
    {
      "cloud_id": "a57f8b7ce54c4afb90ce3755d1e94000:MTkyLjE2OC40NC4xMC5pcC5lcy5pbzo5MjQzJGQzODIwOWU4ZTYwYzRlYTliY2UzMDc1OThhMTljNGI3JDM5ZTRhNjVmYzJiMTQ2NTFiNjY2YWFmZjE4YTEzYjhm",
      "id": "d38209e8e60c4ea9bce307598a19c4b7",
      "kind": "elasticsearch",
      "ref_id": "main-elasticsearch",
      "region": "ece-region"
    }
  ]
}
Cluster [53d104a432a648f68ec76d52ecb521d5][AppSearch]: running step "wait-until-running" (Plan duration 1.38505959s)...
Cluster [39e4a65fc2b14651b666aaff18a13b8f][Kibana]: finished running all the plan steps (Total plan duration: 1.73493053s)
Cluster [d38209e8e60c4ea9bce307598a19c4b7][Elasticsearch]: finished running all the plan steps (Total plan duration: 1.849794895s)
Cluster [53d104a432a648f68ec76d52ecb521d5][AppSearch]: running step "set-maintenance" (Plan duration 11.162178491s)...
Cluster [53d104a432a648f68ec76d52ecb521d5][AppSearch]: finished running all the plan steps (Total plan duration: 16.677195277s)

## Save the definition to a file for later use.
$ ecctl deployment appsearch create --generate-payload --id a57f8b7ce54c4afb90ce3755d1e94000 --zones 2 --size 2048 > appsearch_create_example.json

## Create the deployment piping through the file contents tracking the creation progress
$ cat appsearch_create_example.json | dev-cli deployment appsearch create --track --id a57f8b7ce54c4afb90ce3755d1e94000
[...]
```

### Options

```
      --deployment-template string    Optional deployment template ID, automatically obtained from the current deployment
      --elasticsearch-ref-id string   Optional Elasticsearch ref ID where the AppSearch deployment will connect to
  -f, --file string                   AppSearchPayload file definition. See help for more information
      --generate-payload              Returns the AppSearchPayload without actually creating the deployment resources
  -h, --help                          help for create
      --id string                     Deployment ID where to create the AppSearch deployment
      --name string                   Optional name to set for the AppSearch deployment (Overrides name if present)
      --ref-id string                 RefId for the AppSearch deployment (default "main-appsearch")
      --size int32                    Memory (RAM) in MB that each of the deployment nodes will have (default 2048)
  -t, --track                         Tracks the progress of the performed task
      --version string                Optional version to use. If not specified, it will default to the deployment's stack version
      --zones int32                   Number of zones the deployment will span (default 1)
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

* [ecctl deployment appsearch](ecctl_deployment_appsearch.md)	 - Manages AppSearch deployments

