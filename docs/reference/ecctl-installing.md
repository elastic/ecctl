---
mapped_pages:
  - https://www.elastic.co/guide/en/ecctl/current/ecctl-installing.html
applies_to:
  serverless: ga
  deployment:
    ess: ga
    ece: ga all
---

# Installing [ecctl-installing]

This guide explains how to install the current {{version.ecctl}} release of `ecctl` on supported operating systems. Release artifacts are available from the [release page](https://github.com/elastic/ecctl/releases). If you need changes that are not yet included in the latest release, you can build `ecctl` from source.

## Installation methods [ecctl-installing-methods]

You can install `ecctl` using one of these methods:

- **[macOS](#ecctl-installing-macos)**: Homebrew (recommended) or archive (`.tar.gz`)
- **[Linux](#ecctl-installing-linux)**: `deb`, `rpm`, or archive (`.tar.gz`)
- **[Windows](#ecctl-installing-windows)**: Install with `go install`
- **[Any OS](#ecctl-installing-source)**: Build from source

For shell completions, refer to [Enable shell completions](#ecctl-installing-completions). If you download release artifacts (`.tar.gz`, `.deb`, or `.rpm`), refer to [Verify downloaded artifacts](#ecctl-installing-verify-checksums) to optionally verify their SHA-256 checksums before installation.

## Install on macOS [ecctl-installing-macos]

### Homebrew (recommended) [ecctl-installing-macos-homebrew]

The simplest installation for macOS users is to install `ecctl` with [Homebrew](https://brew.sh/):

```bash
brew tap elastic/tap
brew install elastic/tap/ecctl
```

### Install from archive (.tar.gz) [ecctl-installing-macos-archive]

1. Download the `.tar.gz` archive for your architecture from the [release page](https://github.com/elastic/ecctl/releases):

    ```bash subs=true
    # Apple Silicon
    curl -L -O https://download.elastic.co/downloads/ecctl/{{version.ecctl}}/ecctl_{{version.ecctl}}_darwin_arm64.tar.gz

    # Intel (x86_64)
    curl -L -O https://download.elastic.co/downloads/ecctl/{{version.ecctl}}/ecctl_{{version.ecctl}}_darwin_amd64.tar.gz

    ```

2. Extract the archive and move the `ecctl` binary to a directory in your `PATH`, for example `/usr/local/bin`.
3. Verify the installation:

```bash
ecctl version
```

## Install on Linux [ecctl-installing-linux]

### Install with deb/rpm packages (recommended) [ecctl-installing-linux-packages]

Linux packages are published with every release and are available from the [release page](https://github.com/elastic/ecctl/releases).

- Debian/Ubuntu: `ecctl_<VERSION>_linux_64-bit.deb`, `ecctl_<VERSION>_linux_32-bit.deb`, or `ecctl_<VERSION>_linux_arm64.deb`
- RHEL/CentOS/Fedora: `ecctl_<VERSION>_linux_64-bit.rpm`, `ecctl_<VERSION>_linux_32-bit.rpm`, or `ecctl_<VERSION>_linux_arm64.rpm`

Download the package that matches your system and architecture, then install it with your package manager. For example:

* On Debian/Ubuntu systems:

  ```bash subs=true
  curl -L -O https://download.elastic.co/downloads/ecctl/{{version.ecctl}}/ecctl_{{version.ecctl}}_linux_64-bit.deb
  sudo dpkg -i ecctl_{{version.ecctl}}_linux_64-bit.deb
  ```

* On RHEL/CentOS/Fedora systems:

  ```bash subs=true
  curl -L -O https://download.elastic.co/downloads/ecctl/{{version.ecctl}}/ecctl_{{version.ecctl}}_linux_64-bit.rpm
  sudo rpm -i ecctl_{{version.ecctl}}_linux_64-bit.rpm
  ```

### Install from archive (.tar.gz) [ecctl-installing-linux-archive]

1. Download the Linux `.tar.gz` archive for your architecture from the [release page](https://github.com/elastic/ecctl/releases):

    ```bash subs=true
    # x86_64 (AMD64)
    curl -L -O https://download.elastic.co/downloads/ecctl/{{version.ecctl}}/ecctl_{{version.ecctl}}_linux_amd64.tar.gz

    # ARM64
    curl -L -O https://download.elastic.co/downloads/ecctl/{{version.ecctl}}/ecctl_{{version.ecctl}}_linux_arm64.tar.gz

    # x86 (386)
    curl -L -O https://download.elastic.co/downloads/ecctl/{{version.ecctl}}/ecctl_{{version.ecctl}}_linux_386.tar.gz
    ```

2. Extract the archive and move the `ecctl` binary to a directory in your `PATH`, for example `/usr/local/bin`:

    ```bash subs=true
    tar -xzf ecctl_{{version.ecctl}}_linux_amd64.tar.gz
    sudo cp ecctl /usr/local/bin
    ```

3. Verify the installation:

    ```bash
    ecctl version
    ```

## Install on Windows [ecctl-installing-windows]

Official Windows binaries are not currently published. Install `ecctl` with `go install`:

1. Install [Go](https://go.dev/dl/).
2. Open **Command Prompt** or **PowerShell**.
3. Confirm Go is available:

    ```powershell
    go version
    ```

4. Install `ecctl`:

    ```powershell
    go install github.com/elastic/ecctl@latest
    ```

    To install a different version, replace `@latest` with a specific tag, for example:

    ```powershell subs=true
    go install github.com/elastic/ecctl@v{{version.ecctl}}
    ```

5. Add Go's binary directory to your `PATH` if needed (`%USERPROFILE%\go\bin` by default).
6. Verify the installation:

    ```powershell
    ecctl version
    ```

## Build from source (all operating systems) [ecctl-installing-source]

Use this method if you need changes that are not yet included in an official release or want to build from a specific branch or commit.

### Prerequisites [ecctl-installing-source-prerequisites]

- [Go](https://go.dev/doc/install)
- `git`

### Build steps [ecctl-installing-source-steps]

Clone the repository and build locally:

```bash
git clone https://github.com/elastic/ecctl.git
cd ecctl
go build .
```

This command produces an `ecctl` binary (`ecctl.exe` on Windows) in the current directory.

You can also install directly with Go:

```bash
go install .
```

This installs the binary into `GOBIN` (or `GOPATH/bin` if `GOBIN` is not set). Make sure that location is in your `PATH`.

For contributor-focused setup details, see [Setting up a dev environment](https://github.com/elastic/ecctl/blob/master/CONTRIBUTING.md#setting-up-a-dev-environment).

## Verify downloaded artifacts (optional) [ecctl-installing-verify-checksums]

If you download release artifacts (`.tar.gz`, `.deb`, or `.rpm`), you can optionally verify their SHA-256 checksums before installing them.

1. Download the checksum file:

    ```bash subs=true
    curl -L -O https://download.elastic.co/downloads/ecctl/{{version.ecctl}}/ecctl_{{version.ecctl}}_checksums.txt
    ```

2. Compute the SHA-256 checksum of the downloaded artifact. For example:

    ```bash subs=true
    # macOS
    shasum -a 256 ecctl_{{version.ecctl}}_darwin_arm64.tar.gz

    # Linux
    sha256sum ecctl_{{version.ecctl}}_linux_amd64.tar.gz
    ```

3. Display the expected checksum from the checksums file. For example:

    ```bash subs=true
    grep "ecctl_{{version.ecctl}}_linux_amd64.tar.gz" ecctl_{{version.ecctl}}_checksums.txt
    ```

    Verify that the checksum matches the value returned by the checksum command.

## Enable shell completions [ecctl-installing-completions]

You can enable shell completions after installation. These instructions apply to Bash and Zsh on macOS and Linux.

Load completions for the current shell session:

```bash
source <(ecctl generate completions)
```

Persist completions across sessions:

- Bash: `echo "source <(ecctl generate completions)" >> ~/.bash_profile`
- Zsh: `echo "source <(ecctl generate completions)" >> ~/.zshrc`

## Troubleshooting [ecctl-installing-verify-troubleshooting]

- `ecctl: command not found`: Add the binary location to your `PATH`.
- `permission denied` when copying to system paths: Use elevated privileges (`sudo` on Linux/macOS).
- `go: command not found` on Windows: Ensure Go is installed correctly and restart the terminal after installation.

## Upgrade [ecctl-upgrading]

Use the same installation method you used to install `ecctl`:

- **Homebrew (macOS)**:
  ```bash
  brew upgrade ecctl
  ```
- **Linux package (`deb`/`rpm`)**: Install a newer package from the [release page](https://github.com/elastic/ecctl/releases)
- **Archive (`.tar.gz`) (macOS/Linux)**: Download the latest archive from the [release page](https://github.com/elastic/ecctl/releases) and replace the existing binary.
- **Windows (`go install`)**:
  ```powershell
  go install github.com/elastic/ecctl@latest
  ```

## Next steps [ecctl-installing-next-steps]

After installing `ecctl`, continue with:

- [Configure ecctl](/reference/ecctl-configuring.md)
- [Initialize your first configuration](/reference/ecctl_init.md)
- [Authentication methods](/reference/ecctl-authentication.md)
- [Usage examples](/reference/ecctl-examples.md)
- [Command reference](/reference/ecctl.md)

If you use serverless projects, refer to [Manage serverless projects](/reference/ecctl_project.md)
