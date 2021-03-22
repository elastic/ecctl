## ecctl deployment create

Creates a deployment

### Synopsis

Creates a deployment from a file definition or using the defaults from the deployment template.
Default values are provided by the deployment template, simplifying the getting started experience.
When version is not specified, the latest available stack version will automatically be used.
These are the available options:

  * Using the deployment template defaults with --deployment-template=<deployment template id>
  * Using the deployment template defaults with --deployment-template=<deployment template id> and 
    --generate-payload to use the deployment template, modify certain fields of the definition and
    use the --file flag to create the deployment.
  * File definition: --file=<file path> (shorthand: -f). You can create a definition by using the sample JSON seen here:
    https://elastic.co/guide/en/cloud/current/ec-api-deployment-crud.html#ec_create_a_deployment

As an option "--generate-payload" can be used in order to obtain the generated payload that would be sent as a request. 
Save it, update or extend the topology and create a deployment using the saved payload with the "--file" flag.

```
ecctl deployment create {--file | --es-size <int> --es-zones <int> | --es-node-topology <obj>} [flags]
```

### Examples

```
## Create a deployment with the default values
$ ecctl deployment create --name my-deployment --deployment-template=aws-io-optimized-v2 --region=us-east-1

## To make the command wait until the resources have been created use the "--track" flag, which will output 
the current stage on which the deployment resources are in.
$ ecctl deployment create --name my-deployment --deployment-template=aws-io-optimized-v2 --region=us-east-1 --track
[...]
Deployment [b6ecbea3d5c84124b7dca457f2892086] - [Elasticsearch][b6ecbea3d5c84124b7dca457f2892086]: finished running all the plan steps (Total plan duration: 5m11.s)
Deployment [91c4d60acb804ba0a27651fac02780ec] - [Kibana][8a9d9916cd6e46a7bb0912211d76e2af]: finished running all the plan steps (Total plan duration: 4m29.58s)

## In order to use the "--deployment-template" flag, you'll need to know which deployment templates ara available to you.
You'll need to run the following command to view your deployment templates:
$ ecctl platform deployment-template list

## Use the "--generate-payload" flag to save the definition to a file for later use.
$ ecctl deployment create --name my-deployment --version=7.11.2 --generate-payload > create_example.json

## Create a deployment through the file definition.
$ ecctl deployment create --file create_example.json --track

## To retry a deployment when the previous deployment creation failed, use the request ID provided in the error response of the previous command:
$ ecctl deployment create --request-id=GMZPMRrcMYqHdmxjIQkHbdjnhPIeBElcwrHwzVlhGUSMXrEIzVXoBykSVRsKncNb
```

### Options

```
      --deployment-template string   Deployment template ID on which to base the deployment from
  -f, --file string                  DeploymentCreateRequest file definition. See help for more information
      --generate-payload             Returns the deployment payload without actually creating the deployment resources
  -h, --help                         help for create
      --name string                  Optional name for the deployment
      --request-id string            Optional request ID - Can be found in the Stderr device when a previous deployment creation failed. For more information see the examples in the help command page
  -t, --track                        Tracks the progress of the performed task
      --version string               Version to use, if not specified, the latest available stack version will be used
```

### Options inherited from parent commands

```
      --api-key string        API key to use to authenticate (If empty will look for EC_API_KEY environment variable)
      --config string         Config name, used to have multiple configs in $HOME/.ecctl/<env> (default "config")
      --force                 Do not ask for confirmation
      --format string         Formats the output using a Go template
      --host string           Base URL to use
      --insecure              Skips all TLS validation
      --message string        A message to set on cluster operation
      --output string         Output format [text|json] (default "text")
      --pass string           Password to use to authenticate (If empty will look for EC_PASS environment variable)
      --pprof                 Enables pprofing and saves the profile to pprof-20060102150405
  -q, --quiet                 Suppresses the configuration file used for the run, if any
      --region string         Elasticsearch Service region
      --timeout duration      Timeout to use on all HTTP calls (default 30s)
      --trace                 Enables tracing saves the trace to trace-20060102150405
      --user string           Username to use to authenticate (If empty will look for EC_USER environment variable)
      --verbose               Enable verbose mode
      --verbose-credentials   When set, Authorization headers on the request/response trail will be displayed as plain text
      --verbose-file string   When set, the verbose request/response trail will be written to the defined file
```

### SEE ALSO

* [ecctl deployment](ecctl_deployment.md)	 - Manages deployments

