## ecctl platform enrollment-token create

Creates an enrollment token for role(s)

### Synopsis

Creates an enrollment token for role(s)

```
ecctl platform enrollment-token create --role <ROLE> [flags]
```

### Examples

```
  ecctl [globalFlags] enrollment-token create --role coordinator
  ecctl [globalFlags] enrollment-token create --role coordinator --role proxy
  ecctl [globalFlags] enrollment-token create --role allocator --validity 120s
  ecctl [globalFlags] enrollment-token create --role allocator --validity 2h
```

### Options

```
  -h, --help                help for create
  -r, --role stringArray    Role(s) to associate the tokens with.
  -v, --validity duration   Time in seconds for which this token is valid. Currently this will make the token ephemeral (persistent: false)
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

* [ecctl platform enrollment-token](ecctl_platform_enrollment-token.md)	 - Manages tokens

