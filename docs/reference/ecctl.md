---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl.html
---

# ecctl [ecctl]

Elastic Cloud Control


## Options [_options]

```
      --api-key string        API key to use to authenticate (If empty will look for EC_API_KEY environment variable)
      --config string         Config name, used to have multiple configs in $HOME/.ecctl/<env> (default "config")
      --force                 Do not ask for confirmation
      --format string         Formats the output using a Go template
  -h, --help                  help for ecctl
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


## SEE ALSO [_see_also]

* [ecctl auth](/reference/ecctl_auth.md) - Manages authentication settings
* [ecctl comment](/reference/ecctl_comment.md) - Manages resource comments
* [ecctl deployment](/reference/ecctl_deployment.md) - Manages deployments
* [ecctl generate](/reference/ecctl_generate.md) - Generates completions and docs
* [ecctl init](/reference/ecctl_init.md) - Creates an initial configuration file.
* [ecctl platform](/reference/ecctl_platform.md) - Manages the platform
* [ecctl stack](/reference/ecctl_stack.md) - Manages Elastic StackPacks
* [ecctl user](/reference/ecctl_user.md) - Manages the platform users
* [ecctl version](/reference/ecctl_version.md) - Shows ecctl version

