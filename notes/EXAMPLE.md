# Changelog

## Breaking changes

Removes the previously deprecated `elevate-permissions` command since it does not have a place in ECE.

## Notable changes

### New Commands

Added CRUD commands for proxies filtered group:

```command
$ ecl platform proxy filtered-group

create      Create proxies filtered group
delete      Delete proxies filtered group
show        Show details for proxies filtered group
update      Update proxies filtered group
```

### Text output

All the text outputs have dynamic tab spaces between columns which allow for longer names to be correctly handled.

### Bug fixes

Fixed a bug where the injected Asset Template loader wasn't being used, resulting on the inability to extend the formatter with further 3rd party templates provided by consuming clients.

### Docs

Added a new "Workflow" section which covers all the contributing workflow guidelines for contributors, particularly focusing on commit messages.
Additionally PR and Issue templates have been added and a few re-wording changes have been made.
