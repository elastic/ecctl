---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_comment.html
applies_to:
  deployment:
    ece: all
---

# ecctl comment [ecctl_comment]

Manages resource comments.

```
ecctl comment [flags]
```


## Options [_options_8]

```
  -h, --help   help for comment
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_7]

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


## See also [_see_also_8]

* [ecctl](/reference/ecctl.md)	 - Elastic Cloud Control
* [ecctl comment create](/reference/ecctl_comment_create.md)	 - Creates a new resource comment
* [ecctl comment delete](/reference/ecctl_comment_delete.md)	 - Deletes a resource comment
* [ecctl comment list](/reference/ecctl_comment_list.md)	 - Lists all resource comments
* [ecctl comment show](/reference/ecctl_comment_show.md)	 - Shows information about a resource comment
* [ecctl comment update](/reference/ecctl_comment_update.md)	 - Updates an existing resource comment