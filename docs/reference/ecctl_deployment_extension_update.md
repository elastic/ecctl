---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_deployment_extension_update.html
---

# ecctl deployment extension update [ecctl_deployment_extension_update]

Updates an extension

```
ecctl deployment extension update <extension id> {--file <file-path> | --generate-payload} [--extension-file <file path>] [flags]
```


## Examples [_examples_3]

```
* Return the current extension state as a valid update payload.
  ecctl deployment extension update <extension id> --generate-payload > update.json

* After editing the file with your new values pass it as an argument to the --file flag.
  ecctl deployment extension update <extension id> --file update.json

* The extensions uploaded from a local file will remain unchanged unless the --extension-file flag is used.
  ecctl deployment extension update <extension id> --file update.json --extension-file extension.zip
```


## Options [_options_26]

```
      --extension-file string   Optional flag to upload an extension from a local file path.
      --file string             Path to the file containing the update JSON definition.
      --generate-payload        Outputs JSON which can be used as an argument for the --file flag.
  -h, --help                    help for update
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_25]

:::{include} _snippets/inherited-options.md
:::


## SEE ALSO [_see_also_26]

* [ecctl deployment extension](/reference/ecctl_deployment_extension.md)	 - Manages deployment extensions, such as custom plugins or bundles

