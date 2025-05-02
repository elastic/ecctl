---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_platform_allocator_list.html
---

# ecctl platform allocator list [ecctl_platform_allocator_list]

Returns all allocators that have instances or are connected to the platform. Use --all flag or --output json to show all. Use --query to match any of the allocators properties. ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)


## Synopsis [_synopsis_8]

Returns all allocators that have instances or are connected to the platform. Use --all flag or --output json to show all. Use --query to match any of the allocators properties.

Query examples:

```
* Allocators set to maintenance mode: --query status.maintenance_mode:true

* Allocators with more than 10GB of capacity: --query capacity.memory.total:\>10240
```
Read all the simple query string syntax in [/elasticsearch/docs/reference/query-languages/query-dsl-query-string-query.md#query-string-syntax](elasticsearch://reference/query-languages/query-dsl/query-dsl-query-string-query.md#query-string-syntax)

Filter examples:

```
* Allocators with instance type i3.large : --filter instanceType:i3.large

* Allocators with instance type i3.large AND instance family gcp.highcpu.1 : --filter instanceType:i3.large --filter instanceFamily:gcp.highcpu.1
```
Filter is a post-query action and doesn’t support OR. All filter arguments are applied as AND.

Filter and query flags can be used in combination.

```
ecctl platform allocator list [flags]
```


## Options [_options_66]

```
      --all                  Shows all allocators, including those with no instances or not connected, this is relative to the --size flag.
  -f, --filter stringArray   Post-query filter out allocators based on metadata tags, for instance 'instanceType:i3.large'
  -h, --help                 help for list
      --metadata             Shows allocators metadata
      --query string         Searches allocators using an Elasticsearch search query string query
      --size int             Defines the maximum number of allocators to return (default 100)
      --unhealthy            Searches for unhealthy allocators
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_65]

:::{include} /_snippets/inherited-options.md
:::


## SEE ALSO [_see_also_66]

* [ecctl platform allocator](/reference/ecctl_platform_allocator.md)	 - Manages allocators ![logo cloud ece](https://doc-icons.s3.us-east-2.amazonaws.com/logo_cloud_ece.svg "Supported on {{ece}}") (Available for ECE only)

