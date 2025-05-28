---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl_platform_proxy_settings_update.html
applies_to:
  deployment:
    ece: all
---

# ecctl platform proxy settings update [ecctl_platform_proxy_settings_update]

Updates settings for all proxies.

```
ecctl platform proxy settings update --file settings.json [flags]
```


## Examples [_examples_11]

```
## Update only the defined proxy settings
$ ecctl platform proxy settings update --file settings.json --region us-east-1

## Update by overriding all proxy settings
$ ecctl platform proxy settings update --file settings.json --region us-east-1 --full
```


## Options [_options_102]

```
  -f, --file string      ProxiesSettings file definition. See https://www.elastic.co/guide/en/cloud-enterprise/current/ProxiesSettings.html for more information.
      --full             If set, a full update will be performed and all proxy settings will be overwritten. Any unspecified fields will be deleted.
  -h, --help             help for update
      --version string   If specified, checks for conflicts against the version of the repository configuration
```


## Options inherited from parent commands [_options_inherited_from_parent_commands_101]

:::{include} _snippets/inherited-options.md
:::


## See also [_see_also_102]

* [ecctl platform proxy settings](/reference/ecctl_platform_proxy_settings.md) - Manages proxies settings

