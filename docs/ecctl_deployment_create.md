## ecctl deployment create

Creates a deployment

### Synopsis

Creates a deployment which is defined from a file definition using the --file=<file path> (shorthand: -f) flag.

//You can create a definition by using the sample JSON seen here:
//  https://elastic.co/guide/en/cloud/current/ec-api-deployment-crud.html#ec_create_a_deployment

```
ecctl deployment create {--file | --size <int> --zones <string> | --topology-element <obj>} [flags]
```

### Examples

```
## To make the command wait until the resources have been created use the "--track" flag, which will output 
the current stage on which the deployment resources are in.
$ deployment create --file create_example.json --track
[...]
Deployment [b6ecbea3d5c84124b7dca457f2892086] - [Elasticsearch][b6ecbea3d5c84124b7dca457f2892086]: finished running all the plan steps (Total plan duration: 5m11.s)
Deployment [91c4d60acb804ba0a27651fac02780ec] - [Kibana][8a9d9916cd6e46a7bb0912211d76e2af]: finished running all the plan steps (Total plan duration: 4m29.58s)

## To retry a deployment when the previous deployment creation failed, use the request ID provided in the error response of the previous command:
$ ecctl deployment create --file create_example.json --request-id=GMZPMRrcMYqHdmxjIQkHbdjnhPIeBElcwrHwzVlhGUSMXrEIzVXoBykSVRsKncNb
```

### Options

```
  -f, --file string         DeploymentCreateRequest file definition. See help for more information
  -h, --help                help for create
      --request-id string   Optional request ID - Can be found in the Stderr device when a previous deployment creation failed. For more information see the examples in the help command page
  -t, --track               Tracks the progress of the performed task
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

