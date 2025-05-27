---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_platform_repository_create.html
applies_to:
  deployment:
    ece: all
---

# ecctl platform repository create [ecctl_platform_repository_create]

Creates or updates a snapshot repository.


## Synopsis [_synopsis_11]

Creates / updates a snapshot repository using a set of settings that can be specified as a (yaml|json) file with the --settings flag.

The available settings to set depend on the the --type flag (default s3). A list with the supported settings for each snapshot can be found in the docs: [Snapshot and restore](docs-content://deploy-manage/tools/snapshot-and-restore.md)

The --type flag can be set to any arbitrary value if it’s differs from "s3". Only the S3 available settings are validated.

```
ecctl platform repository create <repository name> --settings <settings file> [flags]
```


## Examples [_examples_12]

```
ecctl platform repository create my-snapshot-repo --settings settings.yml

ecctl platform repository update my-snapshot-repo --settings settings.yml

ecctl platform repository create custom --type fs --settings settings.yml
```


## Options [_options_105]

```
  -h, --help              help for create
      --settings string   Configuration file for the snapshot repository
      --type string       Repository type that will be configured (default "s3")
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_104]

:::{include} _snippets/inherited-options.md
:::


## See also [_see_also_105]

* [ecctl platform repository](/reference/ecctl_platform_repository.md) - Manages snapshot repositories

