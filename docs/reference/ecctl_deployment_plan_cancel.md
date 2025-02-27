---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_deployment_plan_cancel.html
---

# ecctl deployment plan cancel [ecctl_deployment_plan_cancel]

Cancels a resource’s pending plan

```
ecctl deployment plan cancel <deployment id> --kind <kind> [--ref-id <ref-id>] [flags]
```


## Options [_options_29]

```
  -h, --help            help for cancel
      --kind string     Required deployment resource kind (apm, appsearch, kibana, elasticsearch)
      --ref-id string   Optional deployment RefId, if not set, the RefId will be auto-discovered
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_28]

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


## SEE ALSO [_see_also_29]

* [ecctl deployment plan](/reference/ecctl_deployment_plan.md)	 - Manages deployment plans

