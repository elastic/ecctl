## ecctl deployment elasticsearch create

Creates an Elasticsearch cluster

### Synopsis

Creates an Elasticsearch cluster

```
ecctl deployment elasticsearch create [--file] [--capacity|--version|--zones] [name] [flags]
```

### Examples

```

  * Create a single node elasticsearch cluster
  ecctl deployment elasticsearch create --version 5.6.0 --zones 1 --capacity 2048
 
  * Create an elasticsearch cluster from a plan definition.
  ecctl deployment elasticsearch create --file definition.json

```

### Options

```
  -c, --capacity int32   Capacity per node
      --file string      JSON plan definition file location
  -h, --help             help for create
      --plugin strings   Additional plugins to add to the Elasticsearch cluster
  -t, --track            Tracks the progress of the performed task
  -v, --version string   Filter per version
  -z, --zones int32      Number of zones for the cluster
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

