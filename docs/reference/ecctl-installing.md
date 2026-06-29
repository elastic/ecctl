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

The latest stable binaries and Linux packages can be found on the [release page](https://github.com/elastic/ecctl/releases). If you need to use changes that are not part of the latest release, you can build `ecctl` from source.

## Installation methods [ecctl-installing-methods]

You can install `ecctl` using one of these methods:

- **[macOS](#ecctl-installing-macos)**: Homebrew (recommended) or release binary
- **[Linux](#ecctl-installing-linux)**: `deb`, `rpm`, or release binary
- **[Windows](#ecctl-installing-windows)**: install with `go install`
- **[Any OS](#ecctl-installing-source)**: build from source

For shell completions, refer to [Enable shell completions](#ecctl-installing-completions).

## Install on macOS [ecctl-installing-macos]

### Homebrew (recommended) [ecctl-installing-macos-homebrew]

The simplest installation for macOS users is to install `ecctl` with [Homebrew](https://brew.sh/):

```bash
brew tap elastic/tap
brew install elastic/tap/ecctl
```

### Binary from GitHub release [ecctl-installing-macos-binary]

1. Download the archive for your architecture from the [release page](https://github.com/elastic/ecctl/releases), for example:
   - `ecctl_<VERSION>_darwin_amd64.tar.gz`
   - `ecctl_<VERSION>_darwin_arm64.tar.gz`
2. Extract the archive and move the `ecctl` binary to a directory in your `PATH`, for example `/usr/local/bin`.
3. Verify the installation:

```bash
ecctl version
```

## Install on Linux [ecctl-installing-linux]

### Install with deb/rpm packages [ecctl-installing-linux-packages]

Linux packages are published in every release:

- Debian/Ubuntu: `ecctl_<VERSION>_linux_64-bit.deb` or `ecctl_<VERSION>_linux_32-bit.deb`
- RHEL/CentOS/Fedora: `ecctl_<VERSION>_linux_64-bit.rpm` or `ecctl_<VERSION>_linux_32-bit.rpm`

Download the package that matches your system, then install it with your package manager.

Example (`deb`):

```bash
sudo dpkg -i ecctl_<VERSION>_linux_64-bit.deb
```

Example (`rpm`):

```bash
sudo rpm -i ecctl_<VERSION>_linux_64-bit.rpm
```

### Binary from GitHub release [ecctl-installing-linux-binary]

1. Download the archive for your architecture from the [release page](https://github.com/elastic/ecctl/releases), for example:
   - `ecctl_<VERSION>_linux_amd64.tar.gz`
   - `ecctl_<VERSION>_linux_arm64.tar.gz`
   - `ecctl_<VERSION>_linux_386.tar.gz`

2. Extract the archive and move the `ecctl` binary to a directory in your `PATH`, for example `/usr/local/bin`:

```bash
tar -xzf ecctl_<VERSION>_linux_amd64.tar.gz
sudo cp ecctl /usr/local/bin
```

3. Verify the installation:

```bash
ecctl version
```

## Install on Windows [ecctl-installing-windows]

Windows binaries are not currently published as part of the official release artifacts. Use `go install` to build and install `ecctl` from source:

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

5. Add Go's binary directory to your `PATH` if needed (`%USERPROFILE%\go\bin` by default).
6. Verify the installation:

```powershell
ecctl version
```

## Enable shell completions [ecctl-installing-completions]

You can enable shell completions after installation. This setup applies to Bash and Zsh on macOS and Linux.

Load completions for the current shell session:

```bash
source <(ecctl generate completions)
```

Persist completions across sessions:

- Bash: `echo "source <(ecctl generate completions)" >> ~/.bash_profile`
- Zsh: `echo "source <(ecctl generate completions)" >> ~/.zshrc`

## Build from source (all operating systems) [ecctl-installing-source]

Use this method if you need the latest unreleased changes or want to build from a specific branch or commit.

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

## Verify installation [ecctl-installing-verify]

After any installation method, verify that `ecctl` is available:

```bash
ecctl version
```

If the command is not found, check that the binary location is in your `PATH`.

## Troubleshooting [ecctl-installing-verify-troubleshooting]

- `ecctl: command not found`: add the binary location to your `PATH`.
- `permission denied` when copying to system paths: use elevated privileges (`sudo` on Linux/macOS).
- `go: command not found` on Windows: ensure Go is installed correctly and restart the terminal after installation.

## Upgrade [ecctl-upgrading]

Use the same channel you used to install `ecctl`:

- **Homebrew (macOS)**:
  ```bash
  brew upgrade ecctl
  ```
- **Linux package (`deb`/`rpm`)**: install a newer package from the [release page](https://github.com/elastic/ecctl/releases)
- **Release binary (macOS/Linux)**: replace the old binary with the new one
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

If you use serverless projects, refer to:

- [Manage serverless projects](/reference/ecctl_project.md)
