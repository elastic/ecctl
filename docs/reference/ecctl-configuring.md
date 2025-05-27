---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-configuring.html
applies_to:
  deployment:
    ess: all
    ece: all
---

# Configuring [ecctl-configuring]

In order for ecctl to be able to communicate with the RESTful API, it needs to have a set of configuration parameters defined. These parameters can be set in a configuration file, through environment variables, or at runtime using the CLI's global flags.


### Before you begin [_before_you_begin]

The hierarchy for configuration parameters is as follows, from higher precedence to lower:

1. Command line flags `--api-key`, `--region`, `--verbose`
2. Environment variables
3. Shared configuration file `$HOME/.ecctl/config.<json|toml|yaml|hcl>`


## Generate a configuration file [_generate_a_configuration_file]

If it's your first time using ecctl, use the `init` command to assist you in generating a configuration file. The resulting configuration file will be saved under `~/.ecctl/config.json`:

```
$ ecctl init
Welcome to Elastic Cloud Control (ecctl)! This command will guide you through authenticating and setting some default values.

Missing configuration file, would you like to initialise it? [y/n]: y

Select which type of Elastic Cloud offering you will be working with:
  [1] Elasticsearch Service (default).
  [2] Elastic Cloud Enterprise (ECE).
  [3] Elasticsearch Service Private (ESSP).

Please enter your choice: 1

Using "https://api.elastic-cloud.com" as the API endpoint.

Select a region you would like to have as default:

  GCP
  [1] us-central1 (Iowa)
  [2] us-east1 (S. Carolina)
  [3] us-east4 (N. Virginia)
  [4] us-west1 (Oregon)
  [5] northamerica-northeast1 (Montreal)
  [6] southamerica-east1 (São Paulo)
  [7] australia-southeast1 (Sydney)
  [8] europe-west1 (Belgium)
  [9] europe-west2 (London)
  [10] europe-west3 (Frankfurt)
  [11] asia-northeast1 (Tokyo)
  [12] asia-south1 (Mumbai)
  [13] asia-southeast1 (Singapore)

  AWS
  [14] us-east-1 (N. Virginia)
  [15] us-west-1 (N. California)
  [16] us-west-2 (Oregon)
  [17] eu-central-1 (Frankfurt)
  [18] eu-west-2 (London)
  [19] eu-west-1 (Ireland)
  [20] ap-northeast-1 (Tokyo)
  [21] ap-southeast-1 (Singapore)
  [22] ap-southeast-2 (Sydney)
  [23] sa-east-1 (São Paulo)

  Azure
  [24] eastus2 (Virginia)
  [25] westus2 (Washington)
  [26] westeurope (Netherlands)
  [27] uksouth (London)
  [28] japaneast (Tokyo)
  [29] southeastasia (Singapore)

Please enter your choice: 1

Create a new Elasticsearch Service API key (https://cloud.elastic.co/account/keys) and/or
Paste your API Key and press enter: xxxxx

What default output format would you like?
  [1] text - Human-readable output format, commands with no output templates defined will fall back to JSON.
  [2] json - JSON formatted output API responses.

Please enter a choice: 1

Your credentials seem to be valid.

You're all set! Here are some commands to try:
  $ ecctl deployment list

Config written to /home/myuser/.ecctl/config.json
```
