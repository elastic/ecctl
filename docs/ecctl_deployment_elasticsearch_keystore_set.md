## ecctl deployment elasticsearch keystore set

Updates an Elasticsearch cluster keystore with the contents of a file

### Synopsis

Manages the keystore settings of an Elasticsearch cluster.
Note that each operation is add/modify only, unspecified existing keystore values will be unchanged.

```
ecctl deployment elasticsearch keystore set <cluster id> -f <file definition.json> [flags]
```

### Examples

```
$ cat keystore_example.json
{
    "secrets": {
        "s3.client.foobar.access_key": {
            "value": "AKIXAIQFKXPHIFXSILUWPA",
            "as_file": false
        },
        "s3.client.foobar.secret_key": {
            "value": "18qXOpY2zGlApay1237dLXh+LG1X5LUNWjTHq5X1SWjf++m+p0"
        }
    }
}
$ ecctl deployment elasticsearch keystore set 4c052fb17f65467a9b3c36d060106377 --file keystore_example.json
{
  "secrets": {
    "s3.client.foobar.access_key": {
      "as_file": false
    },
    "s3.client.foobar.secret_key": {
      "as_file": false
    }
  }
}
```

### Options

```
  -f, --file string   JSON file that contains JSON-style domain-specific keystore definition
  -h, --help          help for set
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

* [ecctl deployment elasticsearch keystore](ecctl_deployment_elasticsearch_keystore.md)	 - Manages an Elasticsearch cluster's keystore

