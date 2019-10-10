## ecctl auth key

Manages the current authenticated user API keys

### Synopsis

Manages the current authenticated user API keys

```
ecctl auth key [flags]
```

### Options

```
  -h, --help   help for key
```

### Options inherited from parent commands

```
      --apikey string      API key to use to authenticate (If empty will look for EC_APIKEY environment variable)
      --config string      Config name, used to have multiple configs in $HOME/.ecctl/<env> (default "config")
      --force              Do not ask for confirmation
      --format string      Formats the output using a Go template
      --host string        Base URL to use (default "https://api.elastic-cloud.com")
      --insecure           Skips all TLS validation
      --message string     A message to set on cluster operation
      --output string      Output format [text|json] (default "text")
      --pass string        Password to use to authenticate (If empty will look for EC_PASS environment variable)
      --pprof              Enables pprofing and saves the profile to pprof-20060102150405
  -q, --quiet              Suppresses the configuration file used for the run, if any
      --region string      Elastic Cloud region
      --timeout duration   Timeout to use on all HTTP calls (default 30s)
      --trace              Enables tracing saves the trace to trace-20060102150405
      --user string        Username to use to authenticate (If empty will look for EC_USER environment variable)
      --verbose            Enable verbose mode
```

### SEE ALSO

* [ecctl auth](ecctl_auth.md)	 - Manages the platform auth
* [ecctl auth key create](ecctl_auth_key_create.md)	 - Creates a new API key for the current authenticated user
* [ecctl auth key delete](ecctl_auth_key_delete.md)	 - Deletes one or more existing API keys for the specified user
* [ecctl auth key list](ecctl_auth_key_list.md)	 - Lists the API keys for the current authenticated user
* [ecctl auth key show](ecctl_auth_key_show.md)	 - Shows the API key details for the current authenticated user

