## ecctl deployment extension update

Updates an extension

```
ecctl deployment extension update <extension id> {--file <file-path> | --generate-payload} [--extension-file <file path>] [flags]
```

### Examples

```

* Return the current extension state as a valid update payload.
  ecctl deployment extension update <extension id> --generate-payload > update.json

* After editing the file with your new values pass it as an argument to the --file flag.
  ecctl deployment extension update <extension id> --file update.json

* The extensions uploaded from a local file will remain unchanged unless the --extension-file flag is used.
  ecctl deployment extension update <extension id> --file update.json --extension-file extension.zip
```

### Options

```
      --extension-file string   Optional flag to upload an extension from a local file path.
      --file string             Path to the file containing the update JSON definition.
      --generate-payload        Outputs JSON which can be used as an argument for the --file flag.
  -h, --help                    help for update
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

