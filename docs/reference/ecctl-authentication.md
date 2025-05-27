---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-authentication.html
applies_to:
  deployment:
    ess: all
    ece: all
---

# Authentication [ecctl-authentication]

Elastic Cloud uses API keys to authenticate users against its API. Additionally, it supports the usage of [JWT](https://jwt.io/) to validate authenticated clients. The preferred authentication method is API keys.

There are two ways to authenticate against the Elasticsearch Service or the Elastic Cloud Enterprise APIs ecctl:

* By specifying an API key using the `--api-key` flag
* By specifying the `--user` and `--pass` flags

The first method requires the user to already have an API key, if this is the case, all the outgoing API requests will use an Authentication API key header.

The second method uses the `user` and `pass` values to obtain a valid JWT token, that token is then used as the Authentication Bearer header for every API call. A goroutine that refreshes the token every minute is started, so that the token doesn’t expire while we’re performing actions.

