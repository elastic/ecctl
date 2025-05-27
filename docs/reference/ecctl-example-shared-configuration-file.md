---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-example-shared-configuration-file.html
---

# Example: A shared configuration file [ecctl-example-shared-configuration-file]

Below is an example `YAML` configuration file `$HOME/.ecctl/config.yaml` that will effectively point and configure the binary for Elastic Cloud:

```yaml
host: https://api.elastic-cloud.com # URL of your Elastic Cloud or Elastic Cloud Enterprise API endpoint

# Credentials
## api_key is the preferred authentication mechanism.
api_key: bWFyYzo4ZTJmNmZkNjY5ZmQ0MDBkOTQ3ZjI3MTg3ZWI5MWZhYjpOQktHY05jclE0cTBzcUlnTXg3QTd3

## username and password can be used when no API key is available.
user: username
pass: password
```

