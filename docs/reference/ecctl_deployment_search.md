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

:::{include} _snippets/inherited-options.md
:::


## SEE ALSO [_see_also_41]

* [ecctl deployment](/reference/ecctl_deployment.md)	 - Manages deployments

