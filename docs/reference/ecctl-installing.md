---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-installing.html
---

# Installing [ecctl-installing]

The latest stable binaries can be found on the [release page](https://github.com/elastic/ecctl/releases) or compiled from the latest on the master branch to leverage the most recently merged features.

To learn more about building ecctl from the source, see the steps from our [Setting up a dev environment](https://github.com/elastic/ecctl/blob/master/CONTRIBUTING.md#setting-up-a-dev-environment).


## Install on macOS [ecctl-installing-macos]

The simplest installation for macOS users is to install ecctl with [Homebrew](https://brew.sh/):

```
$ brew tap elastic/tap
$ brew install elastic/tap/ecctl

Updating Homebrew...
==> Installing ecctl from elastic/tap
...
==> Caveats
To get autocompletions working make sure to run "source <(ecctl generate completions)".
If you prefer to add to your shell interpreter configuration file run, for bash or zsh respectively:
* `echo "source <(ecctl generate completions)" >> ~/.bash_profile`
* `echo "source <(ecctl generate completions)" >> ~/.zshrc`.
==> Summary
🍺  /usr/local/Cellar/ecctl/1.5.0: 5 files, 22.6MB, built in 4 seconds
```

::::{note}
To get autocompletions working, follow the instructions in the Homebrew output.
::::



## Upgrade on macOS [ecctl-upgrading-macos]

To upgrade ecctl via brew:

```
$ brew upgrade ecctl
```

