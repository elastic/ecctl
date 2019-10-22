## ecctl deployment elasticsearch instances override-capacity

Sets overrides at the instance level, use the --all flag to target all instances

### Synopsis

Sets overrides at the instance level, use the --all flag to target all instances.
If the multiplier flag is not set, the override will be set using the current configured capacity * 2.
Only works for overrides <= 65535

```
ecctl deployment elasticsearch instances override-capacity <cluster id> [flags]
```

### Examples

```

Set all the instances in the cluster to the original (plan) capacity:
    ecctl deployment elasticsearch instances override-capacity 18fc96c491b3d5e10e147463927a5349 --reset

Set all the instances in the cluster to 4096 capacity:
    ecctl deployment elasticsearch instances override-capacity 18fc96c491b3d5e10e147463927a5349 --all --value 4096

Set all the instances in the cluster to 2x multiplier:
    ecctl deployment elasticsearch instances override-capacity 18fc96c491b3d5e10e147463927a5349 --all --multiplier 2

Set the cluster instance to 3x of its current capacity:
    ecctl deployment elasticsearch instances override-capacity 18fc96c491b3d5e10e147463927a5349 --instance instance-0000000003 --multiplier 3
```

### Options

```
      --all                        Applies the override to all of the instances
  -h, --help                       help for override-capacity
  -i, --instance strings           Instances on which to apply the override
      --multiplier uint8           Capacity multiplier
      --reset                      Resets the instance(s) memory to the original value found in the current plan
      --storage-multiplier uint8   Storage multiplier per instance, if not set doesn't override it
      --value uint16               Absolute value of instance override memory (in MBs)
```

### Options inherited from parent commands

```
      --apikey string      API key to use to authenticate (If empty will look for EC_APIKEY environment variable)
      --config string      Config name, used to have multiple configs in $HOME/.ecctl/<env> (default "config")
      --force              Do not ask for confirmation
      --format string      Formats the output using a Go template
      --host string        Base URL to use (default "https://api.elastic-cloud.com")
      --insecure           Skips all TLS validation
      --message string     A message to set on cluster operation
      --output string      Output format [text|json] (default "text")
      --pass string        Password to use to authenticate (If empty will look for EC_PASS environment variable)
      --pprof              Enables pprofing and saves the profile to pprof-20060102150405
  -q, --quiet              Suppresses the configuration file used for the run, if any
      --region string      Elastic Cloud region
      --timeout duration   Timeout to use on all HTTP calls (default 30s)
      --trace              Enables tracing saves the trace to trace-20060102150405
      --user string        Username to use to authenticate (If empty will look for EC_USER environment variable)
      --verbose            Enable verbose mode
```

### SEE ALSO

* [ecctl deployment elasticsearch instances](ecctl_deployment_elasticsearch_instances.md)	 - Manages elasticsearch at the instance level

