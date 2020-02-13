## ecctl user create

Creates a new platform user

### Synopsis

Creates a new platform user

```
ecctl user create --username <username> --role <role> [flags]
```

### Examples

```

  * Create a platform user who has two roles assigned
    ecctl user create --username sam89 --role ece_platform_viewer --role ece_deployment_viewer

```

### Options

```
      --email string               User's email address (must be in a valid email format)
      --fullname string            User's full name
  -h, --help                       help for create
      --insecure-password string   [INSECURE] User's plaintext password
      --role strings               Role or roles assigned to the user. Available roles: 
                                   ece_platform_admin, ece_platform_viewer, ece_deployment_manager, ece_deployment_viewer
      --username string            Unique username for the platform user
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

* [ecctl user](ecctl_user.md)	 - Manages the platform users (Available for ECE only)

