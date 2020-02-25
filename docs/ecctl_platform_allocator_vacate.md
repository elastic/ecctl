## ecctl platform allocator vacate

Moves all the clusters from the specified allocator

### Synopsis

Moves all the clusters from the specified allocator

```
ecctl platform allocator vacate <source> [flags]
```

### Examples

```
  ecctl [globalFlags] allocator vacate i-05e245252362f7f1d
  # Move everything from multiple allocators
  ecctl [globalFlags] allocator vacate i-05e245252362f7f1d i-2362f7f1d252362f7

  # filter by a cluster kind
  ecctl [globalFlags] allocator vacate -k kibana i-05e245252362f7f1d

  # Only move specific cluster IDs
  ecctl [globalFlags] allocator vacate -c f521dedb07194c478fbbc6624f3bbf8f -c f404eea372cc4ea5bd553d47a09705cd i-05e245252362f7f1d

  # Specify multiple allocator targets
  ecctl [globalFlags] allocator vacate -t i-05e245252362f7f2d -t i-2362f7f1d252362f7 i-05e245252362f7f1d
  ecctl [globalFlags] allocator vacate --target i-05e245252362f7f2d --target i-2362f7f1d252362f7 --kind kibana i-05e245252362f7f1d

  # Set the allocators to maintenance mode before vacating them
  ecctl [globalFlags] allocator vacate --maintenance -t i-05e245252362f7f2d -t i-2362f7f1d252362f7 i-05e245252362f7f1d

  # Set the amount of maximum moves to happen at any time
  ecctl [globalFlags] allocator vacate --concurrency 10 i-05e245252362f7f1d

  # Override the allocator health auto discovery
  ecctl [globalFlags] allocator vacate --allocator-down=true i-05e245252362f7f1d

  # Override the skip_snapshot setting
  ecctl [globalFlags] allocator vacate --skip-snapshot=true i-05e245252362f7f1d -c f521dedb07194c478fbbc6624f3bbf8f

  # Override the skip_data_migration setting
  ecctl [globalFlags] allocator vacate --skip-data-migration=true i-05e245252362f7f1d -c f521dedb07194c478fbbc6624f3bbf8f
  
  # Skips tracking the vacate progress which will cause the command to return almost immediately.
  # Not recommended since it can lead to failed vacates without the operator knowing about them.
  ecctl [globalFlags] allocator vacate --skip-tracking i-05e245252362f7f1d

```

### Options

```
      --allocator-down string        Disables the allocator health auto-discovery, setting the allocator-down to either [true|false]
  -c, --cluster stringArray          Cluster IDs to include in the vacate
      --concurrency uint             Maximum number of concurrent moves to perform at any time (default 8)
  -h, --help                         help for vacate
  -k, --kind string                  Kind of workload to vacate (elasticsearch|kibana)
  -m, --maintenance                  Whether to set the allocator(s) in maintenance before performing the vacate
      --max-poll-retries int         Optional maximum plan tracking retries (default 2)
      --move-only                    Keeps the cluster in its current -possibly broken- state and just does the bare minimum to move the requested instances across to another allocator. [true|false] (default true)
      --override-failsafe            If false (the default) then the plan will fail out if it believes the requested sequence of operations can result in data loss - this flag will override some of these restraints. [true|false]
      --poll-frequency duration      Optional polling frequency to check for plan change updates (default 10s)
      --skip-data-migration string   Skips the data-migration operation on the specified cluster IDs. ONLY available when the cluster IDs are specified and --move-only is true. [true|false]
      --skip-snapshot string         Skips the snapshot operation on the specified cluster IDs. ONLY available when the cluster IDs are specified. [true|false]
      --skip-tracking                Skips tracking the vacate progress causing the command to return after the move operation has been executed. Not recommended.
  -t, --target stringArray           Target allocator(s) on which to place the vacated workload
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

* [ecctl platform allocator](ecctl_platform_allocator.md)	 - Manages allocators (Available for ECE only)

