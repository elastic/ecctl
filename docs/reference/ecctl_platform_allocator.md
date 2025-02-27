---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_platform_allocator.html
---

# ecctl platform allocator [ecctl_platform_allocator]

Manages allocators ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)

```
ecctl platform allocator [flags]
```


## Options [_options_65]

```
  -h, --help   help for allocator
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_64]

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


## SEE ALSO [_see_also_65]

* [ecctl platform](/reference/ecctl_platform.md)	 - Manages the platform ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)
* [ecctl platform allocator list](/reference/ecctl_platform_allocator_list.md)	 - Returns all allocators that have instances or are connected to the platform. Use --all flag or --output json to show all. Use --query to match any of the allocators properties. ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)
* [ecctl platform allocator maintenance](/reference/ecctl_platform_allocator_maintenance.md)	 - Sets the allocator in Maintenance mode ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)
* [ecctl platform allocator metadata](/reference/ecctl_platform_allocator_metadata.md)	 - Manages an allocator’s metadata
* [ecctl platform allocator search](/reference/ecctl_platform_allocator_search.md)	 - Performs advanced allocator searching ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)
* [ecctl platform allocator show](/reference/ecctl_platform_allocator_show.md)	 - Returns information about the allocator ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)
* [ecctl platform allocator vacate](/reference/ecctl_platform_allocator_vacate.md)	 - Moves all the resources from the specified allocator ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)

