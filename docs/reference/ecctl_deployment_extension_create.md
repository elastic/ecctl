---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_deployment_extension_create.html
---

# ecctl deployment extension create [ecctl_deployment_extension_create]

Creates an extension

```
ecctl deployment extension create <extension name> --version <version> --type <extension type> {--file <file-path> | --download-url <url>} [--description <description>] [flags]
```


## Options [_options_22]

```
      --description string    Optional flag to add a description to the extension.
      --download-url string   Optional flag to define the URL to download the extension archive.
      --file string           Optional flag to upload an extension from a local file path.
  -h, --help                  help for create
      --type string           Extension type. Can be one of [bundle, plugin].
      --version string        Elastic stack version. Numeric version for plugins, e.g. 7.10.0. Major version e.g. 7.*, or wildcards e.g. * for bundles.
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_21]

:::{include} _snippets/inherited-options.md
:::


## SEE ALSO [_see_also_22]

* [ecctl deployment extension](/reference/ecctl_deployment_extension.md)	 - Manages deployment extensions, such as custom plugins or bundles

