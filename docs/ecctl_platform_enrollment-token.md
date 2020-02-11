## ecctl platform enrollment-token

Manages tokens (Requires platform administration privileges)

### Synopsis

Manages tokens (Requires platform administration privileges)

```
ecctl platform enrollment-token [flags]
```

### Options

```
  -h, --help   help for enrollment-token
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

* [ecctl platform](ecctl_platform.md)	 - Manages the platform
* [ecctl platform enrollment-token create](ecctl_platform_enrollment-token_create.md)	 - Creates an enrollment token for role(s)
* [ecctl platform enrollment-token delete](ecctl_platform_enrollment-token_delete.md)	 - Deletes an enrollment token
* [ecctl platform enrollment-token list](ecctl_platform_enrollment-token_list.md)	 - Retrieves a list of persistent enrollment tokens

