## ecctl user key

Manages the API keys of a platform user

### Synopsis

Manages the API keys of a platform user

```
ecctl user key [flags]
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

* [ecctl user](ecctl_user.md)	 - Manages the platform users (Requires platform administration privileges)
* [ecctl user key delete](ecctl_user_key_delete.md)	 - Deletes an existing API key for the specified user
* [ecctl user key list](ecctl_user_key_list.md)	 - Lists the API keys for the specified user, or all platform users
* [ecctl user key show](ecctl_user_key_show.md)	 - Shows the API key details for the specified user

