## ecctl deployment elasticsearch plan reapply

Reapplies the latest plan attempt, resetting all transient settings

### Synopsis

Reapplies the latest plan attempt, resetting all transient settings

```
ecctl deployment elasticsearch plan reapply <cluster id> [flags]
```

### Options

```
      --default                   Overwrites the strategy to the default one
      --extended-maintenance      Stops routing to the cluster instances after the plan has been applied
      --grow-and-shrink           Overwrites the strategy to grow and shrink
  -h, --help                      help for reapply
      --hide-plan                 Doesn't print the plan before reapplying
      --override-failsafe         Overrides failsafe at the constructor level that prevent bad things from happening
      --reallocate                Forces creation of new instances
      --rolling                   Overwrites the strategy to rolling
      --rolling-all               Overwrites the strategy to apply the change in all the instances at a time (causes downtime)
      --rolling-grow-and-shrink   Overwrites the strategy to rolling grow and shrink (one at a time)
      --skip-data-migration       Bypasses the need to wait for data to be migrated from old instances to new instances before continuing the plan (potentially deleting the old instances and losing data)
      --skip-post-upgrade-steps   Bypasses 2.x->5.x operations for any plan change ending with a 5.x cluster (eg apply a cluster license, ensure Monitoring is configured)
      --skip-snapshot             Skips snapshot on the reapplied plan
      --skip-upgrade-checker      Bypasses issue checks that should be resolved before migration (eg contains old Lucene segments)
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

* [ecctl deployment elasticsearch plan](ecctl_deployment_elasticsearch_plan.md)	 - Manages Elasticsearch plans

