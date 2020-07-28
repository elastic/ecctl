## ecctl deployment create

Creates a deployment

### Synopsis

Creates a deployment which can be defined through flags or from a file definition.
Sane default values are provided, making the command work out of the box even when no parameters are set. 
When version is not specified, the latest available stack version will automatically be used. 
These are the available options:

  * Simplified flags to set size and zone count for each instance type (Elasticsearch, Kibana, APM, Enterprise Search and App Search). 
  * Advanced flag for different Elasticsearch node types: --es-node-topology <json obj> (shorthand: -e).
    Note that the flag can be specified multiple times for complex topologies.
    The JSON object has the following format:
    {
      "node_type": "["data", "master", "ml"]" # type string
      "size": "1g" # type string
      "zone_count": 1 # type int32
    }
  * File definition: --file=<file path> (shorthand: -f). You can create a definition by using the sample JSON seen here:
    https://elastic.co/guide/en/cloud/current/ec-api-deployment-crud.html#ec_create_a_deployment

As an option "--generate-payload" can be used in order to obtain the generated payload that would be sent as a request. 
Save it, update or extend the topology and create a deployment using the saved payload with the "--file" flag.

```
ecctl deployment create {--file | --es-size <int> --es-zones <int> | --es-node-topology <obj>} [flags]
```

### Examples

```
## Create a deployment with the default values for Elasticsearch, a Kibana instance with a modified size, 
and a default APM instance. While Elasticsearch and Kibana come enabled by default, APM, Enterprise Search and App Search need to be 
enabled by using the "--apm", "--enterprise_search" and "--appsearch" flags. The command will exit after the API response has been returned, without 
waiting until the deployment resources have been created. 
$ ecctl deployment create --name my-deployment --zones 2 --kibana-size 2g --apm --apm-size 0.5g

## To make the command wait until the resources have been created use the "--track" flag, which will output 
the current stage on which the deployment resources are in.
$ deployment create --name my-deployment --track
[...]
Deployment [b6ecbea3d5c84124b7dca457f2892086] - [Elasticsearch][b6ecbea3d5c84124b7dca457f2892086]: finished running all the plan steps (Total plan duration: 5m11.s)
Deployment [91c4d60acb804ba0a27651fac02780ec] - [Kibana][8a9d9916cd6e46a7bb0912211d76e2af]: finished running all the plan steps (Total plan duration: 4m29.58s)

## Additionally, a more advanced topology for Elasticsearch can be created through "--es-node-topology" or shorthand "-e".
The following command will create a deployment with two 1GB Elasticsearch instances of the type "data" and 
one 1GB Elasticsearch instance of the type "ml".
$ ecctl deployment create --name my-deployment --es-node-topology '{"size": "1g", "zone_count": 2, "node_type": "data"}' --es-node-topology '{"size": "1g", "zone_count": 1, "node_type": "ml"}'

## In order to use the "--deployment-template" flag, you'll need to know which deployment templates ara available to you.
You'll need to run the following command to view your deployment templates:
$ ecctl platform deployment-template list

## Use the "--generate-payload" flag to save the definition to a file for later use.
$ ecctl deployment create --name my-deployment --size 1g --track --generate-payload --zones 2 > create_example.json

## Create a deployment through the file definition.
$ ecctl deployment create --file create_example.json --track

## To retry a deployment when the previous deployment creation failed, use the request ID provided in the error response of the previous command:
$ ecctl deployment create --request-id=GMZPMRrcMYqHdmxjIQkHbdjnhPIeBElcwrHwzVlhGUSMXrEIzVXoBykSVRsKncNb
```

### Options

```
      --apm                               Enables APM for the deployment
      --apm-ref-id string                 Optional RefId for the APM deployment (default "main-apm")
      --apm-size string                   Memory (RAM) in MB that each of the APM instances will have (default "0.5g")
      --apm-zones int32                   Number of zones the APM instances will span (default 1)
      --appsearch                         Enables App Search for the deployment
      --appsearch-ref-id string           Optional RefId for the App Search deployment (default "main-appsearch")
      --appsearch-size string             Memory (RAM) in MB that each of the App Search instances will have (default "2g")
      --appsearch-zones int32             Number of zones the App Search instances will span (default 1)
      --deployment-template string        Deployment template ID on which to base the deployment from
      --enterprise_search                 Enables Enterprise Search for the deployment
      --enterprise_search-ref-id string   Optional RefId for the Enterprise Search deployment (default "main-enterprise_search")
      --enterprise_search-size string     Memory (RAM) in MB that each of the Enterprise Search instances will have (default "4g")
      --enterprise_search-zones int32     Number of zones the Enterprise Search instances will span (default 1)
  -e, --es-node-topology stringArray      Optional Elasticsearch node topology element definition. See help for more information
      --es-ref-id string                  Optional RefId for the Elasticsearch deployment (default "main-elasticsearch")
      --es-size string                    Memory (RAM) in MB that each of the Elasticsearch instances will have (default "4g")
      --es-zones int32                    Number of zones the Elasticsearch instances will span (default 1)
  -f, --file string                       DeploymentCreateRequest file definition. See help for more information
      --generate-payload                  Returns the deployment payload without actually creating the deployment resources
  -h, --help                              help for create
      --kibana-ref-id string              Optional RefId for the Kibana deployment (default "main-kibana")
      --kibana-size string                Memory (RAM) in MB that each of the Kibana instances will have (default "1g")
      --kibana-zones int32                Number of zones the Kibana instances will span (default 1)
      --name string                       Optional name for the deployment
      --plugin strings                    Additional plugins to add to the Elasticsearch deployment
      --request-id string                 Optional request ID - Can be found in the Stderr device when a previous deployment creation failed. For more information see the examples in the help command page
  -t, --track                             Tracks the progress of the performed task
      --version string                    Version to use, if not specified, the latest available stack version will be used
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

