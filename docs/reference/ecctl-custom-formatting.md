---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-custom-formatting.html
applies_to:
  deployment:
    ess: all
    ece: all
---

# Custom formatting [ecctl-custom-formatting]

ecctl supports a global `--format` flag which can be passed to any existing command or subcommand. Using the `--format` flag allows you to obtain a specific part of a command response that might not have been shown before with the default `--output=text`. The `--format` internally uses Go templates which means that you can use the power of the Go built-in [`text/templates`](https://golang.org/pkg/text/template/) on demand.

