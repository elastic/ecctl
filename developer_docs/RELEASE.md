# Releasing a new version

This guide aims to provide guidance on how to release new versions of the `ecctl` binary as well as updating all the necessary parts to make it successful.

## Prerequisites

Releasing a new version of the binary implies that there have been changes in the source code which are meant to be released for wider consumption. Before releasing a new version there's some prerequisites that have to be checked.

### Verify a release in cloud-sdk-go has been made

Unless this is a patch release, make sure a cloud-sdk-go release with the same version has been made.

Once this is done, the cloud-sdk-go dependency should be updated to that version.

### Make sure the version is updated

Since the source has changed, we need to update the current committed version to a higher version so that the release is published.

The version is currently defined in the [Makefile](./Makefile) as an exported environment variable called `VERSION` in the [SEMVER](https://semver.org) format: `MAJOR.MINOR.BUG`

```Makefile
SHELL := /bin/bash
export VERSION ?= v1.0.0
```

Say we want to perform a minor version release (i.e. no breaking changes and only new features and bug fixes are being included); in which case we'll update the _MINOR_ part of the version:

```Makefile
SHELL := /bin/bash
export VERSION ?= v1.1.0
```

### Generating a changelog for the new version

Once the version is updated, we can then generate the changelog and release notes by calling `make changelog`.

Take a look at one of our previous releases [`v1.0.0-beta2.adoc`](../docs/release_notes/v1.0.0-beta2.adoc) and the [template](../scripts/changelog.tpl.adoc) we use to generate them. The idea is to fill all the applicable sections so that users can consume easily.

After the release notes have been manually curated, a new pull request can be opened with the changelog, release notes and version update changes.

## Executing the release

After the new changelog and version have been merged to master, the only thing remaining is to run `make release`. This is the makefile target which will push the GitHub tag and will trigger the corresponding [GitHub action](.github/workflows/release.yml) which will release ecctl.

## Post release requirements

After a release has been performed there are still a few things we need to do.

### Create documentation specific to the release

In order to have the documentation live for our new release we need to modify the conf.yaml file in the docs repository to [add the release branch](https://github.com/elastic/docs/blob/master/conf.yaml#L837) and have the build [point to our new branch](https://github.com/elastic/docs/blob/master/conf.yaml#L836).

Once the PR for the above changes has been merged, you'll need to run a full doc [rebuild](https://elasticsearch-ci.elastic.co/job/elastic+docs+master+build/build?delay=0sec) to make the new release branch the default docs.

![alt text](docs-rebuild.png "rebuild instructions")
