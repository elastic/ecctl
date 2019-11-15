## ecctl platform deployment-template

Manages deployment templates

### Synopsis

Manages deployment templates

```
ecctl platform deployment-template [flags]
```

### Options

```
  -h, --help   help for deployment-template
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
* [ecctl platform deployment-template create](ecctl_platform_deployment-template_create.md)	 - Creates a platform deployment template
* [ecctl platform deployment-template delete](ecctl_platform_deployment-template_delete.md)	 - Deletes a specific platform deployment template
* [ecctl platform deployment-template list](ecctl_platform_deployment-template_list.md)	 - Lists the platform deployment templates
* [ecctl platform deployment-template pull](ecctl_platform_deployment-template_pull.md)	 - Downloads deployment template into a local folder
* [ecctl platform deployment-template show](ecctl_platform_deployment-template_show.md)	 - Shows information about a specific platform deployment template
* [ecctl platform deployment-template update](ecctl_platform_deployment-template_update.md)	 - Updates a platform deployment template

