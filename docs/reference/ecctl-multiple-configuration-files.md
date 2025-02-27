---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-multiple-configuration-files.html
---

# Multiple configuration files [ecctl-multiple-configuration-files]

ecctl supports having multiple configuration files out of the box. This allows for easy management of multiple environments or specialized targets. By default it will use `$HOME/.ecctl/config.<json|toml|yaml|hcl>`, but when the `--config` flag is specified, it will append the `--config` name to the file:

```
# Default behaviour
$ ecctl version
# will use ~/.ecctl/config.yaml

# When an environment is specified, the configuration file used will change
$ ecctl version --config ece
# will use ~/.ecctl/ece.yaml
```

