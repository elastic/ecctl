## ecctl deployment elasticsearch console

Starts an interactive console with the cluster

### Synopsis

Starts an interactive console that connects to the cluster. If no username or password is specified
it will use the system's credentials to connect to the cluster (needs admin privileges)

```
ecctl deployment elasticsearch console <cluster id> [flags]
```

### Examples

```

	If run without admin credentials, it will need a user and pass specification:
    ecctl elasticsearch console 18fc96c491b3d5e10e147463927a5349 --elasticsearch-user user --elasticsearch-pass pass

```

### Options

```
      --elasticsearch-pass string   Set the elasticsearch password for the interactive console
      --elasticsearch-user string   Set the elasticsearch user for the interactive console
  -h, --help                        help for console
      --query                       Instead of opening a console it will query the cluster with the arguments passed
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

