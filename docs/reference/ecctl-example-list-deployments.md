---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-example-list-deployments.html
---

# List deployments [ecctl-example-list-deployments]

As a first example of using ecctl, run the [ecctl deployment list](/reference/ecctl_deployment_list.md) command to retrieve information about existing deployments. This is a good way to check if ecctl is configured correctly and if you have any deployments already created.

```sh
ecctl deployment list
```

```sh
ID                                 NAME            ELASTICSEARCH                      KIBANA                             APM                                APPSEARCH
00be03849b6a49c1a6541e3ccb5958d2   marvin          00be03849b6a49c1a6541e3ccb5958d2   266e456acf257588a9cde6fb4569d4a0   78c096c22e12408b878083b2d5ff6bcf   -
147cdeace6404c3e4b5018e1401647e4   biggerdata      147cdeace6404c3e4b5018e1401647e4   443a9df7b33952f45921c5823cbad4bc   4678ce52d45547e463455ede663cb4a4   -
```

