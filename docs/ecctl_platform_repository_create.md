## ecctl platform repository create

Creates / updates a snapshot repository (Available for ECE only)

### Synopsis

Creates / updates a snapshot repository using a set of settings that can be
specified as a (yaml|json) file with the --settings flag.

The available settings to set depend on the the --type flag (default s3). A
list with the supported settings for each snapshot can be found in the docs:
https://www.elastic.co/guide/en/elasticsearch/reference/current/modules-snapshots.html#_repository_plugins

The --type flag can be set to any arbitrary value if it's differs from "s3".
Only the S3 available settings are validated.


```
ecctl platform repository create <repository name> --settings <settings file> [flags]
```

### Examples

```
ecctl platform repository create my-snapshot-repo --settings settings.yml

ecctl platform repository update my-snapshot-repo --settings settings.yml

ecctl platform repository create custom --type fs --settings settings.yml

```

### Options

```
  -h, --help              help for create
      --settings string   Configuration file for the snapshot repository
      --type string       Repository type that will be configured (default "s3")
```

### Options inherited from parent commands

```
      --api-key string     API key to use to authenticate (If empty will look for EC_API_KEY environment variable)
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

* [ecctl platform repository](ecctl_platform_repository.md)	 - Manages snapshot repositories (Available for ECE only)

