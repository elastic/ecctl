## ecctl deployment extension create

Creates an extension

```
ecctl deployment extension create <extension name> --version <version> --type <extension type> {--file <file-path> | --download-url <url>} [--description <description>] [flags]
```

### Options

```
      --description string    Optional flag to add a description to the extension.
      --download-url string   Optional flag to define the URL to download the extension archive.
      --file string           Optional flag to upload an extension from a local file path.
  -h, --help                  help for create
      --type string           Extension type. Can be one of [bundle, plugin].
      --version string        Elastic stack version. Numeric version for plugins, e.g. 7.10.0. Major version e.g. 7.*, or wildcards e.g. * for bundles.
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

* [ecctl deployment extension](ecctl_deployment_extension.md)	 - Manages deployment extensions, such as custom plugins or bundles

