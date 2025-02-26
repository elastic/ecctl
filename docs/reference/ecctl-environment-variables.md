---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-environment-variables.html
---

# Environment variables [ecctl-environment-variables]

The same settings can be defined as environment variables instead of a configuration file or to override certain settings of the `YAML` file. If setting environment variables, you’ll need to prefix the configuration parameter with `EC_` and capitalize the setting, i.e. `EC_HOST` or `EC_USER`.

```sh
export EC_API_KEY=bWFyYzo4ZTJmNmZkNjY5ZmQ0MDBkOTQ3ZjI3MTg3ZWI5MWZhYjpOQktHY05jclE0cTBzcUlnTXg3QTd3
```


## Special Environment Variables [ecctl-special-environment-variables]

```sh
export EC_CONFIG=$HOME/.ecctl/cloud.yaml
```

