---
navigation_title: "Elastic Cloud Control (ECCTL)"
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-release-notes.html
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-release-notes-v1.14.3.html
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-release-notes-v1.14.2.html
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-release-notes-v1.14.1.html
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-release-notes-v1.14.0.html
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-release-notes-v1.13.0.html
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-release-notes-v1.12.0.html
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-release-notes-v1.11.0.html
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-release-notes-v1.10.0.html
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-release-notes-v1.9.0.html
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-release-notes-v1.8.0.html
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-release-notes-v1.7.0.html
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-release-notes-v1.6.0.html
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-release-notes-v1.5.0.html
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-release-notes-v1.4.0.html
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-release-notes-v1.3.1.html
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-release-notes-v1.3.0.html
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-release-notes-v1.2.0.html
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-release-notes-v1.1.0.html
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-release-notes-v1.0.0.html
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-release-notes-v1.0.0-beta3.html
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-release-notes-v1.0.0-beta2.html
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-release-notes-v1.0.0-beta1.html
---

# Elastic Cloud Control (ECCTL) release notes

[Elastic Cloud Control (ECCTL)](/reference/index.md) is the command-line interface for the [Elastic Cloud](https://www.elastic.co/docs/api/doc/cloud) and [Elastic Cloud Enterprise](https://www.elastic.co/docs/api/doc/cloud-enterprise/) (ECE) APIs. It wraps typical operations commonly needed by operators within a single command line tool.

Review the changes, fixes, and more in each version of ECCTL.

<!--
Release notes include only features, enhancements, and fixes. Add breaking changes, deprecations, and known issues to the applicable release notes sections.

## version.next [ecctl-next-release-notes]

### Features and enhancements [ecctl-next-features-enhancements]
*

### Fixes [ecctl-next-fixes]
*
-->

## 1.15.0 [ecctl-1.15.0]

Release date: June 13, 2025

[Download the release binaries](https://github.com/elastic/ecctl/releases/tag/v1.15.0)

### Fixes [ecctl-1.15.0-fixes]

* **Bumps `cloud-sdk-go` to 1.24.1**: This version of the Go SDK includes a fix that allows ECCTL to include the integrations server payload when running the `ecctl deployment show` command with the `--generate-update-payload` flag. Previously, the `integrations_server` resource was incorrectly returned as `null`. [#728](https://github.com/elastic/ecctl/pull/728)

## 1.14.3 [ecctl-1.14.3]

Release date: January 17, 2025

[Download the release binaries](https://github.com/elastic/ecctl/releases/tag/v1.14.3)

### Fixes [ecctl-1.14.3-fixes]

* Fix issue with previous release (v1.14.2)

## 1.14.2 [ecctl-1.14.2]

Release date: January 16, 2025

[Download the release binaries](https://github.com/elastic/ecctl/releases/tag/v1.14.2)

### Features and enhancements [ecctl-1.14.2-features-enhancements]

* update module github.com/stretchr/testify to v1.10.0 ([#677](https://github.com/elastic/ecctl/pull/677))
* update module golang.org/x/term to v0.27.0 ([#661](https://github.com/elastic/ecctl/pull/661))

### Changelog

* [1da119d](https://github.com/elastic/ecctl/commit/1da119d) Rename stateful applications -> hosted applications ([#691](https://github.com/elastic/ecctl/pull/691))
* [c248cd2](https://github.com/elastic/ecctl/commit/c248cd2) fix(deps): update module github.com/stretchr/testify to v1.10.0 ([#677](https://github.com/elastic/ecctl/pull/677))
* [95d6fc5](https://github.com/elastic/ecctl/commit/95d6fc5) fix(deps): update module golang.org/x/term to v0.27.0 ([#661](https://github.com/elastic/ecctl/pull/661))

## 1.14.1 [ecctl-1.14.1]

Release date: November 18, 2024

[Download the release binaries](https://github.com/elastic/ecctl/releases/tag/v1.14.1)

### Features and enhancements [ecctl-1.14.1-features-enhancements]

* **Add --config-version flag to instance configuration show command** ([#669](https://github.com/elastic/ecctl/pull/669)): The `platform instance-configuration show` command now also supports the `--config-version` and `--show-deleted` flags to show a specific instance configuration version and allow fetching deleted instance configurations, respectively.

### Changelog 

* [2a6f80f](https://github.com/elastic/ecctl/commit/2a6f80f) feat: add `--config-version` flag to `instance configuration show command` ([#669](https://github.com/elastic/ecctl/pull/669))

## 1.14.0 [ecctl-1.14.0]

Release date: September 26, 2024

[Download the release binaries](https://github.com/elastic/ecctl/releases/tag/v1.14.0)

### Features and enhancements [ecctl-1.14.0-features-enhancements]

* **Deployment search: Add flag to return all matches**:
By default, the `deployment search` command just executes one query and returns the results. The command now also supports the `--all-matches` flag to query and return larger number of results that would exceed the maximum size of a single request. [#664](https://github.com/elastic/ecctl/pull/664)

### Fixes [ecctl-1.14.0-fixes]

* **Clear transients update plan**:
The update payload generated with `--generate-update-payload` used to include transient fields from the latest plan. These are now not included anymore by default. [#649](https://github.com/elastic/ecctl/pull/649)

### Changelog [ecctl-1.14.0-changelog]

* [f2fc756](https://github.com/elastic/ecctl/commit/f2fc756) Deployment search: Add flag to return all matches. ([#664](https://github.com/elastic/ecctl/pull/664))
* [0b631c8](https://github.com/elastic/ecctl/commit/0b631c8) fix(deps): update module golang.org/x/term to v0.22.0 ([#657](https://github.com/elastic/ecctl/pull/657))
* [92d35ea](https://github.com/elastic/ecctl/commit/92d35ea) chore: update the URL fragment for API keys ([#658](https://github.com/elastic/ecctl/pull/658))
* [eaefc6b](https://github.com/elastic/ecctl/commit/eaefc6b) fix(deps): update module github.com/elastic/cloud-sdk-go to v1.20.0 ([#656](https://github.com/elastic/ecctl/pull/656))
* [fd54678](https://github.com/elastic/ecctl/commit/fd54678) fix(deps): update module github.com/spf13/cobra to v1.8.1 ([#655](https://github.com/elastic/ecctl/pull/655))
* [81cceca](https://github.com/elastic/ecctl/commit/81cceca) fix(deps): update module github.com/spf13/viper to v1.19.0 ([#652](https://github.com/elastic/ecctl/pull/652))
* [c9e706f](https://github.com/elastic/ecctl/commit/c9e706f) chore(deps): update goreleaser/goreleaser-action action to v6 ([#654](https://github.com/elastic/ecctl/pull/654))
* [3584bdf](https://github.com/elastic/ecctl/commit/3584bdf) fix(deps): update module golang.org/x/term to v0.21.0 ([#653](https://github.com/elastic/ecctl/pull/653))
* [c47b9d5](https://github.com/elastic/ecctl/commit/c47b9d5) feat: expose clear_transients in ecctl ([#649](https://github.com/elastic/ecctl/pull/649))
* [54341aa](https://github.com/elastic/ecctl/commit/54341aa) fix(deps): update module golang.org/x/term to v0.20.0 ([#647](https://github.com/elastic/ecctl/pull/647))
* [96e6e1e](https://github.com/elastic/ecctl/commit/96e6e1e) fix(deps): update module github.com/elastic/cloud-sdk-go to v1.18.0 ([#648](https://github.com/elastic/ecctl/pull/648))
* [0ff35de](https://github.com/elastic/ecctl/commit/0ff35de) chore(deps): update actions/cache action to v4 ([#640](https://github.com/elastic/ecctl/pull/640))
* [e23d71e](https://github.com/elastic/ecctl/commit/e23d71e) chore(deps): update actions/checkout action to v4 ([#618](https://github.com/elastic/ecctl/pull/618))
* [e332de9](https://github.com/elastic/ecctl/commit/e332de9) fix(deps): update module github.com/asaskevich/govalidator to v0.0.0-20230301143203-a9d515a09cc2 ([#646](https://github.com/elastic/ecctl/pull/646))
* [54cfae1](https://github.com/elastic/ecctl/commit/54cfae1) chore(deps): update actions/setup-go action to v5 ([#632](https://github.com/elastic/ecctl/pull/632))
* [a5cc331](https://github.com/elastic/ecctl/commit/a5cc331) fix(deps): update module golang.org/x/term to v0.19.0 ([#642](https://github.com/elastic/ecctl/pull/642))