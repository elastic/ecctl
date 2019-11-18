## ecctl deployment restore

Restores a previously shut down deployment and all of its associated sub-resources

### Synopsis

Use --restore-snapshot to automatically restore the latest available Elasticsearch snapshot (if applicable)

```
ecctl deployment restore <deployment-id> [flags]
```

### Examples

```
$ ecctl deployment restore 5c17ad7c8df73206baa54b6e2829d9bc
{
  "id": "5c17ad7c8df73206baa54b6e2829d9bc"
}

```

### Options

```
  -h, --help               help for restore
      --restore-snapshot   Restores snapshots for those resources that allow it (Elasticsearch)
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

* [ecctl deployment](ecctl_deployment.md)	 - Manages deployments

