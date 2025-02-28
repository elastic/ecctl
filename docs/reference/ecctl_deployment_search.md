---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_deployment_search.html
---

# ecctl deployment search [ecctl_deployment_search]

Performs advanced deployment search using the Elasticsearch Query DSL


## Synopsis [_synopsis_6]

Read more about [Query DSL](elasticsearch://reference/query-languages/querydsl.md).

```
ecctl deployment search -f <query file.json> [flags]
```


## Examples [_examples_5]

```
$ cat query_string_query.json
{
    "query": {
        "query_string": {
            "query": "name: admin"
        }
    }
}
$ ecctl deployment search -f query_string_query.json
[...]
```

## Options [_options_41]

```
  -a, --all-matches   Uses a cursor to return all matches of the query (ignoring the size in the query). This can be used to query more than 10k results.
  -f, --file string   JSON file that contains JSON-style domain-specific language query
  -h, --help          help for search
      --size int32    Defines the size per request when using the --all-matches option. (default 500)
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_40]

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


## SEE ALSO [_see_also_41]

* [ecctl deployment](/reference/ecctl_deployment.md)	 - Manages deployments

