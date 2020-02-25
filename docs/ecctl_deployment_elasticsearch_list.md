## ecctl deployment elasticsearch list

Returns the list of Elasticsearch clusters

### Synopsis

Returns the list of Elasticsearch clusters

```
ecctl deployment elasticsearch list [flags]
```

### Examples

```
* Returns ES Clusters that contain the string "logging" in the cluster name:
	ecctl elasticsearch list --query cluster_name:logging

* Returns ES Cluster with cluster ID bedd14a60f0871b53636496b5b4b0014:
	ecctl elasticsearch list --query cluster_id:bedd14a60f0871b53636496b5b4b0014

* Returns ES Clusters that contain the string "prod cluster" in the cluster name:
	ecctl elasticsearch list --query cluster_name:'prod cluster'

* To search in both cluster name and ID, ignore the query field:
	ecctl elasticsearch list --query logging

Read all the simple query string syntax in https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-query-string-query.html#query-string-syntax

```

### Options

```
  -h, --help             help for list
  -m, --metadata         Shows deployment metadata
      --query string     Searches clusters using an Elasticsearch search query string query
  -s, --size int         Sets the upper limit of ES clusters to return (default 100)
  -v, --version string   Filters per version
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
      --region string      Elasticsearch Service region
      --timeout duration   Timeout to use on all HTTP calls (default 30s)
      --trace              Enables tracing saves the trace to trace-20060102150405
      --user string        Username to use to authenticate (If empty will look for EC_USER environment variable)
      --verbose            Enable verbose mode
```

### SEE ALSO

* [ecctl deployment elasticsearch](ecctl_deployment_elasticsearch.md)	 - Manages Elasticsearch clusters

