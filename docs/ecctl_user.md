## ecctl user

Manages the platform users

### Synopsis

Manages the platform users

```
ecctl user [flags]
```

### Options

```
  -h, --help   help for user
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

* [ecctl](ecctl.md)	 - Elastic Cloud Control
* [ecctl user create](ecctl_user_create.md)	 - Creates a new platform user
* [ecctl user delete](ecctl_user_delete.md)	 - Deletes one or more platform users
* [ecctl user disable](ecctl_user_disable.md)	 - Disables a platform user
* [ecctl user enable](ecctl_user_enable.md)	 - Enables a previously disabled platform user
* [ecctl user key](ecctl_user_key.md)	 - Manages the API keys of a platform user
* [ecctl user list](ecctl_user_list.md)	 - Lists all platform users
* [ecctl user show](ecctl_user_show.md)	 - Shows details of a specified user
* [ecctl user update](ecctl_user_update.md)	 - Updates a platform user

